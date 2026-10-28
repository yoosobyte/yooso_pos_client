<template>
  <div class="page-container">
    <el-form :inline="true" :model="queryParams" size="default">
      <el-form-item label="关键词">
        <el-input v-model="queryParams.keyword" placeholder="单号/客户" clearable style="width: 200px;" @keyup.enter="loadData" />
      </el-form-item>
      <el-form-item label="状态">
        <el-select v-model="queryParams.status" placeholder="全部" clearable style="width: 200px;">
          <el-option label="已完成" :value="1" />
          <el-option label="已退货" :value="2" />
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
      <el-table-column prop="orderNo" label="销售单号" width="210" />
      <el-table-column label="客户">
        <template #default="{ row }">{{ row.customer?.name || '散客' }}</template>
      </el-table-column>
      <el-table-column label="金额" width="100" align="right">
        <template #default="{ row }">¥{{ row.totalAmount.toFixed(2) }}</template>
      </el-table-column>
      <el-table-column label="优惠" width="80" align="right">
        <template #default="{ row }">¥{{ row.discount.toFixed(2) }}</template>
      </el-table-column>
      <el-table-column label="状态" width="90" align="center">
        <template #default="{ row }">
          <el-tag v-if="row.status === 1" type="success">已完成</el-tag>
          <el-tag v-else-if="row.status === 2" type="warning">已退货</el-tag>
          <el-tag v-else-if="row.status === 3" type="info">已取消</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="创建时间" width="170">
        <template #default="{ row }">{{ parseTime(row.createdAt) }}</template>
      </el-table-column>
      <el-table-column label="操作" width="120" fixed="right" align="center">
        <template #default="{ row }">
          <el-button type="danger" size="small" :disabled="row.status !== 1" @click="handleReturn(row)">
            {{ row.status === 2 ? '已退货' : row.status === 3 ? '已取消' : '退货' }}
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
import { getSalesOrders, returnSalesOrder } from '@/api'
import { parseTime } from '@/utils/ruoyi'


import { ElMessage, ElMessageBox } from 'element-plus'

const loading = ref(false)
const tableData = ref([])
const total = ref(0)
const queryParams = reactive({ page: 1, pageSize: 10, keyword: '', status: '' })

async function loadData() {
  loading.value = true
  try {
    const res = await getSalesOrders(queryParams.page, queryParams.pageSize, queryParams.keyword, queryParams.status || 0)
    if (res && res.code === 0) { tableData.value = res.data.list || []; total.value = res.data.total || 0 }
  } finally { loading.value = false }
}

function resetQuery() {
  queryParams.keyword = ''; queryParams.status = ''; queryParams.page = 1
  loadData()
}

async function handleReturn(row) {
  await ElMessageBox.confirm('退货后将恢复商品库存，确认？', '提示', { type: 'warning' })
  const result = await returnSalesOrder(row.id)
  if (result.code === 0) {
    ElMessage.success(result.msg)
    loadData()
  } else {
    ElMessage.error(result.msg || '操作失败')
  }
}

function handleExpand() {}


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
