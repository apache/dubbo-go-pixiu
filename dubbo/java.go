package dubbo

import (
	"errors"
)

import (
	"github.com/dubbogo/dubbo-go-proxy/common/logger"
)

// java class name elem
const (
	JavaStringClassName      = "java.lang.String"
	JavaLangClassName        = "java.lang.Long"
	JavalangIntegerClassName = "java.lang.Integer"
	JavaStringListClassName  = "[Ljava.lang.String;"
	JavaIntegerListClassName = "[Ljava.lang.Integer;"
	JavaMapClassName         = "java.util.Map"
)

func adapterForJava(ParameterTypes []string, inData []interface{}) ([]interface{}, error) {
	var (
		outData = make([]interface{}, len(ParameterTypes))
		err     error
	)
	defer func() {
		if r := recover(); r != nil {
			logger.Error(r)
			err = errors.New("adapterForJava get err")
			return
		}
	}()

	for i := range ParameterTypes {
		switch ParameterTypes[i] {
		case JavaStringClassName:
			outData[i] = inData[i].(string)
		case JavaLangClassName:
			outData[i] = inData[i].(int)
		default:
			outData[i] = inData[i]
		}
	}
	return outData, err
}
