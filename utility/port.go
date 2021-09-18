package utility

import (
	"net"
	"strconv"
)

func FindPort(start, end int) int {
	if end-start < 10 {
		panic("port range can not less then 10")
	}
	port := RandInt(start, end)
	portStr := strconv.Itoa(port)

	if err := testport(portStr); err == nil {
		return port
	} else {
		return FindPort(start, end)
	}
}

func testport(port string) error {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}
	defer ln.Close()
	return nil
}
