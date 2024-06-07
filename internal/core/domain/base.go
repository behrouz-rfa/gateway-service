package domain

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"time"
)

type Base struct {
	ID        string     `gorm:"type:uuid;primary_key;"  json:"id"` // TODO: switch back to uuid for mariadb
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt *time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (base *Base) BeforeCreate(*gorm.DB) error {

	base.ID = uuid.NewV4().String()

	return nil
}
