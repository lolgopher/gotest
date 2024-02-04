package main

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

func checkError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func main() {
	serverAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:10001")
	checkError(err)

	localAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	checkError(err)

	conn, err := net.DialUDP("udp", localAddr, serverAddr)
	// conn, err := net.DialUDP("udp", nil, serverAddr)
	checkError(err)
	defer conn.Close()

	i := 0
	for {
		msg := strconv.Itoa(i)
		i++
		buf := []byte(msg)
		_, err := conn.Write(buf)
		if err != nil {
			fmt.Println(msg, err)
		}
		time.Sleep(time.Second * 1)
	}

	// var c chan any
	// go func() {
	// 	i := 0
	// 	for {
	// 		// msg := strconv.Itoa(i)
	// 		i++
	// 		// buf := []byte(msg)
	// 		buf := make([]byte, 172)
	// 		rand.Read(buf)
	// 		_, err := conn.Write(buf)
	// 		if err != nil {
	// 			fmt.Println(i, err)
	// 		}
	// 	}
	// 	c <- nil
	// }()

	// <-c
}
