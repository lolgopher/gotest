package main

import (
	"flag"
	"log"
	"math/rand"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const BINARY_SIZE = 320 * 1024

var addr = flag.String("addr", "tawny-vm.kep.k9d.in:8080", "http service address")
var i int
var errCount int

func main() {
	flag.Parse()
	log.SetFlags(0)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	go func() {
		for range ticker.C {
			log.Printf("Error/Client: %d/%d (%.2f%%)", errCount, i, float64(errCount)/float64(i)*100.0)
		}
	}()

	numClients := 10000 // 웹소켓 클라이언트의 개수
	wg := sync.WaitGroup{}
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i = 0; i < numClients; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			connect(id)
		}(i)

		// 300~600ms sleep
		sleepDuration := time.Duration(random.Intn(300)+300) * time.Millisecond
		time.Sleep(sleepDuration)
	}

	wg.Wait()
	log.Printf("Done! Error/Client: %d/%d (%.2f%%)", errCount, i, float64(errCount)/float64(i)*100.0)
	os.Exit(0)
}

func connect(id int) {
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/v1/stt"}
	log.Printf("Client %d connecting to %s", id, u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Printf("Client %d dial: %v", id, err)
		return
	}
	defer c.Close()

	message := `{"type":"recogStart","service":"DICTATION","audioFormat":"RAWPCM/16/16000/1/_/_","recogLongMaxWaitTime":3600000,"requestId":"ibc-stt-proxy-2024-06-03 15:17:03.537538 +0900 KST m=+752.254235626"}`
	if err := c.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		log.Printf("Client %d write: %v", id, err)
		errCount++
		return
	}

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Printf("Client %d read: %v", id, err)
				errCount++
				return
			}
			// log.Printf("Client %d recv: %s", id, message)

			if string(message) == `{"type":"stop"}` {
				return
			}
		}
	}()

	ticker := time.NewTicker(20 * time.Millisecond)
	defer ticker.Stop()
	randomBytes := make([]byte, BINARY_SIZE)

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			rand.Read(randomBytes)
			if err := c.WriteMessage(websocket.BinaryMessage, randomBytes); err != nil {
				log.Printf("Client %d write: %v", id, err)
				return
			}
		}
	}
}

// Client Port Range
