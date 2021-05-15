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

package config

import (
	"go.etcd.io/etcd/clientv3"
	mvccpb2 "go.etcd.io/etcd/mvcc/mvccpb"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

import (
	fc "github.com/dubbogo/dubbo-go-pixiu-filter/pkg/api/config"
	etcdv3 "github.com/dubbogo/gost/database/kv/etcd/v3"
	perrors "github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/yaml"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

var (
	apiConfig *fc.APIConfig
	once      sync.Once
	client    *etcdv3.Client
	listener  APIConfigResourceListener
	lock      sync.RWMutex
)

// APIConfigResourceListener defines api config listener interface
type APIConfigResourceListener interface {
	ResourceChange(new fc.Resource, old fc.Resource) bool // bool is return for interface implement is interesting
	ResourceAdd(res fc.Resource) bool
	ResourceDelete(deleted fc.Resource) bool
	MethodChange(res fc.Resource, method fc.Method, old fc.Method) bool
	MethodAdd(res fc.Resource, method fc.Method) bool
	MethodDelete(res fc.Resource, method fc.Method) bool
}

// LoadAPIConfigFromFile load the api config from file
func LoadAPIConfigFromFile(path string) (*fc.APIConfig, error) {
	if len(path) == 0 {
		return nil, perrors.Errorf("Config file not specified")
	}
	logger.Infof("Load API configuration file form %s", path)
	apiConf := &fc.APIConfig{}
	err := yaml.UnmarshalYMLConfig(path, apiConf)
	if err != nil {
		return nil, perrors.Errorf("unmarshalYmlConfig error %v", perrors.WithStack(err))
	}
	apiConfig = apiConf

	return apiConf, nil
}

func getConfigPath(root string) string {
	return root + "/Resources"
}

// LoadAPIConfig load the api config from config center
func LoadAPIConfig(metaConfig *model.APIMetaConfig) (*fc.APIConfig, error) {
	client, _ = etcdv3.NewConfigClientWithErr(
		etcdv3.WithName(etcdv3.RegistryETCDV3Client),
		etcdv3.WithTimeout(10*time.Second),
		etcdv3.WithEndpoints(strings.Split(metaConfig.Address, ",")...),
	)

	kList, vList, err := client.GetChildren(metaConfig.APIConfigPath)
	if err != nil {
		return nil, perrors.Errorf("Get remote config fail error %v", err)
	}

	if err = initAPIConfigFromKVList(kList, vList); err != nil {
		return nil, err
	}

	// 启动service和plugin group 的监听
	//go listenAPIConfigNodeEvent(metaConfig.APIConfigPath, rev)
	go listenServiceNodeEvent(metaConfig.APIConfigPath)
	//go listenPluginNodeEvent(metaConfig.APIConfigPath, rev)
	return apiConfig, nil
}

func getEtcdServicePath(root string) string {
	return root + "/Resources/"
}

func initAPIConfigFromString(content string) error {
	lock.Lock()
	defer lock.Unlock()

	apiConf := &fc.APIConfig{}
	if len(content) != 0 {
		err := yaml.UnmarshalYML([]byte(content), apiConf)
		if err != nil {
			return perrors.Errorf("unmarshalYmlConfig error %v", perrors.WithStack(err))
		}

		valid := validateAPIConfig(apiConf)
		if !valid {
			return perrors.Errorf("api config not valid error %v", perrors.WithStack(err))
		}

		apiConfig = apiConf
	}
	return nil
}

func initAPIConfigFromKVList(kList, vList []string) error {
	lock.Lock()
	defer lock.Unlock()

	tmpApiConf := &fc.APIConfig{}

	for i, k := range kList {
		v := vList[i]

		re := getCheckBaseInfoRegexp()

		// base info
		if m := re.Match([]byte(k)); m {
			continue
		}

		// resource
		re = getCheckResourceRegexp()
		if m := re.Match([]byte(k)); m {
			resource := &fc.Resource{}
			err := yaml.UnmarshalYML([]byte(v), resource)
			if err != nil {
				logger.Error("unmarshalYmlConfig error %v", err.Error())
				continue
			}

			// for 循环查找
			found := false
			if tmpApiConf.Resources == nil {
				tmpApiConf.Resources = make([]fc.Resource, 0)
			}

			for i, old := range tmpApiConf.Resources {
				if old.Path != resource.Path {
					continue
				}

				// modify one resource
				resource.Methods = old.Methods
				tmpApiConf.Resources[i] = *resource
				found = true
			}

			if !found {
				tmpApiConf.Resources = append(tmpApiConf.Resources, *resource)
			}
			continue
		}

		re = getExtractMethodRegexp()

		if m := re.Match([]byte(k)); m {

			method := &fc.Method{}
			err := yaml.UnmarshalYML([]byte(v), method)
			if err != nil {
				logger.Error("unmarshalYmlConfig error %v", err.Error())
				continue
			}

			found := false
			for r, resource := range tmpApiConf.Resources {
				if method.ResourcePath != resource.Path {
					continue
				}

				for j, old := range resource.Methods {
					if old.HTTPVerb == method.HTTPVerb {
						// modify one method
						resource.Methods[j] = *method
						found = true
					}
				}
				if !found {
					resource.Methods = append(resource.Methods, *method)
					tmpApiConf.Resources[r] = resource
					found = true
				}
			}

			// need add resource first
			if !found {
				resource := &fc.Resource{}
				resource.Methods = append(resource.Methods, *method)
				resource.Path = method.ResourcePath
				tmpApiConf.Resources = append(tmpApiConf.Resources, *resource)
			}

		}
	}
	apiConfig = tmpApiConf
	return nil
}

// validateAPIConfig check api config valid
func validateAPIConfig(conf *fc.APIConfig) bool {
	if conf.Name == "" {
		return false
	}
	if conf.Description == "" {
		return false
	}
	if conf.Resources == nil || len(conf.Resources) == 0 {
		return false
	}
	return true
}

func listenServiceNodeEvent(key string) bool {
	for {
		wc, err := client.WatchWithOption(key, clientv3.WithPrefix())
		if err != nil {
			logger.Warnf("Watch api config {key:%s} = error{%v}", key, err)
			return false
		}

		select {

		// client stopped
		case <-client.Done():
			logger.Warnf("client stopped")
			return false
		// client ctx stop
		// handle etcd events
		case e, ok := <-wc:
			if !ok {
				logger.Warnf("watch-chan closed")
				return false
			}

			if e.Err() != nil {
				logger.Errorf("watch ERR {err: %s}", e.Err())
				continue
			}
			for _, event := range e.Events {
				switch event.Type {
				case mvccpb2.PUT:
					logger.Infof("get event (key{%s}) = event{EventNodePut}", event.Kv.Key)
					handlePutEvent(event.Kv.Key, event.Kv.Value)
					//if err = initAPIConfigFromString(string(event.Kv.Value)); err == nil {
					//	listener.APIConfigChange(GetAPIConf())
					//}
				case mvccpb2.DELETE:
					logger.Infof("get event (key{%s}) = event{EventNodeDeleted}", event.Kv.Key)
					handleDeleteEvent(event.Kv.Key, event.Kv.Value)
					return true
				default:
					return false
				}
			}
		}
	}
}

func getCheckBaseInfoRegexp() *regexp.Regexp {
	return regexp.MustCompile(".+/base$")
}

func getCheckResourceRegexp() *regexp.Regexp {
	return regexp.MustCompile(".+/Resources/[^/]+/?$")
}

func getExtractMethodRegexp() *regexp.Regexp {
	return regexp.MustCompile("Resources/([^/]+)/Method/[^/]+/?$")
}

func handleDeleteEvent(key, val []byte) {

	keyStr := string(key)
	keyStr = strings.TrimSuffix(keyStr, "/")

	re := getCheckResourceRegexp()
	if m := re.Match(key); m {
		pathArray := strings.Split(keyStr, "/")
		if len(pathArray) == 0 {
			logger.Errorf("handleDeleteEvent key format error")
			return
		}
		resourceIdStr := pathArray[len(pathArray)-1]
		ID, err := strconv.Atoi(resourceIdStr)
		if err != nil {
			logger.Error("handleDeleteEvent ID is not int error %v", err)
			return
		}
		deleteApiConfigResource(ID)
		return
	}

	re = getExtractMethodRegexp()
	if m := re.Match(key); m {
		pathArray := strings.Split(keyStr, "/")
		if len(pathArray) < 3 {
			logger.Errorf("handleDeleteEvent key format error")
			return
		}
		resourceIdStr := pathArray[len(pathArray)-3]
		resourceId, err := strconv.Atoi(resourceIdStr)
		if err != nil {
			logger.Error("handleDeleteEvent ID is not int error %v", err)
			return
		}

		methodIdStr := pathArray[len(pathArray)-1]
		methodId, err := strconv.Atoi(methodIdStr)
		if err != nil {
			logger.Error("handleDeleteEvent ID is not int error %v", err)
			return
		}
		deleteApiConfigMethod(resourceId, methodId)
	}
}

func handlePutEvent(key, val []byte) {
	re := getCheckResourceRegexp()
	if m := re.Match(key); m {
		res := &fc.Resource{}
		err := yaml.UnmarshalYML(val, res)
		if err != nil {
			logger.Error("handlePutEvent UnmarshalYML error %v", err)
			return
		}
		mergeApiConfigResource(*res)
		return
	}

	re = getExtractMethodRegexp()
	if m := re.Match(key); m {

		res := &fc.Method{}
		err := yaml.UnmarshalYML(val, res)
		if err != nil {
			logger.Error("handlePutEvent UnmarshalYML error %v", err)
			return
		}

		mergeApiConfigMethod(res.ResourcePath, *res)
	}
}

func deleteApiConfigResource(resourceId int) {
	for i := 0; i < len(apiConfig.Resources); i++ {
		itr := apiConfig.Resources[i]
		if itr.ID == resourceId {
			apiConfig.Resources = append(apiConfig.Resources[:i], apiConfig.Resources[i+1:]...)
			listener.ResourceDelete(itr)
			return
		}
	}
}

func mergeApiConfigResource(val fc.Resource) {
	for i, resource := range apiConfig.Resources {
		if val.ID != resource.ID {
			continue
		}
		// modify one resource
		val.Methods = resource.Methods
		apiConfig.Resources[i] = val

		listener.ResourceChange(val, resource)

		return
	}
	// add one resource
	apiConfig.Resources = append(apiConfig.Resources, val)
	listener.ResourceAdd(val)
}

func deleteApiConfigMethod(resourceId, methodId int) {
	for _, resource := range apiConfig.Resources {
		if resource.ID != resourceId {
			continue
		}

		for i := 0; i < len(resource.Methods); i++ {
			method := resource.Methods[i]

			if method.ID == methodId {
				resource.Methods = append(resource.Methods[:i], resource.Methods[i+1:]...)
				listener.MethodDelete(resource, method)
				return
			}
		}
	}
}

func mergeApiConfigMethod(path string, val fc.Method) {
	for _, resource := range apiConfig.Resources {
		if path != resource.Path {
			continue
		}

		for j, method := range resource.Methods {
			if method.ID == val.ID {
				// modify one method
				resource.Methods[j] = val
				listener.MethodChange(resource, val, method)
				return
			}
		}
		// add one method
		resource.Methods = append(resource.Methods, val)
		listener.MethodAdd(resource, val)
	}
}

// RegisterConfigListener register APIConfigListener
func RegisterConfigListener(li APIConfigResourceListener) {
	listener = li
}

// GetAPIConf returns the init api config
func GetAPIConf() fc.APIConfig {
	return *apiConfig
}
