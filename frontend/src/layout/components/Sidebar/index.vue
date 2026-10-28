<template>
  <div :class="['sidebar-theme-wrapper', {'has-logo':showLogo}, sideTheme]" class="sidebar-container">
    <logo v-if="showLogo" :collapse="isCollapse" />
    <el-scrollbar wrap-class="scrollbar-wrapper">
      <el-menu
          ref="menuRef"
          :default-openeds="defaultOpeneds"
          :default-active="activeMenu"
          :collapse="isCollapse"
          :background-color="getMenuBackground"
          :text-color="getMenuTextColor"
          :unique-opened="true"
          :active-text-color="theme"
          :collapse-transition="false"
          mode="vertical"
          :class="sideTheme"
      >
        <sidebar-item
            v-for="(route, index) in sidebarRoutes"
            :key="route.path + index"
            :item="route"
            :base-path="route.path"
        />
      </el-menu>
    </el-scrollbar>
    <div class="sidebar-version" @click="manualCheck" title="点击检测更新">
      <span v-if="!isCollapse">版本号：v{{ version }}</span>
      <span v-else>v{{ version }}</span>
    </div>
  </div>
</template>
<script setup>

import { ref, computed, watch, nextTick, onMounted,watchEffect } from 'vue'
import { useRoute } from 'vue-router'
import Logo from './Logo'
import SidebarItem from './SidebarItem'
import variables from '@/assets/styles/variables.module.scss'
import useAppStore from '@/store/modules/app'
import useSettingsStore from '@/store/modules/settings'
import { getAppVersion, checkVersion } from '@/api'
import { ElMessage } from 'element-plus'

const route = useRoute()
const appStore = useAppStore()
const settingsStore = useSettingsStore()

const menuRef = ref(null)
const defaultOpeneds = ref([])
const version = ref('')
let hasOpened = false
import { constantRoutes } from '@/router'

// 获取侧边栏菜单路由
const sidebarRoutes = computed(() => {
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
})

const showLogo = computed(() => settingsStore.sidebarLogo)
const sideTheme = computed(() => settingsStore.sideTheme)
const theme = computed(() => settingsStore.theme)
const isCollapse = computed(() => !appStore.sidebar.opened)

const getMenuBackground = computed(() => {
  if (settingsStore.isDark) {
    return 'var(--sidebar-bg)'
  }
  return sideTheme.value === 'theme-dark' ? variables.menuBg : variables.menuLightBg
})

const getMenuTextColor = computed(() => {
  if (settingsStore.isDark) {
    return 'var(--sidebar-text)'
  }
  return sideTheme.value === 'theme-dark' ? variables.menuText : variables.menuLightText
})

const activeMenu = computed(() => {
  const { meta, path } = route
  if (meta.activeMenu) {
    return meta.activeMenu
  }
  return path
})

onMounted(async () => {
  const v = await getAppVersion()
  if (v) version.value = v
})

async function manualCheck() {
  localStorage.removeItem('update_check_date')
  localStorage.removeItem('update_skipped_version')
  const info = await checkVersion()
  if (!info) {
    ElMessage.warning('检测失败，请稍后重试')
    return
  }
  if (info.isNew) {
    window.dispatchEvent(new CustomEvent('show-update-dialog', { detail: info }))
  } else {
    ElMessage.success('已是最新版本')
  }
}

</script>

<style lang="scss" scoped>
.sidebar-container {
  background-color: v-bind(getMenuBackground);
  display: flex;
  flex-direction: column;

  .scrollbar-wrapper {
    background-color: v-bind(getMenuBackground);
    flex: 1;
    overflow: hidden;
  }

  .sidebar-version {
    height: 36px;
    padding: 10px 20px;
    font-size: 12px;
    color: var(--sidebar-text, #999);
    border-top: 1px solid var(--sidebar-border, rgba(0, 0, 0, 0.06));
    box-shadow: 0 -2px 8px rgba(0, 0, 0, 0.08);
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
    gap: 0;
    cursor:pointer;
  }

  .el-menu {
    border: none;
    flex: 1;
    width: 100% !important;

    .el-menu-item {
      height: 44px;
      line-height: 44px;
      color: v-bind(getMenuTextColor);

      &:hover {
        background-color: var(--menu-hover, rgba(0, 0, 0, 0.06)) !important;
      }

      &.is-active {
        color: var(--menu-active-text, #409eff);
        background-color: var(--menu-hover, rgba(0, 0, 0, 0.06)) !important;
      }
    }

    .el-sub-menu__title {
      height: 44px !important;
      line-height: 44px !important;
      color: v-bind(getMenuTextColor);

      &:hover {
        background-color: var(--menu-hover, rgba(0, 0, 0, 0.06)) !important;
      }
    }
  }
}
</style>
