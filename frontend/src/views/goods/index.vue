<template>
  <div class="page-container">
    <!-- 搜索栏 -->
    <el-form :inline="true" :model="queryParams" size="default">
      <el-form-item label="关键词">
        <el-input v-model="queryParams.keyword" placeholder="名称/编号/规格" clearable style="width: 200px;" @keyup.enter="loadData" />
      </el-form-item>
      <el-form-item label="分类">
        <el-tree-select
          v-model="queryParams.categoryId" :data="categoryTree"
          :props="{ label: 'name', value: 'id', children: 'children' }"
          placeholder="全部分类" clearable check-strictly
          style="width: 180px;"
        />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="loadData">查询</el-button>
        <el-button @click="resetQuery">重置</el-button>
        <el-button type="success" @click="handleAdd">新增商品</el-button>
      </el-form-item>
    </el-form>

    <!-- 表格 -->
    <el-table :data="tableData" v-loading="loading" border stripe>
      <el-table-column prop="code" label="编号" width="160" />
      <el-table-column prop="name" label="商品名称" min-width="140" />
      <el-table-column label="分类">
        <template #default="{ row }">{{ row.category?.name || '-' }}</template>
      </el-table-column>
      <el-table-column label="单位" width="70">
        <template #default="{ row }">{{ row.unit?.name || '-' }}</template>
      </el-table-column>
      <el-table-column label="规格">
        <template #default="{ row }">{{ row.spec || '-' }}</template>
      </el-table-column>
      <el-table-column label="进价" width="90" align="right">
        <template #default="{ row }">¥{{ row.costPrice.toFixed(2) }}</template>
      </el-table-column>
      <el-table-column label="售价" width="90" align="right">
        <template #default="{ row }">¥{{ row.salePrice.toFixed(2) }}</template>
      </el-table-column>
      <el-table-column prop="stockQty" label="库存" width="70" align="center" />
      <el-table-column label="状态" width="85" align="center">
        <template #default="{ row }">
          <el-tag v-if="row.status === 1" type="success">已上架</el-tag>
          <el-tag v-else type="info">已下架</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="150" fixed="right" align="center">
        <template #default="{ row }">
          <el-button type="primary" size="small" @click="handleEdit(row)">编辑</el-button>
          <el-button type="danger" size="small" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <pagination
      size="small"
      v-show="total > 0" :total="total"
      v-model:page="queryParams.page" v-model:limit="queryParams.pageSize"
      @pagination="loadData"
    />

    <!-- 新增/编辑弹窗 -->
    <el-dialog :title="dialogTitle" v-model="dialogVisible" width="600px" :close-on-click-modal="false" draggable>
      <el-form :model="form" :rules="rules" ref="formRef" label-position="top">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="商品名称" prop="name">
              <el-input v-model="form.name" placeholder="请输入" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="编号" prop="code">
              <el-input v-model="form.code" placeholder="自动生成" disabled />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="分类">
              <el-tree-select
                v-model="form.categoryId" :data="categoryTree"
                :props="{ label: 'name', value: 'id', children: 'children' }"
                placeholder="选择分类" check-strictly filterable clearable
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="单位">
              <el-select v-model="form.unitId" placeholder="选择单位" filterable clearable>
                <el-option v-for="u in unitList" :key="u.id" :label="u.name" :value="u.id" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="规格">
              <el-input v-model="form.spec" placeholder="如：500ml/红色" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="状态">
              <el-switch v-model="form.statusBool" active-text="上架" inactive-text="下架" class="my-switch"/>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="进价" prop="costPrice">
              <el-input-number v-model="form.costPrice" :min="0" controls-position="right" style="width:100%" placeholder="请输入进价"/>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="售价" prop="salePrice">
              <el-input-number v-model="form.salePrice" :min="0" controls-position="right" style="width:100%" placeholder="请输入售价"/>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="最低库存预警值">
              <el-input-number v-model="form.minStock" :min="0" controls-position="right" style="width:100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="最高库存提示值">
              <el-input-number v-model="form.maxStock" :min="0" controls-position="right" style="width:100%" />
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="商品备注">
              <el-input v-model="form.remark" type="textarea" :rows="2" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitForm">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, nextTick } from 'vue'
