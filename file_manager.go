package main

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"sync"
	"time"
)

var (
	storageNodes = []string{"Server1", "Server2", "Server3"} 
	segmentSize  = getEnvironmentVariableAsInt("CHUNK_SIZE", 5*1024*1024) 
)

func main() {
	sourceFilePath := "/path/to/your/large/file.txt" 
	if err := distributeFileAcrossNodes(sourceFilePath); err != nil {
		log.Fatalf("Failed to distribute file across nodes: %v", err)
	}

	targetFilePath := "/path/to/your/retrieved_file.txt" 
	if err := reconstructFileFromNodes(sourceFilePath, targetFilePath); err != nil {
		log.Fatalf("Failed to reconstruct file from nodes: %v", err)
	}
}

func distributeFileAcrossNodes(sourceFilePath string) error {
	sourceFile, err := os.Open(sourceFilePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	bufferedReader := bufio.NewReader(sourceFile)
	return splitAndDistribute(bufferedReader, sourceFilePath)
}

func splitAndDistribute(bufferedReader *bufio.Reader, sourceFilePath string) error {
	segmentBuffer := make([]byte, segmentSize)
	var batchBuffer []byte 
	segmentIndex := 0

	for {
		bytesRead, readErr := bufferedReader.Read(segmentBuffer)
		if readErr != nil && readErr != io.EOF {
			return readErr
		}
		if bytesRead == 0 {
			break
		}

		segmentData := segmentBuffer[:bytesRead]
		batchBuffer = append(batchBuffer, segmentData...)

		if len(batchBuffer) >= segmentSize*10 || (readErr == io.EOF && len(batchBuffer) > 0) {
			if err := distributeBatchToStorageNode(batchBuffer, sourceFilePath, segmentIndex); err != nil {
				return err
			}
			segmentIndex++
			batchBuffer = []byte{} 
		}
	}

	return nil
}

func distributeBatchToStorageNode(batchBuffer []byte, sourceFilePath string, startIndex int) error {
	segments := splitIntoSegments(batchBuffer, segmentSize)
	for i, segment := range segments {
		segmentIndex := startIndex + i
		node := storageNodes[segmentIndex%len(storageNodes)]
		segmentFilePath := fmt.Sprintf("%s_%d.chunk", sourceFilePath, segmentIndex)

		hasher := sha256.New()
		hasher.Write(segment)
		checksum := hasher.Sum(nil)
		segmentDataWithChecksum := append(checksum, segment...)

		log.Printf("Distributing segment %d (Checksum: %x) to %s\n", segmentIndex, checksum, node)
		if err := os.WriteFile(segmentFilePath, segmentDataWithChecksum, 0644); err != nil {
			return err
		}
	}
	return nil
}

func splitIntoSegments(data []byte, size int) [][]byte {
	var segments [][]byte
	for len(data) > 0 {
		if size > len(data) {
			size = len(data)
		}
		segments = append(segments, data[:size])
		data = data[size:]
	}
	return segments
}

func reconstructFileFromNodes(sourceFilePath, targetFilePath string) error {
	targetFile, err := os.Create(targetFilePath)
	if err != nil {
		return err
	}
	defer targetFile.Close()

	return assembleFileFromSegments(targetFile, sourceFilePath)
}

func assembleFileFromSegments(targetFile *os.File, sourceFilePath string) error {
	segmentIndex := 0
	var assembleWaitGroup sync.WaitGroup
	segmentChannel := make(chan []byte)

	go func() {
		for segment := range segmentChannel {
			if _, err := targetFile.Write(segment); err != nil {
				log.Printf("Error writing segment to file: %v", err)
				continue
			}
		}
	}()

	for {
		assembleWaitGroup.Add(1)
		segmentFilePath := fmt.Sprintf("%s_%d.chunk", sourceFilePath, segmentIndex)
		go func(filePath string) {
			defer assembleWaitGroup.Done()
			segmentData, err := os.ReadFile(filePath)
			if err == os.ErrNotExist {
				close(segmentChannel)
				return
			} else if err != nil {
				log.Printf("Error reading segment: %v", err)
				return
			}

			if len(segmentData) <= sha256.Size {
				log.Printf("Invalid segment size, skipping: %s", filePath)
				return
			}

			storedChecksum, data := segmentData[:sha256.Size], segmentData[sha256.Size:]
			computedChecksum := sha256.Sum256(data)
			if !bytes.Equal(storedChecksum, computedChecksum[:]) {
				log.Printf("Checksum mismatch, skipping segment: %s", filePath)
				return
			}

			segmentChannel <- data
		}(segmentFilePath)

		segmentIndex++
		if segmentIndex%len(storageNodes) == 0 {
			time.Sleep(1 * time.Second) 
		}
	}

	assembleWaitGroup.Wait() 
	return nil
}

func getEnvironmentVariableAsInt(key string, defaultValue int) int {
	valueStr, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	valueInt, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Error parsing %s as integer, using default %d", key, defaultValue)
		return defaultValue
	}

	return valueInt
}