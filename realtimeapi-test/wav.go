package main

import (
	"encoding/binary"
	"os"
)

func writeWAVHeader(w *os.File, numSamples, sampleRate, numChannels int) error {
	// 1. "RIFF" chunk
	_, err := w.Write([]byte("RIFF"))
	if err != nil {
		return err
	}

	// 2. 파일 크기 (Header 제외한 데이터 크기)
	fileSize := 36 + numSamples // 36: 기본 헤더 크기, 1: 8-bit PCM 한 샘플 크기
	err = binary.Write(w, binary.LittleEndian, uint32(fileSize))
	if err != nil {
		return err
	}

	// 3. "WAVE" tag
	_, err = w.Write([]byte("WAVE"))
	if err != nil {
		return err
	}

	// 4. "fmt " chunk
	_, err = w.Write([]byte("fmt "))
	if err != nil {
		return err
	}

	// 5. fmt chunk size (16 for PCM)
	err = binary.Write(w, binary.LittleEndian, uint32(16))
	if err != nil {
		return err
	}

	// 6. Audio format (7: A-law PCM)
	err = binary.Write(w, binary.LittleEndian, uint16(7))
	if err != nil {
		return err
	}

	// 7. Number of channels
	err = binary.Write(w, binary.LittleEndian, uint16(numChannels))
	if err != nil {
		return err
	}

	// 8. Sample rate
	err = binary.Write(w, binary.LittleEndian, uint32(sampleRate))
	if err != nil {
		return err
	}

	// 9. Byte rate (sample rate * numChannels * bytes per sample)
	byteRate := sampleRate * numChannels * 1 // 8-bit PCM => 1 byte per sample
	err = binary.Write(w, binary.LittleEndian, uint32(byteRate))
	if err != nil {
		return err
	}

	// 10. Block align (numChannels * bytes per sample)
	blockAlign := uint16(numChannels * 1)
	err = binary.Write(w, binary.LittleEndian, blockAlign)
	if err != nil {
		return err
	}

	// 11. Bits per sample (8-bit PCM)
	err = binary.Write(w, binary.LittleEndian, uint16(8))
	if err != nil {
		return err
	}

	// 12. "data" chunk
	_, err = w.Write([]byte("data"))
	if err != nil {
		return err
	}

	// 13. Data size (number of samples * numChannels * bytes per sample)
	dataSize := numSamples * numChannels
	err = binary.Write(w, binary.LittleEndian, uint32(dataSize))
	return err
}

func saveAudioToWav(filePath string, data []byte, numSamples, sampleRate, numChannels int) error {
	// Create the file
	wavFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer wavFile.Close()

	// Write WAV header
	err = writeWAVHeader(wavFile, numSamples, sampleRate, numChannels)
	if err != nil {
		return err
	}

	// Write PCMA data (audio data)
	_, err = wavFile.Write(data)
	return err
}
