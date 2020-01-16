package util

import (
	jsoniter "github.com/json-iterator/go"
	"reflect"
	"strings"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func InterfaceTOInterface(in interface{}, out interface{}) error {
	b, e := json.Marshal(in)
	if e != nil {
		return e
	}
	e = json.Unmarshal(b, out)
	if e != nil {
		return e
	}
	return e
}
func StructToJsonString(v interface{}) (string, error) {
	jsonBody, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(jsonBody), err
}
func ParseJsonByStruct(body []byte, v interface{}) error {
	if err := json.Unmarshal(body, v); err != nil {
		return err
	}
	return nil
}
func CheckStringInArray(s string, ss []string) bool {
	for _, v := range ss {
		if s == v {
			return true
		}
	}
	return false
}

func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

func FixPattern(restPattern string) string {
	if len(restPattern) == 0 {
		return restPattern
	}
	if restPattern[:1] != "/" {
		restPattern = "/" + restPattern
	}
	return restPattern
}
func GetServiceName(interfaceName string) string {
	ss := strings.Split(interfaceName, ".")
	if len(ss) <= 3 {
		return ""
	}
	return ss[2]
}
