<template>
  <div class="page-container">
    <el-form :inline="true" :model="queryParams" size="default">
      <el-form-item label="关键词">
        <el-input v-model="queryParams.keyword" placeholder="名称/电话" clearable style="width: 200px;" @keyup.enter="loadData" />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="loadData">查询</el-button>
        <el-button @click="resetQuery">重置</el-button>
        <el-button type="success" @click="handleAdd">新增客户</el-button>
      </el-form-item>
    </el-form>

    <el-table :data="tableData" v-loading="loading" border stripe>
      <el-table-column prop="name" label="客户名称" />
      <el-table-column prop="phone" label="电话" width="140" />
      <el-table-column prop="address" label="地址" min-width="180" />
      <el-table-column label="备注" width="80">
        <template #default="{ row }">{{ row.remark || '-' }}</template>
      </el-table-column>
      <el-table-column label="操作" width="140" fixed="right" align="center">
        <template #default="{ row }">
          <el-button type="primary" size="small" @click="handleEdit(row)">编辑</el-button>
          <el-popconfirm
              title="确定删除客户？"
              icon="SetUp"
              @confirm="handleDel(row.id)"
              popper-class="custom-popconfirm"
          >
            <template #reference>
              <el-button type="danger" size="small">删除</el-button>
            </template>
            <template #actions="{ confirm, cancel }">
              <div class="popconfirm-actions">
                <el-button size="small" type="primary" @click="cancel">点错了</el-button>
                <el-button size="small" type="danger" @click="confirm">确定删除</el-button>
              </div>
            </template>
          </el-popconfirm>
        </template>
      </el-table-column>
    </el-table>

    <pagination
      size="small"
      v-show="total > 0" :total="total"
      v-model:page="queryParams.page" v-model:limit="queryParams.pageSize"
      @pagination="loadData"
    />

    <el-dialog :title="dialogTitle" v-model="dialogVisible" width="450px" draggable>
      <el-form :model="form" ref="formRef" :rules="rules" label-position="top">
        <el-form-item label="客户名称" prop="name">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="客户电话">
          <el-input v-model="form.phone" />
        </el-form-item>
        <el-form-item label="客户地址">
          <el-input v-model="form.address" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="form.remark" type="textarea" :rows="2" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitForm">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, nextTick } from 'vue'
import { getCustomers, createCustomer, updateCustomer, deleteCustomer } from '@/api'


import { ElMessage } from 'element-plus'

const loading = ref(false)
const tableData = ref([])
const total = ref(0)
const queryParams = reactive({ page: 1, pageSize: 10, keyword: '' })

const dialogVisible = ref(false)
const dialogTitle = ref('')
const formRef = ref(null)
const form = reactive({ id: 0, name: '', phone: '', address: '', remark: '' })
const rules = { name: [{ required: true, message: '请输入名称', trigger: 'blur' }] }

async function loadData() {
  loading.value = true
  try {
    const res = await getCustomers(queryParams.page, queryParams.pageSize, queryParams.keyword)
    if (res && res.code === 0) { tableData.value = res.data.list || []; total.value = res.data.total || 0 }
  } finally { loading.value = false }
}

function resetQuery() { queryParams.keyword = ''; queryParams.page = 1; loadData() }

function handleAdd() {
  dialogTitle.value = '新增客户'
  Object.assign(form, { id: 0, name: '', phone: '', address: '', remark: '' })
  dialogVisible.value = true
  nextTick(() => formRef.value?.resetFields())
}

function handleEdit(row) {
  dialogTitle.value = '编辑客户'
  Object.assign(form, { ...row })
  dialogVisible.value = true
}

async function submitForm() {
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return
  const result = form.id
    ? await updateCustomer({ ...form })
    : await createCustomer({ ...form })
  if (result.code === 0) {
    ElMessage.success(result.msg)
    dialogVisible.value = false
    loadData()
  } else {
    ElMessage.error(result.msg || '操作失败')
  }
}

async function handleDel(id) {
  const result = await deleteCustomer(id)
  if (result.code === 0) { ElMessage.success(result.msg); loadData() }
  else { ElMessage.error(result.msg || '删除失败') }
}


loadData()
</script>
<style scoped>
.page-container{
  padding: 20px;
}
</style>