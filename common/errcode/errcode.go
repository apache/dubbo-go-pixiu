package errcode

import "fmt"

var (
	Success    = 0
	NotFind    = 404
	ServerBusy = 500
	ErrData    = 501
)
var errorMap = map[int]string{
	Success:    "success",
	ServerBusy: "server busy",
	NotFind:    "page not find",
	ErrData:    "request parameter error",
}

func GetMsg(ret int) string {
	return fmt.Sprintf(`{"code":%v,"msg":"%v"}`, ret, errorMap[ret])
}
