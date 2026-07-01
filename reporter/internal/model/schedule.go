package model

import "time"

type Schedule struct {
	ID         uint64     `gorm:"column:id;primaryKey"`
	Kind       string     `gorm:"column:kind"`
	ReportCode string     `gorm:"column:report_code"`
	TimeOfDay  string     `gorm:"column:time_of_day"`
	DayOfWeek  int        `gorm:"column:day_of_week"`
	DayOfMonth int        `gorm:"column:day_of_month"`
	Formats    StringList `gorm:"column:formats;type:json"`
	Timezone   string     `gorm:"column:timezone"`
	RunnerType string     `gorm:"column:runner_type"`
	RunnerID   string     `gorm:"column:runner_id"`
	Status     string     `gorm:"column:status"`
	Params     string     `gorm:"column:params;type:json"`
	NextRun    *time.Time `gorm:"column:next_run"`
	CreatedAt  time.Time  `gorm:"column:created_at"`
	UpdatedAt  time.Time  `gorm:"column:updated_at"`
}

func (Schedule) TableName() string {
	return "schedules"
}
