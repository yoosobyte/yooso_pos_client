package main

import (
	"context"
	"fmt"
	"yoosoPosV3/internal/config"
	"yoosoPosV3/internal/model"
	"yoosoPosV3/internal/response"
	"yoosoPosV3/internal/service"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type App struct{}

func NewApp() *App { return &App{} }

func (a *App) ServiceName() string { return "App" }

func (a *App) ServiceStartup(ctx context.Context, _ application.ServiceOptions) error {
	if err := config.InitDB(); err != nil {
		fmt.Printf("数据库初始化失败: %v\n", err)
	}
	// 后台下载更新器
	go service.DownloadUpdater()
	return nil
}

// ==================== 首页仪表盘 ====================

func (a *App) GetDashboard() *response.AjaxResult {
	data, err := service.GetDashboardData()
	if err != nil {
		return response.ErrorWithMsg(err.Error())
	}
	return response.SuccessWithData(data)
}

func (a *App) GetLowStockProducts() *response.AjaxResult {
	list, err := service.GetLowStockProducts()
	if err != nil {
		return response.ErrorWithMsg(err.Error())
	}
	if list == nil {
		list = []model.Product{}
	}
	return response.SuccessWithData(list)
}

// ==================== 商品分类 ====================

func (a *App) GetCategories() *response.AjaxResult {
	list, err := service.GetCategories()
	if err != nil {
		return response.ErrorWithMsg(err.Error())
	}
	if list == nil {
		list = []model.Category{}
	}
	return response.SuccessWithData(list)
}

func (a *App) GetCategoryByID(id uint) *response.AjaxResult {
	c, err := service.GetCategoryByID(id)
	if err != nil {
		return response.ErrorWithMsg(err.Error())
	}
	return response.SuccessWithData(c)
}

func (a *App) CreateCategory(c model.Category) *response.AjaxResult {
	return service.CreateCategory(&c)
}

func (a *App) UpdateCategory(c model.Category) *response.AjaxResult {
	return service.UpdateCategory(&c)
}

func (a *App) DeleteCategory(id uint) *response.AjaxResult {
	return service.DeleteCategory(id)
}

func (a *App) GetMaxCategorySort(parentID uint) *response.AjaxResult {
	sort := service.GetMaxCategorySort(parentID)
	return response.SuccessWithData(sort)
}

// ==================== 商品单位 ====================

func (a *App) GetUnits() *response.AjaxResult {
	list, err := service.GetUnits()
	if err != nil {
		return response.ErrorWithMsg(err.Error())
	}
	if list == nil {
		list = []model.Unit{}
	}
	return response.SuccessWithData(list)
}

func (a *App) CreateUnit(u model.Unit) *response.AjaxResult {
	return service.CreateUnit(&u)
}

func (a *App) UpdateUnit(unit *model.Unit) *response.AjaxResult {
	return service.UpdateUnit(unit)
}

func (a *App) DeleteUnit(id uint) *response.AjaxResult {
	return service.DeleteUnit(id)
}

// ==================== 商品管理 ====================

type ProductListResult struct {
	List  []model.Product `json:"list"`
	Total int64           `json:"total"`
}

func (a *App) GetProducts(page, pageSize int, keyword string, categoryID uint) *response.AjaxResult {
	list, total, err := service.GetProducts(page, pageSize, keyword, categoryID)
	if err != nil {
		return response.ErrorWithMsg(err.Error())
	}
	if list == nil {
		list = []model.Product{}
	}
	return response.SuccessWithData(ProductListResult{List: list, Total: total})
}

func (a *App) GetAllProducts() *response.AjaxResult {
	list, err := service.GetAllProducts()
	if err != nil {
		return response.ErrorWithMsg(err.Error())
	}
	if list == nil {
		list = []model.Product{}
	}
	return response.SuccessWithData(list)
}

func (a *App) GetProductByID(id uint) *response.AjaxResult {
	p, err := service.GetProductByID(id)
	if err != nil {
		return response.ErrorWithMsg(err.Error())
	}
	return response.SuccessWithData(p)
}

func (a *App) CreateProduct(p model.Product) *response.AjaxResult {
	return service.CreateProduct(&p)
}

func (a *App) UpdateProduct(p model.Product) *response.AjaxResult {
	return service.UpdateProduct(&p)
}

func (a *App) DeleteProduct(id uint) *response.AjaxResult {
	return service.DeleteProduct(id)
}

// ==================== 客户管理 ====================

type CustomerListResult struct {
	List  []model.Customer `json:"list"`
	Total int64            `json:"total"`
}

func (a *App) GetCustomers(page, pageSize int, keyword string) *response.AjaxResult {
	list, total, err := service.GetCustomers(page, pageSize, keyword)
	if err != nil {
		return response.ErrorWithMsg(err.Error())
	}
	if list == nil {
		list = []model.Customer{}
	}
	return response.SuccessWithData(CustomerListResult{List: list, Total: total})
}

func (a *App) GetAllCustomers() *response.AjaxResult {
	list, err := service.GetAllCustomers()
	if err != nil {
		return response.ErrorWithMsg(err.Error())
	}
	if list == nil {
		list = []model.Customer{}
	}
	return response.SuccessWithData(list)
}

func (a *App) CreateCustomer(c model.Customer) *response.AjaxResult {
	return service.CreateCustomer(&c)
}

func (a *App) UpdateCustomer(c model.Customer) *response.AjaxResult {
	return service.UpdateCustomer(&c)
}

func (a *App) DeleteCustomer(id uint) *response.AjaxResult {
	return service.DeleteCustomer(id)
}

// ==================== 供应商管理 ====================

type SupplierListResult struct {
	List  []model.Supplier `json:"list"`
	Total int64            `json:"total"`
}

func (a *App) GetSuppliers(page, pageSize int, keyword string) *response.AjaxResult {
	list, total, err := service.GetSuppliers(page, pageSize, keyword)
	if err != nil {
		return response.ErrorWithMsg(err.Error())
	}
	if list == nil {
		list = []model.Supplier{}
	}
	return response.SuccessWithData(SupplierListResult{List: list, Total: total})
}

func (a *App) GetAllSuppliers() *response.AjaxResult {
	list, err := service.GetAllSuppliers()
	if err != nil {
		return response.ErrorWithMsg(err.Error())
	}
	if list == nil {
		list = []model.Supplier{}
	}
	return response.SuccessWithData(list)
}

func (a *App) CreateSupplier(s model.Supplier) *response.AjaxResult {
	return service.CreateSupplier(&s)
}

func (a *App) UpdateSupplier(s model.Supplier) *response.AjaxResult {
	return service.UpdateSupplier(&s)
}

func (a *App) DeleteSupplier(id uint) *response.AjaxResult {
	return service.DeleteSupplier(id)
}

// ==================== 采购管理 ====================

type PurchaseListResult struct {
	List  []model.PurchaseOrder `json:"list"`
	Total int64                 `json:"total"`
}

func (a *App) GetPurchaseOrders(page, pageSize int, keyword string, status int) *response.AjaxResult {
	list, total, err := service.GetPurchaseOrders(page, pageSize, keyword, status)
	if err != nil {
		return response.ErrorWithMsg(err.Error())
	}
	if list == nil {
		list = []model.PurchaseOrder{}
	}
	return response.SuccessWithData(PurchaseListResult{List: list, Total: total})
}

func (a *App) GetPurchaseOrderByID(id uint) *response.AjaxResult {
	po, err := service.GetPurchaseOrderByID(id)
	if err != nil {
		return response.ErrorWithMsg(err.Error())
	}
	return response.SuccessWithData(po)
}

func (a *App) CreatePurchaseOrder(po model.PurchaseOrder) *response.AjaxResult {
	return service.CreatePurchaseOrder(&po)
}

func (a *App) ConfirmPurchaseOrder(id uint) *response.AjaxResult {
	return service.ConfirmPurchaseOrder(id)
}

func (a *App) CancelPurchaseOrder(id uint) *response.AjaxResult {
	return service.CancelPurchaseOrder(id)
}

// ==================== 销售管理 ====================

type SalesListResult struct {
	List  []model.SalesOrder `json:"list"`
	Total int64              `json:"total"`
}

func (a *App) GetSalesOrders(page, pageSize int, keyword string, status int) *response.AjaxResult {
	list, total, err := service.GetSalesOrders(page, pageSize, keyword, status)
	if err != nil {
		return response.ErrorWithMsg(err.Error())
	}
	if list == nil {
		list = []model.SalesOrder{}
	}
	return response.SuccessWithData(SalesListResult{List: list, Total: total})
}

func (a *App) GetSalesOrderByID(id uint) *response.AjaxResult {
	so, err := service.GetSalesOrderByID(id)
	if err != nil {
		return response.ErrorWithMsg(err.Error())
	}
	return response.SuccessWithData(so)
}

func (a *App) CreateSalesOrder(so model.SalesOrder) *response.AjaxResult {
	return service.CreateSalesOrder(&so)
}

func (a *App) ReturnSalesOrder(id uint) *response.AjaxResult {
	return service.ReturnSalesOrder(id)
}

// ==================== 库存管理 ====================

type InventoryLogResult struct {
	List  []model.InventoryLog `json:"list"`
	Total int64                `json:"total"`
}

func (a *App) GetInventoryLogs(page, pageSize int, productName, logType, startDate, endDate string) *response.AjaxResult {
	list, total, err := service.GetInventoryLogs(page, pageSize, productName, logType, startDate, endDate)
	if err != nil {
		return response.ErrorWithMsg(err.Error())
	}
	if list == nil {
		list = []model.InventoryLog{}
	}
	return response.SuccessWithData(InventoryLogResult{List: list, Total: total})
}

func (a *App) CheckInventory(productID uint, actualQty int, remark string) *response.AjaxResult {
	return service.CheckInventory(productID, actualQty, remark)
}

// ==================== 报表分析 ====================

func (a *App) GetSalesReport(startDate, endDate string) *response.AjaxResult {
	items, err := service.GetSalesReport(startDate, endDate)
	if err != nil {
		return response.ErrorWithMsg(err.Error())
	}
	if items == nil {
		items = []service.SalesReportItem{}
	}
	return response.SuccessWithData(items)
}

func (a *App) GetPurchaseReport(startDate, endDate string) *response.AjaxResult {
	items, err := service.GetPurchaseReport(startDate, endDate)
	if err != nil {
		return response.ErrorWithMsg(err.Error())
	}
	if items == nil {
		items = []service.PurchaseReportItem{}
	}
	return response.SuccessWithData(items)
}

// ==================== 系统设置 ====================

func (a *App) GetStoreName() *response.AjaxResult {
	name, err := service.GetStoreName()
	if err != nil {
		return response.ErrorWithMsg(err.Error())
	}
	return response.SuccessWithData(name)
}

func (a *App) SetStoreName(name string) *response.AjaxResult {
	if err := service.SetStoreName(name); err != nil {
		return response.ErrorWithMsg(err.Error())
	}
	return response.SuccessWithMsg("保存成功")
}

func (a *App) GetSetting(key string) string {
	return service.GetSetting(key)
}

func (a *App) SetSetting(key, value string) *response.AjaxResult {
	if err := service.SetSetting(key, value); err != nil {
		return response.ErrorWithMsg(err.Error())
	}
	return response.SuccessWithMsg("保存成功")
}

func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// ==================== 应用信息 ====================
func (a *App) GetAppVersion() string {
	return service.GetVersion()
}

func (a *App) CheckVersion() *service.VersionInfo {
	return service.CheckVersion()
}

func (a *App) DoUpdate(downloadURL, newVersion string) *response.AjaxResult {
	if err := service.DoUpdate(downloadURL, newVersion); err != nil {
		return response.ErrorWithMsg(err.Error())
	}
	return response.SuccessWithMsg("更新已启动")
}

func (a *App) VerifyLockPassword(password string) *response.AjaxResult {
	pwd := service.GetSetting("secondary_password")
	if pwd == "" {
		pwd = "123456"
	}
	if password != pwd {
		return response.ErrorWithMsg("密码错误")
	}
	return response.SuccessWithMsg("验证通过")
}
