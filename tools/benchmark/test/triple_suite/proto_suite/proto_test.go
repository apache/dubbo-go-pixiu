package proto_suite

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

import (
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	tripleConstant "github.com/dubbogo/triple/pkg/common/constant"
	. "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"github.com/onsi/gomega/gmeasure"
)

import (
	triplepb "dubbo-go-pixiu-benchmark/api"
	"dubbo-go-pixiu-benchmark/test"
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

		tripleServerSession = prepareTripleServer()
		pixiuSession = test.PreparePixiu("../../../dist/pixiu", test.CurPath+"/../../../protocol/triple/pb/pixiu/conf/config.yaml")

		time.Sleep(3 * time.Second)

		prepareTripleClient()

	})

	It("triple protocol performance test", func() {

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
				gomega.Expect(resp.Status).To(gomega.Equal("200 OK"))
				reply, err := ioutil.ReadAll(resp.Body)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(string(reply)).NotTo(gomega.MatchRegexp("client call err*"))
				//fmt.Printf("consumer:%+v\n", string(reply))
			})
		}, test.SampleConfig)
	})

	AfterAll(func() {
		tripleServerSession.Terminate().Wait(5 * time.Second)
		pixiuSession.Terminate().Wait(5 * time.Second)
	})
})

func prepareTripleServer() *gexec.Session {
	serverProcess, err := gexec.Build("dubbo-go-pixiu-benchmark/protocol/triple/pb/go-server/cmd")
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	command := exec.Command(serverProcess)
	session, err := gexec.Start(command, ioutil.Discard, ioutil.Discard)
	//session, err := gexec.Start(command, os.Stdout, os.Stderr)

	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	return session
}

func prepareTripleClient() {
	config.SetConsumerService(grpcGreeterImpl)
	err := config.Load(config.WithPath(test.CurPath + "/../../../protocol/triple/pb/go-client/conf/dubbogo.yml"))
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}
