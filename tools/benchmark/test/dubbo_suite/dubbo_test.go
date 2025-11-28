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

package dubbo_suite

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
	"dubbo-go-pixiu-benchmark/protocol/dubbo/go-client/pkg"

	"dubbo-go-pixiu-benchmark/test"

	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"

	hessian "github.com/apache/dubbo-go-hessian2"

	. "github.com/onsi/ginkgo/v2"

	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"github.com/onsi/gomega/gmeasure"
)

var (
	userProvider                     = &pkg.UserProvider{}
	dubboServerSession, pixiuSession *gexec.Session
)

func TestDubboCases(t *testing.T) {
	gomega.RegisterFailHandler(Fail)
	RunSpecs(t, "test")
}

var _ = Describe("test", Ordered, func() {
	BeforeAll(func() {
		var err error
		test.CurPath, err = os.Getwd()
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Wait for ports to be available before starting servers
		waitForPortAvailable("20000", 10*time.Second)
		waitForPortAvailable("8881", 10*time.Second)

		dubboServerSession = prepareDubboServer()
		// Wait for dubbo server to register to Zookeeper
		time.Sleep(5 * time.Second)

		pixiuSession = test.PreparePixiu("../../dist/pixiu", test.CurPath+"/../../protocol/dubbo/pixiu/conf/config.yaml")
		// Wait for pixiu to discover services from Zookeeper
		time.Sleep(5 * time.Second)
	})

	//TODO(kenwaycai): output to external file

	It("pixiu to dubbo protocol performance test", func() {

		urlPrefix := "http://localhost:8881/dubbo.io/org.apache.dubbo.sample.UserProvider/%s"

		experiment := gmeasure.NewExperiment("pixiu to dubbo protocol performance test")
		AddReportEntry(experiment.Name, experiment)

		experiment.Sample(func(idx int) {
			experiment.MeasureDuration("GetUser", func() {
				defer GinkgoRecover()

				url := fmt.Sprintf(urlPrefix, "GetUser")
				data := `
{
    "types": "object",
    "values": {
        "id": "003"
    }
}
`

				resp, err := http.Post(url, "application/json", strings.NewReader(data))
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				defer resp.Body.Close()
				gomega.Expect(resp.Status).To(gomega.Equal("200 OK"))
				_, err = io.ReadAll(resp.Body)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
			})

			experiment.MeasureDuration("GetGender", func() {
				defer GinkgoRecover()

				url := fmt.Sprintf(urlPrefix, "GetGender")
				data := `
{
    "types": "int",
    "values": 1
}
`
				resp, err := http.Post(url, "application/json", strings.NewReader(data))
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				defer resp.Body.Close()
				gomega.Expect(resp.Status).To(gomega.Equal("200 OK"))

				_, err = io.ReadAll(resp.Body)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
			})

			experiment.MeasureDuration("GetUser0", func() {
				defer GinkgoRecover()

				url := fmt.Sprintf(urlPrefix, "GetUser0")
				data := `
{
    "types": "string,string",
    "values": [
        "003",
        "Moorse"
    ]
}
`
				resp, err := http.Post(url, "application/json", strings.NewReader(data))
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				defer resp.Body.Close()
				gomega.Expect(resp.Status).To(gomega.Equal("200 OK"))

				_, err = io.ReadAll(resp.Body)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
			})

			experiment.MeasureDuration("GetUsers", func() {
				defer GinkgoRecover()

				url := fmt.Sprintf(urlPrefix, "GetUsers")
				data := `
{
    "types": "string",
    "values": [
        [
            "003",
            "002"
        ]
    ]
}
`

				resp, err := http.Post(url, "application/json", strings.NewReader(data))
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				defer resp.Body.Close()
				gomega.Expect(resp.Status).To(gomega.Equal("200 OK"))
				_, err = io.ReadAll(resp.Body)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
			})

		}, test.SampleConfig)
	})

	It("dubbo protocol performance test", func() {
		defer GinkgoRecover()

		// Initialize dubbo client only when needed for this test
		prepareDubboClient()

		experiment := gmeasure.NewExperiment("dubbo protocol performance test")
		AddReportEntry(experiment.Name, experiment)

		experiment.Sample(func(idx int) {
			reqUser := &pkg.User{}
			reqUser.ID = "003"
			experiment.MeasureDuration("GetUser", func() {
				_, err := userProvider.GetUser(context.TODO(), reqUser)
				gomega.Expect(err).To(gomega.Succeed())
			})
		}, test.SampleConfig)

		experiment.Sample(func(idx int) {
			experiment.MeasureDuration("GetGender", func() {
				_, err := userProvider.GetGender(context.TODO(), 1)
				gomega.Expect(err).To(gomega.Succeed())
			})
		}, test.SampleConfig)

		experiment.Sample(func(idx int) {
			experiment.MeasureDuration("GetUser0", func() {
				_, err := userProvider.GetUser0("003", "Moorse")
				gomega.Expect(err).To(gomega.Succeed())
			})
		}, test.SampleConfig)

		experiment.Sample(func(idx int) {
			experiment.MeasureDuration("GetUsers", func() {
				_, err := userProvider.GetUsers([]string{"002", "003"})
				gomega.Expect(err).To(gomega.Succeed())
			})
		}, test.SampleConfig)

		experiment.Sample(func(idx int) {
			experiment.MeasureDuration("GetUser2", func() {
				var i int32 = 1
				_, err := userProvider.GetUser2(context.TODO(), i)
				gomega.Expect(err).To(gomega.Succeed())
			})
		}, test.SampleConfig)

		experiment.Sample(func(idx int) {
			experiment.MeasureDuration("GetErr", func() {
				reqUser := &pkg.User{}
				reqUser.ID = "003"
				_, err := userProvider.GetErr(context.TODO(), reqUser)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		}, test.SampleConfig)
	})

	AfterAll(func() {
		// Use Kill instead of Terminate to ensure processes are fully stopped
		if pixiuSession != nil {
			pixiuSession.Kill().Wait(10 * time.Second)
		}
		if dubboServerSession != nil {
			dubboServerSession.Kill().Wait(10 * time.Second)
		}
		gexec.CleanupBuildArtifacts()
		// Wait for ports to be released
		time.Sleep(2 * time.Second)
	})

})

func prepareDubboServer() *gexec.Session {
	serverProcess, err := gexec.Build("dubbo-go-pixiu-benchmark/protocol/dubbo/go-server/cmd")
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	command := exec.Command(serverProcess)
	// Set working directory to test/dubbo_suite so that relative paths in server.go work correctly
	command.Dir = test.CurPath
	session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
	//session, err := gexec.Start(command, os.Stdout, os.Stderr)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	return session
}

func prepareDubboClient() {
	hessian.RegisterJavaEnum(pkg.Gender(pkg.MAN))
	hessian.RegisterJavaEnum(pkg.Gender(pkg.WOMAN))
	hessian.RegisterPOJO(&pkg.User{})

	config.SetConsumerService(userProvider)

	err := config.Load(config.WithPath(test.CurPath + "/../../protocol/dubbo/go-client/conf/dubbogo.yml"))
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

// waitForPortAvailable waits until the port is available (not in use)
func waitForPortAvailable(port string, timeout time.Duration) {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", "127.0.0.1:"+port, 100*time.Millisecond)
		if err != nil {
			// Port is available (connection refused means no one is listening)
			if conn != nil {
				conn.Close()
			}
			return
		}
		conn.Close()
		time.Sleep(500 * time.Millisecond)
	}
	gomega.Expect(fmt.Errorf("port %s is still in use after %v", port, timeout)).NotTo(gomega.HaveOccurred())
}
