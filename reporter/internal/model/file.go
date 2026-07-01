package model

import "time"

type File struct {
	ID          uint64    `gorm:"column:id;primaryKey"`
	JobID       uint64    `gorm:"column:job_id"`
	ReportCode  string    `gorm:"column:report_code"`
	Name        string    `gorm:"column:name"`
	Format      string    `gorm:"column:format"`
	StorageCode string    `gorm:"column:storage_code"`
	FilePath    string    `gorm:"column:file_path"`
	CreatedAt   time.Time `gorm:"column:created_at"`
}

func (File) TableName() string {
	return "files"
}
