package reportqueue

import (
	"context"

	"git.bububla.com/kilogramix/asia/reporter.git/internal/defs"
	"git.bububla.com/kilogramix/asia/reporter.git/internal/model"
	"git.bububla.com/kilogramix/asia/reporter.git/internal/repository"
)

type TaskParams struct {
	ReportCode string
	ScheduleID uint64
	Formats    []string
}

type Publisher interface {
	Publish(ctx context.Context, params TaskParams) error
}

type JobsPublisher struct {
	jobs repository.Jobs
}

func NewJobsPublisher(jobs repository.Jobs) *JobsPublisher {
	return &JobsPublisher{jobs: jobs}
}

func (p *JobsPublisher) Publish(ctx context.Context, params TaskParams) error {
	job := &model.Job{
		ReportCode: params.ReportCode,
		Formats:    params.Formats,
		Status:     defs.JobNew,
	}
	if params.ScheduleID != 0 {
		job.ScheduleID = &params.ScheduleID
	}

	_, err := p.jobs.Create(job)
	return err
}

var _ Publisher = (*JobsPublisher)(nil)
