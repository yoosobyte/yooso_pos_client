<template>
  <div class="page-container">
    <el-form :inline="true" :model="queryParams" size="default">
      <el-form-item label="商品名称">
        <el-input v-model="queryParams.productName" placeholder="商品名称" clearable style="width: 200px;" />
      </el-form-item>
      <el-form-item label="商品类型">
        <el-select v-model="queryParams.logType" placeholder="全部类型" clearable style="width: 200px;">
          <el-option label="采购入库" value="purchase_in" />
          <el-option label="销售出库" value="sales_out" />
          <el-option label="销售退货" value="return_in" />
          <el-option label="盘点调整" value="check" />
        </el-select>
      </el-form-item>
      <el-form-item label="筛选日期">
        <el-date-picker
          v-model="dateRange" type="daterange" range-separator="至"
          start-placeholder="开始日期" end-placeholder="结束日期"
          value-format="YYYY-MM-DD" style="width: 300px;"
        />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="loadData">查询</el-button>
        <el-button @click="resetQuery">重置</el-button>
      </el-form-item>
    </el-form>

    <el-table :data="tableData" v-loading="loading" border stripe>
      <el-table-column label="时间" width="170">
        <template #default="{ row }">{{ parseTime(row.createdAt) }}</template>
      </el-table-column>
      <el-table-column label="商品" min-width="140">
        <template #default="{ row }">{{ row.product?.name || '-' }}</template>
      </el-table-column>
      <el-table-column label="类型" width="100" align="center">
        <template #default="{ row }">
          <el-tag v-if="row.type === 'purchase_in'" type="primary">采购入库</el-tag>
          <el-tag v-else-if="row.type === 'sales_out'" type="danger">销售出库</el-tag>
          <el-tag v-else-if="row.type === 'return_in'" type="success">销售退货</el-tag>
          <el-tag v-else-if="row.type === 'check'" type="warning">盘点调整</el-tag>
          <el-tag v-else>{{ row.type }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="变动数量" width="90" align="center">
        <template #default="{ row }">
          <span :style="{ color: row.qtyChange > 0 ? '#67c23a' : '#f56c6c' }">
            {{ row.qtyChange > 0 ? '+' : '' }}{{ row.qtyChange }}
          </span>
        </template>
      </el-table-column>
      <el-table-column prop="qtyAfter" label="变动后库存" width="100" align="center" />
      <el-table-column prop="orderNo" label="关联单号" width="210" />
      <el-table-column prop="remark" label="备注" min-width="120" />
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
import { getInventoryLogs } from '@/api'
import { parseTime } from '@/utils/ruoyi'



const loading = ref(false)
const tableData = ref([])
const total = ref(0)
const dateRange = ref([])
const queryParams = reactive({ page: 1, pageSize: 10, productName: '', logType: '' })

async function loadData() {
  loading.value = true
  try {
    const start = dateRange.value?.length ? dateRange.value[0] : ''
    const end = dateRange.value?.length ? dateRange.value[1] : ''
    const res = await getInventoryLogs(queryParams.page, queryParams.pageSize, queryParams.productName, queryParams.logType, start, end)
    if (res && res.code === 0) { tableData.value = res.data.list || []; total.value = res.data.total || 0 }
  } finally { loading.value = false }
}

function resetQuery() {
  queryParams.productName = ''; queryParams.logType = ''; dateRange.value = []; queryParams.page = 1
  loadData()
}


loadData()
</script>
<style scoped>
.page-container{
  padding: 20px;
}
</style>