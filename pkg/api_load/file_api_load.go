package api_load

import (
	"encoding/json"
	"io/ioutil"
)

import (
	"github.com/fsnotify/fsnotify"
	"github.com/ghodss/yaml"
)

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/logger"
	"github.com/dubbogo/dubbo-go-proxy/pkg/model"
)

func init() {
	var _ ApiLoad = new(FileApiLoader)
}

type FileApiLoader struct {
	filePath string
	//locker   sync.Mutex
	ApiConfigs []model.Api
}

func (f *FileApiLoader) InitLoad() (err error) {
	//f.locker.Lock()
	//defer f.locker.Unlock()
	content, err := ioutil.ReadFile(f.filePath)
	if err != nil {
		logger.Errorf("fileApiLoader read api config error:%v", err)
		return
	}
	jsonBytes, err := yaml.YAMLToJSON(content)
	if err != nil {
		logger.Errorf("fileApiLoader transfer api config error:%v,is it yaml format?", err)
		return
	}

	err = json.Unmarshal(jsonBytes, &f.ApiConfigs)
	if err != nil {
		logger.Errorf("fileApiLoader read api config error:%v", err)
		return
	}
	if len(f.ApiConfigs) < 1 {
		logger.Warnf("no api loaded!please make sure api config file is located")
	}
	//f.ApiConfigs = apiConfigs
	return
}

func (f *FileApiLoader) HotLoad() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logger.Error(err)
		return err
	}

	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				logger.Debug("event:", event)
				switch event.Op {
				case fsnotify.Write:
					logger.Debug("modified file:", event.Name)
					content, err := ioutil.ReadFile(f.filePath)
					jsonBytes, err := yaml.YAMLToJSON(content)
					if err != nil {
						logger.Warnf("fileApiLoader transfer api config error:%v,is it yaml format?", err)
						break
					}

					err = json.Unmarshal(jsonBytes, &f.ApiConfigs)
					if err != nil {
						logger.Warnf("fileApiLoader read api config error:%v", err)
						break
					}
					break
				case fsnotify.Remove:
					logger.Debug("removed file:", event.Name)
					f.ApiConfigs = nil
					break
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				logger.Error("error:", err)
			}
		}
	}()

	err = watcher.Add(f.filePath)
	if err != nil {
		logger.Error(err)
	}
	return err
}

func (f *FileApiLoader) Clear() error {
	f.ApiConfigs = nil
	return nil
}
