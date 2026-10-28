<template>
  <div class="report-page">
    <el-card class="show-card">
      <template #header>
        <div style="display: flex; justify-content: space-between; align-items: center; flex-wrap: wrap; gap: 8px;">
          <span style="font-weight: bold; white-space: nowrap;"><el-icon :size="18" style="margin-right: 4px; vertical-align: -3px;"><Goods /></el-icon>采购报表</span>
          <div style="display: flex; align-items: center; gap: 8px;">
            <span style="font-size: 14px; font-weight: 600;">日期</span>
            <el-date-picker
              v-model="dateRange" type="daterange" range-separator="至"
              start-placeholder="开始日期" end-placeholder="结束日期"
              value-format="YYYY-MM-DD" style="width: 300px;"
            />
            <span style="font-size: 14px; font-weight: 600;margin-left: 15px;">产品</span>
            <el-select v-model="productFilter" placeholder="全部" filterable clearable style="width: 200px;">
              <el-option v-for="p in productList" :key="p.id" :label="p.name" :value="p.name" />
            </el-select>
            <el-button @click="resetFilter">
              <svg-icon icon-class="clafifter" style="margin-right: 5px;"/> 清除筛选
            </el-button>
          </div>
        </div>
      </template>
      <el-table :data="tableData" border stripe show-summary :summary-method="getSummaries" max-height="calc(100vh - 250px)">
        <el-table-column prop="categoryName" label="分类">
          <template #default="{ row }">{{ row.categoryName || '-' }}</template>
        </el-table-column>
        <el-table-column prop="productName" label="商品" />
        <el-table-column prop="spec" label="规格">
          <template #default="{ row }">{{ row.spec || '-' }}</template>
        </el-table-column>
        <el-table-column prop="totalQty" label="数量" align="center" />
        <el-table-column label="金额" align="right">
          <template #default="{ row }">¥{{ row.totalAmount.toFixed(2) }}</template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { Goods } from '@element-plus/icons-vue'
import { getPurchaseReport, getAllProducts } from '@/api'

const dateRange = ref([])
const productFilter = ref('')
const productList = ref([])
const allData = ref([])
const tableData = computed(() => {
  if (!productFilter.value) return allData.value
  return allData.value.filter(row => row.productName === productFilter.value)
})

function getSummaries(param) {
  const { columns, data } = param
  const sums = []
  const totalAmount = data.reduce((s, r) => s + (r.totalAmount || 0), 0)
  columns.forEach((col, index) => {
    if (index === 0) { sums[index] = '合计'; return }
    if (col.label === '金额') sums[index] = '¥' + totalAmount.toFixed(2)
    else sums[index] = ''
  })
  return sums
}

function resetFilter() {
  dateRange.value = []
  productFilter.value = ''
  loadData()
}

async function loadData() {
  const start = dateRange.value?.length ? dateRange.value[0] : ''
  const end = dateRange.value?.length ? dateRange.value[1] : ''
  const res = await getPurchaseReport(start, end)
  allData.value = (res && res.code === 0) ? (res.data || []) : []
}

watch([dateRange, productFilter], () => loadData())

onMounted(async () => {
  const res = await getAllProducts()
  if (res && res.code === 0) productList.value = res.data || []
})

loadData()
</script>
<style scoped>
.report-page{
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
</style>