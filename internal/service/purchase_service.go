package service

import (
	"yoosoPosV3/internal/config"
	"yoosoPosV3/internal/model"
	"yoosoPosV3/internal/response"
)

// ---------- 客户 ----------

func GetCustomers(page, pageSize int, keyword string) ([]model.Customer, int64, error) {
	var list []model.Customer
	var total int64
	query := config.DB.Model(&model.Customer{})
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("name LIKE ? OR phone LIKE ? OR address LIKE ?", like, like, like)
	}
	query.Count(&total)
	err := query.Order("id desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}

func GetAllCustomers() ([]model.Customer, error) {
	var list []model.Customer
	err := config.DB.Find(&list).Error
	return list, err
}

func CreateCustomer(c *model.Customer) *response.AjaxResult {
	if err := config.DB.Create(c).Error; err != nil {
		return response.ErrorWithMsg("创建失败：" + err.Error())
	}
	return response.SuccessWithData(c)
}

func UpdateCustomer(c *model.Customer) *response.AjaxResult {
	if err := config.DB.Save(c).Error; err != nil {
		return response.ErrorWithMsg("更新失败：" + err.Error())
	}
	return response.SuccessWithMsgData("更新成功", c)
}

func DeleteCustomer(id uint) *response.AjaxResult {
	if err := config.DB.Delete(&model.Customer{}, id).Error; err != nil {
		return response.ErrorWithMsg("删除失败：" + err.Error())
	}
	return response.SuccessWithMsg("删除成功")
}

// ---------- 供应商 ----------

func GetSuppliers(page, pageSize int, keyword string) ([]model.Supplier, int64, error) {
	var list []model.Supplier
	var total int64
	query := config.DB.Model(&model.Supplier{})
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("name LIKE ? OR phone LIKE ? OR contact LIKE ? OR address LIKE ?", like, like, like, like)
	}
	query.Count(&total)
	err := query.Order("id desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}

func GetAllSuppliers() ([]model.Supplier, error) {
	var list []model.Supplier
	err := config.DB.Find(&list).Error
	return list, err
}

func CreateSupplier(s *model.Supplier) *response.AjaxResult {
	if err := config.DB.Create(s).Error; err != nil {
		return response.ErrorWithMsg("创建失败：" + err.Error())
	}
	return response.SuccessWithData(s)
}

func UpdateSupplier(s *model.Supplier) *response.AjaxResult {
	if err := config.DB.Save(s).Error; err != nil {
		return response.ErrorWithMsg("更新失败：" + err.Error())
	}
	return response.SuccessWithMsgData("更新成功", s)
}

func DeleteSupplier(id uint) *response.AjaxResult {
	if err := config.DB.Delete(&model.Supplier{}, id).Error; err != nil {
		return response.ErrorWithMsg("删除失败：" + err.Error())
	}
	return response.SuccessWithMsg("删除成功")
}
