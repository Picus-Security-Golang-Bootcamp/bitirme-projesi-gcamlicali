package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Category struct {
	ID        uuid.UUID `gorm:"primary_key; type:uuid; default:uuid_generate_v4()"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      *string        `gorm:"unique"`
}

func (Category) TableName() string {
	//default table name
	return "categories"
}
