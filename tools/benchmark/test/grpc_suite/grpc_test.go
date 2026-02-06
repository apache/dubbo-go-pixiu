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

package grpc_suite

import (
	"context"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

import (
	. "github.com/onsi/ginkgo/v2"

	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"github.com/onsi/gomega/gmeasure"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

import (
	"github.com/apache/dubbo-go-pixiu/tools/benchmark/api/grpcstub"
	"github.com/apache/dubbo-go-pixiu/tools/benchmark/test"
)

var (
	benchmarkClient                 grpcstub.BenchmarkServiceClient
	grpcServerSession, pixiuSession *gexec.Session
	ctx                             = context.Background()
)

func TestGRPCCases(t *testing.T) {
	gomega.RegisterFailHandler(Fail)
	RunSpecs(t, "gRPC Benchmark Test Suite")
}

var _ = Describe("gRPC protocol performance test", Ordered, func() {
	BeforeAll(func() {
		var err error
		test.CurPath, err = os.Getwd()
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		grpcServerSession = prepareGRPCServer()

		time.Sleep(1 * time.Second)

		pixiuSession = test.PreparePixiu("../../dist/pixiu", test.CurPath+"/../../protocol/grpc/pixiu/conf/config.yaml")

		time.Sleep(3 * time.Second)

		prepareGRPCClient()
	})

	It("gRPC protocol performance test", func() {
		experiment := gmeasure.NewExperiment("gRPC protocol performance test")
		AddReportEntry(experiment.Name, experiment)

		experiment.Sample(func(idx int) {
			defer GinkgoRecover()

			experiment.MeasureDuration("GetUser", func() {
				ctxWithTO, cancel := context.WithTimeout(ctx, 5*time.Second)
				defer cancel()
				resp, err := benchmarkClient.GetUser(ctxWithTO, &grpcstub.GetUserRequest{UserId: 1})
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(resp).NotTo(gomega.BeNil())
				gomega.Expect(len(resp.Users)).To(gomega.Equal(1))
			})
		}, test.SampleConfig)

		experiment.Sample(func(idx int) {
			defer GinkgoRecover()

			experiment.MeasureDuration("GetUsers", func() {
				ctxWithTO, cancel := context.WithTimeout(ctx, 5*time.Second)
				defer cancel()
				resp, err := benchmarkClient.GetUsers(ctxWithTO, &grpcstub.GetUsersRequest{UserIds: []int32{1, 2}})
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(resp).NotTo(gomega.BeNil())
				gomega.Expect(len(resp.Users)).To(gomega.Equal(2))
			})
		}, test.SampleConfig)

		experiment.Sample(func(idx int) {
			defer GinkgoRecover()

			experiment.MeasureDuration("GetUserByName", func() {
				ctxWithTO, cancel := context.WithTimeout(ctx, 5*time.Second)
				defer cancel()
				resp, err := benchmarkClient.GetUserByName(ctxWithTO, &grpcstub.GetUserByNameRequest{Name: "Kenway"})
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(resp).NotTo(gomega.BeNil())
				gomega.Expect(len(resp.Users)).To(gomega.Equal(1))
			})
		}, test.SampleConfig)

		experiment.Sample(func(idx int) {
			defer GinkgoRecover()

			experiment.MeasureDuration("SayHello", func() {
				ctxWithTO, cancel := context.WithTimeout(ctx, 5*time.Second)
				defer cancel()
				resp, err := benchmarkClient.SayHello(ctxWithTO, &grpcstub.HelloRequest{Name: "World"})
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(resp).NotTo(gomega.BeNil())
				gomega.Expect(resp.Message).To(gomega.Equal("Hello World"))
			})
		}, test.SampleConfig)
	})

	It("pixiu to gRPC protocol performance test", func() {
		experiment := gmeasure.NewExperiment("pixiu to gRPC protocol performance test")
		AddReportEntry(experiment.Name, experiment)

		urlPrefix := "http://localhost:8881/api/v1/benchmark.BenchmarkService/"

		experiment.Sample(func(idx int) {
			defer GinkgoRecover()

			experiment.MeasureDuration("GetUser", func() {
				url := urlPrefix + "GetUser"
				data := `{"userId": 1}`
				resp, err := http.Post(url, "application/json", strings.NewReader(data))
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				_, err = io.ReadAll(resp.Body)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(resp.StatusCode).To(gomega.Equal(200))
			})
		}, test.SampleConfig)

		experiment.Sample(func(idx int) {
			defer GinkgoRecover()

			experiment.MeasureDuration("GetUsers", func() {
				url := urlPrefix + "GetUsers"
				data := `{"userIds": [1, 2, 3]}`
				resp, err := http.Post(url, "application/json", strings.NewReader(data))
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				_, err = io.ReadAll(resp.Body)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(resp.StatusCode).To(gomega.Equal(200))
			})
		}, test.SampleConfig)

		experiment.Sample(func(idx int) {
			defer GinkgoRecover()

			experiment.MeasureDuration("GetUserByName", func() {
				url := urlPrefix + "GetUserByName"
				data := `{"name": "Kenway"}`
				resp, err := http.Post(url, "application/json", strings.NewReader(data))
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				_, err = io.ReadAll(resp.Body)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(resp.StatusCode).To(gomega.Equal(200))
			})
		}, test.SampleConfig)

		experiment.Sample(func(idx int) {
			defer GinkgoRecover()

			experiment.MeasureDuration("SayHello", func() {
				url := urlPrefix + "SayHello"
				data := `{"name": "World"}`
				resp, err := http.Post(url, "application/json", strings.NewReader(data))
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				_, err = io.ReadAll(resp.Body)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(resp.StatusCode).To(gomega.Equal(200))
			})
		}, test.SampleConfig)
	})

	AfterAll(func() {
		time.Sleep(5 * time.Second)
		grpcServerSession.Terminate().Wait(5 * time.Second)
		pixiuSession.Terminate().Wait(5 * time.Second)
	})
})

func prepareGRPCServer() *gexec.Session {
	serverProcess, err := gexec.Build("github.com/apache/dubbo-go-pixiu/tools/benchmark/protocol/grpc/go-server/cmd")
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	command := exec.Command(serverProcess)
	session, err := gexec.Start(command, io.Discard, io.Discard)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	return session
}

func prepareGRPCClient() {
	conn, err := grpc.Dial("localhost:50001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	benchmarkClient = grpcstub.NewBenchmarkServiceClient(conn)
}
