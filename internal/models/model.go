package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const NewNotificationStatus = "NEW"
const ReviewedNoificationStatus = "REVIEWED"

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
	NotifyId string  `json:"notify_id" gorm:"type:uuid"`
	FromId   string  `json:"from_id"   gorm:"type:uuid"`
	FromName string  `json:"from_name"`
	TargetId *string `json:"target_id" gorm:"type:uuid;"`
	Type     string  `json:"type"`
	Status   string  `json:"status"`
}

func (n *Notification) BeforeCreate(tx *gorm.DB) (err error) {
	n.Status = NewNotificationStatus
	return n.Base.BeforeCreate(tx)
}
