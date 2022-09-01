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

package mocks

import (
	reflect "reflect"
)

import (
	gomock "github.com/golang/mock/gomock"
)

import (
	model "github.com/apache/dubbo-go-pixiu/pixiu/pkg/model"
	controls "github.com/apache/dubbo-go-pixiu/pixiu/pkg/server/controls"
)

// MockClusterManager is a mock of ClusterManager interface.
type MockClusterManager struct {
	ctrl     *gomock.Controller
	recorder *MockClusterManagerMockRecorder
}

// MockClusterManagerMockRecorder is the mock recorder for MockClusterManager.
type MockClusterManagerMockRecorder struct {
	mock *MockClusterManager
}

// NewMockClusterManager creates a new mock instance.
func NewMockClusterManager(ctrl *gomock.Controller) *MockClusterManager {
	mock := &MockClusterManager{ctrl: ctrl}
	mock.recorder = &MockClusterManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClusterManager) EXPECT() *MockClusterManagerMockRecorder {
	return m.recorder
}

// AddCluster mocks base method.
func (m *MockClusterManager) AddCluster(cluster *model.ClusterConfig) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddCluster", cluster)
}

// AddCluster indicates an expected call of AddCluster.
func (mr *MockClusterManagerMockRecorder) AddCluster(cluster interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCluster", reflect.TypeOf((*MockClusterManager)(nil).AddCluster), cluster)
}

// CloneXdsControlStore mocks base method.
func (m *MockClusterManager) CloneXdsControlStore() (controls.ClusterStore, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloneXdsControlStore")
	ret0, _ := ret[0].(controls.ClusterStore)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CloneXdsControlStore indicates an expected call of CloneXdsControlStore.
func (mr *MockClusterManagerMockRecorder) CloneXdsControlStore() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloneXdsControlStore", reflect.TypeOf((*MockClusterManager)(nil).CloneXdsControlStore))
}

// HasCluster mocks base method.
func (m *MockClusterManager) HasCluster(name string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HasCluster", name)
	ret0, _ := ret[0].(bool)
	return ret0
}

// HasCluster indicates an expected call of HasCluster.
func (mr *MockClusterManagerMockRecorder) HasCluster(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasCluster", reflect.TypeOf((*MockClusterManager)(nil).HasCluster), name)
}

// RemoveCluster mocks base method.
func (m *MockClusterManager) RemoveCluster(names []string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RemoveCluster", names)
}

// RemoveCluster indicates an expected call of RemoveCluster.
func (mr *MockClusterManagerMockRecorder) RemoveCluster(names interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveCluster", reflect.TypeOf((*MockClusterManager)(nil).RemoveCluster), names)
}

// UpdateCluster mocks base method.
func (m *MockClusterManager) UpdateCluster(cluster *model.ClusterConfig) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UpdateCluster", cluster)
}

// UpdateCluster indicates an expected call of UpdateCluster.
func (mr *MockClusterManagerMockRecorder) UpdateCluster(cluster interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCluster", reflect.TypeOf((*MockClusterManager)(nil).UpdateCluster), cluster)
}

// MockListenerManager is a mock of ListenerManager interface.
type MockListenerManager struct {
	ctrl     *gomock.Controller
	recorder *MockListenerManagerMockRecorder
}

// MockListenerManagerMockRecorder is the mock recorder for MockListenerManager.
type MockListenerManagerMockRecorder struct {
	mock *MockListenerManager
}

// NewMockListenerManager creates a new mock instance.
func NewMockListenerManager(ctrl *gomock.Controller) *MockListenerManager {
	mock := &MockListenerManager{ctrl: ctrl}
	mock.recorder = &MockListenerManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockListenerManager) EXPECT() *MockListenerManagerMockRecorder {
	return m.recorder
}

// AddOrUpdateListener mocks base method.
func (m_2 *MockListenerManager) AddListener(m *model.Listener) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "AddListener", m)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddOrUpdateListener indicates an expected call of AddOrUpdateListener.
func (mr *MockListenerManagerMockRecorder) AddOrUpdateListener(m interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddListener", reflect.TypeOf((*MockListenerManager)(nil).AddListener), m)
}

