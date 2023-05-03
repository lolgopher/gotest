package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/lolgopher/synology-filesync/protocol"
	"github.com/lolgopher/synology-filesync/tool"
)

var rootPath string

func main() {
	rootPath, _ = os.Getwd()

	var (
		targetPath      string
		failed          bool
		failedToNotSent bool
		removeZeroSize  bool
	)
	flag.StringVar(&targetPath, "path", rootPath, "Target Path to search metadata files")
	flag.BoolVar(&failed, "failed", false, "Search failed status")
	flag.BoolVar(&failedToNotSent, "retry", false, "Change failed to notsent status")
	flag.BoolVar(&removeZeroSize, "zero", false, "Remove zero size and notsent status")

	flag.Parse()
	targetPath, err := filepath.Abs(filepath.Join(".", targetPath))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("path: %v", targetPath)
	log.Printf("retry: %v", failedToNotSent)
	log.Printf("zero: %v", removeZeroSize)

	if failed {
		if result, err := searchFailedStatus(targetPath, failedToNotSent); err != nil {
			log.Fatal(err)
		} else {
			log.Print(result)
		}
	}

	if removeZeroSize {

	}
}

func searchFailedStatus(targetPath string, failedToNotSent bool) (string, error) {
	result := "\n"

	data, err := tool.GetFailedStatus(targetPath)
	if err != nil {
		return "", err
	}

	rootPath, err = filepath.Abs(filepath.Join(targetPath, ".."))
	if err != nil {
		return "", err
	}

	for key := range data {
		result += key + "\n"

		if failedToNotSent {
			if err := protocol.WriteMetadata(filepath.Join(rootPath, key), 0, protocol.NotSent); err != nil {
				return "", err
			}
		}
	}

	return result, nil
}
