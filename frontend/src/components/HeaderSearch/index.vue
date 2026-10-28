<template>
  <div class="header-search">
    <svg-icon class-name="search-icon" icon-class="search" @click.stop="openDialog" />

    <el-dialog
      v-model="dialogVisible"
      title="菜单搜索"
      width="600px"
      :close-on-click-modal="true"
      top="10vh"
      draggable
      append-to-body
    >
      <div class="search-dialog-content">
        <!-- 搜索输入框 -->
        <el-input
          ref="searchInputRef"
          v-model="query"
          placeholder="请输入菜单名称"
          clearable
          size="large"
          @input="handleQuery"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>

        <!-- 搜索结果区域 -->
        <div class="search-result-container">
          <div class="search-result-empty" v-if="query && searchResults.length === 0">
            <el-empty description="未搜索到相关菜单" :image-size="60" />
          </div>

          <div class="search-result-list" v-if="searchResults.length > 0">
            <div class="search-result-header">
              共搜索到 <b>{{ searchResults.length }}</b> 个结果
            </div>
            <div
              v-for="item in searchResults"
              :key="item.path"
              class="search-result-item"
              @click="navigateTo(item)"
            >
              <div class="result-left">
                <svg-icon :icon-class="item.icon" class="result-icon" />
                <div class="result-text">
                  <div class="result-title">{{ item.title.join(' > ') }}</div>
                  <div class="result-path-text">{{ item.path }}</div>
                </div>
              </div>
              <el-icon class="result-enter"><ArrowRight /></el-icon>
            </div>
          </div>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import Fuse from 'fuse.js'
import { constantRoutes } from '@/router'
import { isHttp } from '@/utils/validate'
import { Search, ArrowRight } from '@element-plus/icons-vue'

const router = useRouter()
const dialogVisible = ref(false)
const query = ref('')
const searchResults = ref([])
const searchInputRef = ref(null)
let fuseInstance = null

// 拼接路径
function resolvePath(base, p) {
  if (isHttp(p) || p.startsWith('/')) return p
  return (base + '/' + p).replace(/\/+/g, '/')
}

// 从路由表中提取可搜索的菜单项
function generateRoutes(routes, basePath = '/', prefixTitle = []) {
  let res = []
  for (const r of routes) {
    if (r.hidden) continue
    if (r.path === '/:pathMatch(.*)*' || r.path === '/401' || r.path === '/lock') continue

    const fullPath = resolvePath(basePath, r.path)
    const data = {
      path: fullPath,
      icon: r.meta?.icon || 'list',
      title: [...prefixTitle]
    }

    if (r.meta?.title) {
      data.title = [...data.title, r.meta.title]

      const hasChildren = r.children && r.children.filter(c => !c.hidden).length > 0
      if (!hasChildren) {
        res.push(data)
      }
    }

    if (r.children) {
      const children = generateRoutes(r.children, fullPath, data.title)
      res = [...res, ...children]
    }
  }
  return res
}

const menuPool = computed(() => generateRoutes(constantRoutes))

// 初始化 Fuse
function initFuse() {
  fuseInstance = new Fuse(menuPool.value, {
    shouldSort: true,
    threshold: 0.3,
    location: 0,
    distance: 100,
    minMatchCharLength: 1,
    keys: [
      { name: 'title', weight: 0.7 },
      { name: 'path', weight: 0.3 }
    ]
  })
}

function openDialog() {
  dialogVisible.value = true
  query.value = ''
  initFuse()
  searchResults.value = [...menuPool.value]
  nextTick(() => {
    if (searchInputRef.value) {
      searchInputRef.value.focus()
    }
  })
}

function handleQuery(val) {
  if (val && val.trim() !== '') {
    searchResults.value = fuseInstance ? fuseInstance.search(val).map(r => r.item) : []
  } else {
    searchResults.value = [...menuPool.value]
  }
}

function navigateTo(item) {
  if (isHttp(item.path)) {
    window.open(item.path, '_blank')
  } else {
    router.push(item.path)
  }
  dialogVisible.value = false
}
</script>

<style lang="scss" scoped>
.header-search {
  display: inline-flex;
  align-items: center;

  .search-icon {
    cursor: pointer;
    font-size: 16px;
    vertical-align: middle;
    transition: transform 0.3s;

    &:hover {
      transform: scale(1.15);
    }
  }
}

.search-dialog-content {
  display: flex;
  flex-direction: column;
  gap: 20px;

  :deep(.el-input__wrapper) {
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  }
}

.search-result-container {
  min-height: 120px;
  max-height: 50vh;
  overflow-y: auto;
}

.search-result-empty {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 20px 0;
}

.search-result-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.search-result-header {
  padding: 4px 0 8px;
  font-size: 13px;
  color: var(--el-text-color-secondary);
  border-bottom: 1px solid var(--el-border-color-lighter);
  margin-bottom: 4px;
}

.search-result-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 12px;
  border-radius: 6px;
  cursor: pointer;
  transition: background-color 0.2s;

  &:hover {
    background-color: var(--el-color-primary-light-9);
  }

  .result-left {
    display: flex;
    align-items: center;
    gap: 12px;
    flex: 1;
    min-width: 0;
  }

  .result-icon {
    font-size: 20px;
    flex-shrink: 0;
    color: var(--el-color-primary);
  }

  .result-text {
    flex: 1;
    min-width: 0;
  }

  .result-title {
    font-size: 14px;
    color: var(--el-text-color-primary);
  }

  .result-path-text {
    font-size: 12px;
    color: var(--el-text-color-secondary);
    margin-top: 2px;
  }

  .result-enter {
    font-size: 14px;
    color: var(--el-text-color-placeholder);
    flex-shrink: 0;
  }
}
</style>
