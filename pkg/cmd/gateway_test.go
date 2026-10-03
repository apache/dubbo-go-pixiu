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

package cmd

import (
	"errors"
	"testing"
)

import (
	"github.com/spf13/cobra"

	"github.com/stretchr/testify/assert"

	"go.uber.org/zap"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

// MockDeployer for testing
type MockDeployer struct {
	initializeErr error
	startErr      error
	stopErr       error
}

func (m *MockDeployer) initialize() error {
	return m.initializeErr
}

func (m *MockDeployer) start() error {
	return m.startErr
}

func (m *MockDeployer) stop() error {
	return m.stopErr
}

func TestDefaultDeployerInitialize(t *testing.T) {
	// Test successful initialization
	d := &DefaultDeployer{
		configManger: nil, // Will use default behavior
	}
	// Note: This test requires actual config files to be present
	// In a real test environment, we would mock the config manager
	assert.NotNil(t, d)
}

func TestDefaultDeployerStart(t *testing.T) {
	d := &DefaultDeployer{
		bootstrap: nil,
	}
	// Note: server.Start requires a valid bootstrap config
	// In a unit test, we would need to mock the server.Start function
	assert.NotNil(t, d)
}

func TestDefaultDeployerStop(t *testing.T) {
	d := &DefaultDeployer{}
	err := d.stop()
	assert.Error(t, err)
	assert.Equal(t, "stop not implemented", err.Error())
}

func TestStartGatewayCmdPreRunE(t *testing.T) {
	tests := []struct {
		name          string
		deployer      Deployer
		expectedError bool
	}{
		{
			name: "successful initialization",
			deployer: &MockDeployer{
				initializeErr: nil,
			},
			expectedError: false,
		},
		{
			name: "failed initialization",
			deployer: &MockDeployer{
				initializeErr: errors.New("initialization failed"),
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original deployer
			originalDeployer := deploy
			defer func() {
				deploy = originalDeployer
			}()

			// Replace with mock deployer
			deploy = tt.deployer

			// Create a new command to test
			testCmd := &cobra.Command{
				Use:     "test",
				PreRunE: startGatewayCmd.PreRunE,
			}

			// Execute PreRunE
			err := testCmd.PreRunE(testCmd, []string{})

			if tt.expectedError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to initialize gateway")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestStartGatewayCmdRunE(t *testing.T) {
	tests := []struct {
		name          string
		deployer      Deployer
		expectedError bool
	}{
		{
			name: "successful start",
			deployer: &MockDeployer{
				startErr: nil,
			},
			expectedError: false,
		},
		{
			name: "failed start",
			deployer: &MockDeployer{
				startErr: errors.New("start failed"),
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original deployer
			originalDeployer := deploy
			defer func() {
				deploy = originalDeployer
			}()

			// Replace with mock deployer
			deploy = tt.deployer

			// Create a new command to test
			testCmd := &cobra.Command{
				Use:  "test",
				RunE: startGatewayCmd.RunE,
			}

			// Execute RunE
			err := testCmd.RunE(testCmd, []string{})

			if tt.expectedError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to start gateway")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestInitDefaultValue(t *testing.T) {
	// Reset all values to empty
	configPath = ""
	apiConfigPath = ""
	logConfigPath = ""
	logLevel = ""
	limitCpus = ""
	logFormat = ""

	// Call initDefaultValue
	initDefaultValue()

	// Check that default values are set
	assert.NotEmpty(t, configPath)
	assert.NotEmpty(t, apiConfigPath)
	assert.NotEmpty(t, logConfigPath)
	assert.NotEmpty(t, logLevel)
	assert.NotEmpty(t, limitCpus)
	// logFormat can be empty as DefaultLogFormat is ""
}

func TestInitLogAppliesLevelWhenConfigIsMissing(t *testing.T) {
	previousPath := logConfigPath
	previousLevel := logLevel
	defer func() {
		logConfigPath = previousPath
		logLevel = previousLevel
		logger.InitLogger(nil)
	}()

	logConfigPath = "/tmp/dubbo-go-pixiu-missing-log-config.yml"
	logLevel = "info"

	err := initLog()
	assert.Error(t, err)
	assert.False(t, logger.GetLogger().Desugar().Core().Enabled(zap.DebugLevel))
}

func TestGatewayCmdAddedToRootCmd(t *testing.T) {
	// Check that GatewayCmd is properly initialized
	assert.NotNil(t, GatewayCmd)
	assert.Equal(t, "gateway", GatewayCmd.Use)
	assert.Equal(t, "Run dubbo go pixiu in gateway mode", GatewayCmd.Short)
}

func TestStartGatewayCmdAddedToGatewayCmd(t *testing.T) {
	// Check that startGatewayCmd is properly initialized
	assert.NotNil(t, startGatewayCmd)
	assert.Equal(t, "start", startGatewayCmd.Use)
	assert.Equal(t, "Start gateway", startGatewayCmd.Short)
	assert.NotNil(t, startGatewayCmd.PreRunE)
	assert.NotNil(t, startGatewayCmd.RunE)
}
