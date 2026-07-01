package report

import (
	"context"

	"git.bububla.com/kilogramix/asia/reporter.git/internal/defs"
)

type Dataset struct {
	name   string
	header []string
	rows   [][]string
}

func NewDataset(name string, header []string, rows [][]string) Dataset {
	return Dataset{name: name, header: header, rows: rows}
}

func (d Dataset) Name() string {
	return d.name
}

func (d Dataset) Header() []string {
	return d.header
}

func (d Dataset) Rows() [][]string {
	return d.rows
}

type ReportType interface {
	Code() defs.ReportCode
	Formatters() []defs.FormatterCode
	StorageCode() defs.StorageCode
	Init(formatters []defs.FormatterCode, storageCode defs.StorageCode) ReportType
	Build(ctx context.Context, params map[string]string) ([]Dataset, error)
}
