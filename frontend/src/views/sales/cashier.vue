<template>
  <div class="cashier-container">
    <!-- 左侧：商品选择 -->
    <div class="cashier-left">
      <el-card shadow="never" style="height: 100%; display: flex; flex-direction: column;">
        <div style="display: flex; gap: 10px; margin-bottom: 12px;">
          <el-input v-model="searchKeyword" placeholder="输入即搜商品名称或编号" clearable
                    style="flex:1;" size="large">
            <template #prefix><el-icon><Search /></el-icon></template>
          </el-input>
        </div>
        <!-- 商品分类标签 -->
        <div class="cat-filter" style="margin-bottom: 12px;">
          <el-button :type="selectedCatId ? '' : 'primary'" size="small" @click="selectedCatId = null; loadProducts()">全部</el-button>
          <el-button v-for="cat in flatCategories" :key="cat.id" :type="selectedCatId === cat.id ? 'primary' : ''"
                     size="small" @click="selectedCatId = cat.id; loadProducts()">{{ cat.name }}</el-button>
        </div>
        <!-- 商品网格 -->
        <div class="product-grid" v-loading="productLoading">
          <div v-for="p in productList" :key="p.id" class="product-card" @click="addToCart(p)" :class="{ 'out-of-stock': p.stockQty <= 0 }">
            <div class="product-type">{{ p.category.name }}</div>
            <div class="product-name">{{ p.name }}</div>
            <div class="product-code">{{ p.spec || p.code }}</div>
            <div class="product-price">¥{{ p.salePrice.toFixed(2) }}</div>
            <div class="product-stock">库存：<span :style="{ color: p.stockQty <= p.minStock && p.minStock > 0 ? 'var(--el-color-danger)' : 'var(--el-color-success)' }">{{ p.stockQty || 0 }} {{ p.unit?.name || '个' }}</span></div>
          </div>
          <el-empty v-if="!productList.length" description="暂无商品" :image-size="60" style="grid-column: 1 / -1; display: flex; justify-content: center; align-items: center;" />
        </div>
      </el-card>
    </div>

    <!-- 右侧：购物车 -->
    <div class="cashier-right">
      <el-card shadow="never" class="cart-card">
        <template #header>
          <div style="display: flex; justify-content: space-between; align-items: center;">
            <div style="display: inline-flex; align-items: center; gap: 4px;">
              <el-icon style="display: flex;align-items: center;justify-content: center;padding-top: 2px;font-size: 17px;"><ShoppingTrolley /></el-icon>
              <span>购物车</span>
            </div>
            <div>
              <el-select v-model="customerId" placeholder="选择客户(不选为散客)" filterable clearable style="width: 220px;" size="default"
                     :filter-method="filterCustomer">
                <el-option v-for="c in filteredCustomerList" :key="c.id" :label="c.phone ? c.name + ' | ' + c.phone : c.name" :value="c.id">
                  <span style="display: flex; justify-content: space-between;">
                    <span>{{ c.name }}</span>
                    <span style="color: var(--el-text-color-secondary); font-size: 12px;">{{ c.phone }}</span>
                  </span>
                </el-option>
              </el-select>
              <el-button size="default" @click="handleClearCart" style="margin-left: 8px;">
                <el-icon><Delete /></el-icon>清空
              </el-button>
            </div>
          </div>
        </template>

        <!-- 购物车列表 -->
        <div class="cart-items">
          <el-table :data="cart" size="small" border v-if="cart.length" max-height="calc(100vh - 390px)">
            <el-table-column label="商品" min-width="90">
              <template #default="{ row }">
                <el-tooltip :content="row.spec || row.code" placement="top" :show-after="400">
                  <div class="cell-name">{{ row.name }}</div>
                </el-tooltip>
              </template>
            </el-table-column>
            <el-table-column label="单价" min-width="52" align="center">
              <template #default="{ row }">¥{{ row.price.toFixed(2) }}</template>
            </el-table-column>
            <el-table-column label="数量" width="115" align="center">
              <template #default="{ row }">
                <el-input-number v-model="row.quantity" :min="1" :max="row.maxStock" size="small" controls-position="right" style="width: 95px" />
              </template>
            </el-table-column>
            <el-table-column label="金额" min-width="70" align="right">
              <template #default="{ row }">¥{{ (row.quantity * row.price).toFixed(2) }}</template>
            </el-table-column>
            <el-table-column label="操作" width="60" align="center">
              <template #default="{ $index }">
                <el-button type="danger" size="small" @click="cart.splice($index, 1)">
                  <el-icon><Delete/></el-icon>
                </el-button>
              </template>
            </el-table-column>
          </el-table>
          <el-empty v-if="!cart.length" description="请选择商品" :image-size="50" />
        </div>

        <!-- 底部结算 -->
        <div class="cart-footer">
          <dl class="bill" v-show="cart.length">
            <div class="bill-row">
              <dt>合计</dt>
              <dd class="bill-amount">¥{{ cartTotalAmount.toFixed(2) }}</dd>
            </div>
            <div class="bill-row">
              <dt>优惠</dt>
              <dd><el-input-number v-model="discount" :min="0" controls-position="right" size="small" style="width: 100px;" placeholder="无" /></dd>
            </div>
            <div class="bill-sep"></div>
            <div class="bill-row bill-total">
              <dt>应收</dt>
              <dd class="bill-receive">¥{{ receiveAmount.toFixed(2) }}</dd>
            </div>
          </dl>
          <el-button :type="cart.length === 0?'info':'primary'" size="large" class="settle-btn"
                     @click="submitOrder" :disabled="cart.length === 0">
            <el-icon><Coin /></el-icon>收款结算(回车↩︎)
          </el-button>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, reactive, onMounted, onBeforeUnmount, watch } from 'vue'
