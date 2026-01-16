/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package store

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

import (
	"github.com/pkg/errors"

	"golang.org/x/crypto/bcrypt"

	"gorm.io/driver/mysql"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

import (
	"github.com/apache/dubbo-go-pixiu/admin/internal/model"
	"github.com/apache/dubbo-go-pixiu/admin/pkg/config"
)

// MySQL handles all MySQL database operations using GORM.
type MySQL struct {
	db *gorm.DB
}

// NewMySQL creates a new MySQL store with GORM.
func NewMySQL(cfg config.MySQLConfig) (*MySQL, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to mysql")
	}

	// Auto migrate tables
	if err := db.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.UserRole{},
		&model.Permission{},
		&model.RolePermission{},
	); err != nil {
		return nil, errors.Wrap(err, "failed to migrate database")
	}

	m := &MySQL{db: db}

	// Initialize default data
	if err := m.initDefaultData(); err != nil {
		return nil, errors.Wrap(err, "failed to initialize default data")
	}

	// Initialize default permissions
	if err := m.InitDefaultPermissions(); err != nil {
		return nil, errors.Wrap(err, "failed to initialize default permissions")
	}

	return m, nil
}

// initDefaultData creates default admin user and roles if they don't exist.
func (m *MySQL) initDefaultData() error {
	// Create default admin role
	adminRole := model.Role{ID: 1, RoleName: "admin", Description: "Administrator"}
	if err := m.db.FirstOrCreate(&adminRole, model.Role{ID: 1}).Error; err != nil {
		return err
	}

	// Create default user role
	userRole := model.Role{ID: 2, RoleName: "user", Description: "Normal User"}
	if err := m.db.FirstOrCreate(&userRole, model.Role{ID: 2}).Error; err != nil {
		return err
	}

	// Create default admin user if not exists (password: admin)
	// Note: Only create if not exists, don't reset password on every startup
	var adminUser model.User
	result := m.db.Where("username = ?", "admin").First(&adminUser)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Create new admin user with bcrypt password
			hashedPassword, err := hashPassword("admin")
			if err != nil {
				return errors.Wrap(err, "failed to hash default admin password")
			}
			adminUser = model.User{
				Username: "admin",
				Password: hashedPassword,
				Role:     1,
				Enabled:  true,
			}
			if err := m.db.Create(&adminUser).Error; err != nil {
				return err
			}
		} else {
			return result.Error
		}
	}
	// Note: Removed automatic password reset on startup for security

	// Assign admin role to admin user
	userRoleAssign := model.UserRole{UserID: adminUser.ID, RoleID: 1}
	if err := m.db.FirstOrCreate(&userRoleAssign, model.UserRole{UserID: adminUser.ID}).Error; err != nil {
		return err
	}

	return nil
}

// Close closes the database connection.
func (m *MySQL) Close() error {
	sqlDB, err := m.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// DB returns the underlying GORM DB instance.
func (m *MySQL) DB() *gorm.DB {
	return m.db
}

// --- User Operations ---

// Login validates user credentials and returns user ID if successful.
// Supports automatic migration from legacy MD5 hashes to bcrypt.
func (m *MySQL) Login(username, password string) (uint, error) {
	var user model.User
	err := m.db.Where("username = ?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, errors.New("invalid username or password")
	}
	if err != nil {
		return 0, errors.Wrap(err, "failed to query user")
	}

	if !user.Enabled {
		return 0, errors.New("user is disabled")
	}

	// Check password
	if !checkPassword(password, user.Password) {
		// If it's a legacy MD5 hash, try MD5 comparison for migration
		if isLegacyMD5Hash(user.Password) {
			if !checkLegacyMD5Password(password, user.Password) {
				return 0, errors.New("invalid username or password")
			}
			// Migrate to bcrypt (ignore error - password will be upgraded on next login)
			_ = m.upgradePasswordHash(user.ID, password)
		} else {
			return 0, errors.New("invalid username or password")
		}
	}

	return user.ID, nil
}

// checkLegacyMD5Password checks password against legacy MD5 hash.
func checkLegacyMD5Password(password, hash string) bool {
	md5Hash := md5.Sum([]byte(password))
	return hex.EncodeToString(md5Hash[:]) == hash
}

// upgradePasswordHash upgrades a user's password hash from MD5 to bcrypt.
func (m *MySQL) upgradePasswordHash(userID uint, password string) error {
	newHash, err := hashPassword(password)
	if err != nil {
		return err
	}
	return m.db.Model(&model.User{}).Where("id = ?", userID).Update("password", newHash).Error
}

// Register creates a new user.
func (m *MySQL) Register(username, password string) error {
	if username == "" {
		return errors.New("username cannot be empty")
	}

	// Check if user already exists
	var count int64
	m.db.Model(&model.User{}).Where("username = ?", username).Count(&count)
	if count > 0 {
		return errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return errors.Wrap(err, "failed to hash password")
	}

	// Create user
	user := model.User{
		Username: username,
		Password: hashedPassword,
		Role:     0,
		Enabled:  true,
	}
	if err := m.db.Create(&user).Error; err != nil {
		return errors.Wrap(err, "failed to create user")
	}

	// Assign default user role
	userRole := model.UserRole{UserID: user.ID, RoleID: 2}
	if err := m.db.Create(&userRole).Error; err != nil {
		return errors.Wrap(err, "failed to assign user role")
	}

	return nil
}

// GetUserInfo returns user information by username.
func (m *MySQL) GetUserInfo(username string) (*model.UserInfo, error) {
	var user model.User
	err := m.db.Where("username = ?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user info")
	}

	return &model.UserInfo{
		ID:       user.ID,
		Username: user.Username,
		Role:     user.Role,
	}, nil
}

