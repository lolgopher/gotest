package main

import (
	"flag"
	"fmt"
	"github.com/lolgopher/synology-filesync/protocol"
	"github.com/lolgopher/synology-filesync/tool"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"path/filepath"
)

var rootPath string

func main() {
	rootPath, _ = os.Getwd()

	var (
		targetPath      string
		failed          bool
		failedToNotSent bool
		notSent         bool
		removeZeroSize  bool
	)
	flag.StringVar(&targetPath, "path", rootPath, "Target Path to search metadata files")
	flag.BoolVar(&failed, "failed", false, "Search failed status")
	flag.BoolVar(&failedToNotSent, "retry", false, "Change failed to notsent status")
	flag.BoolVar(&notSent, "notsent", false, "Search notsent status")
	flag.BoolVar(&removeZeroSize, "zero", false, "Remove zero size to notsent status")

	flag.Parse()
	targetPath, err := filepath.Abs(targetPath)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("path: %v", targetPath)
	log.Printf("retry: %v", failedToNotSent)
	log.Printf("notsent: %v", notSent)
	log.Printf("zero: %v", removeZeroSize)

	if failed {
		if result, err := searchFailedStatus(targetPath, failedToNotSent); err != nil {
			log.Fatal(err)
		} else {
			log.Print(result)
		}
	}

	if notSent {
		if result, err := searchNotSentStatus(targetPath, removeZeroSize); err != nil {
			log.Fatal(err)
		} else {
			log.Print(result)
		}
	}
}

func searchFailedStatus(targetPath string, failedToNotSent bool) (string, error) {
	result := "\n"

	data, err := tool.GetFailedStatus(targetPath)
	if err != nil {
		return "", err
	}

	targetRootPath, err := filepath.Abs(filepath.Join(targetPath, ".."))
	if err != nil {
		return "", err
	}

	for key := range data {
		result += key + "\n"

		if failedToNotSent {
			// 폴더 경로와 메타데이터 파일 경로 설정
			filePath := filepath.Join(targetRootPath, key)
			folderPath := filepath.Dir(filePath)
			metadataFilePath := filepath.Join(folderPath, "metadata.yaml")

			// 메타데이터 파일 읽기
			data, err := os.ReadFile(metadataFilePath)
			if err != nil {
				return "", err
			}
			metadata := make(map[string]protocol.FileMetadata)
			if err := yaml.Unmarshal(data, &metadata); err != nil {
				return "", fmt.Errorf("fail to unmarshal %s metadata file: %v", metadataFilePath, err)
			}

			// 방어 코드
			if _, ok := metadata[key]; !ok {
				return "", fmt.Errorf("fail to find %s in %s metadata file", key, metadataFilePath)
			}

			// 메타데이터 업데이트
			metadata[key] = protocol.FileMetadata{
				Size:   metadata[key].Size,
				Status: string(protocol.NotSent),
			}

			// 메타데이터 파일 쓰기
			metadataData, err := yaml.Marshal(metadata)
			if err != nil {
				return "", fmt.Errorf("fail to marshal %s : %s metadata file: %v", filePath, protocol.NotSent, err)
			}
			if err := os.WriteFile(metadataFilePath, metadataData, 0644); err != nil {
				return "", fmt.Errorf("fail to write %s file: %v", metadataFilePath, err)
			}
		}
	}

	return result, nil
}

func searchNotSentStatus(targetPath string, removeZeroSize bool) (string, error) {
	result := "\n"

	data, err := tool.GetNotSentStatus(targetPath)
	if err != nil {
		return "", err
	}

	targetRootPath, err := filepath.Abs(filepath.Join(targetPath, ".."))
	if err != nil {
		return "", err
	}

	for key := range data {
		result += key + "\n"

		if removeZeroSize {
			// 메타데이터 파일 읽기
			metadataPath := filepath.Dir(filepath.Join(targetRootPath, key))
			//if strings.Count(metadataPath, "/volume2/docker") > 1 {
			//	metadataPath = strings.Replace(metadataPath, "/volume2/docker", "", 1)
			//}
			meta, err := protocol.ReadMetadata(metadataPath)
			if err != nil {
				return "", err
			}

			if _, ok := meta[key]; !ok {
				return "", fmt.Errorf("fail to find %s in %s metadata file", key, metadataPath)
			}
			if meta[key].Size == 0 {
				// 메타데이터 삭제
				delete(meta, key)

				// 메타데이터 파일 쓰기
				metadataData, err := yaml.Marshal(meta)
				if err != nil {
					return "", fmt.Errorf("fail to marshal %v data: %v", meta, err)
				}
				metadataFilePath := filepath.Join(metadataPath, "metadata.yaml")
				if err := os.WriteFile(metadataFilePath, metadataData, 0644); err != nil {
					return "", fmt.Errorf("fail to write %s file: %v", metadataFilePath, err)
				}
			}
		}
	}

	return result, nil
}
