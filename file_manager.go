package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
)

var (
	serverNodes = []string{"Server1", "Server2", "Server3"}
	chunkSize   = getEnvAsInt("CHUNK_SIZE", 5*1024*1024)
)

func main() {
	filePath := "/path/to/your/large/file.txt"
	err := storeFile(filePath)
	if err != nil {
		log.Fatalf("Failed to store file: %v", err)
	}

	retrievedFilePath := "/path/to/your/retrieved_file.txt"
	err = retrieveFile(filePath, retrievedFilePath)
	if err != nil {
		log.Fatalf("Failed to retrieve file: %v", err)
	}
}

func storeFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	return splitAndStore(reader, filePath)
}

func splitAndStore(reader *bufio.Reader, filePath string) error {
	buf := make([]byte, chunkSize)
	chunkIndex := 0

	for {
		n, readErr := reader.Read(buf)
		if readErr != nil && readErr != io.EOF {
			return readErr
		}
		if n == 0 {
			break
		}

		if err := storeChunkToNode(buf[:n], filePath, chunkIndex); err != nil {
			return err
		}
		chunkIndex++
	}

	return nil
}

func storeChunkToNode(data []byte, filePath string, chunkIndex int) error {
	serverNode := serverNodes[chunkIndex%len(serverNodes)]
	chunkFilePath := fmt.Sprintf("%s_%d.chunk", filePath, chunkIndex)
	log.Printf("Storing chunk %d to %s\n", chunkIndex, serverNode)
	return os.WriteFile(chunkFilePath, data, 0644)
}

func retrieveFile(baseFilePath, destinationFilePath string) error {
	destFile, err := os.Create(destinationFilePath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	return reassembleFile(destFile, baseFilePath)
}

func reassembleFile(destFile *os.File, baseFilePath string) error {
	chunkIndex := 0
	for {
		chunkFilePath := fmt.Sprintf("%s_%d.chunk", baseFilePath, chunkIndex)
		chunk, err := os.ReadFile(chunkFilePath)
		if err == os.ErrNotExist {
			break
		} else if err != nil {
			return err
		}

		if _, writeErr := destFile.Write(chunk); writeErr != nil {
			return writeErr
		}
		chunkIndex++
	}

	return nil
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	valueInt, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Error reading %s as integer, using default %d", key, defaultValue)
		return defaultValue
	}

	return valueInt
}