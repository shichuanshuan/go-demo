package main

import (
	"encoding/hex"
	"log"
	"net"
	"time"
)

const (
	address = "239.0.0.0:9999"
)

const (
	maxDatagramSize = 8192
)

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
func ping(addr string) {
	conn, err := NewBroadcaster(addr)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn.Write([]byte("hello, world\n"))
		time.Sleep(1 * time.Second)
	}
}
func main() {
	ping(address)

}

func msgHandler(src *net.UDPAddr, n int, b []byte) {
	log.Println(n, "bytes read from", src)
	log.Println(hex.Dump(b[:n]))
}
