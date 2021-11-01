package utils

import (
	"bufio"
	"crypto/rand"
	"io"
	"log"
	"os"
)

const (
	chunkSize = 10000
)

func GenerateLargeFile(size int) {
	primaryBuffer := make([]byte, chunkSize)
	remainderBuff := make([]byte, size%chunkSize)
	var buffer []byte

	buffers := map[string][]byte{
		"primary":   primaryBuffer,
		"remainder": remainderBuff,
	}

	outputFile, err := os.Create("output.txt")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := outputFile.Close(); err != nil {
			panic(err)
		}
	}()

	writer := bufio.NewWriter(outputFile)
	writtenSize := 0
	var n int
	for {
		if size <= 0 {
			break
		}

		if size < chunkSize {
			buffer = buffers["remainder"]
		} else {
			buffer = buffers["primary"]
		}

		n, err = rand.Read(buffer)
		if err != nil && err != io.EOF {
			panic(err)
		}

		if _, err := writer.Write(buffer[:n]); err != nil {
			panic(err)
		}

		size -= n
		writtenSize += n
		log.Println("=================================")
		log.Printf("SIZE OF THE BUFFER %v\n", len(buffer))
		log.Printf("WRITTEN CHUNK SIZE %v\n", n)
		log.Printf("REMAINING SIZE %v\n", size)
		log.Printf("WRITTEN SIZE %v\n", writtenSize)
		log.Println("=================================\n")
	}
}
