package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	lsnr, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatalf("Error listening to network: %s", err)
	}
	defer lsnr.Close()

	for {
		conn, err := lsnr.Accept()
		if err != nil {
			log.Fatalf("Error awaiting connection: %s", err)
		}

		fmt.Println("Accepted connection from", conn.RemoteAddr())
		strChan := getLinesChannel(conn) // conn implements io...
		for str := range strChan {
			fmt.Println(str)
		}
		fmt.Println("Connection to", conn.RemoteAddr(), "closed")
	}
}
