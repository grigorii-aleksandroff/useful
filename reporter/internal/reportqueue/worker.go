package reportqueue

import (
	"context"
	"fmt"
	"time"

	loggerPkg "git.itechpsp.com/e46/box/platform.git/modules/logger"

	"git.bububla.com/kilogramix/asia/reporter.git/internal/model"
	"git.bububla.com/kilogramix/asia/reporter.git/internal/repository"
	"git.bububla.com/kilogramix/asia/reporter.git/internal/service"
)

type Worker struct {
	logger    loggerPkg.Logger
	jobs      repository.Jobs
	exports   service.ExportService
	interval  time.Duration
	batchSize int
}

func NewWorker(
	logger loggerPkg.Logger,
	jobs repository.Jobs,
	exports service.ExportService,
	interval time.Duration,
	batchSize int,
) *Worker {
	if interval <= 0 {
		interval = 5 * time.Second
	}
	if batchSize <= 0 {
		batchSize = 5
	}

	return &Worker{
		logger:    logger,
		jobs:      jobs,
		exports:   exports,
		interval:  interval,
		batchSize: batchSize,
	}
}

func (w *Worker) Run(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.drain(ctx)
		}
	}
}

func (w *Worker) drain(ctx context.Context) {
	jobs, err := w.jobs.ClaimNext(w.batchSize)
	if err != nil {
		w.logger.Error(fmt.Sprintf("worker: cant claim jobs: %s", err))
		return
	}

	for i := range jobs {
		w.process(ctx, jobs[i])
	}
}

func (w *Worker) process(ctx context.Context, job model.Job) {
	if err := w.exports.Produce(ctx, job); err != nil {
		w.fail(job.ID, err.Error())
		return
	}

	if err := w.jobs.Complete(job.ID); err != nil {
		w.logger.Error(fmt.Sprintf("worker: cant complete job %d: %s", job.ID, err))
	}
}

func (w *Worker) fail(id uint64, reason string) {
	w.logger.Error(fmt.Sprintf("worker: job %d failed: %s", id, reason))
	if err := w.jobs.Fail(id, reason); err != nil {
		w.logger.Error(fmt.Sprintf("worker: cant mark job %d failed: %s", id, err))
	}
}
