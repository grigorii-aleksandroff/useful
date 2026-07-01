package repository

import (
	"time"

	"gorm.io/gorm"

	"git.bububla.com/kilogramix/asia/reporter.git/internal/defs"
	"git.bububla.com/kilogramix/asia/reporter.git/internal/helpers"
	"git.bububla.com/kilogramix/asia/reporter.git/internal/model"
)

type Schedules interface {
	Create(schedule *model.Schedule) (*model.Schedule, error)
	Update(schedule *model.Schedule) (*model.Schedule, error)
	Delete(id uint64) error
	Get(id uint64) (*model.Schedule, error)
	List(limit, offset int, runners []helpers.RunnerFilter) ([]model.Schedule, int64, error)
	FindDue(now time.Time) ([]model.Schedule, error)
	UpdateNextRun(id uint64, nextRun *time.Time) error
}

type gormSchedules struct {
	db *gorm.DB
}

func NewSchedules(db *gorm.DB) Schedules {
	return &gormSchedules{db: db}
}

func (r *gormSchedules) Create(schedule *model.Schedule) (*model.Schedule, error) {
	if err := r.db.Create(schedule).Error; err != nil {
		return nil, err
	}

	return schedule, nil
}

func (r *gormSchedules) Update(schedule *model.Schedule) (*model.Schedule, error) {
	if err := r.db.Save(schedule).Error; err != nil {
		return nil, err
	}

	return schedule, nil
}

func (r *gormSchedules) Delete(id uint64) error {
	return r.db.Model(&model.Schedule{}).
		Where("id = ?", id).
		Update("status", defs.ScheduleDeleted).Error
}

func (r *gormSchedules) Get(id uint64) (*model.Schedule, error) {
	var schedule model.Schedule
	if err := r.db.First(&schedule, id).Error; err != nil {
		return nil, err
	}

	return &schedule, nil
}

func (r *gormSchedules) List(limit, offset int, runners []helpers.RunnerFilter) ([]model.Schedule, int64, error) {
	var count int64
	if err := r.scoped(runners).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	query := r.scoped(runners).Order("id DESC")
	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}

	var schedules []model.Schedule
	if err := query.Find(&schedules).Error; err != nil {
		return nil, 0, err
	}

	return schedules, count, nil
}

func (r *gormSchedules) FindDue(now time.Time) ([]model.Schedule, error) {
	var schedules []model.Schedule

	err := r.db.
		Where("status = ?", defs.ScheduleActive).
		Where("next_run IS NOT NULL").
		Where("next_run <= ?", now).
		Find(&schedules).Error
	if err != nil {
		return nil, err
	}

	return schedules, nil
}

func (r *gormSchedules) UpdateNextRun(id uint64, nextRun *time.Time) error {
	return r.db.Model(&model.Schedule{}).
		Where("id = ?", id).
		Update("next_run", nextRun).Error
}

func (r *gormSchedules) scoped(runners []helpers.RunnerFilter) *gorm.DB {
	query := r.db.Model(&model.Schedule{}).Where("status <> ?", defs.ScheduleDeleted)

	condition, args := helpers.RunnerCondition("runner_type", "runner_id", runners)
	return query.Where(condition, args...)
}
