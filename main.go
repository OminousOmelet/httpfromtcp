package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	f, err := os.Open("messages.txt")
	if err != nil {
		fmt.Printf("Error opening file: %s", err)
		os.Exit(1)
	}
	strChan := getLinesChannel(f)
	for str := range strChan {
		fmt.Printf("read: %s\n", str)
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	strChan := make(chan string)

	go func() {
		defer f.Close()
		defer close(strChan)
		var msg string
		var msgLine string

		// adds 8 bytes at a time to a line, and sends line to channel on linebreaks
		for {
			buff := make([]byte, 8)
			_, err := f.Read(buff)
			if err != nil {
				if !errors.Is(err, io.EOF) {
					fmt.Printf("Error reading file: %s", err)
					return
				}
				break
			}
			msg += string(buff)
			if strings.Contains(msg, "\n") {
				parts := strings.Split(msg, "\n")
				msgLine += parts[0]
				strChan <- msgLine
				msgLine = parts[1]
				msg = ""
			}
		}
	}()
	return strChan
}
