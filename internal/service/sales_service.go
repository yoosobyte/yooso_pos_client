package service

import (
	"fmt"
	"time"
	"yoosoPosV3/internal/config"
	"yoosoPosV3/internal/model"
	"yoosoPosV3/internal/response"

	"gorm.io/gorm"
)

// ---------- 采购单 ----------

func GetPurchaseOrders(page, pageSize int, keyword string, status int) ([]model.PurchaseOrder, int64, error) {
	var list []model.PurchaseOrder
	var total int64
	query := config.DB.Model(&model.PurchaseOrder{}).Preload("Supplier").Preload("Items.Product.Category").Preload("Items.Product.Unit")
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Joins("LEFT JOIN suppliers ON suppliers.id = purchase_orders.supplier_id").
			Where("purchase_orders.order_no LIKE ? OR purchase_orders.remark LIKE ? OR suppliers.name LIKE ?", like, like, like)
	}
	if status > 0 {
		query = query.Where("purchase_orders.status = ?", status)
	}
	query.Count(&total)
	err := query.Order("purchase_orders.id desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}

func GetPurchaseOrderByID(id uint) (*model.PurchaseOrder, error) {
	var po model.PurchaseOrder
	err := config.DB.Preload("Supplier").Preload("Items.Product.Category").Preload("Items.Product.Unit").First(&po, id).Error
	return &po, err
}

func CreatePurchaseOrder(po *model.PurchaseOrder) *response.AjaxResult {
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		po.OrderNo = generateOrderNo("PO")
		var total float64
		for i := range po.Items {
			po.Items[i].Amount = float64(po.Items[i].Quantity) * po.Items[i].Price
			total += po.Items[i].Amount
		}
		po.TotalAmount = total
		return tx.Create(po).Error
	})
	if err != nil {
		return response.ErrorWithMsg(err.Error())
	}
	return response.SuccessWithMsg("采购单已创建")
}

func ConfirmPurchaseOrder(id uint) *response.AjaxResult {
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		var po model.PurchaseOrder
		if err := tx.Preload("Items").First(&po, id).Error; err != nil {
			return err
		}
		if po.Status != 1 {
			return fmt.Errorf("当前状态不可入库")
		}
		for _, item := range po.Items {
			if err := tx.Model(&model.Product{}).Where("id = ?", item.ProductID).
				UpdateColumn("stock_qty", tx.Raw("stock_qty + ?", item.Quantity)).Error; err != nil {
				return err
			}
			var product model.Product
			tx.First(&product, item.ProductID)
			log := model.InventoryLog{
				ProductID: item.ProductID,
				Type:      "purchase_in",
				QtyChange: item.Quantity,
				QtyAfter:  product.StockQty,
				OrderNo:   po.OrderNo,
				Remark:    "采购入库",
			}
			if err := tx.Create(&log).Error; err != nil {
				return err
			}
		}
		return tx.Model(&po).Update("status", 2).Error
	})
	if err != nil {
		return response.ErrorWithMsg(err.Error())
	}
	return response.SuccessWithMsg("入库成功")
}

func CancelPurchaseOrder(id uint) *response.AjaxResult {
	if err := config.DB.Model(&model.PurchaseOrder{}).Where("id = ?", id).Update("status", 3).Error; err != nil {
		return response.ErrorWithMsg(err.Error())
	}
	return response.SuccessWithMsg("已取消")
}

// ---------- 销售单 ----------

func GetSalesOrders(page, pageSize int, keyword string, status int) ([]model.SalesOrder, int64, error) {
	var list []model.SalesOrder
	var total int64
	query := config.DB.Model(&model.SalesOrder{}).Preload("Customer").Preload("Items.Product.Category").Preload("Items.Product.Unit")
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Joins("LEFT JOIN customers ON customers.id = sales_orders.customer_id").
			Where("sales_orders.order_no LIKE ? OR sales_orders.remark LIKE ? OR customers.name LIKE ?", like, like, like)
	}
	if status > 0 {
		query = query.Where("sales_orders.status = ?", status)
	}
	query.Count(&total)
	err := query.Order("sales_orders.id desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}

func GetSalesOrderByID(id uint) (*model.SalesOrder, error) {
	var so model.SalesOrder
	err := config.DB.Preload("Customer").Preload("Items.Product.Category").Preload("Items.Product.Unit").First(&so, id).Error
	return &so, err
}

func CreateSalesOrder(so *model.SalesOrder) *response.AjaxResult {
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		so.OrderNo = generateOrderNo("SO")
		so.Status = 1
		var total float64
		for i := range so.Items {
			so.Items[i].Amount = float64(so.Items[i].Quantity) * so.Items[i].Price
			total += so.Items[i].Amount
			var product model.Product
			if err := tx.First(&product, so.Items[i].ProductID).Error; err != nil {
				return fmt.Errorf("商品不存在")
			}
			if product.StockQty < so.Items[i].Quantity {
				return fmt.Errorf("商品[%s]库存不足，当前库存: %d", product.Name, product.StockQty)
			}
		}
		so.TotalAmount = total - so.Discount
		if so.TotalAmount < 0 {
			so.TotalAmount = 0
		}
		if err := tx.Create(so).Error; err != nil {
			return err
		}
		for _, item := range so.Items {
			if err := tx.Model(&model.Product{}).Where("id = ?", item.ProductID).
				UpdateColumn("stock_qty", tx.Raw("stock_qty - ?", item.Quantity)).Error; err != nil {
				return err
			}
			var product model.Product
			tx.First(&product, item.ProductID)
			log := model.InventoryLog{
				ProductID: item.ProductID,
				Type:      "sales_out",
				QtyChange: -item.Quantity,
				QtyAfter:  product.StockQty,
				OrderNo:   so.OrderNo,
				Remark:    "销售出库",
			}
			if err := tx.Create(&log).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return response.ErrorWithMsg(err.Error())
	}
	return response.SuccessWithMsg("收款成功")
}

func ReturnSalesOrder(id uint) *response.AjaxResult {
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		var so model.SalesOrder
		if err := tx.Preload("Items").First(&so, id).Error; err != nil {
			return err
		}
		if so.Status != 1 {
			return fmt.Errorf("当前状态不可退货")
		}
		for _, item := range so.Items {
			if err := tx.Model(&model.Product{}).Where("id = ?", item.ProductID).
				UpdateColumn("stock_qty", tx.Raw("stock_qty + ?", item.Quantity)).Error; err != nil {
				return err
			}
			var product model.Product
			tx.First(&product, item.ProductID)
			log := model.InventoryLog{
				ProductID: item.ProductID,
				Type:      "return_in",
				QtyChange: item.Quantity,
				QtyAfter:  product.StockQty,
				OrderNo:   so.OrderNo,
				Remark:    "销售退货",
			}
			if err := tx.Create(&log).Error; err != nil {
				return err
			}
		}
		return tx.Model(&so).Update("status", 2).Error
	})
	if err != nil {
		return response.ErrorWithMsg(err.Error())
	}
	return response.SuccessWithMsg("退货成功")
}

func generateOrderNo(prefix string) string {
	now := time.Now()
	return fmt.Sprintf("%s%s%06d", prefix, now.Format("20060102150405"), now.Nanosecond()/1000%1000000)
}
