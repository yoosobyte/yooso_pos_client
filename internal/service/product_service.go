package service

import (
	"fmt"
	"math/rand"
	"time"
	"yoosoPosV3/internal/response"

	"yoosoPosV3/internal/config"
	"yoosoPosV3/internal/model"
)

// ---------- 商品分类 ----------

func GetCategories() ([]model.Category, error) {
	var categories []model.Category
	err := config.DB.Order("sort asc").Find(&categories).Error
	return buildCategoryTree(categories), err
}

func GetCategoryByID(id uint) (*model.Category, error) {
	var c model.Category
	err := config.DB.First(&c, id).Error
	return &c, err
}

func CreateCategory(c *model.Category) *response.AjaxResult {
	if err := config.DB.Create(c).Error; err != nil {
		return response.ErrorWithMsg("创建失败：" + err.Error())
	}
	return response.SuccessWithData(c)
}

func UpdateCategory(c *model.Category) *response.AjaxResult {
	if err := config.DB.Save(c).Error; err != nil {
		return response.ErrorWithMsg("更新失败：" + err.Error())
	}
	return response.SuccessWithMsgData("更新成功", c)
}

func GetMaxCategorySort(parentID uint) int {
	var maxSort int
	query := config.DB.Model(&model.Category{})
	if parentID > 0 {
		query = query.Where("parent_id = ?", parentID)
	} else {
		query = query.Where("parent_id IS NULL")
	}
	query.Select("COALESCE(MAX(sort), 0)").Scan(&maxSort)
	return maxSort + 1
}

func DeleteCategory(id uint) *response.AjaxResult {
	var count int64
	config.DB.Model(&model.Product{}).Where("category_id = ?", id).Count(&count)
	if count > 0 {
		return response.ErrorWithMsg("该分类下有关联商品，无法删除")
	}
	if err := config.DB.Delete(&model.Category{}, id).Error; err != nil {
		return response.ErrorWithMsg("删除失败：" + err.Error())
	}
	return response.SuccessWithMsg("删除成功")
}

func buildCategoryTree(list []model.Category) []model.Category {
	tree := make([]model.Category, 0)
	for i := range list {
		if list[i].ParentID == nil {
			children := findChildren(list, list[i].ID)
			list[i].Children = children
			tree = append(tree, list[i])
		}
	}
	return tree
}

func findChildren(list []model.Category, parentID uint) []*model.Category {
	children := make([]*model.Category, 0)
	for i := range list {
		if list[i].ParentID != nil && *list[i].ParentID == parentID {
			sub := findChildren(list, list[i].ID)
			list[i].Children = sub
			children = append(children, &list[i])
		}
	}
	return children
}

// ---------- 商品单位 ----------

func GetUnits() ([]model.Unit, error) {
	var units []model.Unit
	err := config.DB.Find(&units).Error
	return units, err
}

func CreateUnit(u *model.Unit) *response.AjaxResult {
	// 检查是否存在未删除的同名单位
	var count int64
	config.DB.Model(&model.Unit{}).Where("name = ?", u.Name).Count(&count)
	if count > 0 {
		return response.ErrorWithCode(response.ErrCodeUnitNameExists, fmt.Sprintf("单位名称[%s]已存在", u.Name))
	}
	// 如果存在已逻辑删除的同名单位，先物理删除释放唯一约束
	config.DB.Unscoped().Where("name = ?", u.Name).Delete(&model.Unit{})

	if err := config.DB.Create(u).Error; err != nil {
		return response.ErrorWithMsg("创建失败：" + err.Error())
	}
	return response.SuccessWithData(u)
}

func UpdateUnit(u *model.Unit) *response.AjaxResult {
	// 检查是否存在未删除的同名单位（排除自身）
	var count int64
	config.DB.Model(&model.Unit{}).Where("name = ? AND id != ?", u.Name, u.ID).Count(&count)
	if count > 0 {
		return response.ErrorWithCode(response.ErrCodeUnitNameExists, fmt.Sprintf("单位名称[%s]已存在", u.Name))
	}

	// 如果存在已逻辑删除的同名单位（排除自身），先物理删除释放唯一约束
	config.DB.Unscoped().Where("name = ? AND id != ?", u.Name, u.ID).Delete(&model.Unit{})

	// 保存更新
	if err := config.DB.Save(u).Error; err != nil {
		return response.ErrorWithMsg("更新失败：" + err.Error())
	}

	return response.SuccessWithMsgData("更新成功", u)
}

func DeleteUnit(id uint) *response.AjaxResult {
	var count int64
	config.DB.Model(&model.Product{}).Where("unit_id = ?", id).Count(&count)
	if count > 0 {
		return response.ErrorWithMsg("该单位下有关联商品，无法删除")
	}
	if err := config.DB.Delete(&model.Unit{}, id).Error; err != nil {
		return response.ErrorWithMsg("删除失败：" + err.Error())
	}
	return response.SuccessWithMsg("删除成功")
}

// ---------- 商品 ----------

func GetProducts(page, pageSize int, keyword string, categoryID uint) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64
	query := config.DB.Model(&model.Product{}).Preload("Category").Preload("Unit")
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("name LIKE ? OR code LIKE ? OR spec LIKE ?", like, like, like)
	}
	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}
	query.Count(&total)
	err := query.Order("id desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&products).Error
	return products, total, err
}

func GetAllProducts() ([]model.Product, error) {
	var products []model.Product
	err := config.DB.Preload("Category").Preload("Unit").Order("id desc").Find(&products).Error
	return products, err
}

func GetProductByID(id uint) (*model.Product, error) {
	var p model.Product
	err := config.DB.Preload("Category").Preload("Unit").First(&p, id).Error
	return &p, err
}

func CreateProduct(p *model.Product) *response.AjaxResult {
	if p.Code == "" {
		p.Code = fmt.Sprintf("SP%s%04d%03d", time.Now().Format("20060102150405"), time.Now().UnixMilli()%10000, rand.Intn(1000))
	}
	if err := config.DB.Create(p).Error; err != nil {
		return response.ErrorWithMsg("创建失败：" + err.Error())
	}
	return response.SuccessWithData(p)
}

func UpdateProduct(p *model.Product) *response.AjaxResult {
	if p.Code == "" {
		p.Code = fmt.Sprintf("SP%s%04d%03d", time.Now().Format("20060102150405"), time.Now().UnixMilli()%10000, rand.Intn(1000))
	}
	if err := config.DB.Model(&model.Product{}).Where("id = ?", p.ID).Select(
		"name", "code", "spec", "category_id", "unit_id",
		"sale_price", "cost_price", "min_stock", "max_stock",
		"status", "remark",
	).Updates(p).Error; err != nil {
		return response.ErrorWithMsg("更新失败：" + err.Error())
	}
	return response.SuccessWithMsgData("更新成功", p)
}

func DeleteProduct(id uint) *response.AjaxResult {
	if err := config.DB.Delete(&model.Product{}, id).Error; err != nil {
		return response.ErrorWithMsg("删除失败：" + err.Error())
	}
	return response.SuccessWithMsg("删除成功")
}

func UpdateProductStock(productID uint, qtyChange int) error {
	return config.DB.Model(&model.Product{}).Where("id = ?", productID).
		UpdateColumn("stock_qty", config.DB.Raw("stock_qty + ?", qtyChange)).Error
}