import { getAllProducts, getCategories, getAllCustomers, createSalesOrder } from '@/api'
import { ElMessage, ElMessageBox } from 'element-plus'
import {ShoppingCart, Delete, Search, Coin, Printer, ShoppingTrolley} from '@element-plus/icons-vue'

const searchKeyword = ref('')
const productList = ref([])
const productLoading = ref(false)
const flatCategories = ref([])
const selectedCatId = ref(null)
const customerId = ref(null)
const customerList = ref([])
const customerFilterText = ref('')
const filteredCustomerList = computed(() => {
  if (!customerFilterText.value) return customerList.value
  const kw = customerFilterText.value.toLowerCase()
  return customerList.value.filter(c =>
    c.name.toLowerCase().includes(kw) || (c.phone && c.phone.includes(kw))
  )
})
function filterCustomer(query) {
  customerFilterText.value = query
}
const cart = reactive([])
const discount = ref(null)
const remark = ref('')

const cartTotalQty = computed(() => cart.reduce((s, i) => s + i.quantity, 0))
const cartTotalAmount = computed(() => cart.reduce((s, i) => s + i.quantity * i.price, 0))
const receiveAmount = computed(() => {
  const amt = cartTotalAmount.value - (discount.value || 0)
  return amt < 0 ? 0 : amt
})

async function loadProducts() {
  productLoading.value = true
  try {
    const res = await getAllProducts()
    const all = (res && res.code === 0) ? (res.data || []) : []
    let list = all.filter(p => p.status === 1)
    if (selectedCatId.value) {
      list = list.filter(p => p.categoryId === selectedCatId.value)
    }
    if (searchKeyword.value) {
      const kw = searchKeyword.value.toLowerCase()
      list = list.filter(p => p.name.toLowerCase().includes(kw) || p.code.toLowerCase().includes(kw))
    }
    productList.value = list
  } finally { productLoading.value = false }
}

watch(searchKeyword, () => loadProducts())

function addToCart(p) {
  if (p.stockQty <= 0) {
    ElMessage.warning('该商品已缺货')
    return
  }
  const existing = cart.find(x => x.productId === p.id)
  if (existing) {
    if (existing.quantity >= p.stockQty) {
      ElMessage.warning(`库存不足，当前库存: ${p.stockQty}`)
      return
    }
    existing.quantity++
  } else {
    cart.push({
      productId: p.id,
      name: p.name,
      code: p.code,
      spec: p.spec,
      price: p.salePrice,
      quantity: 1,
      maxStock: p.stockQty
    })
  }
}

