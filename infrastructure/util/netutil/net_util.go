package netutil

import (
	"errors"
	"fmt"
	"net"
)

func GetLocalIp() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", errors.New(fmt.Sprintf("get all address error: %s", err.Error()))
	}
	// 检查ip地址判断是否回环地址
	for _, address := range addrs {
		if netIp, ok := address.(*net.IPNet); ok && !netIp.IP.IsLoopback() {
			if netIp.IP.To4() != nil {
				return netIp.IP.String(), nil
			}
		}
	}
	return "", errors.New("can not get local ip")
}
