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
	defer f.Close()

	var msg string
	var msgLines []string
	step := 0
	msgLines = append(msgLines, "")

	// Prints each line as soon as they are ready
	for {
		buff := make([]byte, 8)
		_, err := f.Read(buff)
		if err != nil {
			if !errors.Is(err, io.EOF) {
				fmt.Printf("Error reading file: %s", err)
			}
			break
		}
		msg += string(buff)

		if strings.Contains(msg, "\n") {
			part := strings.Split(msg, "\n")
			msgLines[step] += part[0]
			fmt.Printf("read: %s\n", msgLines[step])

			msgLines = append(msgLines, part[1])
			step += 1
			msg = ""
		}
	}
}
