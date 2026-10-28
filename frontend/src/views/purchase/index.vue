<template>
  <div class="page-container">
    <!-- 顶部操作栏 -->
    <el-card shadow="never" style="margin-bottom: 12px;">
      <div style="display: flex; align-items: center; gap: 16px; flex-wrap: wrap;">
        <span style="font-weight: 600; white-space: nowrap;">供应商</span>
        <el-select v-model="supplierId" placeholder="选择供应商" filterable clearable style="width: 200px;">
          <el-option v-for="s in supplierList" :key="s.id" :label="s.name" :value="s.id" />
        </el-select>
        <el-divider direction="vertical" style="height: 24px;" />
        <span style="font-weight: 600; white-space: nowrap;">备注</span>
        <el-input v-model="remark" placeholder="选填" clearable style="width: 200px;" />
        <div style="flex: 1;" />
        <el-statistic class="my-statistic" title="合计金额" :value="totalAmount" prefix="¥" :precision="2" :value-style="{ color: '#409eff', fontWeight: 'bold' }" />
        <el-button type="primary" size="large" @click="submitOrder" :disabled="orderItems.length === 0">
          <el-icon style="margin-right: 4px;"><Plus /></el-icon>提交采购单
        </el-button>
      </div>
    </el-card>

    <el-card shadow="never">
      <el-input v-model="searchKeyword" placeholder="搜索名称/编号/规格" clearable style="width:100%;margin-bottom: 20px;" @input="filterProducts">
        <template #prefix><el-icon><Search /></el-icon></template>
      </el-input>

      <el-table :data="displayList" border stripe max-height="calc(100vh - 350px)" show-summary :summary-method="getSummaries">
        <el-table-column label="编号" min-width="140">
          <template #default="{ row }">{{ row.code }}</template>
        </el-table-column>
        <el-table-column label="商品名称" min-width="120">
          <template #default="{ row }">{{ row.name }}</template>
        </el-table-column>
        <el-table-column label="规格" min-width="80">
          <template #default="{ row }">{{ row.spec || '-' }}</template>
        </el-table-column>
        <el-table-column label="现库存" min-width="70" align="center">
          <template #default="{ row }">{{ row.unit ? row.stockQty + ' ' + row.unit.name : row.stockQty }}</template>
        </el-table-column>
        <el-table-column label="采购数量" width="140" align="center">
          <template #default="{ row }">
            <el-input-number v-model="qtyMap[row.id]" :min="0" :max="99999" controls-position="right" size="small" style="width:100%;" placeholder="采购数量"/>
          </template>
        </el-table-column>
        <el-table-column label="采购单价" width="140" align="right">
          <template #default="{ row }">
            <el-input-number v-model="priceMap[row.id]" :min="0" controls-position="right" size="small" style="width:100%;" />
          </template>
        </el-table-column>
        <el-table-column label="小计" width="120" align="right">
          <template #default="{ row }">
            ¥{{ ((qtyMap[row.id] || 0) * (priceMap[row.id] || 0)).toFixed(2) }}
          </template>
        </el-table-column>
      </el-table>
    </el-card>

  </div>
</template>

<script setup>
import { ref, computed, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { getAllProducts, getAllSuppliers, createPurchaseOrder } from '@/api'
import { ElMessage } from 'element-plus'
import { Search, Plus } from '@element-plus/icons-vue'

const router = useRouter()

const supplierId = ref(null)
const remark = ref('')
const searchKeyword = ref('')
const allProducts = ref([])
const displayList = ref([])
const supplierList = ref([])

// 用 reactive 对象存储每个产品的数量和价格
const qtyMap = reactive({})
const priceMap = reactive({})

const orderItems = computed(() => {
  return allProducts.value
      .filter(p => (qtyMap[p.id] || 0) > 0)
      .map(p => ({
        productId: p.id,
        product: p,
        quantity: qtyMap[p.id],
        price: priceMap[p.id] || p.costPrice || 0,
        amount: qtyMap[p.id] * (priceMap[p.id] || p.costPrice || 0)
      }))
})

const totalAmount = computed(() => {
  return orderItems.value.reduce((sum, item) => sum + item.amount, 0)
})

function filterProducts() {
  const kw = searchKeyword.value.toLowerCase()
  if (!kw) {
    displayList.value = allProducts.value
  } else {
    displayList.value = allProducts.value.filter(p =>
        p.name.toLowerCase().includes(kw) || (p.code && p.code.toLowerCase().includes(kw)) || (p.spec && p.spec.toLowerCase().includes(kw))
    )
  }
}

// 表格底部合计
function getSummaries(param) {
  const { columns, data } = param
  const sums = []
  columns.forEach((column, index) => {
    if (index === 0) {
      sums[index] = '合计'
      return
    }
    if (column.property === '小计' || column.label === '小计') {
      // 计算小计列的总和
      let total = 0
      data.forEach(row => {
        const qty = qtyMap[row.id] || 0
        const price = priceMap[row.id] || 0
        total += qty * price
      })
      sums[index] = '¥' + total.toFixed(2)
    } else {
      sums[index] = ''
    }
  })
  return sums
}

async function submitOrder() {
  if (orderItems.value.length === 0) return
  const data = {
    supplierId: supplierId.value || null,
    remark: remark.value,
    items: orderItems.value.map(item => ({
      productId: item.productId,
      quantity: item.quantity,
      price: item.price,
      amount: item.amount
    }))
  }
  const result = await createPurchaseOrder(data)
  if (result.code === 0) {
    ElMessage.success(result.msg)
    allProducts.value.forEach(p => {
      qtyMap[p.id] = null
      priceMap[p.id] = p.costPrice || 0
    })
    supplierId.value = null
    remark.value = ''
    router.push('/purchase/record')
  } else {
    ElMessage.error(result.msg || '提交失败')
  }
}

getAllProducts().then(res => {
  if (res && res.code === 0) {
    const prods = res.data || []
    allProducts.value = prods
    displayList.value = prods
    prods.forEach(p => {
      qtyMap[p.id] = null
      priceMap[p.id] = p.costPrice || 0
    })
  }
}).catch(e => console.error('获取商品列表失败:', e))

getAllSuppliers().then(res => {
  if (res && res.code === 0) supplierList.value = res.data || []
}).catch(e => console.error('获取供应商列表失败:', e))
</script>

<style scoped>
.page-container {
  padding: 20px;
}

:deep(.el-table__footer-wrapper) {
  box-shadow: 0 -4px 8px rgba(0, 0, 0, 0.08);
  position: relative;
  z-index: 1;
}

:deep(.el-table__footer-wrapper .el-table__cell) {
  font-weight: 700;
  background-color: var(--el-fill-color-light);
}

:deep(.my-statistic .el-statistic__prefix) {
  font-weight: 550;
  font-size: 18px;
}
</style>