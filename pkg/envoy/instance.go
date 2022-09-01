// Copyright Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package envoy

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"

	envoyAdmin "github.com/envoyproxy/go-control-plane/envoy/admin/v3"

	"istio.io/pkg/log"
)

const (
	defaultName        = "envoy"
	defaultLiveTimeout = 20 * time.Second
)

// Config for an Envoy Instance.
type Config struct {
	// Options provides the command-line options to be passed to Envoy.
	Options Options

	// Name of the envoy instance, used for logging only. If not provided, defaults to "envoy".
	Name string

	// BinaryPath the path to the Envoy binary.
	BinaryPath string

	// WorkingDir to be used when running Envoy. If not set, the current working directory is used.
	WorkingDir string

	// AdminPort specifies the administration port for the Envoy server. If not set, will
	// be determined by parsing the Envoy bootstrap configuration file.
	AdminPort uint32

	// SkipBaseIDClose skips calling Close on the BaseID, if one was specified.
	SkipBaseIDClose bool
}

// Waitable specifies a waitable operation.
type Waitable interface {
	// WithTimeout specifies an upper bound on the wait time.
	WithTimeout(timeout time.Duration) Waitable

	// Do performs the wait. By default, waits indefinitely. To specify an upper bound on
	// the wait time, use WithTimeout. If the wait times out, returns the last known error
	// for retried operations or context.DeadlineExceeded if no previous error was encountered.
	Do() error
}

// Instance of an Envoy process.
type Instance interface {
	// Config returns the configuration for this Instance.
	Config() Config

	// BaseID used to start Envoy. If not set, returns InvalidBaseID.
	BaseID() BaseID

	// Epoch used to start Envoy. If it was not set, defaults to 0.
	Epoch() Epoch

	// Start the Envoy Instance. The process will be killed if the given context is canceled.
	//
	// If this Instance was created via NewInstanceForHotRestart, this method will block until the parent
	// Instance terminates or goes "live". This is due to the fact that hot restart will fail if the previous
	// envoy process is still initializing.
	Start(ctx context.Context) Instance

	// NewInstanceForHotRestart creates a new Envoy Instance that is configured for a hot restart of this
	// Instance (i.e. epoch is incremented). During a hot restart of Envoy, the old process is drained and
	// traffic is shifted over to the new process.
	//
	// The caller must Start the returned instance to initiate the hot restart.
	//
	// If a new Instance is successfully created, it assumes ownership of the Envoy shared memory segment
	// used for hot restart. This means that when this Instance exits, it will no longer destroy
	// the shared memory segment, regardless of the value of SkipBaseIDClose.
	//
	// If this Instance hasn't been started, calling this method does nothing and simply returns
	// this Instance since there is nothing to restart.
	//
	// This method may only be called once on a given Instance. Subsequent calls will return an error.
	NewInstanceForHotRestart() (Instance, error)

	// WaitUntilLive polls the Envoy ServerInfo endpoint and waits for it to transition to "live". If the
	// wait times out, returns the last known error or context.DeadlineExceeded if no error occurred within the
	// specified duration.
	WaitLive() Waitable

	// AdminPort gets the administration port for Envoy.
	AdminPort() uint32

	// GetServerInfo returns a structure representing a call to /server_info
	GetServerInfo() (*envoyAdmin.ServerInfo, error)

	// GetConfigDump polls Envoy admin port for the config dump and returns the response.
	GetConfigDump() (*envoyAdmin.ConfigDump, error)

	// Wait for the Instance to terminate.
	Wait() Waitable

	// Kill the process, if running.
	Kill() error

	// KillAndWait is a helper that calls Kill and then waits for the process to terminate.
	KillAndWait() Waitable

	// Shutdown initiates the graceful termination of Envoy. Returns immediately and does not
	// wait for the process to exit.
	Shutdown() error

	// ShutdownAndWait is a helper that calls Shutdown and waits for the process to terminate.
	ShutdownAndWait() Waitable

	// DrainListeners drains listeners of Envoy so that inflight requests
	// can gracefully finish and even continue making outbound calls as needed.
	DrainListeners() error
}

// FactoryFunc is a function that manufactures Envoy Instances.
type FactoryFunc func(cfg Config) (Instance, error)

var _ FactoryFunc = New

