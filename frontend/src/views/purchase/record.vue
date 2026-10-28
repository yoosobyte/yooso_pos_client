<template>
  <div class="page-container">
    <el-form :inline="true" :model="queryParams" size="default">
      <el-form-item label="关键词">
        <el-input v-model="queryParams.keyword" placeholder="单号/备注" clearable style="width: 200px;" @keyup.enter="loadData" />
      </el-form-item>
      <el-form-item label="状态">
        <el-select v-model="queryParams.status" placeholder="全部" clearable style="width: 200px;">
          <el-option label="草稿" :value="1" />
          <el-option label="已入库" :value="2" />
          <el-option label="已取消" :value="3" />
        </el-select>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="loadData">查询</el-button>
        <el-button @click="resetQuery">重置</el-button>
      </el-form-item>
    </el-form>

    <el-table :data="tableData" v-loading="loading" border stripe @expand-change="handleExpand">
      <el-table-column type="expand">
        <template #default="{ row }">
          <el-table :data="row.items" size="small" border style="margin: 8px 0;">
            <el-table-column label="分类">
              <template #default="{ row: r }">{{ r.product?.category?.name || '-' }}</template>
            </el-table-column>
            <el-table-column label="商品">
              <template #default="{ row: r }">{{ r.product?.name || '-' }}</template>
            </el-table-column>
            <el-table-column label="规格" align="center">
              <template #default="{ row: r }">{{ r.product?.spec || '-' }}</template>
            </el-table-column>
            <el-table-column prop="quantity" label="数量" align="center" />
            <el-table-column label="单位" align="center">
              <template #default="{ row: r }">{{ r.product?.unit?.name || '-' }}</template>
            </el-table-column>
            <el-table-column label="单价" align="right">
              <template #default="{ row: r }">¥{{ r.price.toFixed(2) }}</template>
            </el-table-column>
            <el-table-column label="金额" align="right">
              <template #default="{ row: r }">¥{{ r.amount.toFixed(2) }}</template>
            </el-table-column>
          </el-table>
        </template>
      </el-table-column>
      <el-table-column prop="orderNo" label="采购单号" width="210" />
      <el-table-column label="供应商">
        <template #default="{ row }">{{ row.supplier?.name || '-' }}</template>
      </el-table-column>
      <el-table-column label="金额" width="120" align="right">
        <template #default="{ row }">¥{{ row.totalAmount.toFixed(2) }}</template>
      </el-table-column>
      <el-table-column label="状态" width="90" align="center">
        <template #default="{ row }">
          <el-tag v-if="row.status === 1" type="warning">草稿</el-tag>
          <el-tag v-else-if="row.status === 2" type="success">已入库</el-tag>
          <el-tag v-else-if="row.status === 3" type="info">已取消</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="创建时间" width="170">
        <template #default="{ row }">{{ parseTime(row.createdAt) }}</template>
      </el-table-column>
      <el-table-column label="备注" width="80">
        <template #default="{ row }">
          <el-tooltip :content="row.remark" placement="top" :disabled="!row.remark" :show-after="300">
            <span class="text-ellipsis">{{ row.remark || '-' }}</span>
          </el-tooltip>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="170" fixed="right" align="center">
        <template #default="{ row }">
          <el-button type="success" size="small" :disabled="row.status !== 1" @click="handleConfirm(row)">
            {{ row.status === 2 ? '已入库' : '确认入库' }}
          </el-button>
          <el-button type="danger" size="small" :disabled="row.status !== 1" @click="handleCancel(row)">
            {{ row.status === 3 ? '已取消' : row.status === 2 ? '无法取消' : '取消' }}
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <pagination
      size="small"
      v-show="total > 0" :total="total"
      v-model:page="queryParams.page" v-model:limit="queryParams.pageSize"
      @pagination="loadData"
    />
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { getPurchaseOrders, confirmPurchaseOrder, cancelPurchaseOrder } from '@/api'
import { parseTime } from '@/utils/ruoyi'


import { ElMessage, ElMessageBox } from 'element-plus'

const loading = ref(false)
const tableData = ref([])
const total = ref(0)
const queryParams = reactive({ page: 1, pageSize: 10, keyword: '', status: '' })

async function loadData() {
  loading.value = true
  try {
    const res = await getPurchaseOrders(queryParams.page, queryParams.pageSize, queryParams.keyword, queryParams.status || 0)
    if (res && res.code === 0) { tableData.value = res.data.list || []; total.value = res.data.total || 0 }
  } finally { loading.value = false }
}

function resetQuery() {
  queryParams.keyword = ''; queryParams.status = ''; queryParams.page = 1
  loadData()
}

async function handleConfirm(row) {
  await ElMessageBox.confirm('确认入库后将更新商品库存，确认？', '提示', { type: 'info' })
  const result = await confirmPurchaseOrder(row.id)
  if (result.code === 0) { ElMessage.success(result.msg); loadData() }
  else { ElMessage.error(result.msg || '操作失败') }
}

async function handleCancel(row) {
  await ElMessageBox.confirm('确定取消该采购单？', '提示', { type: 'warning' })
  const result = await cancelPurchaseOrder(row.id)
  if (result.code === 0) { ElMessage.success(result.msg); loadData() }
  else { ElMessage.error(result.msg || '操作失败') }
}

function handleExpand(row, rows) {
  // 展开时自然显示 items 数据
}


loadData()
</script>
<style scoped>
.page-container{
  padding: 20px;
}
.text-ellipsis {
  display: block;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>