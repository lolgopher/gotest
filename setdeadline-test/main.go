package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	// UDP 서버 생성
	serverAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println(err)
		return
	}
	serverConn, err := net.ListenUDP("udp", serverAddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer serverConn.Close()

	// UDP 클라이언트 생성
	clientAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	if err != nil {
		fmt.Println(err)
		return
	}
	clientConn, err := net.DialUDP("udp", clientAddr, serverAddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer clientConn.Close()

	timeout := time.Now().Add(time.Second)
	// SetDeadline 테스트
	fmt.Println("SetDeadline test started.")
	clientConn.SetDeadline(timeout)
	for i := 0; i < 5; i++ {
		_, err = clientConn.Write([]byte("aaaaa"))
		if err != nil {
			fmt.Println("Write error:", err)
		}
		time.Sleep(time.Millisecond * 500)
	}

	fmt.Println("SetDeadline test finished.")
}
