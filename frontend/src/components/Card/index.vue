<template>
  <div class="card">
    <div class="card-header" style="position: relative;">
      <h2>{{ picName }}</h2>
      <div style="position: absolute;bottom: 10px;right: 3px;color:#757575;">
        排序:{{ picSort }}
      </div>
    </div>
    <div class="card-body">
      <img style="width: 100%;height: auto;" :src="picUrl" :alt="picName">
    </div>
    <div class="card-footer">
      <el-button icon="Edit" size="small" type="primary" style="width: 49%;" @click="editItem">修改</el-button>
      <el-button icon="Delete" size="small" type="danger" style="width: 49%;" @click="removeItem">删除</el-button>
    </div>
  </div>
</template>

<script setup>
import { Edit, Delete } from '@element-plus/icons-vue'

defineProps({
  picId: {
    type: Number,
    default: null
  },
  picName: {
    type: String,
    default: '图片名称'
  },
  picUrl: {
    type: String,
    default: null
  },
  picSort: {
    type: Number,
    default: 0
  }
})

const emit = defineEmits(['editItem', 'removeItem'])

const editItem = () => {
  emit('editItem', picId)
}

const removeItem = () => {
  emit('removeItem', picId)
}
</script>

<style scoped>
.card {
  width: 200px;
  background-color: var(--el-bg-color);
  border-radius: 5px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  overflow: hidden;
  margin: 5px 10px 10px 5px;
  transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.card:hover {
  transform: scale(1.05);
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.2);
}

.card-header {
  background-color: #409EFF;
  color: white;
  padding: 10px;
  text-align: center;
  border-top-left-radius: 5px;
  border-top-right-radius: 5px;
}

.card-header h2 {
  margin: 0;
  font-size: 16px;
}

.card-body {
  height: 170px;
  padding: 1px;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: var(--el-fill-color-light);
}

.card-footer {
  background-color: var(--el-fill-color-light);
  padding: 10px;
  text-align: center;
  border-bottom-left-radius: 5px;
  border-bottom-right-radius: 5px;
  display: flex;
  justify-items: center;
  align-self: center;
}

/* ========== 暗黑模式适配 ========== */
html.dark .card {
  background-color: var(--el-bg-color);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.3);
}

html.dark .card:hover {
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.5);
}

html.dark .card-body {
  background-color: var(--el-fill-color);
}

html.dark .card-footer {
  background-color: var(--el-fill-color);
}
</style>