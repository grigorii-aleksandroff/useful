package scheduler

import (
	"context"
	"fmt"
	"time"

	loggerPkg "git.itechpsp.com/e46/box/platform.git/modules/logger"

	"git.bububla.com/kilogramix/asia/reporter.git/internal/model"
	"git.bububla.com/kilogramix/asia/reporter.git/internal/reportqueue"
	"git.bububla.com/kilogramix/asia/reporter.git/internal/repository"
	"git.bububla.com/kilogramix/asia/reporter.git/internal/service"
)

type Scheduler struct {
	logger    loggerPkg.Logger
	schedules repository.Schedules
	publisher reportqueue.Publisher
	interval  time.Duration
}

func New(
	logger loggerPkg.Logger,
	schedules repository.Schedules,
	publisher reportqueue.Publisher,
	interval time.Duration,
) *Scheduler {
	if interval <= 0 {
		interval = 30 * time.Second
	}

	return &Scheduler{
		logger:    logger,
		schedules: schedules,
		publisher: publisher,
		interval:  interval,
	}
}

func (s *Scheduler) Run(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.tick(ctx)
		}
	}
}

func (s *Scheduler) tick(ctx context.Context) {
	now := time.Now().UTC()

	due, err := s.schedules.FindDue(now)
	if err != nil {
		s.logger.Error(fmt.Sprintf("scheduler: cant read due schedules: %s", err))
		return
	}

	for i := range due {
		s.fire(ctx, due[i], now)
	}
}

func (s *Scheduler) fire(ctx context.Context, schedule model.Schedule, now time.Time) {
	err := s.publisher.Publish(ctx, reportqueue.TaskParams{
		ReportCode: schedule.ReportCode,
		ScheduleID: schedule.ID,
		Formats:    schedule.Formats,
	})
	if err != nil {
		s.logger.Error(fmt.Sprintf("scheduler: cant publish schedule %d: %s", schedule.ID, err))
		return
	}

	s.rearm(schedule, now)
}

func (s *Scheduler) rearm(schedule model.Schedule, now time.Time) {
	next, _ := service.NextRun(schedule.Kind, schedule.TimeOfDay, schedule.DayOfWeek, schedule.DayOfMonth, nil, now, schedule.Timezone)
	if err := s.schedules.UpdateNextRun(schedule.ID, next); err != nil {
		s.logger.Error(fmt.Sprintf("scheduler: cant update next_run for schedule %d: %s", schedule.ID, err))
	}
}
