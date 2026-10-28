package model

import (
	"time"

	"gorm.io/gorm"
)

// Customer 客户
type Customer struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`
	Name      string         `gorm:"size:200;not null" json:"name"`
	Phone     string         `gorm:"size:50" json:"phone"`
	Address   string         `gorm:"size:500" json:"address"`
	Remark    string         `gorm:"size:500" json:"remark"`
}