// GetUserRole returns the user's role information.
func (m *MySQL) GetUserRole(username string) (*model.Role, error) {
	var role model.Role
	err := m.db.
		Joins("JOIN pixiu_user_role ur ON pixiu_role.id = ur.role_id").
		Joins("JOIN pixiu_user u ON u.id = ur.user_id").
		Where("u.username = ?", username).
		First(&role).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("role not found")
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user role")
	}

	return &role, nil
}

// IsAdmin checks if the user is an administrator.
func (m *MySQL) IsAdmin(username string) (bool, error) {
	var count int64
	err := m.db.Model(&model.User{}).
		Where("username = ? AND role = 1", username).
		Count(&count).Error

	if err != nil {
		return false, errors.Wrap(err, "failed to check admin status")
	}

	return count > 0, nil
}

// ChangePassword updates the user's password.
func (m *MySQL) ChangePassword(username, oldPassword, newPassword string) error {
	// Get user first
	var user model.User
	err := m.db.Where("username = ?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("user not found")
	}
	if err != nil {
		return errors.Wrap(err, "failed to query user")
	}

	// Verify old password (support both bcrypt and legacy MD5)
	passwordValid := checkPassword(oldPassword, user.Password)
	if !passwordValid && isLegacyMD5Hash(user.Password) {
		passwordValid = checkLegacyMD5Password(oldPassword, user.Password)
	}
	if !passwordValid {
		return errors.New("invalid old password")
	}

	// Hash new password with bcrypt
	newHashed, err := hashPassword(newPassword)
	if err != nil {
		return errors.Wrap(err, "failed to hash new password")
	}

	// Update password
	result := m.db.Model(&model.User{}).
		Where("id = ?", user.ID).
		Update("password", newHashed)

	if result.Error != nil {
		return errors.Wrap(result.Error, "failed to update password")
	}

	return nil
}

// hashPassword hashes a password using bcrypt.
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.Wrap(err, "failed to hash password")
	}
	return string(bytes), nil
}

// checkPassword verifies a password against a hash.
// Supports both bcrypt (new) and legacy MD5 hashes for migration.
func checkPassword(password, hash string) bool {
	// bcrypt hashes are always 60 characters
	if len(hash) == 60 {
		err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
		return err == nil
	}
	// Legacy MD5 hash (32 characters hex) - should not happen after migration
	// but kept for backwards compatibility
	return false
}

// isLegacyMD5Hash checks if the hash is a legacy MD5 hash.
func isLegacyMD5Hash(hash string) bool {
	// MD5 produces 32 character hex string, bcrypt produces 60 character string starting with $2
	return len(hash) == 32
}

// --- Permission Operations ---

// HasPermission checks if a user has permission to perform an action on a resource.
func (m *MySQL) HasPermission(username, resource, action string) (bool, error) {
	// Admin users have all permissions
	isAdmin, err := m.IsAdmin(username)
	if err != nil {
		return false, err
	}
	if isAdmin {
		return true, nil
	}

	// Check specific permission through role
	var count int64
	err = m.db.Table("pixiu_permission p").
		Joins("JOIN pixiu_role_permission rp ON p.id = rp.permission_id").
		Joins("JOIN pixiu_user_role ur ON rp.role_id = ur.role_id").
		Joins("JOIN pixiu_user u ON ur.user_id = u.id").
		Where("u.username = ? AND p.resource = ? AND p.action = ?", username, resource, action).
		Count(&count).Error

	if err != nil {
		return false, errors.Wrap(err, "failed to check permission")
	}

	return count > 0, nil
}

