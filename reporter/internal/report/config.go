package report

import "git.bububla.com/kilogramix/asia/reporter.git/internal/defs"

type Config struct {
	Reports map[defs.ReportCode]ReportConfig
}

type ReportConfig struct {
	Formatters []defs.FormatterCode
	Storage    defs.StorageCode
}