// RemoveListener mocks base method.
func (m *MockListenerManager) RemoveListener(names []string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RemoveListener", names)
}

// RemoveListener indicates an expected call of RemoveListener.
func (mr *MockListenerManagerMockRecorder) RemoveListener(names interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveListener", reflect.TypeOf((*MockListenerManager)(nil).RemoveListener), names)
}

// MockDynamicResourceManager is a mock of DynamicResourceManager interface.
type MockDynamicResourceManager struct {
	ctrl     *gomock.Controller
	recorder *MockDynamicResourceManagerMockRecorder
}

// MockDynamicResourceManagerMockRecorder is the mock recorder for MockDynamicResourceManager.
type MockDynamicResourceManagerMockRecorder struct {
	mock *MockDynamicResourceManager
}

// NewMockDynamicResourceManager creates a new mock instance.
func NewMockDynamicResourceManager(ctrl *gomock.Controller) *MockDynamicResourceManager {
	mock := &MockDynamicResourceManager{ctrl: ctrl}
	mock.recorder = &MockDynamicResourceManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDynamicResourceManager) EXPECT() *MockDynamicResourceManagerMockRecorder {
	return m.recorder
}

// GetCds mocks base method.
func (m *MockDynamicResourceManager) GetCds() *model.ApiConfigSource {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCds")
	ret0, _ := ret[0].(*model.ApiConfigSource)
	return ret0
}

// GetCds indicates an expected call of GetCds.
func (mr *MockDynamicResourceManagerMockRecorder) GetCds() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCds", reflect.TypeOf((*MockDynamicResourceManager)(nil).GetCds))
}

// GetLds mocks base method.
func (m *MockDynamicResourceManager) GetLds() *model.ApiConfigSource {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLds")
	ret0, _ := ret[0].(*model.ApiConfigSource)
	return ret0
}

// GetLds indicates an expected call of GetLds.
func (mr *MockDynamicResourceManagerMockRecorder) GetLds() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLds", reflect.TypeOf((*MockDynamicResourceManager)(nil).GetLds))
}

// GetNode mocks base method.
func (m *MockDynamicResourceManager) GetNode() *model.Node {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNode")
	ret0, _ := ret[0].(*model.Node)
	return ret0
}

// GetNode indicates an expected call of GetNode.
func (mr *MockDynamicResourceManagerMockRecorder) GetNode() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNode", reflect.TypeOf((*MockDynamicResourceManager)(nil).GetNode))
}

// MockClusterStore is a mock of ClusterStore interface.
type MockClusterStore struct {
	ctrl     *gomock.Controller
	recorder *MockClusterStoreMockRecorder
}

// MockClusterStoreMockRecorder is the mock recorder for MockClusterStore.
type MockClusterStoreMockRecorder struct {
	mock *MockClusterStore
}

// NewMockClusterStore creates a new mock instance.
func NewMockClusterStore(ctrl *gomock.Controller) *MockClusterStore {
	mock := &MockClusterStore{ctrl: ctrl}
	mock.recorder = &MockClusterStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClusterStore) EXPECT() *MockClusterStoreMockRecorder {
	return m.recorder
}

// Config mocks base method.
func (m *MockClusterStore) Config() []*model.ClusterConfig {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Config")
	ret0, _ := ret[0].([]*model.ClusterConfig)
	return ret0
}

// Config indicates an expected call of Config.
func (mr *MockClusterStoreMockRecorder) Config() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Config", reflect.TypeOf((*MockClusterStore)(nil).Config))
}

// HasCluster mocks base method.
func (m *MockClusterStore) HasCluster(name string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HasCluster", name)
	ret0, _ := ret[0].(bool)
	return ret0
}

// HasCluster indicates an expected call of HasCluster.
func (mr *MockClusterStoreMockRecorder) HasCluster(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasCluster", reflect.TypeOf((*MockClusterStore)(nil).HasCluster), name)
}
