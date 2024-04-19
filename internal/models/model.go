package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const newNotificationStatus = "NEW"

type Base struct {
	ID        string         `gorm:"type:uuid" json:"id"`
	CreatedAt time.Time      `                 json:"created_at"`
	UpdatedAt time.Time      `                 json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"     json:"deleted_at"`
}

func (b *Base) BeforeCreate(tx *gorm.DB) (err error) {
	if b.ID == "" {
		b.ID = uuid.New().String()
	}
	return
}

type Notification struct {
	Base
	NotifyId string `gorm:"type:uuid"`
	FromId   string `gorm:"type:uuid"`
	FromName string
	TargetId *string `gorm:"type:uuid"`
	Type     string
	Status   string
}

func (n *Notification) BeforeCreate(tx *gorm.DB) (err error) {
	n.Status = newNotificationStatus
	return n.Base.BeforeCreate(tx)
}
