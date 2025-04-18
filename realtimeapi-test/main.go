package main

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

type SessionUpdate struct {
	EventId string      `json:"event_id"`
	Type    string      `json:"type"`
	Session SessionType `json:"session"`
}

type SessionType struct {
	Modalities        []string `json:"modalities,omitempty"`
	Instructions      string   `json:"instructions,omitempty"`
	InputAudioFormat  string   `json:"input_audio_format,omitempty"`
	OutputAudioFormat string   `json:"output_audio_format,omitempty"`
}

type InputAudioBuffer struct {
	EventId string `json:"event_id"`
	Type    string `json:"type"`
	Audio   string `json:"audio,omitempty"`
}

type ResponseAudioDelta struct {
	EventId      string `json:"event_id"`
	Type         string `json:"type"`
	ResponseId   string `json:"response_id"`
	ItemId       string `json:"item_id"`
	OutputIndex  int    `json:"output_index"`
	ContentIndex int    `json:"content_index"`
	Delta        string `json:"delta"`
}

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{
		Scheme:   "wss",
		Host:     "api.openai.com",
		Path:     "/v1/realtime",
		RawQuery: "model=gpt-4o-mini-realtime-preview-2024-12-17",
	}
	log.Printf("connecting to %s", u.String())

	h := http.Header{}
	h.Add("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY"))
	h.Add("OpenAI-Beta", "realtime=v1")

	c, _, err := websocket.DefaultDialer.Dial(u.String(), h)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		var pcmData []byte

		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}

			// JSON 형식의 메시지를 파싱
			var responseAudioDelta ResponseAudioDelta
			if err := json.Unmarshal(message, &responseAudioDelta); err != nil {
				fmt.Println("Error unmarshalling JSON:", err)
				return
			}
			if responseAudioDelta.Type == "response.audio.delta" {
				audioData, err := b64.StdEncoding.DecodeString(responseAudioDelta.Delta)
				if err != nil {
					fmt.Println("Error decoding base64:", err)
					return
				}
				pcmData = append(pcmData, audioData...)
			} else if responseAudioDelta.Type == "response.audio.done" {
				// PCM 데이터를 WAV 파일로 저장
				if err := saveAudioToWav("output.wav", pcmData, len(pcmData)/2, 8000, 1); err != nil {
					fmt.Println("Error saving audio to WAV:", err)
					return
				}
			} else {
				log.Printf("recv: %s", message)
			}
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	var count int
	var flag bool

	// 청크 전송 완료 메시지 전송
	jsonMessage, err := json.Marshal(SessionUpdate{
		EventId: fmt.Sprintf("event_%d", count),
		Type:    "session.update",
		Session: SessionType{
			Modalities:        []string{"audio", "text"},
			Instructions:      "한국어로 대답해줘.",
			InputAudioFormat:  "g711_alaw",
			OutputAudioFormat: "g711_alaw",
		},
	})
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}
	if err := c.WriteMessage(websocket.TextMessage, jsonMessage); err != nil {
		log.Println("write:", err)
		return
	}

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			log.Println("tick", t)

			if flag {
				continue
			} else {
				flag = true
			}

			filename := "test.wav"

			// ChunkReader 객체 생성
			chunkReader, err := NewChunkReader(filename)
			if err != nil {
				fmt.Println("Error opening file:", err)
				return
			}
			defer chunkReader.Close()

			// 청크 단위로 파일 읽기
			for {
				chunk, err := chunkReader.ReadChunk()
				if err != nil {
					fmt.Println("Error reading chunk:", err)
					break
				}

				if chunk == nil {
					// EOF에 도달한 경우
					fmt.Println("All chunks read successfully.")
					break
				}
				// fmt.Printf("Read chunk of size %d bytes\n", len(chunk))

				// 청크를 JSON으로 변환
				audioBufferMessage := InputAudioBuffer{
					EventId: fmt.Sprintf("event_%d", count),
					Type:    "input_audio_buffer.append",
					Audio:   b64.StdEncoding.EncodeToString(chunk),
				}
				jsonMessage, err := json.Marshal(audioBufferMessage)
				if err != nil {
					fmt.Println("Error marshalling JSON:", err)
					break
				}

				// 청크 전송
				if err := c.WriteMessage(websocket.TextMessage, jsonMessage); err != nil {
					log.Println("write:", err)
					return
				}

				// 20ms 대기
				time.Sleep(20 * time.Millisecond)
			}

			// 묵음 전송
			for i := 0; i < 50; i++ {
				chunk := make([]byte, 640)

				// 청크를 JSON으로 변환
				audioBufferMessage := InputAudioBuffer{
					EventId: fmt.Sprintf("event_%d", count),
					Type:    "input_audio_buffer.append",
					Audio:   b64.StdEncoding.EncodeToString(chunk),
				}
				jsonMessage, err := json.Marshal(audioBufferMessage)
				if err != nil {
					fmt.Println("Error marshalling JSON:", err)
					break
				}

				// 청크 전송
				if err := c.WriteMessage(websocket.TextMessage, jsonMessage); err != nil {
					log.Println("write:", err)
					return
				}

				// 20ms 대기
				time.Sleep(20 * time.Millisecond)
			}

			// 1초 대기
			// time.Sleep(1 * time.Second)

			// 청크 전송 완료 메시지 전송
			// jsonMessage, err := json.Marshal(InputAudioBuffer{
			// 	EventId: fmt.Sprintf("event_%d", count),
			// 	Type:    "input_audio_buffer.commit",
			// })
			// if err != nil {
			// 	fmt.Println("Error marshalling JSON:", err)
			// 	break
			// }
			// if err := c.WriteMessage(websocket.TextMessage, jsonMessage); err != nil {
			// 	log.Println("write:", err)
			// 	return
			// }

			// err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			// if err != nil {
			// 	log.Println("write:", err)
			// 	return
			// }

		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
				// case <-time.After(time.Second):
			}
			return
		}
	}
}
