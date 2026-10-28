package model

import (
	"time"

	"gorm.io/gorm"
)

// Category 商品分类
type Category struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`
	Name      string         `gorm:"size:100;not null" json:"name"`
	ParentID  *uint          `gorm:"default:null" json:"parentId"`
	Sort      int            `gorm:"default:0" json:"sort"`
	Children  []*Category    `gorm:"-" json:"children,omitempty"`
}

// Unit 商品单位
type Unit struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`
	Name      string         `gorm:"size:50;not null;unique" json:"name"`
}

// Product 商品
type Product struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deletedAt"`
	Name       string         `gorm:"size:200;not null" json:"name"`
	Code       string         `gorm:"size:50;unique" json:"code"`
	CategoryID *uint          `gorm:"default:null" json:"categoryId"`
	Category   *Category      `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	UnitID     *uint          `gorm:"default:null" json:"unitId"`
	Unit       *Unit          `gorm:"foreignKey:UnitID" json:"unit,omitempty"`
	SalePrice  float64        `gorm:"type:decimal(10,2);default:0" json:"salePrice"`
	CostPrice  float64        `gorm:"type:decimal(10,2);default:0" json:"costPrice"`
	MinStock   int            `gorm:"default:0" json:"minStock"`
	MaxStock   int            `gorm:"default:0" json:"maxStock"`
	StockQty   int            `gorm:"default:0" json:"stockQty"`
	Spec       string         `gorm:"size:200" json:"spec"`
	Remark     string         `gorm:"size:500" json:"remark"`
	Status     int            `gorm:"default:1" json:"status"` // 1上架 0下架
}
