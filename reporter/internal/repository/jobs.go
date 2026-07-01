package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"git.bububla.com/kilogramix/asia/reporter.git/internal/defs"
	"git.bububla.com/kilogramix/asia/reporter.git/internal/model"
)

type Jobs interface {
	Create(job *model.Job) (*model.Job, error)
	ClaimNext(limit int) ([]model.Job, error)
	Complete(id uint64) error
	Fail(id uint64, reason string) error
}

type gormJobs struct {
	db *gorm.DB
}

func NewJobs(db *gorm.DB) Jobs {
	return &gormJobs{db: db}
}

func (r *gormJobs) Create(job *model.Job) (*model.Job, error) {
	if err := r.db.Create(job).Error; err != nil {
		return nil, err
	}

	return job, nil
}

func (r *gormJobs) ClaimNext(limit int) ([]model.Job, error) {
	var jobs []model.Job

	err := r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.
			Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
			Where("status = ?", defs.JobNew).
			Order("id ASC").
			Limit(limit).
			Find(&jobs).Error
		if err != nil {
			return err
		}

		if len(jobs) == 0 {
			return nil
		}

		return tx.Model(&model.Job{}).
			Where("id IN ?", jobIDs(jobs)).
			Update("status", defs.JobProcessing).Error
	})

	return jobs, err
}

func (r *gormJobs) Complete(id uint64) error {
	return r.update(id, map[string]interface{}{"status": defs.JobComplete})
}

func (r *gormJobs) Fail(id uint64, reason string) error {
	return r.update(id, map[string]interface{}{
		"status": defs.JobFailed,
		"error":  reason,
	})
}

func (r *gormJobs) update(id uint64, values map[string]interface{}) error {
	return r.db.Model(&model.Job{}).Where("id = ?", id).Updates(values).Error
}

func jobIDs(jobs []model.Job) []uint64 {
	ids := make([]uint64, len(jobs))
	for i := range jobs {
		ids[i] = jobs[i].ID
	}
	return ids
}
