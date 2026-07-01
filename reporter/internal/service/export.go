package service

import (
	"context"
	"fmt"

	"git.bububla.com/kilogramix/asia/reporter.git/internal/defs"
	"git.bububla.com/kilogramix/asia/reporter.git/internal/formatter"
	"git.bububla.com/kilogramix/asia/reporter.git/internal/model"
	"git.bububla.com/kilogramix/asia/reporter.git/internal/report"
	"git.bububla.com/kilogramix/asia/reporter.git/internal/repository"
	"git.bububla.com/kilogramix/asia/reporter.git/internal/storage"
)

type ExportService interface {
	Produce(ctx context.Context, job model.Job) error
}

type exportService struct {
	registry   report.Registry
	files      repository.Files
	formatters map[defs.FormatterCode]formatter.Formatter
	storages   map[defs.StorageCode]storage.Storage
}

func NewExportService(
	registry report.Registry,
	files repository.Files,
	formatters map[defs.FormatterCode]formatter.Formatter,
	storages map[defs.StorageCode]storage.Storage,
) ExportService {
	return &exportService{registry: registry, files: files, formatters: formatters, storages: storages}
}

func (s *exportService) Produce(ctx context.Context, job model.Job) error {
	reportType, ok := s.registry.Get(defs.ReportCode(job.ReportCode))
	if !ok {
		return fmt.Errorf("unknown report code %q", job.ReportCode)
	}

	datasets, err := reportType.Build(ctx, nil)
	if err != nil {
		return err
	}

	store := s.storage(reportType.StorageCode())

	for _, code := range jobFormats(job) {
		current, ok := s.formatters[defs.FormatterCode(code)]
		if !ok {
			return fmt.Errorf("unknown format %q", code)
		}

		for _, data := range datasets {
			if err := s.writeFile(job, store, current, data); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *exportService) writeFile(job model.Job, store storage.Storage, current formatter.Formatter, data report.Dataset) error {
	content, err := current.Format(data)
	if err != nil {
		return err
	}

	fileName := fmt.Sprintf("report_%d_%s.%s", job.ID, data.Name(), current.Extension())
	path, err := store.Save(fileName, content)
	if err != nil {
		return err
	}

	_, err = s.files.Create(&model.File{
		JobID:       job.ID,
		ReportCode:  job.ReportCode,
		Name:        data.Name(),
		Format:      current.Code().ToString(),
		StorageCode: store.Code().ToString(),
		FilePath:    path,
	})
	return err
}

func (s *exportService) storage(code defs.StorageCode) storage.Storage {
	if store, ok := s.storages[code]; ok {
		return store
	}

	return s.storages[defs.StorageFile]
}

func jobFormats(job model.Job) []string {
	if len(job.Formats) == 0 {
		return []string{defs.FormatCSV.ToString()}
	}

	return job.Formats
}
