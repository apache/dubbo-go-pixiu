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
	"dubbo.apache.org/dubbo-go/v3/client"
	_ "dubbo.apache.org/dubbo-go/v3/imports"

	. "github.com/onsi/ginkgo/v2"

	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"github.com/onsi/gomega/gmeasure"
)

import (
	"github.com/apache/dubbo-go-pixiu/tools/benchmark/api"
	"github.com/apache/dubbo-go-pixiu/tools/benchmark/test"
)

var (
	benchmarkService                  api.BenchmarkService
	tripleServerSession, pixiuSession *gexec.Session
)

func TestTripleCases(t *testing.T) {
	gomega.RegisterFailHandler(Fail)
	RunSpecs(t, "Triple Benchmark Test Suite")
}

var _ = Describe("triple protocol performance test", Ordered, func() {
	BeforeAll(func() {
		var err error
		test.CurPath, err = os.Getwd()
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		waitForPortAvailable("20000", 10*time.Second)
		waitForPortAvailable("8881", 10*time.Second)

		tripleServerSession = prepareTripleServer()
		time.Sleep(8 * time.Second)

		pixiuSession = test.PreparePixiu("../../../dist/pixiu", test.CurPath+"/../../../protocol/triple/pixiu/conf/config.yaml")
		time.Sleep(6 * time.Second)
	})

	It("pixiu to triple protocol performance test", func() {
		defer GinkgoRecover()

		urlPrefix := "http://localhost:8881/dubbo.io/benchmark.BenchmarkService/%s"

		experiment := gmeasure.NewExperiment("pixiu to triple protocol performance test")
		AddReportEntry(experiment.Name, experiment)

		experiment.Sample(func(idx int) {
			experiment.MeasureDuration("GetUser", func() {
				defer GinkgoRecover()

				url := fmt.Sprintf(urlPrefix, "GetUser")
				data := `{"userId": 1}`
				resp, err := http.Post(url, "application/json", strings.NewReader(data))
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				defer resp.Body.Close()
				gomega.Expect(resp.Status).To(gomega.Equal("200 OK"))
				reply, err := io.ReadAll(resp.Body)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(string(reply)).NotTo(gomega.MatchRegexp("client call err*"))
			})
		}, test.SampleConfig)

		experiment.Sample(func(idx int) {
			experiment.MeasureDuration("GetUsers", func() {
				defer GinkgoRecover()

				url := fmt.Sprintf(urlPrefix, "GetUsers")
				data := `{"userIds": [1, 2]}`
				resp, err := http.Post(url, "application/json", strings.NewReader(data))
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				defer resp.Body.Close()
				gomega.Expect(resp.Status).To(gomega.Equal("200 OK"))
				reply, err := io.ReadAll(resp.Body)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(string(reply)).NotTo(gomega.MatchRegexp("client call err*"))
			})
		}, test.SampleConfig)

		experiment.Sample(func(idx int) {
			experiment.MeasureDuration("GetUserByName", func() {
				defer GinkgoRecover()

				url := fmt.Sprintf(urlPrefix, "GetUserByName")
				data := `{"name": "Kenway"}`
				resp, err := http.Post(url, "application/json", strings.NewReader(data))
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				defer resp.Body.Close()
				gomega.Expect(resp.Status).To(gomega.Equal("200 OK"))
				reply, err := io.ReadAll(resp.Body)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(string(reply)).NotTo(gomega.MatchRegexp("client call err*"))
			})
		}, test.SampleConfig)

		experiment.Sample(func(idx int) {
			experiment.MeasureDuration("SayHello", func() {
				defer GinkgoRecover()

				url := fmt.Sprintf(urlPrefix, "SayHello")
				data := `{"name": "test"}`
				resp, err := http.Post(url, "application/json", strings.NewReader(data))
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				defer resp.Body.Close()
				gomega.Expect(resp.Status).To(gomega.Equal("200 OK"))
				reply, err := io.ReadAll(resp.Body)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(string(reply)).NotTo(gomega.MatchRegexp("client call err*"))
			})
		}, test.SampleConfig)
	})

	It("triple protocol performance test", func() {
		prepareTripleClient()

		experiment := gmeasure.NewExperiment("triple protocol performance test")
		AddReportEntry(experiment.Name, experiment)

		experiment.Sample(func(idx int) {
			defer GinkgoRecover()

			experiment.MeasureDuration("GetUser", func() {
				req := &api.GetUserRequest{UserId: 1}
				_, err := benchmarkService.GetUser(context.Background(), req)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
			})
		}, test.SampleConfig)

		experiment.Sample(func(idx int) {
			defer GinkgoRecover()

			experiment.MeasureDuration("GetUsers", func() {
				req := &api.GetUsersRequest{UserIds: []int32{1, 2}}
				_, err := benchmarkService.GetUsers(context.Background(), req)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
			})
		}, test.SampleConfig)

		experiment.Sample(func(idx int) {
			defer GinkgoRecover()

			experiment.MeasureDuration("GetUserByName", func() {
				req := &api.GetUserByNameRequest{Name: "Kenway"}
				_, err := benchmarkService.GetUserByName(context.Background(), req)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
			})
		}, test.SampleConfig)

		experiment.Sample(func(idx int) {
			defer GinkgoRecover()

			experiment.MeasureDuration("SayHello", func() {
				req := &api.HelloRequest{Name: "laurence"}
				_, err := benchmarkService.SayHello(context.Background(), req)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
			})
		}, test.SampleConfig)
	})

	AfterAll(func() {
		if pixiuSession != nil {
			pixiuSession.Kill().Wait(10 * time.Second)
		}
		if tripleServerSession != nil {
			tripleServerSession.Kill().Wait(10 * time.Second)
		}
		gexec.CleanupBuildArtifacts()
		time.Sleep(2 * time.Second)
	})
})

func prepareTripleServer() *gexec.Session {
	serverProcess, err := gexec.Build("github.com/apache/dubbo-go-pixiu/tools/benchmark/protocol/triple/go-server/cmd")
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	command := exec.Command(serverProcess)
	command.Dir = test.CurPath
	session, err := gexec.Start(command, io.Discard, io.Discard)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	return session
}

func prepareTripleClient() {
	// Create client using new API
	cli, err := client.NewClient(
		client.WithClientURL("tri://127.0.0.1:20000"),
	)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	// Create BenchmarkService instance
	benchmarkService, err = api.NewBenchmarkService(cli)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

func waitForPortAvailable(port string, timeout time.Duration) {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", "127.0.0.1:"+port, 100*time.Millisecond)
		if err != nil {
			return
		}
		conn.Close()
		time.Sleep(500 * time.Millisecond)
	}
}
