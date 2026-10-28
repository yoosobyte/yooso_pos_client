<template>
  <el-dialog v-model="visible" width="460px" :close-on-click-modal="false" draggable class="update-dialog">
    <template #header>
      <span class="dlg-title">发现新版本</span>
    </template>

    <div class="body">
      <div class="row">
        <span class="label">当前版本</span>
        <el-tag size="small" type="info">v{{ info.current }}</el-tag>
      </div>
      <div class="row">
        <span class="label">最新版本</span>
        <el-tag size="small" type="primary">v{{ info.newVersion }}</el-tag>
      </div>
      <div v-if="info.updateContent" class="changelog">
        <div class="cl-title">更新说明</div>
        <div class="cl-text">{{ info.updateContent }}</div>
      </div>
    </div>

    <template #footer>
      <div class="footer">
        <el-button size="small" @click="handleToday">今日不再弹出</el-button>
        <el-button size="small" @click="handleIgnore">忽略此版本</el-button>
        <el-button size="small" type="primary" @click="handleUpdate">
          立即更新 <el-icon style="margin-left:2px"><Download /></el-icon>
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref } from 'vue'
import { Download } from '@element-plus/icons-vue'
import { doUpdate } from '@/api'
import { ElMessage } from 'element-plus'

const visible = ref(false)
const info = ref({})

const TODAY_KEY = 'update_check_date'
const VERSION_KEY = 'update_skipped_version'

function shouldShow() {
  const today = new Date().toDateString()
  return localStorage.getItem(TODAY_KEY) !== today && localStorage.getItem(VERSION_KEY) !== info.value.newVersion
}

function open(versionInfo) {
  if (versionInfo.forceUpdate) {
    ElMessage.info('正在下载更新...')
    doUpdate(versionInfo.downloadUrl, versionInfo.newVersion)
    return
  }
  info.value = versionInfo
  if (!shouldShow()) return
  visible.value = true
}

function handleToday() {
  localStorage.setItem(TODAY_KEY, new Date().toDateString())
  visible.value = false
}

function handleIgnore() {
  localStorage.setItem(VERSION_KEY, info.value.newVersion)
  visible.value = false
}

async function handleUpdate() {
  localStorage.removeItem(VERSION_KEY)
  visible.value = false
  await doUpdate(info.value.downloadUrl, info.value.newVersion)
}

defineExpose({ open })
</script>

<style scoped>
.dlg-title {
  font-size: 16px;
  font-weight: 600;
}

.update-dialog :deep(.el-dialog__header) {
  padding: 12px 20px 0;
}

.update-dialog :deep(.el-dialog__body) {
  padding: 6px 20px;
}

.update-dialog :deep(.el-dialog__footer) {
  padding: 8px 20px 12px;
}

.body {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.row {
  display: flex;
  align-items: center;
  gap: 10px;
  line-height: 1;
}

.label {
  font-size: 13px;
  color: var(--el-text-color-secondary);
  width: 56px;
  flex-shrink: 0;
}

.row :deep(.el-tag) {
  padding: 0 8px;
  height: 24px;
  line-height: 24px;
}

.changelog {
  margin-top: 0;
}

.cl-title {
  font-size: 13px;
  color: var(--el-text-color-secondary);
  margin: 0 0 1px 0;
}

.cl-text {
  margin-top: 10px;
  font-size: 12px;
  line-height: 1.4;
  color: var(--el-text-color-regular);
  white-space: pre-wrap;
  max-height: 80px;
  overflow-y: auto;
  padding: 4px 6px;
  background: var(--el-fill-color-lighter);
  border-radius: 4px;
}

.footer {
  display: flex;
  justify-content: flex-end;
}
</style>