// New creates a new Envoy Instance with the given options.
func New(cfg Config) (Instance, error) {
	if cfg.Name == "" {
		cfg.Name = defaultName
	}

	// Process the binary path.
	if cfg.BinaryPath == "" {
		return nil, errors.New("must specify an Envoy binary")
	}

	// Create the config object from the options.
	ctx := newConfigContext()
	if err := cfg.Options.validate(ctx); err != nil {
		return nil, err
	}

	// Extract the admin port from the configuration.
	adminPort := cfg.AdminPort
	if adminPort == 0 {
		var err error
		adminPort, err = ctx.getAdminPort()
		if err != nil {
			return nil, err
		}
	}

	// Create a new command with the specified options.
	args := cfg.Options.ToArgs()
	cmd := exec.Command(cfg.BinaryPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if cfg.WorkingDir != "" {
		cmd.Dir = cfg.WorkingDir
	}

	return &instance{
		name:      cfg.Name,
		config:    cfg,
		cmd:       cmd,
		adminPort: adminPort,
		baseID:    ctx.baseID,
		epoch:     ctx.epoch,
		waitCh:    make(chan struct{}, 1),
	}, nil
}

type instance struct {
	config     Config
	name       string
	waitErr    error
	cmd        *exec.Cmd
	waitCh     chan struct{}
	adminPort  uint32
	baseID     BaseID
	epoch      Epoch
	started    bool
	hotRestart Instance
	mux        sync.Mutex
}

func (i *instance) Config() Config {
	return i.config
}

func (i *instance) Start(ctx context.Context) Instance {
	i.mux.Lock()
	defer i.mux.Unlock()

	log.Infof("%s starting with command: %v", i.logID(), i.cmd.Args)

	// Make sure we haven't already started.
	if i.started {
		log.Infof("%s was already started, skipping Start", i.logID())
		return i
	}
	i.started = true

	// Start Envoy.
	if err := i.cmd.Start(); err != nil {
		i.waitErr = err
		i.close()
		return i
	}

	// Asynchronously wait for the command to terminate.
	doneCh := make(chan error, 1)
	go func() {
		// Send the
		doneCh <- i.cmd.Wait()
		close(doneCh)
	}()

	go func() {
		// Free all resources and close doneCh when we exit.
		defer i.close()

		select {
		case <-ctx.Done():
			log.Infof("%s Aborting: %v", i.logID(), ctx.Err())
			i.waitErr = ctx.Err()

			// Context aborted ... kill the process.
			if err := i.Kill(); err != nil {
				log.Warnf("%s kill failed: %v", i.logID(), err)
			}
			return
		case i.waitErr = <-doneCh:
			log.Infof("%s exited with error: %v", i.logID(), i.waitErr)
			return
		}
	}()

	return i
}

func (i *instance) NewInstanceForHotRestart() (Instance, error) {
	i.mux.Lock()
	defer i.mux.Unlock()

	if i.hotRestart != nil {
		return nil, fmt.Errorf("%s already created a hot restart Instance", i.logID())
	}

	if !i.started {
		// This instance hasn't been started yet, no restart required.
		return i, nil
	}

	// If this is a hot restart, wait for the parent process to be live before creating the new
	// instance.
	if err := i.WaitLive().WithTimeout(defaultLiveTimeout).Do(); err != nil {
		log.Warnf("%s failed to go live: %v. Proceeding with hot restart",
			i.logID(), err)
	}

	// Copy the configuration, but replace the epoch.
	cfg := i.config
	cfg.Options = make(Options, 0, len(cfg.Options))
	for _, o := range i.config.Options {
		if o.FlagName() != Epoch(0).FlagName() {
			cfg.Options = append(cfg.Options, o)
		}
	}

	// Increment the epoch on the new Instance.
	cfg.Options = append(cfg.Options, i.epoch+1)

	// Create the new instance.
	hotRestart, err := New(cfg)
	if err != nil {
		return nil, err
	}
	i.hotRestart = hotRestart

	return hotRestart, nil
}

func (i *instance) WaitLive() Waitable {
	return &waitableImpl{
		instance:    i,
		retryPeriod: 200 * time.Millisecond,
		retryHandler: func() (bool, error) {
			info, err := GetServerInfo(i.adminPort)
			if err != nil {
				return true, err
			}

			switch info.State {
			case envoyAdmin.ServerInfo_LIVE:
				// We're live!
				return false, nil
			case envoyAdmin.ServerInfo_DRAINING:
				// Don't retry, it'll never happen.
				return false, errors.New("envoy will never go live, it's already draining")
			default:
				// Retry.
				return true, fmt.Errorf("envoy not live. Server State: %s", info.State)
			}
		},
	}
}

func (i *instance) AdminPort() uint32 {
	return i.adminPort
}

func (i *instance) BaseID() BaseID {
	return i.baseID
}

func (i *instance) Epoch() Epoch {
	return i.epoch
}

func (i *instance) GetServerInfo() (*envoyAdmin.ServerInfo, error) {
	return GetServerInfo(i.adminPort)
}

func (i *instance) GetConfigDump() (*envoyAdmin.ConfigDump, error) {
	return GetConfigDump(i.adminPort)
}

func (i *instance) Wait() Waitable {
	return &waitableImpl{
		instance: i,
	}
}

func (i *instance) Kill() error {
	if i.cmd.Process == nil {
		return errors.New("envoy process was not started")
	}

	return i.cmd.Process.Kill()
}

func (i *instance) KillAndWait() Waitable {
	return &waitableImpl{
		instance:    i,
		creationErr: i.Kill(),
	}
}

func (i *instance) Shutdown() error {
	return Shutdown(i.adminPort)
}

func (i *instance) ShutdownAndWait() Waitable {
	return &waitableImpl{
		instance:    i,
		creationErr: i.Shutdown(),
	}
}

func (i *instance) DrainListeners() error {
	return DrainListeners(i.adminPort, true)
}

func (i *instance) close() {
	// Delete the shared memory segment (used for hot restart) if configured to do and no
	// further hot-restarts were initiated. If another restart was initiated, we hand off
	// ownership of the shared memory to that Instance.
	if !i.config.SkipBaseIDClose && i.hotRestart == nil {
		if err := i.baseID.Close(); err != nil {
			log.Infof("Failed freeing BaseID for %s: %v", i.logID(), err)
		}
	}

	close(i.waitCh)
}

func (i *instance) logID() string {
	return fmt.Sprintf("Envoy '%s' (epoch %d)", i.name, i.epoch)
}

var _ Waitable = &waitableImpl{}

type waitableImpl struct {
	*instance

	creationErr error

	retryPeriod  time.Duration
	retryHandler func() (bool, error)

	timeout time.Duration
}

func (w *waitableImpl) Do() error {
	if w.creationErr != nil {
		return w.creationErr
	}

	// Create a dummy time channel to be used if needed.
	dummyTimeCh := make(chan time.Time, 1)
	defer close(dummyTimeCh)

	// Create the timeout channel
	var timeoutCh <-chan time.Time
	if w.timeout > 0 {
		// Create a timer for the specified duration.
		timer := time.NewTimer(w.timeout)
		defer timer.Stop()
		timeoutCh = timer.C
	} else {
		// No timeout was specified, just create dummy channel.
		timeoutCh = dummyTimeCh
	}

	var retryCh <-chan time.Time
	isRetrying := w.retryPeriod > 0
	if isRetrying {
		// Create a ticker for the specified duration.
		ticker := time.NewTicker(w.retryPeriod)
		defer ticker.Stop()
		retryCh = ticker.C
	} else {
		// No ticker was specified, just use a dummy channel.
		retryCh = dummyTimeCh
	}

	var lastErr error
	for {
		select {
		case <-w.waitCh:
			if w.waitErr != nil {
				return w.waitErr
			}
			if lastErr != nil {
				return lastErr
			}
			if isRetrying {
				// Envoy exited before the retry operation was successful,
				return errors.New("envoy process exited before wait completed")
			}
			return nil
		case <-retryCh:
			shouldRetry, err := w.retryHandler()
			if !shouldRetry {
				return err
			}

			// We're retrying, save the last error.
			lastErr = err
		case <-timeoutCh:
			if lastErr != nil {
				return lastErr
			}
			// The timeout occurred before any other events.
			return context.DeadlineExceeded
		}
	}
}

func (w *waitableImpl) WithTimeout(d time.Duration) Waitable {
	out := *w
	out.timeout = d
	return &out
}
