package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
)

const (
	BUFFER_SIZE   = 1000
	NAME_SIZE     = 64
	FILENAME_SIZE = 256
)

type st_PACKET_HEADER struct {
	DwPacketCode uint32 // 0x11223344 pre-define value
	SzName       [NAME_SIZE]byte
	SzFileName   [FILENAME_SIZE]byte
	IFileSize    int
}

var wg sync.WaitGroup

func main() {
	// target address
	address := "127.0.0.1:10010"

	// target file
	filePath := "./golang.jpg"

	wg = sync.WaitGroup{}
	wg.Add(1)
	go startServer(address)

	// connect to the server
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("Error connecting to the server:", err)
		return
	}
	defer conn.Close()

	// get file info
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		fmt.Println("Fail to get file info:", err)
		return
	}

	// get file size
	fileSize := fileInfo.Size()

	// make message
	var name [NAME_SIZE]byte
	copy(name[:], "홍길동")

	var fileName [FILENAME_SIZE]byte
	copy(fileName[:], "golang.jpg")

	message := st_PACKET_HEADER{
		DwPacketCode: 0x11223344, // 58851161
		SzName:       name,
		SzFileName:   fileName,
		IFileSize:    int(fileSize),
	}
	buf := structToBytes(message)

	// send a message to the server
	_, err = conn.Write(buf)
	if err != nil {
		fmt.Println("Error sending data:", err)
		return
	}
	fmt.Println("Message sent:", message)

	// file open
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Fail to open file:", err)
		return
	}
	defer file.Close()

	// read file data
	buffer := make([]byte, BUFFER_SIZE)
	for {
		// read file data by BUFFER_SIZE byte
		n, err := file.Read(buffer)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Fail to read file:", err)
			} else {
				fmt.Println("EOF: Client")
			}
			break
		}

		// send a message to the server
		_, err = conn.Write(buffer[:n])
		if err != nil {
			fmt.Println("Error sending data:", err)
			return
		}
	}

	// receive response from the server
	response := make([]byte, 4)
	n, err := conn.Read(response)
	if err != nil {
		fmt.Println("Error reading data:", err)
		return
	}
	fmt.Println("Response from the server:", response[:n])

	wg.Wait()
}

func structToBytes(packet st_PACKET_HEADER) []byte {
	data := make([]byte, 4+NAME_SIZE+FILENAME_SIZE+4)

	binary.LittleEndian.PutUint32(data[:4], packet.DwPacketCode)
	copy(data[4:], packet.SzName[:])
	copy(data[NAME_SIZE+4:], packet.SzFileName[:])
	binary.LittleEndian.PutUint32(data[FILENAME_SIZE+NAME_SIZE+4:], uint32(packet.IFileSize))

	return data
}

func startServer(address string) {
	defer wg.Done()

	// listen for incoming connections
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	// accept incoming connection
	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("Error accepting connection:", err)
		return
	}

	handleConnection(conn)
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// read data from the client
	buffer := make([]byte, 4+NAME_SIZE+FILENAME_SIZE+4)
	_, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading data from client:", err)
		return
	}

	header := bytesToStruct(buffer)
	fmt.Println("Received message from client:", header)

	var result []byte
	buf := make([]byte, 10)
	size := 0

	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading data from client:", err)
			return
		}

		if n == 0 {
			fmt.Println("EOF: Server")
			break
		}

		result = append(result, buf[:n]...)

		size += n
		if size >= header.IFileSize {
			fmt.Println("Done: Server")
			break
		}
	}
	if err := os.WriteFile("./output.jpg", result, 0644); err != nil {
		fmt.Println("Fail to save file:", err)
	}

	// send a response back to the client
	response := make([]byte, 4)
	binary.LittleEndian.PutUint32(response, 0xdddddddd)
	_, err = conn.Write(response)
	if err != nil {
		fmt.Println("Error sending data to client:", err)
		return
	}
	fmt.Println("Response sent to client:", response)
}

func bytesToStruct(data []byte) st_PACKET_HEADER {
	packet := st_PACKET_HEADER{}

	packet.DwPacketCode = binary.LittleEndian.Uint32(data[:4])
	copy(packet.SzName[:], data[4:NAME_SIZE+4])
	copy(packet.SzFileName[:], data[NAME_SIZE+4:FILENAME_SIZE+NAME_SIZE+4])
	packet.IFileSize = int(binary.LittleEndian.Uint32(data[FILENAME_SIZE+NAME_SIZE+4:]))

	return packet
}
