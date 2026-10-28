<template>
  <div class="page-container">
    <el-row :gutter="20">
      <!-- 分类树 -->
      <el-col :span="12">
        <el-card header="商品分类" shadow="hover" class="show-card">
          <div style="margin-bottom: 10px;">
            <el-button type="primary" size="small" style="width:100%" @click="handleAdd(null)">新增根分类</el-button>
          </div>
          <el-tree
              ref="treeRef"
              class="category-tree"
              :data="treeData"
              :props="{ label: 'name', children: 'children' }"
              node-key="id"
              default-expand-all
              highlight-current
          >
            <template #default="{ node, data }">
              <span class="tree-node">
                <span>{{ node.label + ' ▸' + data.sort }}</span>
                <span class="tree-actions">
                  <el-button link type="primary" size="small" @click.stop="handleAdd(data)">添加子分类</el-button>
                  <el-button link type="warning" size="small" @click.stop="handleEdit(data)">编辑</el-button>
                  <el-popconfirm
                      title="确定删除分类？"
                      icon="SetUp"
                      @confirm="handleDelete(data.id)"
                      popper-class="custom-popconfirm"
                  >
                    <template #reference>
                      <el-button link type="danger" size="small">删除</el-button>
                    </template>
                    <template #actions="{ confirm, cancel }">
                      <div class="popconfirm-actions">
                        <el-button size="small" type="primary" @click="cancel">点错了</el-button>
                        <el-button size="small" type="danger" @click="confirm">确定删除</el-button>
                      </div>
                    </template>
                  </el-popconfirm>
                </span>
              </span>
            </template>
          </el-tree>
        </el-card>
      </el-col>

      <!-- 单位管理 -->
      <el-col :span="12">
        <el-card header="商品单位" shadow="hover" class="show-card">
          <div style="margin-bottom: 10px;">
            <el-button type="primary" size="small" style="width:100%" @click="handleUnitAdd">新增单位</el-button>
          </div>
          <el-table :data="unitList" size="small" border>
            <el-table-column prop="name" label="单位名称" />
            <el-table-column label="操作" width="130" align="center">
              <template #default="{ row }">
                <el-button type="warning" size="small" @click="handleUnitEdit(row)">编辑</el-button>
                <el-popconfirm
                    title="确定删除单位？"
                    icon="SetUp"
                    @confirm="handleUnitDelete(row.id)"
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
        </el-card>
      </el-col>
    </el-row>

    <!-- 分类弹窗 -->
    <el-dialog :title="catDialogTitle" v-model="catDialogVisible" width="400px" draggable>
      <el-form ref="catFormRef" :model="catForm" :rules="catRules" label-position="top">
        <el-form-item label="分类名称" prop="name">
          <el-input v-model="catForm.name" />
        </el-form-item>
        <el-form-item label="分类排序">
          <el-input-number v-model="catForm.sort" :min="0" controls-position="right" style="width: 100%;" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="catDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitCatForm">确定</el-button>
      </template>
    </el-dialog>

    <!-- 单位弹窗 -->
    <el-dialog :title="unitDialogTitle" v-model="unitDialogVisible" width="350px">
      <el-form ref="unitFormRef" :model="unitForm" :rules="unitRules" label-position="top">
        <el-form-item label="单位名称" prop="name">
          <el-input v-model="unitForm.name" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="unitDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitUnitForm">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, nextTick } from 'vue'
import { getCategories, createCategory, updateCategory, deleteCategory, getMaxCategorySort, getUnits, createUnit, updateUnit, deleteUnit } from '@/api'
import { ElMessage } from 'element-plus'

const treeRef = ref(null)
const treeData = ref([])
const unitList = ref([])

// 分类
const catFormRef = ref(null)
const catDialogVisible = ref(false)
const catDialogTitle = ref('')
const catForm = reactive({ id: 0, name: '', parentId: null, sort: 0 })
const catRules = { name: [{ required: true, message: '请输入分类名称', trigger: 'blur' }] }

async function loadCategories() {
  const res = await getCategories()
  if (res && res.code === 0) {
    treeData.value = res.data
    nextTick(() => {
      showExpandIcons()
    })
  }
}

