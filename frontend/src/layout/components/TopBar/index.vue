<template>
  <el-menu class="topbar-menu" :ellipsis="false" :default-active="activeMenu" :active-text-color="theme" mode="horizontal">
    <sidebar-item :key="route.path + index" v-for="(route, index) in topMenus" :item="route" :base-path="route.path" />

    <el-sub-menu index="more" class="el-sub-menu__hide-arrow" v-if="moreRoutes.length > 0">
      <template #title>
        <span>更多菜单</span>
      </template>
      <sidebar-item :key="route.path + index" v-for="(route, index) in moreRoutes" :item="route" :base-path="route.path" />
    </el-sub-menu>
  </el-menu>
</template>
<script setup>
import { ref, computed, onMounted, onBeforeUnmount, watch, watchEffect } from 'vue'

import { useRoute } from 'vue-router'
import SidebarItem from '../Sidebar/SidebarItem'
import useAppStore from '@/store/modules/app'
import useSettingsStore from '@/store/modules/settings'

const route = useRoute()
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

// 获取顶部菜单路由
function staticMenus() {
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

const topMenus = computed(() => {
  return staticMenus().filter((f) => !f.hidden).slice(0, visibleNumber.value)
})

const moreRoutes = computed(() => {
  return staticMenus().filter((f) => !f.hidden).slice(visibleNumber.value)
})

function setVisibleNumber() {
  const width = document.body.getBoundingClientRect().width / 3
  visibleNumber.value = Math.max(1, parseInt(width / 85))
}

onMounted(() => {
  window.addEventListener('resize', setVisibleNumber)
  setVisibleNumber()
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', setVisibleNumber)
})
</script>

<style lang="scss">
/* menu item */
.topbar-menu.el-menu--horizontal .el-submenu__title, .topbar-menu.el-menu--horizontal .el-menu-item {
  padding: 0 10px !important;
}

.topbar-menu.el-menu--horizontal > .el-menu-item {
  float: left;
  height: 50px !important;
  line-height: 50px !important;
  color: #303133 !important;
  padding: 0 5px !important;
  margin: 0 10px !important;
}

.el-sub-menu.is-active .svg-icon, .el-menu-item.is-active .svg-icon + span, .el-sub-menu.is-active .svg-icon + span, .el-sub-menu.is-active .el-sub-menu__title span {
  color: v-bind(theme);
}

/* sub-menu item */
.topbar-menu.el-menu--horizontal > .el-sub-menu .el-sub-menu__title {
  float: left;
  line-height: 50px !important;
  color: #303133 !important;
  margin: 0 15px -3px!important;
}

/* topbar more arrow */
.topbar-menu .el-sub-menu .el-sub-menu__icon-arrow {
  position: static;
  margin-left: 8px;
  margin-top: 0px;
  display: block !important;
}

/* menu__title el-menu-item */
.topbar-menu.el-menu--horizontal .el-sub-menu__title, .topbar-menu.el-menu--horizontal .el-menu-item {
  height: 60px;
}
</style>
