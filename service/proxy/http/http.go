package http

import (
	"io"
	"net/http"
	"time"
)

import (
	"github.com/dubbogo/dubbo-go-proxy/common/config"
	"github.com/dubbogo/dubbo-go-proxy/common/errcode"
	"github.com/dubbogo/dubbo-go-proxy/common/logger"
	"github.com/dubbogo/dubbo-go-proxy/common/util"
	"github.com/dubbogo/dubbo-go-proxy/dubbo"
	ct "github.com/dubbogo/dubbo-go-proxy/service/context"
	"github.com/dubbogo/dubbo-go-proxy/service/metadata/redis"
)

var srv http.Server
var mc = redis.NewRedisMetaDataCenter()

func Run() {
	dubbo.Client.Init()
	startHttpServer()
}
func startHttpServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", commonHandle)
	srv = http.Server{
		Addr:           config.Config.HttpListenAddr,
		Handler:        mux,
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   20 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	err := srv.ListenAndServe()
	if err != nil {
		if err != http.ErrServerClosed {
			logger.Errorf("http.ListenAndServe err:%v", err)
		}
	}
}
func commonHandle(w http.ResponseWriter, r *http.Request) {
	setJsHeader(w, r)
	if r.Method == http.MethodOptions {
		return
	}
	var (
		ret         int
		response    interface{}
		err         error
		responseStr string
		ctx         = ct.NewHttpContext(w, r)
	)
	defer func() {
		// to http response
		if ret == errcode.Success {
			io.WriteString(w, responseStr)
		} else {
			io.WriteString(w, getErrRsp(ret))
		}
	}()
	filter := ctx.NextFilter()
	for filter != nil {
		ret = filter.OnRequest(ctx)
		if ret != errcode.Success {
			return
		}
		filter = ctx.NextFilter()
	}
	reqData := ctx.InvokeData()
	if reqData == nil {
		ret = errcode.ErrData
		return
	}
	if response, ret = dubbo.Client.Call(*reqData); ret != errcode.Success {
		return
	}
	responseStr, err = util.StructToJsonString(response)
	if err != nil {
		ret = errcode.ServerBusy
	}
	return

}

func getErrRsp(ret int) string {
	return errcode.GetMsg(ret)
}

func setJsHeader(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"))
	w.Header().Add("Access-Control-Allow-Headers", "access_token,Content-Type")

	w.Header().Add("Access-Control-Allow-Methods", "POST")
	w.Header().Add("Access-Control-Allow-Methods", "OPTIONS")
	w.Header().Add("Access-Control-Allow-Methods", "GET")
	w.Header().Add("Access-Control-Allow-Methods", "DELETE")
	w.Header().Add("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
}
