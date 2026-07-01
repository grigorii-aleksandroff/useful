package formatter

import "git.bububla.com/kilogramix/asia/reporter.git/internal/defs"

type Dataset interface {
	Name() string
	Header() []string
	Rows() [][]string
}

type Formatter interface {
	Code() defs.FormatterCode
	Extension() string
	Format(data Dataset) ([]byte, error)
}
