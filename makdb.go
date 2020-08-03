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

type row struct {
	id       uint32
	username string
	email    string
}

type page struct {
	rows []row
}

func (p *page) appendToPage(r *row) {
	p.rows = append(p.rows, *r)
}

func newPage() page {
	p := page{}
	p.rows = make([]row, 0, numRowsPerPage)
	return p
}

type table struct {
	numRows      int
	numFullPages int
	pages        []page
}

func newTable() *table {
	pTable := new(table)
	pTable.numRows      = 0
	pTable.numFullPages = 0
	pTable.pages = make([]page, maxTableRows)

	for i := 0; i < maxTableRows; i++ {
		pTable.pages[i] = newPage()
	}

	return pTable
}

func (t *table) appendRow(r *row) error {
	if t.numRows >= maxTableRows {
		return errors.New("table is full")
	}

	currentPage := &t.pages[t.numFullPages]
	currentPage.appendToPage(r)
	t.numRows++
	if len(currentPage.rows) >= numRowsPerPage {
		t.numFullPages++
	}

	return nil
}

func (t *table) getRow(rowNum int) *row {
	pageInd, rowInd := rowNumToIndex(rowNum)
	return &(t.pages[pageInd].rows[rowInd])
}

func rowNumToIndex(rowNum int) (int, int) {
	rowIndex := rowNum % numRowsPerPage
	pageIndex := rowNum / numRowsPerPage
	return pageIndex, rowIndex
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
	return newStatement, nil
}

func executeStatement(st *statement, t *table) error {
	switch st.sType {
	case selectStatement:
		err := executeSelect(st, t)
		if err != nil {
			return err
		}
	case insertStatement:
		err := executeInsert(st, t)
		if err != nil {
			return err
		}
	}

	return nil
}

func executeInsert(st *statement, t *table) error {
	err := t.appendRow(&st.rowToInsert)
	if err != nil {
		errString := fmt.Sprintf("Error appending to table: %v\n", err)
		return errors.New(errString)
	}

	return nil
}

func executeSelect(st *statement, t *table) error {
	fmt.Printf("%10s %32s %32s\n", "ID", "Username", "Email")
	fmt.Println(strings.Repeat("-", 76))
	for i := 0; i < t.numRows; i++ {
		pRow := t.getRow(i)
		fmt.Printf("%10d %32s %32s\n", pRow.id, pRow.username, pRow.email)
	}
	return nil
}

func main() {
	var inputReader *os.File
	var err error
	interactive := true

	if len(os.Args) > 1 {
		inputReader, err = os.Open(os.Args[1])
		interactive = false
		if err != nil {
			fmt.Printf("Failed to open file: %v\n", err)
		}
	} else {
		inputReader = os.Stdin
	}

	stdinScanner := bufio.NewScanner(inputReader)
	prompt := "maksql>"

	tab := newTable()

	for {
		if interactive{
			fmt.Print(prompt)
		}

		success := stdinScanner.Scan()
		if !success {
			break
		}

		inputText := strings.TrimSpace(stdinScanner.Text())
		inputText = strings.ToLower(inputText)
		if len(inputText) == 0 {
			continue
		}

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
			executeStatement(&inputStatement, tab)
		}
	}

	if err := stdinScanner.Err(); err != nil {
		fmt.Printf("\nError: %v\n", err)
	}
}
