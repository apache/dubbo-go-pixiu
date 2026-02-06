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

package pkg

import (
	"context"
	"fmt"
)

import (
	"github.com/dubbogo/gost/log/logger"

	perrors "github.com/pkg/errors"
)

// UserProvider implements the Dubbo protocol benchmark service with Hessian serialization
type UserProvider struct {
}

func (u *UserProvider) getUser(userID string) (*User, error) {
	if user, ok := userMap[userID]; ok {
		return &user, nil
	}
	return nil, fmt.Errorf("invalid user id: %s", userID)
}

func (u *UserProvider) GetUser(ctx context.Context, req *User) (*User, error) {
	logger.Infof("Dubbo UserProvider GetUser, req: %#v", req)
	user, err := u.getUser(req.ID)
	if err == nil {
		logger.Infof("Dubbo UserProvider GetUser, rsp: %#v", user)
	}
	return user, err
}

func (u *UserProvider) GetUser0(id string, name string) (User, error) {
	logger.Infof("Dubbo UserProvider GetUser0, id: %s, name: %s", id, name)
	user, err := u.getUser(id)
	if err != nil {
		return User{}, err
	}
	if user.Name != name {
		return User{}, perrors.New("name is not " + user.Name)
	}
	return *user, nil
}

func (u *UserProvider) GetUsers(req []string) ([]*User, error) {
	logger.Infof("Dubbo UserProvider GetUsers, req: %v", req)
	users := make([]*User, 0, len(req))
	for _, id := range req {
		user, err := u.getUser(id)
		if err != nil {
			continue
		}
		users = append(users, user)
	}
	return users, nil
}

func (u *UserProvider) GetGender(i int32) (Gender, error) {
	if i == 1 {
		return Gender(WOMAN), nil
	}
	return Gender(MAN), nil
}
