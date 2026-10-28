package service

import (
	"fmt"
	"time"

	"yoosoPosV3/internal/config"
	"yoosoPosV3/internal/model"
	"yoosoPosV3/internal/response"

	"gorm.io/gorm"
)

// ---------- 库存变动记录 ----------

func GetInventoryLogs(page, pageSize int, productName string, logType string, startDate, endDate string) ([]model.InventoryLog, int64, error) {
	var list []model.InventoryLog
	var total int64
	query := config.DB.Model(&model.InventoryLog{}).Preload("Product")
	if productName != "" {
		query = query.Joins("JOIN products ON products.id = inventory_logs.product_id").
			Where("products.name LIKE ?", "%"+productName+"%")
	}
	if logType != "" {
		query = query.Where("inventory_logs.type = ?", logType)
	}
	if startDate != "" {
		query = query.Where("inventory_logs.created_at >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("inventory_logs.created_at <= ?", endDate+" 23:59:59")
	}
	query.Count(&total)
	err := query.Order("inventory_logs.id desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}

// ---------- 库存盘点 ----------

func CheckInventory(productID uint, actualQty int, remark string) *response.AjaxResult {
	var product model.Product
	if err := config.DB.First(&product, productID).Error; err != nil {
		return response.ErrorWithMsg("商品不存在：" + err.Error())
	}
	diff := actualQty - product.StockQty
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&product).Update("stock_qty", actualQty).Error; err != nil {
			return err
		}
		log := model.InventoryLog{
			ProductID: productID,
			Type:      "check",
			QtyChange: diff,
			QtyAfter:  actualQty,
			OrderNo:   fmt.Sprintf("PD%s", time.Now().Format("20060102150405")),
			Remark:    remark,
		}
		return tx.Create(&log).Error
	})
	if err != nil {
		return response.ErrorWithMsg("盘点失败：" + err.Error())
	}
	return response.SuccessWithMsg("盘点完成")
}
