package grpc

import (
	"fmt"

	"git.bububla.com/kilogramix/asia/reporter.git/internal/defs"
	"git.bububla.com/kilogramix/asia/reporter.git/internal/service"
)

func validateScheduleInput(in service.ScheduleInput) error {
	if in.ReportCode == "" {
		return fmt.Errorf("report_code is required")
	}

	switch in.Kind {
	case defs.KindOnce:
		if in.RunAt == nil {
			return fmt.Errorf("run_at is required for once schedule")
		}
	case defs.KindDaily:
		return requireTimeOfDay(in.TimeOfDay)
	case defs.KindWeekly:
		if err := requireTimeOfDay(in.TimeOfDay); err != nil {
			return err
		}
		if in.DayOfWeek < 1 || in.DayOfWeek > 7 {
			return fmt.Errorf("day_of_week must be between 1 and 7")
		}
	case defs.KindMonthly:
		if err := requireTimeOfDay(in.TimeOfDay); err != nil {
			return err
		}
		if in.DayOfMonth < 1 || in.DayOfMonth > 31 {
			return fmt.Errorf("day_of_month must be between 1 and 31")
		}
	default:
		return fmt.Errorf("unknown schedule kind %q", in.Kind)
	}

	return nil
}

func requireTimeOfDay(value string) error {
	if value == "" {
		return fmt.Errorf("time_of_day is required")
	}
	return nil
}
