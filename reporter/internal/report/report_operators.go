package report

import (
	"context"

	"git.bububla.com/kilogramix/asia/reporter.git/internal/defs"
)

func init() {
	registerType(&Operators{})
}

type Operators struct {
	formatters  []defs.FormatterCode
	storageCode defs.StorageCode
}

func (o *Operators) Code() defs.ReportCode {
	return defs.ReportOperators
}

func (o *Operators) Formatters() []defs.FormatterCode {
	return o.formatters
}

func (o *Operators) StorageCode() defs.StorageCode {
	return o.storageCode
}

func (o *Operators) Init(formatters []defs.FormatterCode, storageCode defs.StorageCode) ReportType {
	o.formatters = formatters
	o.storageCode = storageCode
	return o
}

func (o *Operators) Build(ctx context.Context, params map[string]string) ([]Dataset, error) {
	return []Dataset{
		NewDataset(
			"operators",
			[]string{"id", "name", "created_at"},
			[][]string{
				{"1", "stub-operator-1", "2026-01-01"},
				{"2", "stub-operator-2", "2026-01-02"},
			},
		),
	}, nil
}
