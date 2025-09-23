package test

import (
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"github.com/onsi/gomega/gmeasure"
	"io/ioutil"
	"os/exec"
	"time"
)

var (
	CurPath      string
	SampleConfig = gmeasure.SamplingConfig{
		N:           500,
		Duration:    300 * time.Second,
		NumParallel: 10,
	}
)

func PreparePixiu(pixiu, path string) *gexec.Session {
	command := exec.Command(pixiu, "gateway", "start", "-c", path)
	session, err := gexec.Start(command, ioutil.Discard, ioutil.Discard)
	//session, err := gexec.Start(command, os.Stdout, os.Stderr)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	return session
}
