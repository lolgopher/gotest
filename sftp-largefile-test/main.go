package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

const (
	ip       = ""
	port     = ""
	username = ""
	password = ""

	capacity    = 1024 * 1024
	maxCapacity = 3 * 1024 * 1024 * 1024
)

func main() {
	var filePath string
	flag.StringVar(&filePath, "path", "", "temp file path")
	flag.Parse()

	var tempFile *os.File
	var err error
	if filePath == "" {
		tempFile, err = makeTempFile(maxCapacity)

	} else {
		tempFile, err = os.Open(filePath)
	}
	if err != nil {
		log.Fatalf("fail to init temp file: %v", err)
	}
	defer func() {
		_ = tempFile.Close()

		// remove temp file
		if err = os.Remove(tempFile.Name()); err != nil {
			log.Fatalf("fail to remove temp file: %v", err)
		}
	}()

	// set ssh connection info
	sshConfig := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// make ssh client
	addr := fmt.Sprintf("%s:%s", ip, port)
	sshClient, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer func() {
		if err := sshClient.Close(); err != nil {
			log.Printf("fail to close ssh client: %v", err)
		}
	}()

	// make sftp client
	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		log.Fatalf("fail to create SFTP client: %v", err)
	}
	defer func() {
		if err := sftpClient.Close(); err != nil {
			log.Printf("fail to close sftp client: %v", err)
		}
	}()

	localFilePath := tempFile.Name()
	remoteFilePath := localFilePath

	// open local file
	localFile, err := os.Open(localFilePath)
	if err != nil {
		log.Fatalf("fail to open local file: %v", err)
	}
	defer func() {
		if err := localFile.Close(); err != nil {
			log.Printf("fail to close %s file: %v", localFilePath, err)
		}
	}()

	localFileContent, err := io.ReadAll(localFile)
	if err != nil {
		log.Fatalf("fail to read local file: %v", err)
	}

	// make dir
	dir := filepath.Dir(remoteFilePath)
	if _, err := sftpClient.Stat(dir); os.IsNotExist(err) {
		if err := sftpClient.MkdirAll(dir); err != nil {
			log.Fatalf("fail to create remote dir: %v", err)
		}
	}

	// send file
	newFile, err := sftpClient.OpenFile(remoteFilePath, os.O_CREATE|os.O_WRONLY|os.O_EXCL)
	if err != nil {
		log.Fatalf("fail to create remote file: %v", err)
	}
	defer func() {
		if err := newFile.Close(); err != nil {
			log.Printf("fail to close %s file: %v", remoteFilePath, err)
		}
	}()
	size, err := newFile.Write(localFileContent)
	if err != nil {
		log.Fatalf("fail to write to remote file: %v", err)
	}

	fileSize := getFileSize(tempFile)
	log.Printf("temp file size = %d, send file size = %d, is same = %v", fileSize, size, fileSize == size)
}

func makeTempFile(size int64) (tempFile *os.File, err error) {
	// make temp file
	tempFile, err = os.CreateTemp(".", "temp")
	if err != nil {
		return nil, errors.Wrap(err, "fail to make temp file")
	}

	// print temp file path
	log.Printf("temp file path: %s", tempFile.Name())

	// write to temp file
	start := time.Now()
	buf := make([]byte, capacity)

	for getFileSize(tempFile) < size {
		if _, err := rand.Read(buf); err != nil {
			return nil, errors.Wrap(err, "fail to make rand content")
		}

		if _, err = tempFile.Write(buf); err != nil {
			return nil, errors.Wrap(err, "fail to write temp file")
		}
	}
	end := time.Since(start)

	log.Printf("make temp file: %v", end)
	log.Printf("write to temp file %d size", getFileSize(tempFile))

	return tempFile, err
}

func getFileSize(file *os.File) int64 {
	info, err := os.Stat(file.Name())
	if err != nil {
		log.Fatalf("fail to get file %s info: %v", file.Name(), err)
	}

	return info.Size()
}
