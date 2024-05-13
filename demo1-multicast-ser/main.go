package main

import (
	"encoding/hex"
	"fmt"
	"golang.org/x/net/ipv4"
	"log"
	"net"
)

const (
	address = "239.0.0.0:9999"
)

const (
	maxDatagramSize = 8192
)

// Listen binds to the UDP address and port given and writes packets received
// from that address to a buffer which is passed to a hander
func Listen(address string, handler func(*net.UDPAddr, int, []byte)) {
	// Parse the string address
	addr, err := net.ResolveUDPAddr("udp4", address)
	if err != nil {
		log.Fatal(err)
	}

	iface, _ := net.InterfaceByName("以太网")
	// Open up a connection
	conn, err := net.ListenMulticastUDP("udp", iface, addr)
	if err != nil {
		log.Fatal(err)
	}

	// 设置IGMP版本为V2
	p := ipv4.NewPacketConn(conn)
	defer p.Close()
	if err := p.SetControlMessage(ipv4.FlagSrc, true); err != nil {
		log.Fatal(err)
	}
	if err := p.SetControlMessage(ipv4.FlagDst, true); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("addr %v %v \n", addr.IP, addr.Port)

	conn.SetReadBuffer(maxDatagramSize)
	fmt.Printf("con %v %v \n", conn.LocalAddr(), conn.RemoteAddr())

	// Loop forever reading from the socket
	for {
		fmt.Printf("read\n")
		buffer := make([]byte, maxDatagramSize)
		numBytes, src, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Fatal("ReadFromUDP failed:", err)
		}

		handler(src, numBytes, buffer)
	}
}

func NewBroadcaster(address string) (*net.UDPConn, error) {
	addr, err := net.ResolveUDPAddr("udp4", address)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		return nil, err
	}

	return conn, nil

}

func main() {
	//go ping(address)
	Listen(address, msgHandler)
}

func msgHandler(src *net.UDPAddr, n int, b []byte) {
	log.Println(n, "bytes read from", src)
	log.Println(hex.Dump(b[:n]))
}
