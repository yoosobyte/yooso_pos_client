// Wails v3 绑定适配层 — 从自动生成的 bindings 导入，导出 camelCase 别名
import { App } from "../../bindings/yoosoPosV3"

// 仪表盘
export const getDashboard = () => App.GetDashboard()
export const getLowStockProducts = () => App.GetLowStockProducts()
// 分类
export const getCategories = () => App.GetCategories()
export const getCategoryByID = (id) => App.GetCategoryByID(id)
export const createCategory = (d) => App.CreateCategory(d)
export const updateCategory = (d) => App.UpdateCategory(d)
export const deleteCategory = (id) => App.DeleteCategory(id)
export const getMaxCategorySort = (parentId) => App.GetMaxCategorySort(parentId || 0)
// 单位
export const getUnits = () => App.GetUnits()
export const createUnit = (d) => App.CreateUnit(d)
export const updateUnit = (d) => App.UpdateUnit(d)
export const deleteUnit = (id) => App.DeleteUnit(id)
// 商品
export const getProducts = (p, s, kw, cid) => App.GetProducts(p, s, kw, cid)
export const getAllProducts = () => App.GetAllProducts()
export const getProductByID = (id) => App.GetProductByID(id)
export const createProduct = (d) => App.CreateProduct(d)
export const updateProduct = (d) => App.UpdateProduct(d)
export const deleteProduct = (id) => App.DeleteProduct(id)
// 客户
export const getCustomers = (p, s, kw) => App.GetCustomers(p, s, kw)
export const getAllCustomers = () => App.GetAllCustomers()
export const createCustomer = (d) => App.CreateCustomer(d)
export const updateCustomer = (d) => App.UpdateCustomer(d)
export const deleteCustomer = (id) => App.DeleteCustomer(id)
// 供应商
export const getSuppliers = (p, s, kw) => App.GetSuppliers(p, s, kw)
export const getAllSuppliers = () => App.GetAllSuppliers()
export const createSupplier = (d) => App.CreateSupplier(d)
export const updateSupplier = (d) => App.UpdateSupplier(d)
export const deleteSupplier = (id) => App.DeleteSupplier(id)
// 采购
export const getPurchaseOrders = (p, s, kw, st) => App.GetPurchaseOrders(p, s, kw, st)
export const getPurchaseOrderByID = (id) => App.GetPurchaseOrderByID(id)
export const createPurchaseOrder = (d) => App.CreatePurchaseOrder(d)
export const confirmPurchaseOrder = (id) => App.ConfirmPurchaseOrder(id)
export const cancelPurchaseOrder = (id) => App.CancelPurchaseOrder(id)
// 销售
export const getSalesOrders = (p, s, kw, st) => App.GetSalesOrders(p, s, kw, st)
export const getSalesOrderByID = (id) => App.GetSalesOrderByID(id)
export const createSalesOrder = (d) => App.CreateSalesOrder(d)
export const returnSalesOrder = (id) => App.ReturnSalesOrder(id)
// 库存
export const getInventoryLogs = (p, s, pn, lt, sd, ed) => App.GetInventoryLogs(p, s, pn, lt, sd, ed)
export const checkInventory = (pid, qty, rk) => App.CheckInventory(pid, qty, rk)
// 报表
export const getSalesReport = (sd, ed) => App.GetSalesReport(sd, ed)
export const getPurchaseReport = (sd, ed) => App.GetPurchaseReport(sd, ed)
// 应用信息
export const getAppVersion = () => App.GetAppVersion()
export const checkVersion = () => App.CheckVersion()
export const doUpdate = (url, version) => App.DoUpdate(url, version)
export const verifyLockPassword = (pwd) => App.VerifyLockPassword(pwd)
// 系统设置
export const getStoreName = () => App.GetStoreName()
export const setStoreName = (name) => App.SetStoreName(name)
export const getSetting = (key) => App.GetSetting(key)
export const setSetting = (key, value) => App.SetSetting(key, value)
