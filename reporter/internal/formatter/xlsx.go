package formatter

import (
	"git.bububla.com/kilogramix/asia/reporter.git/internal/defs"
	"github.com/xuri/excelize/v2"
)

const xlsxSheet = "Sheet1"

type XLSX struct{}

func NewXLSX() XLSX {
	return XLSX{}
}

func (XLSX) Code() defs.FormatterCode {
	return defs.FormatXLSX
}

func (XLSX) Extension() string {
	return "xlsx"
}

func (XLSX) Format(data Dataset) ([]byte, error) {
	file := excelize.NewFile()
	defer file.Close()

	sheet := xlsxSheet
	if name := data.Name(); name != "" {
		if err := file.SetSheetName(sheet, name); err != nil {
			return nil, err
		}
		sheet = name
	}

	rowIndex := 1

	if header := data.Header(); len(header) > 0 {
		if err := writeRow(file, sheet, rowIndex, header); err != nil {
			return nil, err
		}
		rowIndex++
	}

	for _, row := range data.Rows() {
		if err := writeRow(file, sheet, rowIndex, row); err != nil {
			return nil, err
		}
		rowIndex++
	}

	buffer, err := file.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func writeRow(file *excelize.File, sheet string, rowIndex int, values []string) error {
	cell, err := excelize.CoordinatesToCellName(1, rowIndex)
	if err != nil {
		return err
	}

	return file.SetSheetRow(sheet, cell, &values)
}

var _ Formatter = XLSX{}
