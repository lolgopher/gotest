package main

import (
	"io"
	"os"
)

const chunkSize = 320 // 320 Byte

// ChunkReader 구조체: MP3 파일을 청크 단위로 읽는 객체
type ChunkReader struct {
	file    *os.File
	current int64
}

// NewChunkReader: 새로운 ChunkReader 객체 생성
func NewChunkReader(filename string) (*ChunkReader, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	return &ChunkReader{
		file:    file,
		current: 0,
	}, nil
}

// ReadChunk: MP3 파일에서 청크로 읽기
func (cr *ChunkReader) ReadChunk() ([]byte, error) {
	buf := make([]byte, chunkSize)
	n, err := cr.file.Read(buf)
	if err == io.EOF {
		if n == 0 {
			return nil, nil // 파일 끝에 도달
		}
		return buf[:n], nil // 마지막 청크
	} else if err != nil {
		return nil, err // 읽기 오류
	}

	cr.current += int64(n)
	return buf[:n], nil
}

// Close: 파일을 닫는 메서드
func (cr *ChunkReader) Close() error {
	return cr.file.Close()
}
