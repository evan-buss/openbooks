package model

import (
	"context"

	"gorm.io/gorm"
)

type Schedule struct {
	gorm.Model
	Server   string    `json:"server" validate:"required"`
	Schedule string    `json:"schedule" validate:"required"`
	History  []History `json:"history"`
}

func FindAllSchedules(db *gorm.DB, ctx context.Context) ([]Schedule, error) {
	var schedules []Schedule
	tx := db.WithContext(ctx).Find(&schedules)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return schedules, nil
}

func CreateSchedule(db *gorm.DB, ctx context.Context, schedule *Schedule) error {
	result := db.WithContext(ctx).Create(&schedule)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func GetScheduleById(db *gorm.DB, ctx context.Context, id uint) (Schedule, error) {
	var schedule Schedule
	tx := db.WithContext(ctx).First(&schedule, id)
	if tx.Error != nil {
		return Schedule{}, tx.Error
	}
	return schedule, nil
}

func UpdateSchedule(db *gorm.DB, ctx context.Context, schedule *Schedule) error {
	result := db.WithContext(ctx).Save(&schedule)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func DeleteSchedule(db *gorm.DB, ctx context.Context, id uint) error {
	result := db.WithContext(ctx).Delete(&Schedule{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
