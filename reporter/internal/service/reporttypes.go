package service

import (
	"context"

	"git.bububla.com/kilogramix/asia/reporter.git/internal/report"
)

type ReportTypeService interface {
	ListReportTypes(ctx context.Context) []report.ReportType
}

type reportTypeService struct {
	registry report.Registry
}

var _ ReportTypeService = (*reportTypeService)(nil)

func (s *reportTypeService) ListReportTypes(ctx context.Context) []report.ReportType {
	return s.registry.List()
}
