package healthcheck

import (
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"net"
	"strconv"
	"strings"
	"time"
)

func TcpConn(tarAddr string, port int, timeout time.Duration) bool {
	if _, _, err := net.SplitHostPort(tarAddr); err != nil {
		if strings.Contains(err.Error(), "missing port in address") {
			tarAddr = tarAddr + ":" + strconv.Itoa(port)
		} else {
			logger.Infof("[health check] invalid address format: %s", tarAddr)
			return false
		}
	}
	conn, err := net.DialTimeout("tcp", tarAddr, timeout)
	if err != nil {
		logger.Infof("[health check] http checker for host %s error: %v", tarAddr, err)
		return false
	}
	conn.Close()
	return true
}
