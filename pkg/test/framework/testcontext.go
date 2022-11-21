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

package framework

import (
	"fmt"
	"io"
	"os"
	"path"
	"testing"
	"time"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/test"
	"github.com/apache/dubbo-go-pixiu/pkg/test/framework/components/cluster"
	"github.com/apache/dubbo-go-pixiu/pkg/test/framework/errors"
	"github.com/apache/dubbo-go-pixiu/pkg/test/framework/label"
	"github.com/apache/dubbo-go-pixiu/pkg/test/framework/resource"
	"github.com/apache/dubbo-go-pixiu/pkg/test/scopes"
	"github.com/apache/dubbo-go-pixiu/pkg/test/util/yml"
)

// TestContext is a test-level context that can be created as part of test executing tests.
type TestContext interface {
	resource.Context
	test.Failer

	// NewSubTest creates a new sub-test under the current running Test. The lifecycle of a sub-Test is scoped to the
	// parent. Calls to Done() will block until all children are also Done(). When Run, sub-Tests will automatically
	// create their own Golang *testing.T with the name provided.
	//
	// If this TestContext was not created by a Test or if that Test is not running, this method will panic.
	NewSubTest(name string) Test
	NewSubTestf(format string, a ...interface{}) Test

	// WorkDir allocated for this test.
	WorkDir() string

	// CreateDirectoryOrFail creates a new sub directory with the given name in the workdir, or fails the test.
	CreateDirectoryOrFail(name string) string

	// CreateTmpDirectoryOrFail creates a new temporary directory with the given prefix in the workdir, or fails the test.
	CreateTmpDirectoryOrFail(prefix string) string

	// SkipDumping will skip dumping debug logs/configs/etc for this scope only (child scopes are not skipped).
	SkipDumping()

	// Done should be called when this context is no longer needed. It triggers the asynchronous cleanup of any
	// allocated resources.
	Done()

	// Methods for interacting with the underlying *testing.T.
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Failed() bool
	Name() string
	Skip(args ...interface{})
	SkipNow()
	Skipf(format string, args ...interface{})
	Skipped() bool
}

var (
	_ TestContext = &testContext{}
	_ test.Failer = &testContext{}
)

// testContext for the currently executing test.
type testContext struct {
	yml.FileWriter

	// The id of the test context. useful for debugging the test framework itself.
	id string

	// The currently running Test. Non-nil if this context was created by a Test
	test *testImpl

	// The underlying Go testing.T for this context.
	*testing.T

	// suite-level context
	suite *suiteContext

	// resource scope for this context.
	scope *scope

	// The workDir for this particular context
	workDir string
}

// Before executing a new context, we should wait for existing contexts to terminate if they are NOT parents of this context.
// This is to workaround termination of functions run with RunParallel. When this is used, child tests will not run until the parent
// has terminated. This means that the parent cannot synchronously cleanup, or it would block its children. However, if we do async cleanup,
// then new tests can unexpectedly start during the cleanup of another. This may lead to odd results, like a test cleanup undoing the setup of a future test.
// To workaround this, we maintain a set of all contexts currently terminating. Before starting the context, we will search this set;
// if any non-parent contexts are found, we will wait.
func waitForParents(test *testImpl) {
	iterations := 0
	for {
		iterations++
		done := true
		globalParentLock.Range(func(key, value interface{}) bool {
			k := key.(*testImpl)
			current := test
			for current != nil {
				if current == k {
					return true
				}
				current = current.parent

			}
			// We found an item in the list, and we are *not* a child of it. This means another test hierarchy has exclusive access right now
			// Wait until they are finished before proceeding
			done = false
			return true
		})
		if done {
			return
		}
		time.Sleep(time.Millisecond * 50)
		// Add some logging in case something locks up so we can debug
		if iterations%10 == 0 {
			globalParentLock.Range(func(key, value interface{}) bool {
				scopes.Framework.Warnf("Stuck waiting for parent test suites to terminate... %v is blocking", key.(*testImpl).goTest.Name())
				return true
			})
		}
	}
}

