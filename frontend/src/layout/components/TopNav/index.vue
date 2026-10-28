<template>
  <el-menu
    :default-active="activeMenu"
    mode="horizontal"
    @select="handleSelect"
    :ellipsis="false"
  >
    <template v-for="(item, index) in topMenus">
      <el-menu-item :style="{'--theme': theme}" :index="item.path" :key="index" v-if="index < visibleNumber">
        <svg-icon
        v-if="item.meta && item.meta.icon && item.meta.icon !== '#'"
        :icon-class="item.meta.icon"/>
        {{ item.meta.title }}
      </el-menu-item>
    </template>

    <!-- 顶部菜单超出数量折叠 -->
    <el-sub-menu :style="{'--theme': theme}" index="more" v-if="topMenus.length > visibleNumber">
      <template #title>更多菜单</template>
      <template v-for="(item, index) in topMenus">
        <el-menu-item :style="{'--theme': theme}" :index="item.path" :key="index" v-if="index >= visibleNumber">
          <svg-icon
          v-if="item.meta && item.meta.icon && item.meta.icon !== '#'"
          :icon-class="item.meta.icon"/>
          {{ item.meta.title }}
        </el-menu-item>
      </template>
    </el-sub-menu>
  </el-menu>
</template>
<script setup>
import { ref, computed, onMounted, onBeforeUnmount, watch, watchEffect } from 'vue'

import { useRoute, useRouter } from 'vue-router'
import useAppStore from '@/store/modules/app'
import useSettingsStore from '@/store/modules/settings'

const route = useRoute()
const router = useRouter()
const appStore = useAppStore()
const settingsStore = useSettingsStore()

const theme = computed(() => settingsStore.theme)
const device = computed(() => appStore.device)
const activeMenu = computed(() => {
  const { meta, path } = route
  if (meta.activeMenu) {
    return meta.activeMenu
  }
  return path
})

const visibleNumber = ref(5)
import { constantRoutes } from '@/router'

// 获取路由数据
function routers() {
  const result = []
  
  // 遍历所有路由
  constantRoutes.forEach(route => {
    // 跳过隐藏的路由
    if (route.hidden) return
    
    // 如果路由有 children，处理 children
    if (route.children && route.children.length > 0) {
      // 创建一个父级菜单项（如果有 meta.title）
      if (route.meta?.title) {
        result.push({
          ...route,
          children: route.children.filter(child => !child.hidden && child.meta?.title)
        })
      } else {
        // 如果没有父级 meta.title，直接添加 children 作为顶级菜单
        route.children.forEach(child => {
          if (!child.hidden && child.meta?.title) {
            result.push({
              ...child,
              children: []
            })
          }
        })
      }
    } else if (route.meta?.title) {
      // 单个路由（没有 children）
      result.push(route)
    }
  })
  
  return result
}

// 顶部显示菜单
const topMenus = computed(() => {
  return routers().filter((f) => !f.hidden).slice(0, visibleNumber.value)
})

const moreRoutes = computed(() => {
  return routers().filter((f) => !f.hidden).slice(visibleNumber.value)
})

function setVisibleNumber() {
  const width = document.body.getBoundingClientRect().width / 3
  visibleNumber.value = Math.max(1, parseInt(width / 85))
}

function handleSelect(key, keyPath) {
  router.push(key)
}

onMounted(() => {
  window.addEventListener('resize', setVisibleNumber)
  setVisibleNumber()
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', setVisibleNumber)
})
</script>
