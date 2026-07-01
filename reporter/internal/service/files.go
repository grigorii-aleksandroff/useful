package service

import (
	"context"
	"fmt"

	"git.bububla.com/kilogramix/asia/reporter.git/internal/defs"
	"git.bububla.com/kilogramix/asia/reporter.git/internal/helpers"
	"git.bububla.com/kilogramix/asia/reporter.git/internal/model"
	"git.bububla.com/kilogramix/asia/reporter.git/internal/repository"
	"git.bububla.com/kilogramix/asia/reporter.git/internal/storage"
)

type FileService interface {
	ListFiles(ctx context.Context, limit, offset int, runners []helpers.RunnerFilter) ([]model.File, int64, error)
	DownloadFile(ctx context.Context, id uint64, runners []helpers.RunnerFilter) (model.File, storage.Download, error)
}

type fileService struct {
	files    repository.Files
	storages map[defs.StorageCode]storage.Storage
}

func (s *fileService) ListFiles(ctx context.Context, limit, offset int, runners []helpers.RunnerFilter) ([]model.File, int64, error) {
	return s.files.List(limit, offset, runners)
}

func (s *fileService) DownloadFile(ctx context.Context, id uint64, runners []helpers.RunnerFilter) (model.File, storage.Download, error) {
	file, err := s.files.Get(id, runners)
	if err != nil {
		return model.File{}, storage.Download{}, err
	}

	store, ok := s.storages[defs.StorageCode(file.StorageCode)]
	if !ok {
		return model.File{}, storage.Download{}, fmt.Errorf("unknown storage code %q", file.StorageCode)
	}

	download, err := store.Download(file.FilePath)
	if err != nil {
		return model.File{}, storage.Download{}, err
	}

	return *file, download, nil
}

var _ FileService = (*fileService)(nil)