func newTestContext(test *testImpl, goTest *testing.T, s *suiteContext, parentScope *scope, labels label.Set) *testContext {
	waitForParents(test)
	id := s.allocateContextID(goTest.Name())

	allLabels := s.suiteLabels.Merge(labels)
	if !s.settings.Selector.Selects(allLabels) {
		goTest.Skipf("Skipping: label mismatch: labels=%v, filter=%v", allLabels, s.settings.Selector)
	}

	if s.settings.SkipMatcher.MatchTest(goTest.Name()) {
		goTest.Skipf("Skipping: test %v matched -istio.test.skip regex", goTest.Name())
	}

	scopes.Framework.Debugf("Creating New test context")
	workDir := path.Join(s.settings.RunDir(), goTest.Name(), "_test_context")
	if _, err := os.Stat(path.Join(s.settings.RunDir(), goTest.Name())); !os.IsNotExist(err) {
		// Folder already exist. This can happen due to --istio.test.retries. Switch to using "id", which
		// is globally unique. We do not due this all the time since it breaks the structure of subtests
		// a bit. When we use id we end up with "Parent-0". However, subtests end up with
		// "Parent/Child-0", which is not in the same folder. As a compromise, we only append the id in
		// the rare case of retry. This ensures we always have all data, and in the common cases the data
		// is more readable.
		workDir = path.Join(s.settings.RunDir(), id, "_test_context")
	}
	if err := os.MkdirAll(workDir, os.ModePerm); err != nil {
		goTest.Fatalf("Error creating work dir %q: %v", workDir, err)
	}

	scopes.Framework.Debugf("Creating new testContext: %q", id)

	if parentScope == nil {
		parentScope = s.globalScope
	}

	scopeID := fmt.Sprintf("[%s]", id)
	return &testContext{
		id:         id,
		test:       test,
		T:          goTest,
		suite:      s,
		scope:      newScope(scopeID, parentScope),
		workDir:    workDir,
		FileWriter: yml.NewFileWriter(workDir),
	}
}

func (c *testContext) Settings() *resource.Settings {
	return c.suite.settings
}

func (c *testContext) TrackResource(r resource.Resource) resource.ID {
	id := c.suite.allocateResourceID(c.id, r)
	rid := &resourceID{id: id}
	c.scope.add(r, rid)
	return rid
}

func (c *testContext) GetResource(ref interface{}) error {
	return c.scope.get(ref)
}

func (c *testContext) WorkDir() string {
	return c.workDir
}

func (c *testContext) Environment() resource.Environment {
	return c.suite.environment
}

func (c *testContext) Clusters() cluster.Clusters {
	if c == nil || c.Environment() == nil {
		return nil
	}
	return c.Environment().Clusters()
}

func (c *testContext) AllClusters() cluster.Clusters {
	if c == nil || c.Environment() == nil {
		return nil
	}
	return c.Environment().AllClusters()
}

func (c *testContext) CreateDirectory(name string) (string, error) {
	dir, err := os.MkdirTemp(c.workDir, name)
	if err != nil {
		scopes.Framework.Errorf("Error creating dir: runID='%v', prefix='%s', workDir='%v', err='%v'",
			c.suite.settings.RunID, name, c.workDir, err)
	} else {
		scopes.Framework.Debugf("Created a dir: runID='%v', name='%s'", c.suite.settings.RunID, dir)
	}
	return dir, err
}

func (c *testContext) CreateDirectoryOrFail(name string) string {
	tmp, err := c.CreateDirectory(name)
	if err != nil {
		c.Fatalf("Error creating  directory with name %q: %v", name, err)
	}
	return tmp
}

func (c *testContext) CreateTmpDirectory(prefix string) (string, error) {
	dir, err := os.MkdirTemp(c.workDir, prefix)
	if err != nil {
		scopes.Framework.Errorf("Error creating temp dir: runID='%v', prefix='%s', workDir='%v', err='%v'",
			c.suite.settings.RunID, prefix, c.workDir, err)
	} else {
		scopes.Framework.Debugf("Created a temp dir: runID='%v', name='%s'", c.suite.settings.RunID, dir)
	}

	return dir, err
}

func (c *testContext) SkipDumping() {
	c.scope.skipDumping()
}

func (c *testContext) ConfigKube(clusters ...cluster.Cluster) resource.ConfigManager {
	return newConfigManager(c, clusters)
}

func (c *testContext) ConfigIstio() resource.ConfigManager {
	return newConfigManager(c, c.Clusters().Configs())
}

