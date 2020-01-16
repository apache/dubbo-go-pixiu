package logger

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

func TestInitLog(t *testing.T) {
	var (
		err  error
		path string
	)

	err = InitLog("")
	assert.EqualError(t, err, "log configure file name is nil")

	path, err = filepath.Abs("./log.xml")
	assert.NoError(t, err)
	err = InitLog(path)
	assert.EqualError(t, err, "log configure file name{"+path+"} suffix must be .yml")

	path, err = filepath.Abs("./logger.yml")
	assert.NoError(t, err)
	err = InitLog(path)
	var errMsg string
	if runtime.GOOS == "windows" {
		errMsg = fmt.Sprintf("open %s: The system cannot find the file specified.", path)
	} else {
		errMsg = fmt.Sprintf("open %s: no such file or directory", path)
	}
	assert.EqualError(t, err, fmt.Sprintf("ioutil.ReadFile(file:%s) = error:%s", path, errMsg))

	err = InitLog("./log.yml")
	assert.NoError(t, err)

	Debug("debug")
	Info("info")
	Warn("warn")
	Error("error")
	Debugf("%s", "debug")
	Infof("%s", "info")
	Warnf("%s", "warn")
	Errorf("%s", "error")
}

func TestSetLevel(t *testing.T) {
	err := InitLog("./log.yml")
	assert.NoError(t, err)
	Debug("debug")
	Info("info")

	assert.True(t, SetLoggerLevel("info"))
	Debug("debug")
	Info("info")

	SetLogger(GetLogger().(*DubboLogger).Logger)
	assert.False(t, SetLoggerLevel("debug"))
	Debug("debug")
	Info("info")
}
