package logger

import (
	"io/ioutil"
	"log"
	"path"
	"sync"
)

import (
	perrors "github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v2"
)

var (
	logger Logger
)

type DubbogoProxyLogger struct {
	mutex sync.Mutex
	Logger
	dynamicLevel zap.AtomicLevel
	// disable presents the logger state. if disable is true, the logger will write nothing
	// the default value is false
	disable bool
}

type Logger interface {
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Debug(args ...interface{})

	Infof(fmt string, args ...interface{})
	Warnf(fmt string, args ...interface{})
	Errorf(fmt string, args ...interface{})
	Debugf(fmt string, args ...interface{})
}

func init() {
	logConfFile := "./conf/log.yml"
	err := InitLog(logConfFile)
	if err != nil {
		log.Printf("[InitLog] warn: %v", err)
	}
}

func InitLog(logConfFile string) error {
	if logConfFile == "" {
		InitLogger(nil)
		return perrors.New("log configure file name is nil")
	}
	if path.Ext(logConfFile) != ".yml" {
		InitLogger(nil)
		return perrors.Errorf("log configure file name{%s} suffix must be .yml", logConfFile)
	}

	confFileStream, err := ioutil.ReadFile(logConfFile)
	if err != nil {
		InitLogger(nil)
		return perrors.Errorf("ioutil.ReadFile(file:%s) = error:%v", logConfFile, err)
	}

	conf := &zap.Config{}
	err = yaml.Unmarshal(confFileStream, conf)
	if err != nil {
		InitLogger(nil)
		return perrors.Errorf("[Unmarshal]init logger error: %v", err)
	}

	InitLogger(conf)

	return nil
}

func InitLogger(conf *zap.Config) {
	var zapLoggerConfig zap.Config
	if conf == nil {
		zapLoggerConfig = zap.NewDevelopmentConfig()
		zapLoggerEncoderConfig := zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}
		zapLoggerConfig.EncoderConfig = zapLoggerEncoderConfig
	} else {
		zapLoggerConfig = *conf
	}
	zapLogger, _ := zapLoggerConfig.Build(zap.AddCallerSkip(1))
	//logger = zapLogger.Sugar()
	logger = &DubbogoProxyLogger{Logger: zapLogger.Sugar(), dynamicLevel: zapLoggerConfig.Level}
}

func SetLogger(log Logger) {
	logger = log
}

func GetLogger() Logger {
	return logger
}

func SetLoggerLevel(level string) bool {
	if l, ok := logger.(OpsLogger); ok {
		l.SetLoggerLevel(level)
		return true
	}
	return false
}

type OpsLogger interface {
	Logger
	SetLoggerLevel(level string)
}

func (dpl *DubbogoProxyLogger) SetLoggerLevel(level string) {
	l := new(zapcore.Level)
	l.Set(level)
	dpl.dynamicLevel.SetLevel(*l)
}