func (c *testContext) CreateTmpDirectoryOrFail(prefix string) string {
	tmp, err := c.CreateTmpDirectory(prefix)
	if err != nil {
		c.Fatalf("Error creating temp directory with prefix %q: %v", prefix, err)
	}
	return tmp
}

func (c *testContext) newChildContext(test *testImpl) *testContext {
	return newTestContext(test, test.goTest, c.suite, c.scope, label.NewSet(test.labels...))
}

func (c *testContext) NewSubTest(name string) Test {
	if c.test == nil {
		panic(fmt.Sprintf("Attempting to create subtest %s from a TestContext with no associated Test", name))
	}

	if c.test.goTest == nil || c.test.ctx == nil {
		panic(fmt.Sprintf("Attempting to create subtest %s before running parent", name))
	}

	return &testImpl{
		name:          name,
		parent:        c.test,
		s:             c.test.s,
		featureLabels: c.test.featureLabels,
	}
}

func (c *testContext) NewSubTestf(format string, a ...interface{}) Test {
	return c.NewSubTest(fmt.Sprintf(format, a...))
}

func (c *testContext) ConditionalCleanup(fn func()) {
	c.scope.addCloser(&closer{fn: func() error {
		fn()
		return nil
	}, noskip: true})
}

func (c *testContext) Cleanup(fn func()) {
	c.scope.addCloser(&closer{fn: func() error {
		fn()
		return nil
	}})
}

func (c *testContext) Done() {
	if c.Failed() && c.Settings().CIMode {
		scopes.Framework.Debugf("Begin dumping testContext: %q", c.id)
		// make sure we dump suite-level resources, but don't dump sibling tests or their children
		rt.DumpShallow(c)
		c.scope.dump(c, true)
		scopes.Framework.Debugf("Completed dumping testContext: %q", c.id)
	}

	scopes.Framework.Debugf("Begin cleaning up testContext: %q", c.id)
	if err := c.scope.done(c.suite.settings.NoCleanup); err != nil {
		c.Logf("error scope cleanup: %v", err)
		if c.Settings().FailOnDeprecation {
			if errors.IsOrContainsDeprecatedError(err) {
				c.Error("Using deprecated Envoy features. Failing due to -istio.test.deprecation_failure flag.")
			}
		}
	}
	scopes.Framework.Debugf("Completed cleaning up testContext: %q", c.id)
}

func (c *testContext) Error(args ...interface{}) {
	c.Helper()
	c.T.Error(args...)
}

func (c *testContext) Errorf(format string, args ...interface{}) {
	c.Helper()
	c.T.Errorf(format, args...)
}

func (c *testContext) Fail() {
	c.Helper()
	c.T.Fail()
}

func (c *testContext) FailNow() {
	c.Helper()
	c.T.FailNow()
}

func (c *testContext) Failed() bool {
	c.Helper()
	return c.T.Failed()
}

func (c *testContext) Fatal(args ...interface{}) {
	c.Helper()
	c.T.Fatal(args...)
}

func (c *testContext) Fatalf(format string, args ...interface{}) {
	c.Helper()
	c.T.Fatalf(format, args...)
}

func (c *testContext) Log(args ...interface{}) {
	c.Helper()
	c.T.Log(args...)
}

func (c *testContext) Logf(format string, args ...interface{}) {
	c.Helper()
	c.T.Logf(format, args...)
}

func (c *testContext) Name() string {
	c.Helper()
	return c.T.Name()
}

func (c *testContext) Skip(args ...interface{}) {
	c.Helper()
	c.T.Skip(args...)
}

func (c *testContext) SkipNow() {
	c.Helper()
	c.T.SkipNow()
}

func (c *testContext) Skipf(format string, args ...interface{}) {
	c.Helper()
	c.T.Skipf(format, args...)
}

func (c *testContext) Skipped() bool {
	c.Helper()
	return c.T.Skipped()
}

var _ io.Closer = &closer{}

type closer struct {
	fn     func() error
	noskip bool
}

func (c *closer) Close() error {
	return c.fn()
}

func (c *testContext) RecordTraceEvent(key string, value interface{}) {
	// Currently, only supported at suite level.
	panic("TODO: implement tracing in test context")
}
