package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/tealeg/xlsx/v3"
)

// SheetNamesTemplate define name's for new created sheets.
var SheetNamesTemplate = "Sheet %d"

func main() {
	initCommandLine(os.Args)
}

func writeAllSheets(xlFile *xlsx.File, dataFiles []string, sheetNames []string, exampleRowNumber int) (err error) {

	for i, dataFileName := range dataFiles {

		sheet, err := getSheet(xlFile, sheetNames, i)
		if err != nil {
			return err
		}

		var exampleRow *xlsx.Row
		if exampleRowNumber != 0 && exampleRowNumber <= sheet.MaxRow {
			// example row counting from 1
			exampleRow, _ = sheet.Row(exampleRowNumber - 1)

			sheet.RemoveRowAtIndex(exampleRowNumber - 1)
		}

		err = writeSheet(dataFileName, sheet, exampleRow)

		if err != nil {
			return err
		}
	}

	return nil
}

func getSheet(xlFile *xlsx.File, sheetNames []string, i int) (sheet *xlsx.Sheet, err error) {

	var sheetName string
	if len(sheetNames) > i {
		sheetName = sheetNames[i]
	} else {
		sheetName = fmt.Sprintf(SheetNamesTemplate, i+1)
	}

	sheet, ok := xlFile.Sheet[sheetName]
	if ok != true {
		sheet, err = xlFile.AddSheet(sheetName)

		if err != nil {
			return nil, err
		}
	}
	return sheet, nil
}

func writeSheet(dataFileName string, sheet *xlsx.Sheet, exampleRow *xlsx.Row) error {

	data, err := getCsvData(dataFileName)

	if err != nil {
		return err
	}

	var i int
	for {
		record, err := data.Read()

		if err == io.EOF || record == nil {
			break
		} else if err != nil {
			return err
		}

		// if i > 5000 {
		//	break
		// }

		// if i%500 == 0 {
		// 	fmt.Println(i)
		// }

		i++

		writeRowToXls(sheet, record, exampleRow)
	}

	return nil
}

func buildXls(p *params) (err error) {

	var xlFile *xlsx.File
	if p.xlsxTemplate == "" {
		xlFile = xlsx.NewFile()
	} else {
		xlFile, err = xlsx.OpenFile(p.xlsxTemplate)
		if err != nil {
			return err
		}
	}

	writeAllSheets(xlFile, p.input, p.sheets, p.row)

	return xlFile.Save(p.output)
}

func writeRowToXls(sheet *xlsx.Sheet, record []string, exampleRow *xlsx.Row) {

	var row *xlsx.Row
	var cell *xlsx.Cell

	row = sheet.AddRow()

	var cellsLen int
	if exampleRow != nil {
		cellsLen = exampleRow.Sheet.MaxCol
	}

	for k, v := range record {
		cell = row.AddCell()

		cell.Value = v

		if exampleRow != nil && cellsLen > k {
			style := exampleRow.GetCell(k).GetStyle()

			cell.SetStyle(style)
		}
	}
}


// getCsvData read's data from CSV file.
func getCsvData(dataFileName string) (*csv.Reader, error) {
	dataFile, err := os.Open(dataFileName)
	if err != nil {
		return nil, errors.New("Problem with reading data from " + dataFileName)
	}

	return csv.NewReader(dataFile), nil
}
