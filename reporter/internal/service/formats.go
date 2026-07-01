package service

import (
	"context"

	"git.bububla.com/kilogramix/asia/reporter.git/internal/defs"
	"git.bububla.com/kilogramix/asia/reporter.git/internal/formatter"
	"git.bububla.com/kilogramix/asia/reporter.git/internal/report"
)

type FormatService interface {
	Formatters(ctx context.Context, reportCode string) []formatter.Formatter
}

type formatService struct {
	registry   report.Registry
	formatters map[defs.FormatterCode]formatter.Formatter
}

func (s *formatService) Formatters(ctx context.Context, reportCode string) []formatter.Formatter {
	reportType, ok := s.registry.Get(defs.ReportCode(reportCode))
	if !ok {
		return nil
	}

	codes := reportType.Formatters()
	result := make([]formatter.Formatter, 0, len(codes))
	for _, code := range codes {
		if current, ok := s.formatters[code]; ok {
			result = append(result, current)
		}
	}
	return result
}

var _ FormatService = (*formatService)(nil)