// InitDefaultPermissions initializes default permissions and assigns them to admin role.
func (m *MySQL) InitDefaultPermissions() error {
	// Create default permissions
	for _, perm := range model.DefaultPermissions {
		var existing model.Permission
		result := m.db.Where("resource = ? AND action = ?", perm.Resource, perm.Action).First(&existing)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			if err := m.db.Create(&perm).Error; err != nil {
				return errors.Wrap(err, "failed to create permission")
			}
		}
	}

	// Assign all permissions to admin role (role_id = 1)
	var permissions []model.Permission
	if err := m.db.Find(&permissions).Error; err != nil {
		return errors.Wrap(err, "failed to fetch permissions")
	}

	for _, perm := range permissions {
		rp := model.RolePermission{RoleID: 1, PermissionID: perm.ID}
		m.db.FirstOrCreate(&rp, model.RolePermission{RoleID: 1, PermissionID: perm.ID})
	}

	// Assign read permissions to user role (role_id = 2)
	var readPermissions []model.Permission
	if err := m.db.Where("action = ?", model.ActionRead).Find(&readPermissions).Error; err != nil {
		return errors.Wrap(err, "failed to fetch read permissions")
	}

	for _, perm := range readPermissions {
		rp := model.RolePermission{RoleID: 2, PermissionID: perm.ID}
		m.db.FirstOrCreate(&rp, model.RolePermission{RoleID: 2, PermissionID: perm.ID})
	}

	return nil
}

// --- User Management Operations ---

// ListUsers returns all users with pagination.
func (m *MySQL) ListUsers(page, pageSize int) (*model.UserListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	var total int64
	if err := m.db.Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, errors.Wrap(err, "failed to count users")
	}

	var users []model.User
	if err := m.db.Order("id ASC").Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, errors.Wrap(err, "failed to list users")
	}

	// Convert to UserInfo to hide password
	items := make([]model.UserListItem, len(users))
	for i, u := range users {
		items[i] = model.UserListItem{
			ID:          u.ID,
			Username:    u.Username,
			Role:        u.Role,
			Enabled:     u.Enabled,
			DateCreated: u.DateCreated,
			DateUpdated: u.DateUpdated,
		}
	}

	return &model.UserListResponse{
		Items:    items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// GetUserByID returns user by ID.
func (m *MySQL) GetUserByID(id uint) (*model.UserListItem, error) {
	var user model.User
	if err := m.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, errors.Wrap(err, "failed to get user")
	}

	return &model.UserListItem{
		ID:          user.ID,
		Username:    user.Username,
		Role:        user.Role,
		Enabled:     user.Enabled,
		DateCreated: user.DateCreated,
		DateUpdated: user.DateUpdated,
	}, nil
}