function showExpandIcons() {
  const treeEl = document.querySelector('.category-tree')
  if (!treeEl) return

  const nodes = treeEl.querySelectorAll('.el-tree-node__children > .el-tree-node > .el-tree-node__content')
  nodes.forEach(node => {
    const icon = node.querySelector('.el-tree-node__expand-icon')
    if (icon) {
      icon.style.visibility = 'visible'
      const treeNode = node.parentElement
      if (treeNode.classList.contains('is-leaf')) {
        icon.style.cursor = 'pointer'
        icon.addEventListener('click', (e) => {
          e.stopPropagation()
          ElMessage.info('该分类下暂无子分类')
        })
      }
    }
  })
}

async function handleAdd(parent) {
  catDialogTitle.value = parent ? `添加「${parent.name}」的子分类` : '新增根分类'
  const sortRes = await getMaxCategorySort(parent?.id)
  const defaultSort = (sortRes && sortRes.code === 0) ? sortRes.data : 0
  Object.assign(catForm, { id: 0, name: '', parentId: parent?.id || null, sort: defaultSort })
  catDialogVisible.value = true
}

function handleEdit(data) {
  catDialogTitle.value = '编辑分类'
  Object.assign(catForm, { id: data.id, name: data.name, parentId: data.parentId || null, sort: data.sort })
  catDialogVisible.value = true
}

async function submitCatForm() {
  const valid = await catFormRef.value.validate().catch(() => false)
  if (!valid) return
  const result = catForm.id
    ? await updateCategory({ ...catForm })
    : await createCategory({ ...catForm })
  if (result.code === 0) {
    ElMessage.success(result.msg)
    catDialogVisible.value = false
    loadCategories()
  } else {
    ElMessage.error(result.msg || '操作失败')
  }
}

async function handleDelete(id) {
  const result = await deleteCategory(id)
  if (result.code === 0) {
    ElMessage.success(result.msg)
    loadCategories()
  } else {
    ElMessage.error(result.msg || '删除失败')
  }
}

// 单位
const unitFormRef = ref(null)
const unitDialogVisible = ref(false)
const unitDialogTitle = ref('')
const unitForm = reactive({ id: 0, name: '' })
const unitRules = { name: [{ required: true, message: '请输入单位名称', trigger: 'blur' }] }

async function loadUnits() {
  const res = await getUnits()
  if (res && res.code === 0) unitList.value = res.data
}

function handleUnitAdd() {
  unitDialogTitle.value = '新增单位'
  Object.assign(unitForm, { id: 0, name: '' })
  unitDialogVisible.value = true
}

function handleUnitEdit(row) {
  unitDialogTitle.value = '编辑单位'
  Object.assign(unitForm, { id: row.id, name: row.name })
  unitDialogVisible.value = true
}

async function submitUnitForm() {
  const valid = await unitFormRef.value.validate().catch(() => false)
  if (!valid) return
  const result = unitForm.id
      ? await updateUnit({ ...unitForm })
      : await createUnit({ ...unitForm })

  if (result.code === 0) {
    ElMessage.success('操作成功')
    unitDialogVisible.value = false
    loadUnits()
  } else {
    ElMessage.error(result.msg || '操作失败')
  }
}

async function handleUnitDelete(id) {
  const result = await deleteUnit(id)
  if (result.code === 0) {
    ElMessage.success(result.msg)
    loadUnits()
  } else {
    ElMessage.error(result.msg || '删除失败')
  }
}

loadCategories()
loadUnits()
</script>


<style scoped>
.page-container {
  padding: 20px;
}

.tree-node {
  flex: 1;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-right: 8px;
}

.tree-actions {
  display: inline-flex;
  gap: 4px;
}

.tree-actions .el-button {
  padding: 0 4px;
}

.category-tree {
  border: 1px solid var(--el-border-color);
  border-radius: 4px;
  padding: 8px;
}

/* ===== 树节点展开图标样式 ===== */
:deep(.category-tree .el-tree-node__children > .el-tree-node > .el-tree-node__content .el-tree-node__expand-icon) {
  visibility: visible !important;
}

:deep(.category-tree .el-tree-node__children > .el-tree-node.is-leaf .el-tree-node__expand-icon) {
  color: var(--el-text-color-placeholder);
  cursor: default;
}

:deep(.category-tree .el-tree-node__children > .el-tree-node:not(.is-leaf) .el-tree-node__expand-icon) {
  color: var(--el-text-color-secondary);
  cursor: pointer;
}

:deep(.category-tree .el-tree-node__expand-icon) {
  transition: transform 0.3s ease;
}

:deep(.category-tree .el-tree-node__expand-icon.expanded) {
  transform: rotate(90deg);
}
</style>