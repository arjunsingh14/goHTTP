package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	f, err := os.Open("message.txt")
	if err != nil {
		log.Fatalf("error reading file: %s", err)
	}
	ch := getLinesChannel(f)
	for line := range ch {
		fmt.Printf("read: %s", line)
	}
	err = f.Close()
}

func getLinesChannel(f io.ReadCloser) <- chan string {
	currentLine := ""
	ch := make(chan string)
	go func() { 
		defer close(ch)
		for {
			buffer := make([]byte, 8)
			bytesRead, err := f.Read(buffer)
			if err != nil {
				if currentLine != "" {
					ch <- currentLine
					currentLine = ""
				}
				if errors.Is(err, io.EOF){
					break
				}
				log.Fatalf("error reading chunk: %s", err)
			}
			text := string(buffer[:bytesRead])
			parts := strings.Split(text, "\n")
			for i := 0; i < len(parts) - 1; i++ {
				ch <- fmt.Sprintf("%s%s\n", currentLine, parts[i])
				currentLine = ""
			}
			currentLine += parts[len(parts) - 1]
		}
	}()
	return ch
}