import { getProducts, getAllProducts, createProduct, updateProduct, deleteProduct, getCategories, getUnits } from '@/api'


import { ElMessage, ElMessageBox } from 'element-plus'

const loading = ref(false)
const tableData = ref([])
const total = ref(0)
const queryParams = reactive({ page: 1, pageSize: 10, keyword: '', categoryId: null })
const categoryTree = ref([])
const unitList = ref([])

const dialogVisible = ref(false)
const dialogTitle = ref('')
const formRef = ref(null)
const initForm = { id: 0, name: '', code: '', spec: '', categoryId: null, unitId: null, costPrice: null, salePrice: null, minStock: 0, maxStock: 0, statusBool: true, remark: '' }
const form = reactive({ ...initForm })
const rules = {
  name: [{ required: true, message: '请输入商品名称', trigger: 'blur' }],
  costPrice: [{ required: true, message: '请输入进价', trigger: 'blur' }],
  salePrice: [{ required: true, message: '请输入售价', trigger: 'blur' }]
}

async function loadData() {
  loading.value = true
  try {
    const res = await getProducts(queryParams.page, queryParams.pageSize, queryParams.keyword, queryParams.categoryId || 0)
    if (res && res.code === 0) { tableData.value = res.data.list || []; total.value = res.data.total || 0 }
  } finally { loading.value = false }
}

async function loadCategories() {
  const res = await getCategories()
  if (res && res.code === 0) categoryTree.value = res.data
}

async function loadUnits() {
  const res = await getUnits()
  if (res && res.code === 0) unitList.value = res.data
}

function resetQuery() {
  queryParams.keyword = ''; queryParams.categoryId = null; queryParams.page = 1
  loadData()
}

function handleAdd() {
  dialogTitle.value = '新增商品'
  // 清除上次编辑遗留的 createdAt/updatedAt/deletedAt/category/unit 等脏属性
  for (const key of Object.keys(form)) {
    if (!(key in initForm)) delete form[key]
  }
  Object.assign(form, initForm)
  dialogVisible.value = true
  nextTick(() => {
    formRef.value?.resetFields()
    form.code = 'SP' + crypto.randomUUID().replace(/-/g, '').slice(0, 12).toUpperCase()
  })
}

function handleEdit(row) {
  dialogTitle.value = '编辑商品'
  Object.assign(form, {
    ...row,
    statusBool: row.status === 1,
    categoryId: row.categoryId || null,
    unitId: row.unitId || null
  })
  dialogVisible.value = true
}

async function submitForm() {
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return
  // 只提取后端需要的字段，避免发送 createdAt/updatedAt/deletedAt 等脏数据
  const data = {
    id: form.id || 0,
    name: form.name,
    code: form.code || 'SP' + crypto.randomUUID().replace(/-/g, '').slice(0, 12).toUpperCase(),
    spec: form.spec,
    categoryId: form.categoryId || null,
    unitId: form.unitId || null,
    costPrice: form.costPrice,
    salePrice: form.salePrice,
    minStock: form.minStock,
    maxStock: form.maxStock,
    status: form.statusBool ? 1 : 0,
    remark: form.remark
  }
  const result = data.id
    ? await updateProduct(data)
    : await createProduct(data)
  if (result.code === 0) {
    ElMessage.success(result.msg)
    dialogVisible.value = false
    loadData()
  } else {
    ElMessage.error(result.msg || '操作失败')
  }
}

async function handleDelete(row) {
  await ElMessageBox.confirm(`确定删除商品"${row.name}"吗？`, '提示', { type: 'warning' })
  const result = await deleteProduct(row.id)
  if (result.code === 0) {
    ElMessage.success(result.msg)
    loadData()
  } else {
    ElMessage.error(result.msg || '删除失败')
  }
}

loadCategories()
loadUnits()
loadData()
</script>

<style scoped>
.page-container{
  padding: 20px;
}
:deep(.my-switch .el-switch__core) {
  min-width: 58px !important;
}
</style>
