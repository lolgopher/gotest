package main

import (
	"fmt"
	"net"
	"strings"
	"time"
)

func sendUDPMessage(leftIP, rightIP string, leftPort, rightPort int, message string) error {
	// Resolve UDP address
	rtpLeftAddr := &net.UDPAddr{
		IP:   net.ParseIP(leftIP),
		Port: leftPort,
	}
	rtpRightAddr := &net.UDPAddr{
		IP:   net.ParseIP(rightIP),
		Port: rightPort,
	}

	// Create UDP connection
	conn, err := net.DialUDP("udp", rtpLeftAddr, rtpRightAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	ticker := time.NewTicker(2 * time.Second)
	for {
		select {
		case <-ticker.C:
			break
		default:
			// Write message to connection
			_, err := conn.Write([]byte(message))
			if err != nil {
				if !strings.Contains(err.Error(), "connection refused") {
					fmt.Println("Error:", err)
				}
				return err
			}
		}
	}

	return nil
}

func main() {
	// 고루틴에서 stacktrace를 계속 저장, err가 리턴되어서 channel에 전달될 때 까지 저장
	// stacktrace의 내용(stack)이 가장 많은 내용만 저장(가장 하위 stack)
	// pprof를 찾아보자
	// https://cs.opensource.google/go/go/+/refs/tags/go1.22.1:src/net/http/pprof/pprof.go

	// Repeat as many times as you want
	for i := 0; i < 30000; i++ {
		i := i
		go func(int) {
			time.Sleep(time.Duration(100*i) * time.Millisecond)
			port := 30000 + (i % 400)
			_ = sendUDPMessage("127.0.0.1", "127.0.0.1", port, 8080, fmt.Sprintf("H%d\n", i))
			//if err != nil {
			//	fmt.Println("Error:", err)
			//}
		}(i)
	}

	// Wait for goroutines to finish
	fmt.Scanln()
}
