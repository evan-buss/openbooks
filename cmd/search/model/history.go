package model

import (
	"time"

	"gorm.io/gorm"
)

type History struct {
	gorm.Model
	ScheduleId   uint      `json:"schedule_id"`
	StartedAt    time.Time `json:"started_at"`
	EndedAt      time.Time `json:"ended_at"`
	ErrorCount   uint64    `json:"error_count"`
	SuccessCount uint64    `json:"success_count"`
}
