package service

import (
	"context"
	"time"

	"git.bububla.com/kilogramix/asia/reporter.git/internal/defs"
	"git.bububla.com/kilogramix/asia/reporter.git/internal/helpers"
	"git.bububla.com/kilogramix/asia/reporter.git/internal/model"
	"git.bububla.com/kilogramix/asia/reporter.git/internal/repository"
)

type ScheduleInput struct {
	Kind       string
	ReportCode string
	TimeOfDay  string
	DayOfWeek  int
	DayOfMonth int
	RunAt      *time.Time
	Enabled    bool
	Formats    []string
	Timezone   string
	RunnerType string
	RunnerID   string
	Params     string
}

type ScheduleService interface {
	CreateSchedule(ctx context.Context, in ScheduleInput) (*model.Schedule, error)
	UpdateSchedule(ctx context.Context, id uint64, in ScheduleInput) (*model.Schedule, error)
	DeleteSchedule(ctx context.Context, id uint64) error
	ListSchedules(ctx context.Context, limit, offset int, runners []helpers.RunnerFilter) ([]model.Schedule, int64, error)
}

type scheduleService struct {
	schedules repository.Schedules
}

func (s *scheduleService) CreateSchedule(ctx context.Context, in ScheduleInput) (*model.Schedule, error) {
	schedule := &model.Schedule{}
	applyInput(schedule, in)
	schedule.RunnerType = in.RunnerType
	schedule.RunnerID = in.RunnerID

	return s.schedules.Create(schedule)
}

func (s *scheduleService) UpdateSchedule(ctx context.Context, id uint64, in ScheduleInput) (*model.Schedule, error) {
	schedule, err := s.schedules.Get(id)
	if err != nil {
		return nil, err
	}

	applyUpdate(schedule, in)

	return s.schedules.Update(schedule)
}

func (s *scheduleService) DeleteSchedule(ctx context.Context, id uint64) error {
	return s.schedules.Delete(id)
}

func (s *scheduleService) ListSchedules(ctx context.Context, limit, offset int, runners []helpers.RunnerFilter) ([]model.Schedule, int64, error) {
	return s.schedules.List(limit, offset, runners)
}

func applyInput(schedule *model.Schedule, in ScheduleInput) {
	timezone := tzOrUTC(in.Timezone)
	next, _ := NextRun(in.Kind, in.TimeOfDay, in.DayOfWeek, in.DayOfMonth, in.RunAt, time.Now().UTC(), timezone)

	schedule.Kind = in.Kind
	schedule.ReportCode = in.ReportCode
	schedule.TimeOfDay = in.TimeOfDay
	schedule.DayOfWeek = in.DayOfWeek
	schedule.DayOfMonth = in.DayOfMonth
	schedule.Status = statusFor(in.Enabled)
	schedule.Formats = in.Formats
	schedule.Timezone = timezone
	schedule.Params = orEmptyJSON(in.Params)
	schedule.NextRun = next
}

func applyUpdate(schedule *model.Schedule, in ScheduleInput) {
	preserved := schedule.NextRun

	applyInput(schedule, in)

	if in.Kind == defs.KindOnce && in.RunAt == nil {
		schedule.NextRun = preserved
	}
}

func statusFor(enabled bool) string {
	if enabled {
		return defs.ScheduleActive
	}
	return defs.ScheduleDisabled
}

func tzOrUTC(value string) string {
	if value == "" {
		return "UTC"
	}
	return value
}

func orEmptyJSON(value string) string {
	if value == "" {
		return "{}"
	}
	return value
}

var _ ScheduleService = (*scheduleService)(nil)
