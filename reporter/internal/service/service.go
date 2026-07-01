package service

import (
	"git.bububla.com/kilogramix/asia/reporter.git/internal/defs"
	"git.bububla.com/kilogramix/asia/reporter.git/internal/formatter"
	"git.bububla.com/kilogramix/asia/reporter.git/internal/report"
	"git.bububla.com/kilogramix/asia/reporter.git/internal/repository"
	"git.bububla.com/kilogramix/asia/reporter.git/internal/storage"
)

type ReportService interface {
	ScheduleService
	ReportTypeService
	FormatService
	FileService
}

type Service struct {
	*scheduleService
	*reportTypeService
	*formatService
	*fileService
}

func NewReporterServices(
	repo *repository.ReporterRepository,
	registry report.Registry,
	formatters map[defs.FormatterCode]formatter.Formatter,
	storages map[defs.StorageCode]storage.Storage,
) *Service {
	return &Service{
		scheduleService:   &scheduleService{schedules: repo.Schedules},
		reportTypeService: &reportTypeService{registry: registry},
		formatService:     &formatService{registry: registry, formatters: formatters},
		fileService:       &fileService{files: repo.Files, storages: storages},
	}
}

var _ ReportService = (*Service)(nil)
