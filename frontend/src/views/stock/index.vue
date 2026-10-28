<template>
  <div class="page-container">
    <el-form :inline="true" :model="queryParams" size="default">
      <el-form-item label="关键词">
        <el-input v-model="queryParams.keyword" placeholder="名称/编号/规格" clearable style="width: 200px;" @keyup.enter="loadData" />
      </el-form-item>
      <el-form-item label="分类">
        <el-tree-select
          v-model="queryParams.categoryId" :data="categories"
          :props="{ label: 'name', value: 'id', children: 'children' }"
          placeholder="全部分类" clearable check-strictly style="width: 200px;"
        />
      </el-form-item>
      <el-form-item>
        <el-checkbox v-model="queryParams.onlyLowStock" @change="loadData">仅看低库存</el-checkbox>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="loadData">查询</el-button>
        <el-button @click="resetQuery">重置</el-button>
        <el-button type="warning" @click="openBatchCheck">批量盘库</el-button>
      </el-form-item>
    </el-form>

    <el-table :data="tableData" v-loading="loading" border stripe>
      <el-table-column label="编号" align="center">
        <template #default="{ row }">{{ row.code || '-' }}</template>
      </el-table-column>
      <el-table-column label="商品名称" align="center">
        <template #default="{ row }">{{ row.name || '-' }}</template>
      </el-table-column>
      <el-table-column label="分类" align="center">
        <template #default="{ row }">{{ row.category?.name || '-' }}</template>
      </el-table-column>
      <el-table-column label="单位" width="70" align="center">
        <template #default="{ row }">{{ row.unit?.name || '-' }}</template>
      </el-table-column>
      <el-table-column label="规格" align="center">
        <template #default="{ row }">{{ row.spec || '-' }}</template>
      </el-table-column>
      <el-table-column label="当前库存" width="90" align="center">
        <template #default="{ row }">
      <span :class="{ 'text-danger': row.stockQty <= row.minStock && row.minStock > 0 }" style="font-weight:700;font-size:16px;">
        {{ row.stockQty ?? '-' }}
      </span>
        </template>
      </el-table-column>
      <el-table-column label="最低库存" width="80" align="center">
        <template #default="{ row }">{{ row.minStock ?? '-' }}</template>
      </el-table-column>
      <el-table-column label="状态" width="95" align="center">
        <template #default="{ row }">
          <el-tag v-if="row.stockQty <= 0 && row.minStock > 0" type="danger">缺货</el-tag>
          <el-tag v-else-if="row.stockQty <= row.minStock && row.minStock > 0" type="warning">偏低</el-tag>
          <el-tag v-else type="success">库存充足</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="130" fixed="right" align="center">
        <template #default="{ row }">
          <div style="display: flex; justify-content: center;">
            <el-button type="primary" size="small" @click="handleCheck(row)">盘点</el-button>
            <el-button type="danger" size="small" @click="handleLog(row)">记录</el-button>
          </div>
        </template>
      </el-table-column>
    </el-table>

    <pagination
      size="small"
      v-show="total > 0" :total="total"
      v-model:page="queryParams.page" v-model:limit="queryParams.pageSize"
      @pagination="loadData"
    />

    <!-- 盘点弹窗 -->
    <el-dialog title="库存盘点" v-model="checkDialogVisible" width="400px" draggable>
      <el-form ref="checkFormRef" :model="checkForm" :rules="checkRules" label-position="top">
        <el-form-item label="商品名称" prop="name">
          <el-input v-model="checkProduct.name" style="width:100%" disabled/>
        </el-form-item>
        <el-form-item label="系统库存" prop="stockQty">
          <el-input-number v-model="checkProduct.stockQty" :min="0" controls-position="right" style="width:100%" disabled/>
        </el-form-item>
        <el-form-item label="实际数量" prop="actualQty">
          <el-input-number v-model="checkForm.actualQty" :min="0" controls-position="right" style="width:100%" placeholder="填写真实数量" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="checkForm.remark" type="textarea" :rows="2" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="checkDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitCheck">确认盘点</el-button>
      </template>
    </el-dialog>

    <!-- 批量盘库弹窗 -->
    <el-dialog title="批量盘库" v-model="batchCheckVisible" width="950px" draggable :close-on-click-modal="false">
      <el-table :data="tableData" border stripe max-height="450">
        <el-table-column label="编号" width="150">
          <template #default="{ row }">{{ row.code || '-' }}</template>
        </el-table-column>
        <el-table-column label="商品名称" min-width="120">
          <template #default="{ row }">{{ row.name }}</template>
        </el-table-column>
        <el-table-column label="规格" min-width="120">
          <template #default="{ row }">{{ row.spec || '-' }}</template>
        </el-table-column>
        <el-table-column label="当前库存" width="85" align="center">
          <template #default="{ row }">{{ row.stockQty }}</template>
        </el-table-column>
        <el-table-column label="实际数量" width="140" align="center">
          <template #default="{ row }">
            <el-input-number v-model="batchQtyMap[row.id]" :min="0" controls-position="right" size="small" style="width:100%;" placeholder="填写真实数量" />
          </template>
        </el-table-column>
        <el-table-column label="差异" width="70" align="center">
          <template #default="{ row }">
            <span v-if="batchQtyMap[row.id] != null" :style="{ color: batchQtyMap[row.id] - row.stockQty >= 0 ? '#67c23a' : '#f56c6c', fontWeight: 'bold' }">
              {{ batchQtyMap[row.id] - row.stockQty >= 0 ? '+' : '' }}{{ batchQtyMap[row.id] - row.stockQty }}
            </span>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column label="盘点备注" min-width="80">
          <template #default="{ row }">
            <el-input v-model="batchRemarkMap[row.id]" size="small" placeholder="选填" />
          </template>
        </el-table-column>
      </el-table>
      <template #footer>
        <el-button @click="batchCheckVisible = false">取消</el-button>
        <el-button type="primary" @click="submitBatchCheck" :disabled="batchCount === 0" :loading="batchSubmitting">
          确认盘点（{{ batchCount }} 项）
        </el-button>
      </template>
    </el-dialog>

    <!-- 库存记录弹窗 -->
    <el-dialog :title="`出入库明细 - ${logProduct?.name || ''}`" v-model="logDialogVisible" width="850px" draggable>
      <el-table :data="logData" size="small" border max-height="350">
        <el-table-column label="时间" width="160">
          <template #default="{ row }">{{ parseTime(row.createdAt) }}</template>
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
        <el-table-column prop="qtyAfter" label="变动后" width="70" align="center" />
        <el-table-column prop="orderNo" label="单号" width="180" />
        <el-table-column prop="remark" label="备注" min-width="120" />
      </el-table>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { getProducts, getCategories, checkInventory, getInventoryLogs } from '@/api'
