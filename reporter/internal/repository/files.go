package repository

import (
	"gorm.io/gorm"

	"git.bububla.com/kilogramix/asia/reporter.git/internal/helpers"
	"git.bububla.com/kilogramix/asia/reporter.git/internal/model"
)

type Files interface {
	Create(file *model.File) (*model.File, error)
	Get(id uint64, runners []helpers.RunnerFilter) (*model.File, error)
	List(limit, offset int, runners []helpers.RunnerFilter) ([]model.File, int64, error)
}

type gormFiles struct {
	db *gorm.DB
}

func NewFiles(db *gorm.DB) Files {
	return &gormFiles{db: db}
}

func (r *gormFiles) Create(file *model.File) (*model.File, error) {
	if err := r.db.Create(file).Error; err != nil {
		return nil, err
	}

	return file, nil
}

func (r *gormFiles) Get(id uint64, runners []helpers.RunnerFilter) (*model.File, error) {
	var file model.File
	if err := r.scoped(runners).Where("files.id = ?", id).First(&file).Error; err != nil {
		return nil, err
	}

	return &file, nil
}

func (r *gormFiles) List(limit, offset int, runners []helpers.RunnerFilter) ([]model.File, int64, error) {
	var count int64
	if err := r.scoped(runners).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	query := r.scoped(runners).Order("files.id DESC")
	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}

	var files []model.File
	if err := query.Find(&files).Error; err != nil {
		return nil, 0, err
	}

	return files, count, nil
}

func (r *gormFiles) scoped(runners []helpers.RunnerFilter) *gorm.DB {
	condition, args := helpers.RunnerCondition("schedules.runner_type", "schedules.runner_id", runners)

	return r.db.Model(&model.File{}).
		Joins("JOIN jobs ON jobs.id = files.job_id").
		Joins("JOIN schedules ON schedules.id = jobs.schedule_id").
		Where(condition, args...)
}
