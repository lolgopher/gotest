package main

import (
	"flag"
	"github.com/lolgopher/synology-filesync/protocol"
	"log"
)

func main() {
	var targetPath string
	flag.StringVar(&targetPath, "path", "/", "path")
	flag.Parse()

	sc, err := protocol.NewSFTPClient(&protocol.ConnectionInfo{
		IP:       "192.168.0.200",
		Port:     "2222",
		Username: "u0_a192",
		Password: "admin",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := sc.Close(); err != nil {
			log.Fatalf("fail to close sftp client: %v", err)
		}
	}()

	// 남은 용량 확인
	stat1, err := sc.Client.Stat(targetPath)
	if err != nil {
		log.Fatalf("failed to get root directory information: %v", err)
	}

	availableSpace1 := stat1.Size()
	log.Printf("size: %v", availableSpace1)

	stat2, err := sc.Client.StatVFS(targetPath)
	if err != nil {
		log.Fatalf("failed to get root directory information: %v", err)
	}
	totalSpace := stat2.TotalSpace()
	availableSpace2 := stat2.FreeSpace()
	log.Printf("total size: %v", totalSpace/1024/1024)
	log.Printf("size: %v", availableSpace2/1024/1024)
}
