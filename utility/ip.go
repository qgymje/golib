package utility

import (
	"encoding/binary"
	"net"
	"net/http"
)

var localIP string
var localIntIP uint32

func GetLocalIP() string {
	if localIP != "" {
		return localIP
	}
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				localIP = ipnet.IP.String()
				return localIP
			}
		}
	}
	return ""
}

func GetRemoteIP(req *http.Request) string {
	var ip string
	if ipProxy := req.Header.Get("X-FORWARDED-FOR"); len(ipProxy) > 0 {
		ip = ipProxy
	} else {
		ip, _, _ = net.SplitHostPort(req.RemoteAddr)
	}
	return ip
}

func IP2Int(ip net.IP) uint32 {
	if len(ip) == 16 {
		return binary.BigEndian.Uint32(ip[12:16])
	}
	return binary.BigEndian.Uint32(ip)
}

func Int2IP(n uint32) net.IP {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, n)
	return ip
}

func GetLocalIPToInt() uint32 {
	if localIntIP != 0 {
		return localIntIP
	}

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return 0
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return IP2Int(ipnet.IP)
			}
		}
	}
	return 0
}
