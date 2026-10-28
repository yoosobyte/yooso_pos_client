package model

import (
	"time"

	"gorm.io/gorm"
)

// PurchaseOrder 采购单
type PurchaseOrder struct {
	ID          uint                `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time           `json:"createdAt"`
	UpdatedAt   time.Time           `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt      `gorm:"index" json:"deletedAt"`
	OrderNo     string              `gorm:"size:50;not null;uniqueIndex" json:"orderNo"`
	SupplierID  *uint               `gorm:"default:null" json:"supplierId"`
	Supplier    *Supplier           `gorm:"foreignKey:SupplierID" json:"supplier,omitempty"`
	TotalAmount float64             `gorm:"type:decimal(10,2);default:0" json:"totalAmount"`
	Status      int                 `gorm:"default:1" json:"status"` // 1草稿 2已入库 3已取消
	Remark      string              `gorm:"size:500" json:"remark"`
	Items       []PurchaseOrderItem `gorm:"foreignKey:OrderID" json:"items,omitempty"`
}

// PurchaseOrderItem 采购单明细
type PurchaseOrderItem struct {
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
