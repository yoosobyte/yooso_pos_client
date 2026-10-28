package model

import (
	"time"
)

// InventoryLog 库存变动记录
type InventoryLog struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	ProductID uint      `gorm:"not null;index" json:"productId"`
	Product   *Product  `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Type      string    `gorm:"size:50;not null" json:"type"` // purchase_in:采购入库, sales_out:销售出库, check:盘点调整, return_in:退货入库, return_out:采购退货出库
	QtyChange int       `gorm:"not null" json:"qtyChange"`
	QtyAfter  int       `gorm:"not null" json:"qtyAfter"`
	OrderNo   string    `gorm:"size:50" json:"orderNo"`
	Remark    string    `gorm:"size:500" json:"remark"`
}
