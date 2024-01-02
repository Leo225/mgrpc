package mgrpc

import (
	"errors"
	"net"
)

func LocalIPv4(name ...string) (ipAddr string, err error) {
	var adds []net.Addr
	if len(name) > 0 && name[0] != "" {
		var inet *net.Interface
		inet, err = net.InterfaceByName(name[0])
		if err != nil {
			return
		}
		adds, err = inet.Addrs()
	} else {
		adds, err = net.InterfaceAddrs()
	}

	for _, a := range adds {
		if ipNet, ok := a.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ipAddr = ipNet.IP.String()
				return
			}
		}
	}
	err = errors.New("not found local ip")
	return
}
