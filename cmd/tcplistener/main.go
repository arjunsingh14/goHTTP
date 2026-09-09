package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"net"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		fmt.Println("error listening on port 42069", err.Error())
		os.Exit(1)
	}
	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println("error connecting")
			continue
		}
		
		fmt.Println("connection accepted")
		ch := getLinesChannel(connection)
		for line := range ch {
			fmt.Printf("%s", line)
		}
		err = connection.Close()
		if err != nil {
			fmt.Println("failed to close listener")
		}
	}
}

func getLinesChannel(f io.ReadCloser) <- chan string {
	currentLine := ""
	ch := make(chan string)
	go func() { 
		defer close(ch)
		for {
			buffer := make([]byte, 8)
			bytesRead, err := f.Read(buffer)

			if bytesRead > 0 {
				text := string(buffer[:bytesRead])
				parts := strings.Split(text, "\n")
				for i := 0; i < len(parts) - 1; i++ {
					ch <- fmt.Sprintf("%s%s\n", currentLine, parts[i])
					currentLine = ""
				}
				currentLine += parts[len(parts) - 1]
			}

			if err != nil {
				if currentLine != "" {
					ch <- currentLine
					currentLine = ""
				}
				if errors.Is(err, io.EOF){
					break
				}
				fmt.Println("error reading chunk:", err)
				return
			}
		}
	}()
	return ch
}