// UpdateUser updates user information.
func (m *MySQL) UpdateUser(id uint, req *model.UpdateUserRequest) error {
	updates := make(map[string]any)

	if req.Role != nil {
		updates["role"] = *req.Role
	}
	if req.Enabled != nil {
		updates["enabled"] = *req.Enabled
	}

	if len(updates) == 0 {
		return nil
	}

	result := m.db.Model(&model.User{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return errors.Wrap(result.Error, "failed to update user")
	}
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

// DeleteUser deletes a user by ID.
func (m *MySQL) DeleteUser(id uint) error {
	// Prevent deleting admin user (ID=1)
	if id == 1 {
		return errors.New("cannot delete admin user")
	}

	// Delete user roles first
	if err := m.db.Where("user_id = ?", id).Delete(&model.UserRole{}).Error; err != nil {
		return errors.Wrap(err, "failed to delete user roles")
	}

	// Delete user
	result := m.db.Delete(&model.User{}, id)
	if result.Error != nil {
		return errors.Wrap(result.Error, "failed to delete user")
	}
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

// ResetUserPassword resets a user's password (admin operation).
func (m *MySQL) ResetUserPassword(id uint, newPassword string) error {
	hashedPassword, err := hashPassword(newPassword)
	if err != nil {
		return errors.Wrap(err, "failed to hash password")
	}

	result := m.db.Model(&model.User{}).Where("id = ?", id).Update("password", hashedPassword)
	if result.Error != nil {
		return errors.Wrap(result.Error, "failed to reset password")
	}
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

// --- Role Management Operations ---

// ListRoles returns all roles.
func (m *MySQL) ListRoles() ([]model.Role, error) {
	var roles []model.Role
	if err := m.db.Order("id ASC").Find(&roles).Error; err != nil {
		return nil, errors.Wrap(err, "failed to list roles")
	}
	return roles, nil
}

// GetRole returns a role by ID.
func (m *MySQL) GetRole(id uint) (*model.Role, error) {
	var role model.Role
	if err := m.db.First(&role, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("role not found")
		}
		return nil, errors.Wrap(err, "failed to get role")
	}
	return &role, nil
}

// AssignUserRole assigns a role to a user.
func (m *MySQL) AssignUserRole(userID, roleID uint) error {
	// Check if user exists
	var user model.User
	if err := m.db.First(&user, userID).Error; err != nil {
		return errors.New("user not found")
	}

	// Check if role exists
	var role model.Role
	if err := m.db.First(&role, roleID).Error; err != nil {
		return errors.New("role not found")
	}

	// Update user's role field
	if err := m.db.Model(&model.User{}).Where("id = ?", userID).Update("role", roleID).Error; err != nil {
		return errors.Wrap(err, "failed to update user role")
	}

	// Update user_role table
	userRole := model.UserRole{UserID: userID, RoleID: roleID}
	if err := m.db.Where("user_id = ?", userID).Delete(&model.UserRole{}).Error; err != nil {
		return errors.Wrap(err, "failed to clear old user role")
	}
	if err := m.db.Create(&userRole).Error; err != nil {
		return errors.Wrap(err, "failed to assign user role")
	}

	return nil
}

// CreateRole creates a new role.
func (m *MySQL) CreateRole(role *model.Role) error {
	if err := m.db.Create(role).Error; err != nil {
		return errors.Wrap(err, "failed to create role")
	}
	return nil
}

// UpdateRole updates a role.
func (m *MySQL) UpdateRole(id uint, roleName, description string) error {
	result := m.db.Model(&model.Role{}).Where("id = ?", id).Updates(map[string]any{
		"role_name":   roleName,
		"description": description,
	})
	if result.Error != nil {
		return errors.Wrap(result.Error, "failed to update role")
	}
	if result.RowsAffected == 0 {
		return errors.New("role not found")
	}
	return nil
}

// DeleteRole deletes a role by ID.
func (m *MySQL) DeleteRole(id uint) error {
	// Check if any users have this role
	var count int64
	if err := m.db.Model(&model.User{}).Where("role = ?", id).Count(&count).Error; err != nil {
		return errors.Wrap(err, "failed to check role usage")
	}
	if count > 0 {
		return errors.New("cannot delete role: role is assigned to users")
	}

	// Delete role permissions first
	if err := m.db.Where("role_id = ?", id).Delete(&model.RolePermission{}).Error; err != nil {
		return errors.Wrap(err, "failed to delete role permissions")
	}

	// Delete the role
	result := m.db.Delete(&model.Role{}, id)
	if result.Error != nil {
		return errors.Wrap(result.Error, "failed to delete role")
	}
	if result.RowsAffected == 0 {
		return errors.New("role not found")
	}
	return nil
}

// --- Permission Management Operations ---

// ListPermissions returns all permissions.
func (m *MySQL) ListPermissions() ([]model.Permission, error) {
	var permissions []model.Permission
	if err := m.db.Order("resource ASC, action ASC").Find(&permissions).Error; err != nil {
		return nil, errors.Wrap(err, "failed to list permissions")
	}
	return permissions, nil
}

// GetRolePermissions returns permissions for a role.
func (m *MySQL) GetRolePermissions(roleID uint) ([]model.Permission, error) {
	var permissions []model.Permission
	err := m.db.Joins("JOIN pixiu_role_permission rp ON rp.permission_id = pixiu_permission.id").
		Where("rp.role_id = ?", roleID).
		Find(&permissions).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to get role permissions")
	}
	return permissions, nil
}

// UpdateRolePermissions updates the permissions for a role.
func (m *MySQL) UpdateRolePermissions(roleID uint, permissionIDs []uint) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		// Delete existing role permissions
		if err := tx.Where("role_id = ?", roleID).Delete(&model.RolePermission{}).Error; err != nil {
			return errors.Wrap(err, "failed to clear role permissions")
		}

		// Insert new role permissions
		for _, permID := range permissionIDs {
			rp := model.RolePermission{RoleID: roleID, PermissionID: permID}
			if err := tx.Create(&rp).Error; err != nil {
				return errors.Wrap(err, "failed to assign permission to role")
			}
		}
		return nil
	})
}

// CreateUserByAdmin creates a new user (admin operation).
func (m *MySQL) CreateUserByAdmin(username, password string, roleID uint) (*model.User, error) {
	// Check if username exists
	var existingUser model.User
	if err := m.db.Where("username = ?", username).First(&existingUser).Error; err == nil {
		return nil, errors.New("username already exists")
	}

	// Hash password
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, errors.Wrap(err, "failed to hash password")
	}

	user := &model.User{
		Username: username,
		Password: hashedPassword,
		Role:     int(roleID),
		Enabled:  true,
	}

	if err := m.db.Create(user).Error; err != nil {
		return nil, errors.Wrap(err, "failed to create user")
	}

	// Assign role
	if roleID > 0 {
		userRole := model.UserRole{UserID: user.ID, RoleID: roleID}
		if err := m.db.Create(&userRole).Error; err != nil {
			return nil, errors.Wrap(err, "failed to assign role to user")
		}
	}

	return user, nil
}
