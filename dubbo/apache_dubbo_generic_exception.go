package dubbo

import (
	hessian "github.com/apache/dubbo-go-hessian2"
	"github.com/apache/dubbo-go-hessian2/java_exception"
)

func init() {
	hessian.RegisterPOJO(&ApacheDubboGenericException{})
}

type ApacheDubboGenericException struct {
	SerialVersionUID     int64
	DetailMessage        string
	SuppressedExceptions []java_exception.Throwabler
	StackTrace           []java_exception.StackTraceElement
	Cause                java_exception.Throwabler
	ExceptionClass       string
	ExceptionMessage     string
}

func NewApacheDubboGenericException(exceptionClass, exceptionMessage string) *ApacheDubboGenericException {
	return &ApacheDubboGenericException{ExceptionClass: exceptionClass, ExceptionMessage: exceptionMessage}
}

func (e ApacheDubboGenericException) Error() string {
	return e.DetailMessage
}

func (ApacheDubboGenericException) JavaClassName() string {
	return "org.apache.dubbo.rpc.service.GenericException"
}
