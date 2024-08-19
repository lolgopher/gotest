package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Message struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	done := make(chan interface{})
	startFlag := false
	for {
		select {
		case <-done:
			return
		default:
			conn.SetReadDeadline(time.Now().Add(1 * time.Second))
			_, _, err := conn.ReadMessage()
			if err != nil {
				if err, ok := err.(net.Error); !ok || !err.Timeout() {
					log.Println(err)
					return

				}
			}

			// log.Printf("Received message: %s\n", message)

			if !startFlag {
				startFlag = true
				go sendMessage(conn, done)
			}
		}
	}
}

func sendMessage(conn *websocket.Conn, done chan<- interface{}) {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	defer func() {
		done <- struct{}{}
		conn.Close()
	}()

	messages := []string{
		`{"type":"ready","sessionId":"32e42215e41bf38901b2cdad7f61702dd532b902"}`,
		`{"type":"beginPointDetection","value":"BPD"}`,
		`{"type":"partialResult","value":"아"}`,
		`{"type":"partialResult","value":"아아"}`,
		`{"type":"partialResult","value":"아아 마이크"}`,
		`{"type":"partialResult","value":"아아 마이크 테스트"}`,
		`{"type":"endPointDetection","value":"EPD"}`,
		`{"type":"finalResult","value":"아아 마이크 테스트","durationMS":3300,"x-metering-count":3,"nBest":[{"value":"아아 마이크 테스트","resultInfo":null,"score":2}],"voiceProfile":{"authenticated":false},"gender":0}`,
		`{"type":"partialResult","value":"아"}`,
		`{"type":"partialResult","value":"아아"}`,
		`{"type":"partialResult","value":"아아 마이크"}`,
		`{"type":"partialResult","value":"아아 마이크 테스트"}`,
		`{"type":"endPointDetection","value":"EPD"}`,
		`{"type":"finalResult","value":"아아 마이크 테스트","durationMS":3300,"x-metering-count":3,"nBest":[{"value":"아아 마이크 테스트","resultInfo":null,"score":2}],"voiceProfile":{"authenticated":false},"gender":0}`,
		`{"type":"partialResult","value":"아"}`,
		`{"type":"partialResult","value":"아아"}`,
		`{"type":"partialResult","value":"아아 마이크"}`,
		`{"type":"partialResult","value":"아아 마이크 테스트"}`,
		`{"type":"endPointDetection","value":"EPD"}`,
		`{"type":"finalResult","value":"아아 마이크 테스트","durationMS":3300,"x-metering-count":3,"nBest":[{"value":"아아 마이크 테스트","resultInfo":null,"score":2}],"voiceProfile":{"authenticated":false},"gender":0}`,
		`{"type":"endLongRecognition","value":"ELR"}`,
	}

	for _, message := range messages {
		err := conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Println(err)
		}

		// 300~600ms sleep
		sleepDuration := time.Duration(random.Intn(300)+300) * time.Millisecond
		time.Sleep(sleepDuration)
	}
}

func main() {
	http.HandleFunc("/ws", handleWebSocket)

	fmt.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
