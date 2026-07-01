package formatter

import (
	"bytes"
	"encoding/csv"

	"git.bububla.com/kilogramix/asia/reporter.git/internal/defs"
)

type CSV struct{}

func NewCSV() CSV {
	return CSV{}
}

func (CSV) Code() defs.FormatterCode {
	return defs.FormatCSV
}

func (CSV) Extension() string {
	return "csv"
}

func (CSV) Format(data Dataset) ([]byte, error) {
	var buffer bytes.Buffer
	writer := csv.NewWriter(&buffer)

	if header := data.Header(); len(header) > 0 {
		if err := writer.Write(header); err != nil {
			return nil, err
		}
	}

	if err := writer.WriteAll(data.Rows()); err != nil {
		return nil, err
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

var _ Formatter = CSV{}
