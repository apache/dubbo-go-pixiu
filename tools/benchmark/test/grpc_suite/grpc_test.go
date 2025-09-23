package grpc_suite

import (
	"context"
	"dubbo-go-pixiu-benchmark/test"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gmeasure"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega/gexec"

	pb "dubbo-go-pixiu-benchmark/protocol/grpc/proto"
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
				_, err = ioutil.ReadAll(resp.Body)
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
				_, err = ioutil.ReadAll(resp.Body)
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
				_, err = ioutil.ReadAll(resp.Body)
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
	session, err := gexec.Start(command, ioutil.Discard, ioutil.Discard)
	//session, err := gexec.Start(command, os.Stdout, os.Stderr)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	return session
}

func prepareGRPCClient() {
	conn, err := grpc.Dial("localhost:50001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	userProviderClient = pb.NewUserProviderClient(conn)
}
