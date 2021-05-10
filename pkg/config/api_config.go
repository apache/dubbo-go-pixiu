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

// LoadAPIConfig load the api config from config center
func LoadAPIConfig(metaConfig *model.APIMetaConfig) (*fc.APIConfig, error) {
	client, _ = etcdv3.NewConfigClientWithErr(
		etcdv3.WithName(etcdv3.RegistryETCDV3Client),
		etcdv3.WithTimeout(10*time.Second),
		etcdv3.WithEndpoints(strings.Split(metaConfig.Address, ",")...),
	)

	content, rev, err := client.GetValAndRev(metaConfig.APIConfigPath)
	if err != nil {
		return nil, perrors.Errorf("Get remote config fail error %v", err)
	}

	if err = initAPIConfigFromString(content); err != nil {
		return nil, err
	}

	// 启动service和plugin group 的监听
	//go listenAPIConfigNodeEvent(metaConfig.APIConfigPath, rev)
	go listenServiceNodeEvent(getEtcdServicePath(metaConfig.APIConfigPath), rev)
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

func listenServiceNodeEvent(key string, rev int64) bool {
	for {
		wc, err := client.WatchWithOption(key, clientv3.WithRev(rev), clientv3.WithPrefix())
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
					handlePutEvent(event.Kv.Key, event.Kv.Value)
					//if err = initAPIConfigFromString(string(event.Kv.Value)); err == nil {
					//	listener.APIConfigChange(GetAPIConf())
					//}
				case mvccpb2.DELETE:
					logger.Warnf("get event (key{%s}) = event{EventNodeDeleted}", event.Kv.Key)
					handleDeleteEvent(event.Kv.Key, event.Kv.Value)
					return true
				default:
					return false
				}
			}
		}
	}
}

func getCheckResourceRegexp() *regexp.Regexp {
	return regexp.MustCompile(".+/Resources/[^/]+/?$")
}

func getExtractMethodRegexp() *regexp.Regexp {
	return regexp.MustCompile("Resources/([^/]+)/Method")
}


func handleDeleteEvent(key, val []byte) {
	re := getCheckResourceRegexp()
	matchResource := re.Match(key)

	if matchResource {
		res := fc.Resource{}
		err := yaml.UnmarshalYML(val, res)
		if err != nil {
			logger.Error("handlePutEvent UnmarshalYML error %v", err)
			return
		}
		mergeApiConfigResource(res)
	} else {
		res := fc.Method{}
		err := yaml.UnmarshalYML(val, res)
		if err != nil {
			logger.Error("handlePutEvent UnmarshalYML error %v", err)
			return
		}

		reExtract := getExtractMethodRegexp()

		result := reExtract.FindStringSubmatch(strings.Replace(string(key), "[pixiu]", "/", -1))
		if len(result) != 2 {
			return
		}
		fullPath := result[1]
		mergeApiConfigMethod(fullPath, res)
	}
}


func handlePutEvent(key, val []byte) {
	re := getCheckResourceRegexp()
	matchResource := re.Match(key)

	if matchResource {
		res := fc.Resource{}
		err := yaml.UnmarshalYML(val, res)
		if err != nil {
			logger.Error("handlePutEvent UnmarshalYML error %v", err)
			return
		}
		mergeApiConfigResource(res)
	} else {

		res := fc.Method{}
		err := yaml.UnmarshalYML(val, res)
		if err != nil {
			logger.Error("handlePutEvent UnmarshalYML error %v", err)
			return
		}

		reExtract := getExtractMethodRegexp()

		result := reExtract.FindStringSubmatch(strings.Replace(string(key), "[pixiu]", "/", -1))
		if len(result) != 2 {
			return
		}
		fullPath := result[1]
		mergeApiConfigMethod(fullPath, res)
	}
}

func mergeApiConfigResource(val fc.Resource) {
	for i, resource := range apiConfig.Resources {
		if val.Path != resource.Path {
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

func mergeApiConfigMethod(path string, val fc.Method) {
	for _, resource := range apiConfig.Resources {
		if path != resource.Path {
			continue
		}

		for j, method := range resource.Methods {
			if method.HTTPVerb == val.HTTPVerb {
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
