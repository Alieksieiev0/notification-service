package services

import (
	"context"

	"github.com/Alieksieiev0/notification-service/internal/models"
	"gorm.io/gorm"
)

type Service interface {
	GetByNotifyId(c context.Context, id string) ([]models.Notification, error)
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

func (ns *service) GetByNotifyId(c context.Context, id string) ([]models.Notification, error) {
	notifications := []models.Notification{}
	return notifications, ns.db.Find(&notifications, "notify_id = ?", id).Error
}

func (ns *service) Save(c context.Context, notification *models.Notification) error {
	return ns.db.Save(notification).Error
}