import { parseTime } from '@/utils/ruoyi'


import { ElMessage } from 'element-plus'

const loading = ref(false)
const tableData = ref([])
const total = ref(0)
const categories = ref([])
const queryParams = reactive({ page: 1, pageSize: 10, keyword: '', categoryId: null, onlyLowStock: false })

async function loadData() {
  loading.value = true
  try {
    let cid = queryParams.categoryId || 0
    const res = await getProducts(queryParams.page, queryParams.pageSize, queryParams.keyword, cid)
    if (res && res.code === 0) {
      let list = res.data.list || []
      if (queryParams.onlyLowStock) {
        list = list.filter(item => item.stockQty <= item.minStock && item.minStock > 0)
      }
      tableData.value = list
      total.value = res.data.total || 0
    }
  } finally { loading.value = false }
}

async function loadCategories() {
  const res = await getCategories()
  if (res && res.code === 0) categories.value = res.data
}

function resetQuery() {
  queryParams.keyword = ''; queryParams.categoryId = null; queryParams.onlyLowStock = false; queryParams.page = 1
  loadData()
}

// 盘点
const checkFormRef = ref(null)
const checkDialogVisible = ref(false)
const checkProduct = ref(null)
const checkForm = reactive({ actualQty: null, remark: '' })
const checkRules = {
  actualQty: [{ required: true, message: '请填写真实数量', trigger: 'blur' }]
}

function handleCheck(row) {
  checkProduct.value = row
  checkForm.actualQty = null
  checkForm.remark = ''
  checkDialogVisible.value = true
}

async function submitCheck() {
  const valid = await checkFormRef.value.validate().catch(() => false)
  if (!valid) return
  const diff = checkForm.actualQty - checkProduct.value.stockQty
  let prefix = ''
  if (diff > 0) prefix = '盘库增加 '
  else if (diff < 0) prefix = '盘库减少 '
  else prefix = '盘库持平 '
  const fullRemark = prefix + (checkForm.remark || '')
  const result = await checkInventory(checkProduct.value.id, checkForm.actualQty, fullRemark)
  if (result.code === 0) {
    ElMessage.success(result.msg)
    checkDialogVisible.value = false
    loadData()
  } else {
    ElMessage.error(result.msg || '盘点失败')
  }
}

// 批量盘库
const batchCheckVisible = ref(false)
const batchQtyMap = reactive({})
const batchRemarkMap = reactive({})
const batchSubmitting = ref(false)

const batchCount = computed(() => {
  return tableData.value.filter(row => batchQtyMap[row.id] != null).length
})

function openBatchCheck() {
  tableData.value.forEach(row => {
    batchQtyMap[row.id] = null
    batchRemarkMap[row.id] = ''
  })
  batchCheckVisible.value = true
}

async function submitBatchCheck() {
  const items = tableData.value.filter(row => batchQtyMap[row.id] != null)
  if (items.length === 0) return
  batchSubmitting.value = true
  let done = 0
  for (const row of items) {
    const qty = batchQtyMap[row.id]
    const remark = batchRemarkMap[row.id] || ''
    const diff = qty - row.stockQty
    let prefix = ''
    if (diff > 0) prefix = '盘库增加 '
    else if (diff < 0) prefix = '盘库减少 '
    else prefix = '盘库持平 '
    const result = await checkInventory(row.id, qty, prefix + remark)
    if (result.code !== 0) {
      ElMessage.error(`第 ${done + 1} 项盘点失败：${result.msg}`)
      batchSubmitting.value = false
      return
    }
    done++
  }
  ElMessage.success(`盘点完成，共处理 ${done} 项`)
  batchCheckVisible.value = false
  loadData()
  batchSubmitting.value = false
}

// 库存记录
const logDialogVisible = ref(false)
const logProduct = ref(null)
const logData = ref([])

async function handleLog(row) {
  logProduct.value = row
  try {
    const res = await getInventoryLogs(1, 200, row.name, '', '', '')
    console.log('库存记录查询结果:', row.name, res)
    logData.value = (res && res.code === 0) ? (res.data.list || []) : []
  } catch(e) {
    console.error('查询库存记录失败:', e)
    logData.value = []
  }
  logDialogVisible.value = true
}

loadCategories()
loadData()
</script>

<style scoped>
.page-container{
  padding: 20px;
}
.text-danger { color: #f56c6c; }
</style>
