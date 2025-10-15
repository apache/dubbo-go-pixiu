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
	pb "dubbo-go-pixiu-benchmark/protocol/grpc/proto"

	"dubbo-go-pixiu-benchmark/test"

	. "github.com/onsi/ginkgo/v2"

	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"github.com/onsi/gomega/gmeasure"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	userProviderClient              pb.UserProviderClient
	grpcServerSession, pixiuSession *gexec.Session
	ctx                             = context.Background()
)

func TestGRPCCases(t *testing.T) {
	gomega.RegisterFailHandler(Fail)
	RunSpecs(t, "test")
}

var _ = Describe("grpc protocol performance test", Ordered, func() {
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

	It("grpc protocol performance test", func() {
		experiment := gmeasure.NewExperiment("grpc protocol performance test")
		AddReportEntry(experiment.Name, experiment)

		experiment.Sample(func(idx int) {
			defer GinkgoRecover()

			experiment.MeasureDuration("GetUser", func() {
				ctxWithTO, cancel := context.WithTimeout(ctx, 5*time.Second)
				defer cancel()
				resp, err := userProviderClient.GetUser(ctxWithTO, &pb.GetUserRequest{UserId: 1})
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(resp).NotTo(gomega.BeNil())
				gomega.Expect(len(resp.Users), 1)
				//fmt.Printf("consumer:%+v\n", resp.Users)
			})
		}, test.SampleConfig)

		experiment.Sample(func(idx int) {
			defer GinkgoRecover()

			experiment.MeasureDuration("GetUsers", func() {
				ctxWithTO, cancel := context.WithTimeout(ctx, 5*time.Second)
				defer cancel()
				resp, err := userProviderClient.GetUsers(ctxWithTO, &pb.GetUsersRequest{UserId: []int32{1, 2}})
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(resp).NotTo(gomega.BeNil())
				gomega.Expect(len(resp.Users), 2)
				//fmt.Printf("consumer:%+v\n", resp.Users)
			})
		}, test.SampleConfig)

		experiment.Sample(func(idx int) {
			defer GinkgoRecover()

			experiment.MeasureDuration("GetUserByName", func() {
				ctxWithTO, cancel := context.WithTimeout(ctx, 5*time.Second)
				defer cancel()
				resp, err := userProviderClient.GetUserByName(ctxWithTO, &pb.GetUserByNameRequest{Name: "Kenway"})
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(resp).NotTo(gomega.BeNil())
				gomega.Expect(len(resp.Users), 1)
				//fmt.Printf("consumer:%+v\n", resp.Users)
			})
		}, test.SampleConfig)
	})

	It("pixiu to grpc protocol performance test", func() {
		experiment := gmeasure.NewExperiment("pixiu to grpc protocol performance test")
		AddReportEntry(experiment.Name, experiment)

		urlPrefix := "http://localhost:8881/api/v1/provider.UserProvider/"

		experiment.Sample(func(idx int) {
			defer GinkgoRecover()

			experiment.MeasureDuration("GetUser", func() {
				url := urlPrefix + "GetUser"
				data := `
{
	"userId": 1
}
`
				resp, err := http.Post(url, "application/json", strings.NewReader(data))
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				_, err = io.ReadAll(resp.Body)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(resp.Status, 200)
				//println(string(respBytes))
			})
		}, test.SampleConfig)

		experiment.Sample(func(idx int) {
			defer GinkgoRecover()

			experiment.MeasureDuration("GetUsers", func() {
				url := urlPrefix + "GetUsers"
				data := `
{
	"userId": [1, 2, 3]
}
`
				resp, err := http.Post(url, "application/json", strings.NewReader(data))
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				_, err = io.ReadAll(resp.Body)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(resp.Status, 200)
				//println(string(respBytes))
			})
		}, test.SampleConfig)

		experiment.Sample(func(idx int) {
			defer GinkgoRecover()

			experiment.MeasureDuration("GetUserByName", func() {
				url := urlPrefix + "GetUserByName"
				data := `
{
	"name": "Kenway"
}
`
				resp, err := http.Post(url, "application/json", strings.NewReader(data))
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				_, err = io.ReadAll(resp.Body)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(resp.Status, 200)
				//println(string(respBytes))
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
	serverProcess, err := gexec.Build("dubbo-go-pixiu-benchmark/protocol/grpc/go-server/cmd")
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	command := exec.Command(serverProcess)
	session, err := gexec.Start(command, io.Discard, io.Discard)
	//session, err := gexec.Start(command, os.Stdout, os.Stderr)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	return session
}

func prepareGRPCClient() {
	conn, err := grpc.Dial("localhost:50001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	userProviderClient = pb.NewUserProviderClient(conn)
}
