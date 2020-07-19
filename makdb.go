package main

import (
	"fmt"
	"os"
	"bufio"
	"strings"
)

func main() {
	stdinScanner := bufio.NewScanner(os.Stdin)
	prompt := "maksql>"

	for {
		fmt.Print(prompt)
		success := stdinScanner.Scan()
		if !success {
			break
		}
		
		inputText := strings.TrimSpace(stdinScanner.Text())

		if inputText == ".exit" {
			os.Exit(0)
		} else {
			fmt.Printf("Unrecognised command: %v\n", inputText)
		}
	}

	if err := stdinScanner.Err(); err != nil {
		fmt.Printf("\nError: %v\n", err)
	}
}