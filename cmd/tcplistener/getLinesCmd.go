package main

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	strChan := make(chan string)

	go func() {
		defer f.Close()
		defer close(strChan)
		var msgLine string

		// reads 8 bytes at a time, adds to msgLine, sends line to channel on linebreaks
		for {
			buff := make([]byte, 8)
			n, err := f.Read(buff)
			if err != nil { //need to check for nil first so no panic
				if msgLine != "" {
					strChan <- msgLine
				}
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Printf("Error reading file: %s", err)
				return
			}
			msg := string(buff[:n])
			parts := strings.Split(msg, "\n")
			for i := 0; i < len(parts)-1; i++ {
				strChan <- fmt.Sprintf("%s%s", msgLine, parts[i])
				msgLine = ""
			}
			msgLine += parts[len(parts)-1]
		}
	}()
	return strChan
}
