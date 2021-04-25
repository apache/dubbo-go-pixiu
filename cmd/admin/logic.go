package main

import (
	"strings"
)

import (
	fc "github.com/dubbogo/dubbo-go-pixiu-filter/pkg/api/config"
	perrors "github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/yaml"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

// BizSetBaseInfo business layer create base info
func BizSetBaseInfo(info *BaseInfo, created bool) error {
	// validate the api config

	data, _ := yaml.MarshalYML(info)

	if created {
		setErr := Client.Create(getRootPath(Base), string(data))

		if setErr != nil {
			logger.Warnf("update etcd error, %v\n", setErr)
			return perrors.WithMessage(setErr, "BizSetBaseInfo error")
		}
	} else {
		setErr := Client.Update(getRootPath(Base), string(data))

		if setErr != nil {
			logger.Warnf("update etcd error, %v\n", setErr)
			return perrors.WithMessage(setErr, "BizSetBaseInfo error")
		}
	}

	return nil
}

// BizGetBaseInfo business layer get base info
func BizGetBaseInfo() (*BaseInfo, error) {
	config, err := Client.Get(getRootPath(Base))
	if err != nil {
		logger.Errorf("GetBaseInfo err, %v\n", err)
		return nil, perrors.WithMessage(err, "BizGetBaseInfo error")
	}
	data := &BaseInfo{}
	_ = yaml.UnmarshalYML([]byte(config), data)

	return data, nil
}

// BizGetResourceDetail business layer get resource detail
func BizGetResourceDetail(path string) (string, error) {
	key := getResourceKey(path)
	detail, err := Client.Get(key)
	if err != nil {
		logger.Errorf("BizGetResourceDetail err, %v\n", err)
		return "", perrors.WithMessage(err, "BizGetResourceDetail error")
	}
	return detail, nil
}

// BizGetMethodDetail business layer get method detail
func BizGetMethodDetail(path string, method string) (string, error) {
	key := getMethodKey(path, method)
	detail, err := Client.Get(key)
	if err != nil {
		logger.Errorf("BizGetResourceDetail err, %v\n", err)
		return "", perrors.WithMessage(err, "BizGetResourceDetail error")
	}
	return detail, nil
}

// BizGetResourceList business layer get resource list
func BizGetResourceList() ([]fc.Resource, error) {
	_, vList, err := Client.GetChildrenKVList(getRootPath(Resources))
	if err != nil {
		logger.Errorf("GetResourceList err, %v\n", err)
		return nil, perrors.WithMessage(err, "BizGetResourceList error")
	}

	var ret []fc.Resource
	for _, v := range vList {
		res := &fc.Resource{}
		err := yaml.UnmarshalYML([]byte(v), res)
		if err != nil {
			logger.Errorf("UnmarshalYML err, %v\n", err)
		}
		ret = append(ret, *res)
	}

	return ret, nil
}

// BizGetMethodList business layer get method list
func BizGetMethodList(path string) ([]fc.Method, error) {
	key := getResourceMethodPrefixKey(path)

	_, vList, err := Client.GetChildrenKVList(key)
	if err != nil {
		logger.Errorf("GetResourceList err, %v\n", err)
		return nil, perrors.WithMessage(err, "BizGetResourceList error")
	}

	var ret []fc.Method
	for _, v := range vList {
		res := &fc.Method{}
		err := yaml.UnmarshalYML([]byte(v), res)
		if err != nil {
			logger.Errorf("UnmarshalYML err, %v\n", err)
		}
		ret = append(ret, *res)
	}

	return ret, nil
}

// BizGetPluginGroupList business layer get plugin group list
func BizGetPluginGroupList() ([]fc.PluginsGroup, error) {
	key := getPluginGroupPrefixKey()

	_, vList, err := Client.GetChildrenKVList(key)
	if err != nil {
		logger.Errorf("GetResourceList err, %v\n", err)
		return nil, perrors.WithMessage(err, "BizGetResourceList error")
	}

	var ret []fc.PluginsGroup
	for _, v := range vList {
		res := &fc.PluginsGroup{}
		err := yaml.UnmarshalYML([]byte(v), res)
		if err != nil {
			logger.Errorf("UnmarshalYML err, %v\n", err)
		}
		ret = append(ret, *res)
	}

	return ret, nil
}

// BizGetPluginGroupDetail business layer get plugin group detail
func BizGetPluginGroupDetail(name string) (string, error) {
	key := getPluginGroupKey(name)
	detail, err := Client.Get(key)
	if err != nil {
		logger.Errorf("BizGetResourceDetail err, %v\n", err)
		return "", perrors.WithMessage(err, "BizGetResourceDetail error")
	}
	return detail, nil
}

func getResourceKey(path string) string {
	return getRootPath(Resources) + "/" + strings.Replace(path, "/", "_", -1)
}

func getPluginGroupKey(name string) string {
	return getPluginGroupPrefixKey() + "/" + name
}

func getPluginGroupPrefixKey() string {
	return getRootPath(PluginGroup)
}

func getResourceMethodPrefixKey(path string) string {
	return getResourceKey(path) + "/" + "method"
}

func getMethodKey(path string, method string) string {
	return getResourceMethodPrefixKey(path) + "/" + method
}

// BizSetResourceInfo business layer create resource
func BizSetResourceInfo(res *fc.Resource, created bool) error {

	methods := res.Methods
	// 创建 methods
	BizCreateResourceMethod(getResourceMethodPrefixKey(res.Path), methods)
	// 创建resource
	res.Methods = nil
	data, _ := yaml.MarshalYML(res)
	if created {
		setErr := Client.Create(getResourceKey(res.Path), string(data))

		if setErr != nil {
			logger.Warnf("update etcd error, %v\n", setErr)
			return perrors.WithMessage(setErr, "BizSetResourceInfo error")
		}
	} else {
		setErr := Client.Update(getResourceKey(res.Path), string(data))

		if setErr != nil {
			logger.Warnf("update etcd error, %v\n", setErr)
			return perrors.WithMessage(setErr, "BizSetResourceInfo error")
		}
	}

	return nil
}

// BizSetPluginGroupInfo create plugin group
func BizSetPluginGroupInfo(res *fc.PluginsGroup, created bool) error {

	data, _ := yaml.MarshalYML(res)
	if created {
		setErr := Client.Create(getPluginGroupKey(res.GroupName), string(data))

		if setErr != nil {
			logger.Warnf("create etcd error, %v\n", setErr)
			return perrors.WithMessage(setErr, "BizSetPluginGroupInfo error")
		}
	} else {
		setErr := Client.Update(getPluginGroupKey(res.GroupName), string(data))

		if setErr != nil {
			logger.Warnf("update etcd error, %v\n", setErr)
			return perrors.WithMessage(setErr, "BizSetPluginGroupInfo error")
		}
	}

	return nil
}

// BizDeleteResourceInfo business layer delete resource
func BizDeleteResourceInfo(path string) error {
	key := getResourceKey(path)
	err := Client.Delete(key)
	if err != nil {
		logger.Warnf("BizDeleteResourceInfo, %v\n", err)
		return perrors.WithMessage(err, "BizDeleteResourceInfo error")
	}
	return nil
}

// BizDeleteMethodInfo business layer delete method
func BizDeleteMethodInfo(path string, method string) error {
	key := getMethodKey(path, method)
	err := Client.Delete(key)
	if err != nil {
		logger.Warnf("BizDeleteMethodInfo, %v\n", err)
		return perrors.WithMessage(err, "BizDeleteMethodInfo error")
	}
	return nil
}

// BizDeletePluginGroupInfo business layer delete plugin group
func BizDeletePluginGroupInfo(name string) error {
	key := getPluginGroupKey(name)
	err := Client.Delete(key)
	if err != nil {
		logger.Warnf("BizDeletePluginGroupInfo, %v\n", err)
		return perrors.WithMessage(err, "BizDeletePluginGroupInfo error")
	}
	return nil
}

// BizCreateResourceMethod batch create method below specific path
func BizCreateResourceMethod(path string, methods []fc.Method) error {

	if len(methods) == 0 {
		return nil
	}

	var kList, vList []string

	for _, method := range methods {
		kList = append(kList, getMethodKey(path, string(method.HTTPVerb)))
		data, _ := yaml.MarshalYML(method)
		vList = append(vList, string(data))
	}

	err := Client.BatchCreate(kList, vList)
	if err != nil {
		logger.Warnf("update etcd error, %v\n", err)
		return perrors.WithMessage(err, "BizCreateResourceMethod error")
	}
	return nil
}

// BizSetResourceMethod batch create method below specific path
func BizSetResourceMethod(path string, method *fc.Method, created bool) error {

	key := getMethodKey(path, string(method.HTTPVerb))
	data, _ := yaml.MarshalYML(method)
	if created {
		// 需要判断是创建还是修改
		err := Client.Create(key, string(data))
		if err != nil {
			logger.Warnf("BizSetResourceMethod etcd error, %v\n", err)
			return perrors.WithMessage(err, "BizSetResourceMethod error")
		}
	} else {
		// 需要判断是创建还是修改
		err := Client.Update(key, string(data))
		if err != nil {
			logger.Warnf("BizSetResourceMethod etcd error, %v\n", err)
			return perrors.WithMessage(err, "BizSetResourceMethod error")
		}
	}

	return nil
}

func getRootPath(key string) string {
	return Bootstrap.EtcdConfig.Path + "/" + key
}
