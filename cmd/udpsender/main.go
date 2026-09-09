package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		fmt.Println("failed to resolve udp addr: ", err)
		return
	}
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Println("failed to dial udp connection", err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("> ")
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("error getting stdin: ", err)
			continue
		}
		_, err = conn.Write([]byte(msg))
		if err != nil {
			fmt.Println("error writing message: ", err)
			continue
		}
	}

}