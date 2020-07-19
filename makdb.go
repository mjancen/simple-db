package main

import (
	"fmt"
	"os"
	"bufio"
	"strings"
	"errors"
)

const (
	selectStatement = iota
	insertStatement
	invalidStatement
)

type statementType int

type statement struct {
	sType statementType
}

func doMetaCommand(inputCommand string) error {
	if inputCommand == ".exit" {
		os.Exit(0)
		return nil
	} else {
		return nil
	}
}

func prepareStatement(inputCommand string) (statement, error) {
	var newStatement statement
	if strings.HasPrefix(inputCommand, "select") {
		newStatement.sType = selectStatement
	} else if strings.HasPrefix(inputCommand, "insert") {
		newStatement.sType = insertStatement
	} else {
		errString := fmt.Sprintf("unrecognised command: %v\n", inputCommand)
		return statement{sType : invalidStatement}, errors.New(errString)
	}
	return newStatement, nil
}

func executeStatement(st statement) {
	switch st.sType {
	case selectStatement:
		fmt.Println("This is a select statement.")
	case insertStatement:
		fmt.Println("This is an insert statement.")
	}
}

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
		inputText = strings.ToLower(inputText)

		if inputText[0] == '.' {
			err := doMetaCommand(inputText)
			if err != nil {
				fmt.Printf("Unrecognised meta-command: %v\n", inputText)
			}
		}
		
		inputStatement, err := prepareStatement(inputText)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			executeStatement(inputStatement)
		}

	}

	if err := stdinScanner.Err(); err != nil {
		fmt.Printf("\nError: %v\n", err)
	}
}