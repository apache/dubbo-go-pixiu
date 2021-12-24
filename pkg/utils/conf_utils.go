package utils

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

const (
	// ipAddress regex :  ((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})(\.((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})){3}
	REGISTRY_ADDRESS_REGEX = "(zookeeper|nacos)://([0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3})\\:([0-9]{1,5})"
	ZK_ADDRESS             = "127.0.0.1:2181"
	NACOS_ADDRESS          = "127.0.0.1:8848"
)

func replaceConf(filePath, origin, target string) error {
	output, needHandle, err := readFile(filePath, origin, target)
	if err != nil {
		logger.Errorf("read file error!")
		return err
	}
	if needHandle {
		err = writeToFile("conf/conf.yaml", output)
		if err != nil {
			logger.Errorf("replaceConf error!")
			return err
		}
	}
	logger.Warnf("not need do something.")
	return nil
}

func adapterConf(address string) error {
	matched, _ := regexp.MatchString(REGISTRY_ADDRESS_REGEX, address)
	if matched {
		filePath := ""
		strArr := strings.Split(address, "://")
		switch strArr[0] {
		case "zookeeper":
			filePath = "conf/conf-zk.yaml"
			return replaceConf(filePath, ZK_ADDRESS, strArr[1])
		case "nacos":
			filePath = "conf/conf-nacos.yaml"
			return replaceConf(filePath, NACOS_ADDRESS, strArr[1])
		default:
			logger.Errorf("config address=%s error", address)
			return errors.New("registry address error")
		}
	}
	return errors.New("parse address error")
}

func getFiles(path string) []string {
	files := make([]string, 0)
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
	}
	return files
}

func readFile(filePath, origin, target string) ([]byte, bool, error) {
	f, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil {
		return nil, false, err
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	needHandle := false
	output := make([]byte, 0)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				return output, needHandle, nil
			}
			return nil, needHandle, err
		}
		if ok, _ := regexp.Match(origin, line); ok {
			reg := regexp.MustCompile(target)
			newByte := reg.ReplaceAll(line, []byte(target))
			output = append(output, newByte...)
			output = append(output, []byte("\n")...)
			if !needHandle {
				needHandle = true
			}
		} else {
			output = append(output, line...)
			output = append(output, []byte("\n")...)
		}
	}
	return output, needHandle, nil
}

func writeToFile(filePath string, outPut []byte) error {
	_, b := IsFile(filePath)
	var f *os.File
	var err error
	if b {
		f, _ = os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	} else {
		logger.Errorf("file [%s] not exists, create.", filePath)
		f, err = os.Create(filePath)
	}

	writer := bufio.NewWriter(f)
	_, err = writer.Write(outPut)
	if err != nil {
		return err
	}
	writer.Flush()
	return nil
}

func WriteFile(path string, str string) {
	_, b := IsFile(path)
	var f *os.File
	var err error
	if b {
		f, _ = os.OpenFile(path, os.O_APPEND, 0666)
	} else {
		// if not exists, create
		f, err = os.Create(path)
	}

	// close when end
	defer func() {
		err = f.Close()
		if err != nil {
			fmt.Println("err = ", err)
		}
	}()

	if err != nil {
		fmt.Println("err = ", err)
		return
	}
	_, err = f.WriteString(str)
	if err != nil {
		fmt.Println("err = ", err)
	}
}

func IsExists(path string) (os.FileInfo, bool) {
	f, err := os.Stat(path)
	return f, err == nil || os.IsExist(err)
}

func IsDir(path string) (os.FileInfo, bool) {
	f, flag := IsExists(path)
	return f, flag && f.IsDir()
}

func IsFile(path string) (os.FileInfo, bool) {
	f, flag := IsExists(path)
	return f, flag && !f.IsDir()
}
