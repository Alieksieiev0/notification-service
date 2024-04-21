package services

import (
	"context"
	"fmt"

	"github.com/Alieksieiev0/notification-service/internal/models"
	"gorm.io/gorm"
)

type Service interface {
	Get(c context.Context, params ...Param) ([]models.Notification, error)
	Save(c context.Context, notification *models.Notification) error
}

func NewService(db *gorm.DB) Service {
	return &service{
		db: db,
	}
}

type service struct {
	db *gorm.DB
}

func (s *service) Get(c context.Context, params ...Param) ([]models.Notification, error) {
	notifications := []models.Notification{}
	db := ApplyParams(s.db, params...)
	return notifications, db.Find(&notifications).Error
}

func (s *service) Save(c context.Context, notification *models.Notification) error {
	return s.db.Save(notification).Error
}

type Param func(db *gorm.DB) *gorm.DB

func Limit(limit int) Param {
	return func(db *gorm.DB) *gorm.DB {
		return db.Limit(limit)
	}
}

func Offset(offset int) Param {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(offset)
	}
}

func Filter(name string, value string, isStrict bool) Param {
	return func(db *gorm.DB) *gorm.DB {
		if isStrict {
			return db.Where(name+"= ?", value)
		}
		return db.Where(
			fmt.Sprintf("LOWER(%s) LIKE LOWER(?)", name),
			fmt.Sprintf("%%%s%%", value),
		)
	}
}

func ApplyParams(db *gorm.DB, params ...Param) *gorm.DB {
	for _, param := range params {
		db = param(db)
	}
	return db
}
