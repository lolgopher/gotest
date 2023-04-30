package main

import (
	"flag"
	"github.com/lolgopher/synology-filesync/tool"
	"log"
	"os"
)

func main() {
	rootPath, _ := os.Getwd()

	var targetPath string
	flag.StringVar(&targetPath, "targetPath", rootPath, "Target Path to search metadata files")

	data, err := tool.GetFailedStatus(targetPath)
	if err != nil {
		log.Fatal(err)
	}

	for key := range data {
		log.Println(key)
	}
}
