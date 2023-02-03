package network

import (
	"log"
	"net"
)

// Get preferred outbound ip of this machine
// https://stackoverflow.com/a/37382208/12221657
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
