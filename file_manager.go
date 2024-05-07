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
	var batchData []byte // Simulated batching data container
	chunkIndex := 0

	for {
		n, readErr := reader.Read(buf)
		if readErr != nil && readErr != io.EOF {
			return readErr
		}
		if n == 0 {
			break
		}

		// Prepare and accumulate batch data. This is a simplified simulation.
		// In real applications, you would accumulate data in a structure best suited for your storage API.
		dataToStore := buf[:n]
		batchData = append(batchData, dataToStore...)

		if len(batchData) >= chunkSize*10 || (readErr == io.EOF && len(batchData) > 0) {
			if err := storeBatchToNode(batchData, filePath, chunkIndex); err != nil {
				return err
			}
			chunkIndex++
			batchData = []byte{} // Reset batch data after storing
		}
	}

	return nil
}

func storeBatchToNode(data []byte, filePath string, startIndex int) error {
	// In a real-world scenario, this method would make a bulk request to an API instead of writing to the disk.
	// The logic for splitting data back into chunks, computing checksums, and assigning to servers is kept for consistency with your initial design.
	chunks := splitIntoChunks(data, chunkSize)
	for i, chunk := range chunks {
		chunkIndex := startIndex + i
		serverNode := serverNodes[chunkIndex%len(serverNodes)]
		chunkFilePath := fmt.Sprintf("%s_%d.chunk", filePath, chunkIndex)

		// Calculate checksum for the chunk
		hasher := sha256.New()
		hasher.Write(chunk)
		checksum := hasher.Sum(nil)
		chunkDataWithChecksum := append(checksum, chunk...)

		log.Printf("Storing chunk %d (Checksum: %x) to %s\n", chunkIndex, checksum, serverNode)
		if err := os.WriteFile(chunkFilePath, chunkDataWithChecksum, 0644); err != nil {
			return err
		}
	}
	return nil
}

func splitIntoChunks(data []byte, size int) [][]byte {
	var chunks [][]byte
	for len(data) > 0 {
		if size > len(data) {
			size = len(data)
		}
		chunks = append(chunks, data[:size])
		data = data[size:]
	}
	return chunks
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
	var wg sync.WaitGroup
	chunkChan := make(chan []byte)

	go func() {
		for chunk := range chunkChan {
			if _, writeErr := destFile.Write(chunk); writeErr != nil {
				log.Printf("Error writing chunk to file: %v", writeErr)
				continue
			}
		}
	}()

	for {
		wg.Add(1)
		chunkFilePath := fmt.Sprintf("%s_%d.chunk", baseFilePath, chunkIndex)
		go func(chunkFilePath string, index int) {
			defer wg.Done()
			chunkDataWithChecksum, err := os.ReadFile(chunkFilePath)
			if err == os.ErrNotExist {
				close(chunkChan)
				return
			} else if err != nil {
				log.Printf("Error reading chunk: %v", err)
				return
			}

			if len(chunkDataWithChecksum) <= sha256.Size {
				log.Printf("Invalid chunk size, skipping: %s", chunkFilePath)
				return
			}

			checksum, data := chunkDataWithChecksum[:sha256.Size], chunkDataWithChecksum[sha256.Size:]
			calculatedChecksum := sha256.Sum256(data)
			if !bytes.Equal(checksum, calculatedChecksum[:]) {
				log.Printf("Checksum mismatch, skipping chunk: %s", chunkFilePath)
				return
			}

			chunkChan <- data
		}(chunkFilePath, chunkIndex)

		chunkIndex++
		if chunkIndex%len(serverNodes) == 0 { // Assume each file chunk retrieval starts after the previous one
			time.Sleep(1 * time.Second) // Simulate network/disk latency
		}
	}

	wg.Wait() // Wait for all go routines to finish
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