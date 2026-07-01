package model

import "time"

type Job struct {
	ID         uint64     `gorm:"column:id;primaryKey"`
	ScheduleID *uint64    `gorm:"column:schedule_id"`
	ReportCode string     `gorm:"column:report_code"`
	Formats    StringList `gorm:"column:formats;type:json"`
	Status     string     `gorm:"column:status"`
	Error      string     `gorm:"column:error"`
	CreatedAt  time.Time  `gorm:"column:created_at"`
	UpdatedAt  time.Time  `gorm:"column:updated_at"`
}

func (Job) TableName() string {
	return "jobs"
}
