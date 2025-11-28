/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package proto_suite

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

import (
	triplepb "dubbo-go-pixiu-benchmark/api"

	"dubbo-go-pixiu-benchmark/test"

	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"

	tripleConstant "github.com/dubbogo/triple/pkg/common/constant"

	. "github.com/onsi/ginkgo/v2"

	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"github.com/onsi/gomega/gmeasure"
)

var (
	grpcGreeterImpl                   = new(triplepb.GreeterClientImpl)
	tripleServerSession, pixiuSession *gexec.Session
)

func TestTripleCases(t *testing.T) {
	gomega.RegisterFailHandler(Fail)
	RunSpecs(t, "test")
}

var _ = Describe("triple protocol performance test", Ordered, func() {
	BeforeAll(func() {
		var err error
		test.CurPath, err = os.Getwd()
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Wait for ports to be available before starting servers
		waitForPortAvailable("20000", 10*time.Second)
		waitForPortAvailable("8881", 10*time.Second)

		tripleServerSession = prepareTripleServer()
		// Wait for dubbo server to register to Zookeeper
		time.Sleep(5 * time.Second)

		pixiuSession = test.PreparePixiu("../../../dist/pixiu", test.CurPath+"/../../../protocol/triple/pb/pixiu/conf/config.yaml")
		// Wait for pixiu to discover services from Zookeeper
		time.Sleep(5 * time.Second)
	})

	It("pixiu to triple protocol performance test", func() {
		defer GinkgoRecover()

		urlPrefix := "http://localhost:8881/dubbo.io/org.apache.dubbogo.samples.api.Greeter/%s"

		experiment := gmeasure.NewExperiment("pixiu to triple protocol performance test")
		AddReportEntry(experiment.Name, experiment)

		experiment.Sample(func(idx int) {
			experiment.MeasureDuration("SayHello", func() {
				defer GinkgoRecover()

				url := fmt.Sprintf(urlPrefix, "SayHello")
				data := `
{
    "name":"test"
}
`
				resp, err := http.Post(url, "application/json", strings.NewReader(data))
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				defer resp.Body.Close()
				gomega.Expect(resp.Status).To(gomega.Equal("200 OK"))
				reply, err := io.ReadAll(resp.Body)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(string(reply)).NotTo(gomega.MatchRegexp("client call err*"))
				//fmt.Printf("consumer:%+v\n", string(reply))
			})
		}, test.SampleConfig)
	})

	It("triple protocol performance test", func() {
		// Initialize dubbo client only when needed for this test
		prepareTripleClient()

		experiment := gmeasure.NewExperiment("triple protocol performance test")
		AddReportEntry(experiment.Name, experiment)

		experiment.Sample(func(idx int) {
			defer GinkgoRecover()

			experiment.MeasureDuration("SayHello", func() {
				req := &triplepb.HelloRequest{
					Name: "laurence",
				}
				ctx := context.WithValue(context.Background(), tripleConstant.TripleCtxKey("tri-req-id"), "test_value_XXXXXXXX")
				_, err := grpcGreeterImpl.SayHello(ctx, req)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				//fmt.Printf("consumer:%+v\n", reply)
			})
		}, test.SampleConfig)
	})

	AfterAll(func() {
		// Use Kill instead of Terminate to ensure processes are fully stopped
		if pixiuSession != nil {
			pixiuSession.Kill().Wait(10 * time.Second)
		}
		if tripleServerSession != nil {
			tripleServerSession.Kill().Wait(10 * time.Second)
		}
		gexec.CleanupBuildArtifacts()
		// Wait for ports to be released
		time.Sleep(2 * time.Second)
	})
})

func prepareTripleServer() *gexec.Session {
	serverProcess, err := gexec.Build("dubbo-go-pixiu-benchmark/protocol/triple/pb/go-server/cmd")
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	command := exec.Command(serverProcess)
	// Set working directory
	command.Dir = test.CurPath
	session, err := gexec.Start(command, io.Discard, io.Discard)
	//session, err := gexec.Start(command, os.Stdout, os.Stderr)

	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	return session
}

func prepareTripleClient() {
	config.SetConsumerService(grpcGreeterImpl)
	err := config.Load(config.WithPath(test.CurPath + "/../../../protocol/triple/pb/go-client/conf/dubbogo.yml"))
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

// waitForPortAvailable waits until the port is available (not in use)
func waitForPortAvailable(port string, timeout time.Duration) {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", "127.0.0.1:"+port, 100*time.Millisecond)
		if err != nil {
			// Port is available (connection refused means no one is listening)
			return
		}
		conn.Close()
		time.Sleep(500 * time.Millisecond)
	}
}
