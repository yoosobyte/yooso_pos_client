package model

import (
	"time"

	"gorm.io/gorm"
)

// SalesOrder 销售单
type SalesOrder struct {
	ID          uint             `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time        `json:"createdAt"`
	UpdatedAt   time.Time        `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt   `gorm:"index" json:"deletedAt"`
	OrderNo     string           `gorm:"size:50;not null;uniqueIndex" json:"orderNo"`
	CustomerID  *uint            `gorm:"default:null" json:"customerId"`
	Customer    *Customer        `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	TotalAmount float64          `gorm:"type:decimal(10,2);default:0" json:"totalAmount"`
	Discount    float64          `gorm:"type:decimal(10,2);default:0" json:"discount"`
	PaidAmount  float64          `gorm:"type:decimal(10,2);default:0" json:"paidAmount"`
	Status      int              `gorm:"default:1" json:"status"` // 1已完成 2已退货 3已取消
	Remark      string           `gorm:"size:500" json:"remark"`
	Items       []SalesOrderItem `gorm:"foreignKey:OrderID" json:"items,omitempty"`
}

// SalesOrderItem 销售单明细
type SalesOrderItem struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`
	OrderID   uint           `gorm:"not null;index" json:"orderId"`
	ProductID uint           `gorm:"not null" json:"productId"`
	Product   *Product       `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Quantity  int            `gorm:"not null" json:"quantity"`
	Price     float64        `gorm:"type:decimal(10,2);default:0" json:"price"`
	Amount    float64        `gorm:"type:decimal(10,2);default:0" json:"amount"`
	Remark    string         `gorm:"size:255" json:"remark"`
}
