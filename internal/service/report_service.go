package service

import (
	"yoosoPosV3/internal/config"
	"yoosoPosV3/internal/model"
)

// SalesReportItem 销售报表项
type SalesReportItem struct {
	CategoryName string  `json:"categoryName"`
	ProductName  string  `json:"productName"`
	Spec         string  `json:"spec"`
	TotalQty     int     `json:"totalQty"`
	TotalAmount  float64 `json:"totalAmount"`
}

// PurchaseReportItem 采购报表项
type PurchaseReportItem struct {
	CategoryName string  `json:"categoryName"`
	ProductName  string  `json:"productName"`
	Spec         string  `json:"spec"`
	TotalQty     int     `json:"totalQty"`
	TotalAmount  float64 `json:"totalAmount"`
}

// DashboardData 首页仪表盘数据
type DashboardData struct {
	TotalProducts    int64   `json:"totalProducts"`
	TotalCustomers   int64   `json:"totalCustomers"`
	TotalSuppliers   int64   `json:"totalSuppliers"`
	TodaySalesCount  int64   `json:"todaySalesCount"`
	TodaySalesAmount float64 `json:"todaySalesAmount"`
	MonthSalesAmount float64 `json:"monthSalesAmount"`
	LowStockCount    int64   `json:"lowStockCount"`
}

// GetDashboardData 获取首页仪表盘数据
func GetDashboardData() (*DashboardData, error) {
	d := &DashboardData{}
	config.DB.Model(&model.Product{}).Count(&d.TotalProducts)
	config.DB.Model(&model.Customer{}).Count(&d.TotalCustomers)
	config.DB.Model(&model.Supplier{}).Count(&d.TotalSuppliers)

	// 总单数
	config.DB.Model(&model.SalesOrder{}).
		Where("status = 1").
		Count(&d.TodaySalesCount)
	config.DB.Model(&model.SalesOrder{}).
		Where("date(created_at) = date('now', 'localtime') AND status = 1").
		Select("COALESCE(SUM(total_amount),0)").Scan(&d.TodaySalesAmount)

	// 本月销售
	config.DB.Model(&model.SalesOrder{}).
		Where("strftime('%Y-%m', created_at) = strftime('%Y-%m', 'now', 'localtime') AND status = 1").
		Select("COALESCE(SUM(total_amount),0)").Scan(&d.MonthSalesAmount)

	// 低库存预警
	config.DB.Model(&model.Product{}).
		Where("stock_qty <= min_stock AND min_stock > 0").
		Count(&d.LowStockCount)

	return d, nil
}

// GetSalesReport 销售报表
func GetSalesReport(startDate, endDate string) ([]SalesReportItem, error) {
	var items []SalesReportItem
	query := config.DB.Table("sales_order_items soi").
		Select("COALESCE(c.name,'') as category_name, p.name as product_name, COALESCE(p.spec,'') as spec, SUM(soi.quantity) as total_qty, SUM(soi.amount) as total_amount").
		Joins("JOIN sales_orders so ON so.id = soi.order_id").
		Joins("JOIN products p ON p.id = soi.product_id").
		Joins("LEFT JOIN categories c ON c.id = p.category_id").
		Where("so.status = 1 AND soi.deleted_at IS NULL")
	if startDate != "" {
		query = query.Where("so.created_at >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("so.created_at <= ?", endDate+" 23:59:59")
	}
	err := query.Group("p.name").Order("total_amount desc").Scan(&items).Error
	return items, err
}

// GetPurchaseReport 采购报表
func GetPurchaseReport(startDate, endDate string) ([]PurchaseReportItem, error) {
	var items []PurchaseReportItem
	query := config.DB.Table("purchase_order_items poi").
		Select("COALESCE(c.name,'') as category_name, p.name as product_name, COALESCE(p.spec,'') as spec, SUM(poi.quantity) as total_qty, SUM(poi.amount) as total_amount").
		Joins("JOIN purchase_orders po ON po.id = poi.order_id").
		Joins("JOIN products p ON p.id = poi.product_id").
		Joins("LEFT JOIN categories c ON c.id = p.category_id").
		Where("po.status = 2 AND poi.deleted_at IS NULL")
	if startDate != "" {
		query = query.Where("po.created_at >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("po.created_at <= ?", endDate+" 23:59:59")
	}
	err := query.Group("p.name").Order("total_amount desc").Scan(&items).Error
	return items, err
}

// GetLowStockProducts 获取低库存商品
func GetLowStockProducts() ([]model.Product, error) {
	var products []model.Product
	err := config.DB.Preload("Unit").Where("stock_qty <= min_stock AND min_stock > 0").Find(&products).Error
	return products, err
}
