package model

import (
	"time"

	"gorm.io/gorm"
)

// Supplier 供应商
type Supplier struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`
	Name      string         `gorm:"size:200;not null" json:"name"`
	Contact   string         `gorm:"size:50" json:"contact"`
	Phone     string         `gorm:"size:50" json:"phone"`
	Address   string         `gorm:"size:500" json:"address"`
	Remark    string         `gorm:"size:500" json:"remark"`
}