function handleClearCart() {
  cart.length = 0
  customerId.value = null
  discount.value = null
}

async function submitOrder() {
  if (cart.length === 0) return
  const data = {
    customerId: customerId.value || null,
    discount: discount.value || 0,
    remark: remark.value || '',
    items: cart.map(item => ({
      productId: item.productId,
      quantity: item.quantity,
      price: item.price,
      amount: item.quantity * item.price
    }))
  }
  const result = await createSalesOrder(data)
  if (result.code === 0) {
    ElMessage.success(result.msg)
    handleClearCart()
    loadProducts()
  } else {
    ElMessage.error(result.msg || '操作失败')
    loadProducts()
  }
}

async function loadCategories() {
  const res = await getCategories()
  if (res && res.code === 0) {
    function flatten(list) {
      let result = []
      for (const item of list) {
        result.push({ id: item.id, name: item.name })
        if (item.children) result = result.concat(flatten(item.children))
      }
      return result
    }
    flatCategories.value = flatten(res.data)
  }
}

loadProducts()
loadCategories()
getAllCustomers().then((res) => {
  if (res && res.code === 0) customerList.value = res.data
})

function onKeyDown(e) {
  if (e.key === 'Enter') {
    e.preventDefault()
    if (cart.length) submitOrder()
  }
}

onMounted(() => window.addEventListener('keydown', onKeyDown))
onBeforeUnmount(() => window.removeEventListener('keydown', onKeyDown))
</script>

<style scoped>
.cashier-container {
  display: flex;
  gap: 12px;
  height: calc(100vh - 120px);
  padding: 20px;
  background-color: var(--el-bg-color-page);
}

.cashier-left {
  flex: 1;
  overflow: hidden;
}

.cashier-right {
  width: 520px;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
}

.cart-card {
  height: 100%;
  display: flex;
  flex-direction: column;
  background-color: var(--el-bg-color);
  border-color: var(--el-border-color-light);
}

.cat-filter {
  display: flex;
  gap: 4px 4px;
  flex-wrap: wrap;
}

.cat-filter .el-button {
  margin-left: 0;
}

.cart-items {
  flex: 1;
  overflow: visible;
  margin-bottom: 4px;
}

.product-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(155px, 1fr));
  gap: 12px;
  overflow-y: auto;
  flex: 1;
  align-content: start;
}

.product-card {
  position: relative;
  border: 1px solid var(--el-border-color-light);
  border-radius: 10px;
  padding: 16px;
  cursor: pointer;
  transition: all 0.2s;
  background-color: var(--el-bg-color);
  color: var(--el-text-color-primary);
  user-select: none;
}

.product-card:hover {
  border-color: var(--el-color-primary);
  box-shadow: 0 2px 8px rgba(64, 158, 255, 0.2);
}

.product-card.out-of-stock {
  opacity: 0.5;
  background-color: var(--el-fill-color-light);
}

.product-type{
  position:absolute;
  top: 0;
  right: 0;
  background-color: #fff3e0;
  height: 20px;
  padding: 0 7px;
  border-top-right-radius: 10px;
  border-bottom-left-radius: 10px;
  color: #a5a5a5;
  font-size: 10px;
  display: flex;
  justify-content: center;
  align-items: center;
}

html.dark .product-type {
  background-color: #FFFFFF;
  color: #333333;
}

.product-name {
  font-weight: 600;
  font-size: 15px;
  margin-bottom: 8px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: var(--el-text-color-primary);
}

.product-code {
  font-size: 12px;
  color: var(--el-text-color-secondary);
  margin-bottom: 6px;
}

.product-price {
  color: var(--el-color-danger);
  font-weight: 700;
  font-size: 18px;
  margin-bottom: 4px;
}

.product-stock {
  font-size: 11px;
  color: var(--el-text-color-secondary);
  margin-top: 6px;
}

