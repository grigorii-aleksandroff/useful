package repository

import "gorm.io/gorm"

type ReporterRepository struct {
	Schedules Schedules
	Jobs      Jobs
	Files     Files
}

func NewReporterRepository(db *gorm.DB) *ReporterRepository {
	return &ReporterRepository{
		Schedules: NewSchedules(db),
		Jobs:      NewJobs(db),
		Files:     NewFiles(db),
	}
}
