package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

const (
	selectStatement = iota
	insertStatement
	invalidStatement
)

const (
	columnUsernameSize = 32
	columnEmailSize    = 255
)

const numRowsPerPage = 10
const maxTableRows = 100

type page struct {
	rows []row
}

func newPage() page {
	p := page{}
	p.rows = make([]row, 0, numRowsPerPage)
	return p
}

func getRow(t *table, rowNum int) *row {
	pageNum := rowNum / numRowsPerPage
	rowInPage := rowNum % numRowsPerPage
	return &(t.pages[pageNum].rows[rowInPage])
}

type table struct {
	numRows uint32
	pages   []page
}

type row struct {
	id       uint32
	username string
	email    string
}

type statementType int

type statement struct {
	sType       statementType
	rowToInsert row
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
		n, _ := fmt.Sscanf(
			inputCommand,
			"insert %d %s %s",
			&newStatement.rowToInsert.id,
			&newStatement.rowToInsert.username,
			&newStatement.rowToInsert.email)
		fmt.Printf("row to insert: %v\n", newStatement.rowToInsert)

		if n != 3 {
			newStatement.sType = invalidStatement
			errString := fmt.Sprintf("syntax error: %v\n", inputCommand)
			return newStatement, errors.New(errString)
		}
	} else {
		newStatement.sType = invalidStatement
		errString := fmt.Sprintf("unrecognised command: %v\n", inputCommand)
		return newStatement, errors.New(errString)
	}
	fmt.Printf("New statement: %v\n", newStatement)
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