.cell-name {
  font-weight: 600;
  font-size: 13px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.cart-card :deep(.el-card__body) {
  padding: 6px 8px !important;
}

.cart-card :deep(.el-card__header) {
  padding: 6px 10px !important;
}

.cart-footer {
  margin-top: auto;
}

.bill {
  margin: 8px 0 8px 0;
  padding: 0;
  border: 1px solid var(--el-border-color);
}

.bill-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 8px;
  min-height: 36px;
}

.bill-row dt {
  font-size: 13px;
  color: var(--el-text-color-secondary);
  margin: 0;
}

.bill-row dd {
  font-size: 15px;
  font-weight: 600;
  color: var(--el-text-color-primary);
  margin: 0;
}

.bill-amount {
  font-size: 17px !important;
  font-weight: 600 !important;
  color: var(--el-color-primary) !important;
}

.bill-sep {
  display: none;
}

.bill-total {
  border-top: 1px solid var(--el-border-color);
}

.bill-total dt {
  font-size: 13px !important;
  font-weight: 600 !important;
  color: var(--el-text-color-primary) !important;
}

.bill-receive {
  font-size: 17px !important;
  font-weight: 600 !important;
  color: var(--el-color-danger) !important;
}

.settle-btn {
  width: 100%;
  height: 42px;
  font-size: 16px;
  letter-spacing: 2px;
}

/* 暗黑模式额外适配 */
:deep(.el-input__inner) {
  background-color: var(--el-bg-color);
  border-color: var(--el-border-color);
  color: var(--el-text-color-primary);
}

:deep(.el-input__inner:hover) {
  border-color: var(--el-border-color-hover);
}

:deep(.el-input__inner:focus) {
  border-color: var(--el-color-primary);
}

:deep(.el-tag) {
  background-color: var(--el-fill-color-light);
  border-color: var(--el-border-color-light);
  color: var(--el-text-color-primary);
}

:deep(.el-tag--primary) {
  background-color: var(--el-color-primary-light-9);
  border-color: var(--el-color-primary);
  color: var(--el-color-primary);
}

:deep(.el-empty__description) {
  color: var(--el-text-color-secondary);
}

:deep(.el-select .el-input__inner) {
  background-color: var(--el-bg-color);
}

:deep(.el-select-dropdown__item) {
  color: var(--el-text-color-primary);
}

:deep(.el-select-dropdown__item.hover) {
  background-color: var(--el-fill-color-light);
}

:deep(.el-select-dropdown__item.selected) {
  color: var(--el-color-primary);
}

:deep(.el-button--primary) {
  background-color: var(--el-color-primary);
  border-color: var(--el-color-primary);
}

:deep(.el-button--primary.is-disabled) {
  background-color: var(--el-color-primary-light-5);
  border-color: var(--el-color-primary-light-5);
}

:deep(.el-input-number .el-input__inner) {
  background-color: var(--el-bg-color);
  color: var(--el-text-color-primary);
}

:deep(.el-input-number.is-controls-right .el-input-number__decrease),
:deep(.el-input-number.is-controls-right .el-input-number__increase) {
  background-color: var(--el-bg-color);
  border-color: var(--el-border-color);
  color: var(--el-text-color-primary);
}

:deep(.el-input-number.is-controls-right .el-input-number__decrease:hover),
:deep(.el-input-number.is-controls-right .el-input-number__increase:hover) {
  background-color: var(--el-fill-color-light);
}

/* 滚动条暗黑适配 */
.cart-items::-webkit-scrollbar,
.product-grid::-webkit-scrollbar {
  width: 4px;
}

.cart-items::-webkit-scrollbar-track,
.product-grid::-webkit-scrollbar-track {
  background: var(--el-fill-color-light);
  border-radius: 4px;
}

.cart-items::-webkit-scrollbar-thumb,
.product-grid::-webkit-scrollbar-thumb {
  background: var(--el-border-color);
  border-radius: 4px;
}

.cart-items::-webkit-scrollbar-thumb:hover,
.product-grid::-webkit-scrollbar-thumb:hover {
  background: var(--el-border-color-hover);
}

/* 响应式适配 */
@media (max-width: 768px) {
  .cashier-container {
    flex-direction: column;
    height: auto;
    padding: 12px;
  }
  .cashier-right {
    width: 100%;
  }
  .product-grid {
    grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
  }
}
</style>
