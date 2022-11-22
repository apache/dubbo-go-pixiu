//  Copyright Istio Authors
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package stackdriver

import (
	cloudtracepb "google.golang.org/genproto/googleapis/devtools/cloudtrace/v1"
	loggingpb "google.golang.org/genproto/googleapis/logging/v2"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
)

import (
	md "github.com/apache/dubbo-go-pixiu/pkg/bootstrap/platform"
	"github.com/apache/dubbo-go-pixiu/pkg/test"
	"github.com/apache/dubbo-go-pixiu/pkg/test/framework"
	"github.com/apache/dubbo-go-pixiu/pkg/test/framework/components/cluster"
	"github.com/apache/dubbo-go-pixiu/pkg/test/framework/resource"
)

// Instance represents a deployed Stackdriver app instance in a Kubernetes cluster.
type Instance interface {
	Address() string
	// Gets the namespace in which stackdriver is deployed.
	GetStackdriverNamespace() string
	ListTimeSeries(namespace string) ([]*monitoringpb.TimeSeries, error)
	ListLogEntries(lt LogType, namespace string) ([]*loggingpb.LogEntry, error)
	ListTraces(namespace string) ([]*cloudtracepb.Trace, error)
}

type Config struct {
	// Cluster to be used in a multicluster environment
	Cluster cluster.Cluster
}

// New returns a new instance of stackdriver.
func New(ctx resource.Context, c Config) (i Instance, err error) {
	if UseRealStackdriver() {
		return newRealStackdriver(ctx, c)
	}
	return newKube(ctx, c)
}

// NewOrFail returns a new Stackdriver instance or fails test.
func NewOrFail(t test.Failer, ctx resource.Context, c Config, realSD bool) Instance {
	t.Helper()
	i, err := New(ctx, c)
	if err != nil {
		t.Fatalf("stackdriver.NewOrFail: %v", err)
	}

	return i
}

func UseRealStackdriver() bool {
	// Use real stackdriver only if the test intends to AND the test is running on GCP.
	// Currently real stackdriver only works if the test runs on GCP.
	return framework.UseRealStackdriver && md.IsGCP()
}
