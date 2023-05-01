package main

import (
	"flag"
	"github.com/lolgopher/synology-filesync/protocol"
	"github.com/lolgopher/synology-filesync/tool"
	"log"
	"os"
	"path/filepath"
)

func main() {
	rootPath, _ := os.Getwd()

	var targetPath string
	var failedToNotSent bool
	flag.StringVar(&targetPath, "path", rootPath, "Target Path to search metadata files")
	flag.BoolVar(&failedToNotSent, "retry", false, "Change failed to notsent status")

	flag.Parse()
	targetPath, err := filepath.Abs(filepath.Join(".", targetPath))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("path: %v", targetPath)
	log.Printf("retry: %v", failedToNotSent)

	data, err := tool.GetFailedStatus(targetPath)
	if err != nil {
		log.Fatal(err)
	}

	rootPath, err = filepath.Abs(filepath.Join(targetPath, ".."))
	if err != nil {
		log.Fatal(err)
	}

	for key := range data {
		log.Println(key)

		if failedToNotSent {
			if err := protocol.WriteMetadata(filepath.Join(rootPath, key), 0, protocol.NotSent); err != nil {
				log.Fatal(err)
			}
		}
	}
}
