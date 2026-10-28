<template>
  <div id="tags-view-container" class="tags-view-container" :class="{ 'tags-view-container--chrome': tagsViewStyle === 'chrome' }">
    <!-- 左切换箭头 -->
    <span class="tags-nav-btn tags-nav-btn--left" :class="{ disabled: !canScrollLeft }" @click="scrollLeft">
      <el-icon><arrow-left /></el-icon>
    </span>

    <!-- 标签滚动区 -->
    <scroll-pane ref="scrollPaneRef" class="tags-view-wrapper" @scroll="handleScroll" @update-arrows="updateArrowState">
      <router-link
        v-for="tag in visitedViews"
        :key="tag.path"
        :data-path="tag.path"
        :class="{ 'active': isActive(tag), 'has-icon': tagsIcon }"
        :to="{ path: tag.path, query: tag.query, fullPath: tag.fullPath }"
        class="tags-view-item"
        :style="tagActiveStyle(tag)"
        @click.middle="!isAffix(tag) ? closeSelectedTag(tag) : ''"
        @contextmenu.prevent="openMenu(tag, $event)"
      >
        <svg-icon v-if="tagsIcon && tag.meta && tag.meta.icon && tag.meta.icon !== '#'" :icon-class="tag.meta.icon" style="margin-right: 3px;" />
        {{ tag.title }}
        <span v-if="!isAffix(tag)" @click.prevent.stop="closeSelectedTag(tag)" class="tags-close-btn">
          <close class="el-icon-close" />
        </span>
      </router-link>
    </scroll-pane>

    <!-- 右切换箭头 -->
    <span class="tags-nav-btn tags-nav-btn--right" :class="{ disabled: !canScrollRight }" @click="scrollRight">
      <el-icon><arrow-right /></el-icon>
    </span>

    <!-- 下拉操作菜单 -->
    <el-dropdown class="tags-action-dropdown" trigger="click" placement="bottom-end" @command="handleDropdownCommand">
      <span class="tags-action-btn">
        <el-icon><arrow-down /></el-icon>
      </span>
      <template #dropdown>
        <el-dropdown-menu class="tags-dropdown-menu">
          <el-dropdown-item v-if="!isAffix(selectedDropdownTag)" command="close"><close style="width: 1em; height: 1em;" />关闭当前</el-dropdown-item>
          <el-dropdown-item command="closeOthers"><circle-close style="width: 1em; height: 1em;" />关闭其他</el-dropdown-item>
          <el-dropdown-item command="closeLeft" :disabled="isFirstView()"><back style="width: 1em; height: 1em;" />关闭左侧</el-dropdown-item>
          <el-dropdown-item command="closeRight" :disabled="isLastView()"><right style="width: 1em; height: 1em;" />关闭右侧</el-dropdown-item>
          <el-dropdown-item command="closeAll"><circle-close style="width: 1em; height: 1em;" />全部关闭</el-dropdown-item>
        </el-dropdown-menu>
      </template>
    </el-dropdown>

    <!-- 右键菜单 -->
    <ul v-show="visible" :style="{ left: left + 'px', top: top + 'px' }" class="contextmenu">
      <li @click="refreshSelectedTag(selectedTag)"><refresh-right style="width: 1em; height: 1em;" />刷新页面</li>
      <li v-if="!isAffix(selectedTag)" @click="closeSelectedTag(selectedTag)"><close style="width: 1em; height: 1em;" />关闭当前</li>
      <li @click="closeOthersTags"><circle-close style="width: 1em; height: 1em;" />关闭其他</li>
      <li @click="closeLeftTags" :disabled="isFirstView()"><back style="width: 1em; height: 1em;" />关闭左侧</li>
      <li @click="closeRightTags" :disabled="isLastView()"><right style="width: 1em; height: 1em;" />关闭右侧</li>
      <li @click="closeAllTags(selectedTag)"><circle-close style="width: 1em; height: 1em;" />全部关闭</li>
    </ul>
  </div>
</template>
<script setup>
import { ref, computed, watch, nextTick, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import ScrollPane from './ScrollPane'
import useTagsViewStore from '@/store/modules/tagsView'
import useSettingsStore from '@/store/modules/settings'
import useAppStore from '@/store/modules/app'
import { ElMessage } from 'element-plus'
import {
  ArrowLeft,
  ArrowRight,
  ArrowDown,
  Close,
  CircleClose,
  Back,
  Right,
  RefreshRight
} from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const tagsViewStore = useTagsViewStore()
const settingsStore = useSettingsStore()
import { constantRoutes } from '@/router'


const visitedViews = computed(() => useTagsViewStore().visitedViews)

// 获取可搜索的路由
function routes() {
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

const theme = computed(() => useSettingsStore().theme)
const tagsIcon = computed(() => useSettingsStore().tagsIcon)
const tagsViewStyle = computed(() => useSettingsStore().tagsViewStyle)

const scrollPaneRef = ref(null)
const left = ref(0)
const top = ref(0)
const visible = ref(false)
const selectedTag = ref({})
const selectedDropdownTag = ref({})
const canScrollLeft = ref(false)
const canScrollRight = ref(false)

const filterRoutes = computed(() => {
  return routes()
})

const isActive = (tag) => {
  return tag.path === route.path
}

const isAffix = (tag) => {
  return tag.meta && tag.meta.affix
}

const isFirstView = () => {
  try {
    return selectedDropdownTag.value.fullPath === filterRoutes.value[0].fullPath
  } catch (err) {
    return false
  }
}

const isLastView = () => {
  try {
    return selectedDropdownTag.value.fullPath === filterRoutes.value[filterRoutes.value.length - 1].fullPath
  } catch (err) {
    return false
  }
}

const tagActiveStyle = (tag) => {
  if (!isActive(tag)) return {}
  return {
    'background-color': theme.value,
    'border-color': theme.value
  }
}

const toLastView = (visitedViews, view) => {
  const latestView = visitedViews.slice(-1)[0]
  if (latestView) {
    router.push(latestView.fullPath)
  } else {
    // 如果 tags 栏为空，重定向到默认页
    if (view.name === 'Dashboard') {
      router.replace({ path: '/redirect' + view.fullPath })
    } else {
      router.push('/')
    }
  }
}

const openMenu = (tag, e) => {
  const menuMinWidth = 105
  const offsetLeft = scrollPaneRef.value?.$el.getBoundingClientRect().left
  const offsetWidth = scrollPaneRef.value?.$el.offsetWidth
  const maxLeft = offsetWidth - menuMinWidth
  const l = e.clientX - offsetLeft + 15

  if (l > maxLeft) {
    left.value = maxLeft
  } else {
    left.value = l
  }

  top.value = e.clientY
  visible.value = true
  selectedTag.value = tag
}

const closeMenu = () => {
  visible.value = false
}

const refreshSelectedTag = (view) => {
  tagsViewStore.delCachedView(view)
  const { fullPath } = view
  nextTick(() => {
    router.replace({ path: '/redirect' + fullPath })
  })
}

const closeSelectedTag = (view) => {
  tagsViewStore.delView(view).then(({ visitedViews }) => {
    if (isActive(view)) {
      toLastView(visitedViews, view)
    }
  })
}

const closeOthersTags = () => {
  if (selectedDropdownTag.value.fullPath !== undefined && selectedDropdownTag.value.fullPath !== '') {
    router.push(selectedDropdownTag.value.fullPath)
  }
  tagsViewStore.delOthersViews(selectedDropdownTag.value).then(() => {
    moveToCurrentTag()
  })
}

const closeLeftTags = () => {
  if (selectedDropdownTag.value.fullPath !== undefined && selectedDropdownTag.value.fullPath !== '') {
    router.push(selectedDropdownTag.value.fullPath)
  }
  tagsViewStore.delLeftViews(selectedDropdownTag.value).then(() => {
    moveToCurrentTag()
  })
}

const closeRightTags = () => {
  if (selectedDropdownTag.value.fullPath !== undefined && selectedDropdownTag.value.fullPath !== '') {
    router.push(selectedDropdownTag.value.fullPath)
  }
  tagsViewStore.delRightViews(selectedDropdownTag.value).then(() => {
    moveToCurrentTag()
  })
}

const closeAllTags = (view) => {
  tagsViewStore.delAllViews().then(({ visitedViews }) => {
    if (visitedViews.length === 0) {
      toLastView(visitedViews, view)
    }
  })
}

const handleDropdownCommand = (command) => {
  if (command === 'close') {
    closeSelectedTag(selectedDropdownTag.value)
  } else if (command === 'closeOthers') {
    closeOthersTags()
  } else if (command === 'closeLeft') {
    closeLeftTags()
  } else if (command === 'closeRight') {
    closeRightTags()
  } else if (command === 'closeAll') {
    closeAllTags(selectedDropdownTag.value)
  } else if (command === 'fullscreen') {
    // 已移除全屏功能
  }
}

const handleScroll = () => {
  updateArrowState()
}

const updateArrowState = () => {
  const scrollWrapper = scrollPaneRef.value?.wrapRef
  if (!scrollWrapper) return

  const { scrollLeft, scrollWidth, clientWidth } = scrollWrapper
  canScrollLeft.value = scrollLeft > 0
  canScrollRight.value = scrollLeft < scrollWidth - clientWidth - 1
}

const scrollLeft = () => {
  scrollPaneRef.value?.scrollByOffset(-200)
}

const scrollRight = () => {
  scrollPaneRef.value?.scrollByOffset(200)
}

const moveToCurrentTag = () => {
  nextTick(() => {
    for (const r of filterRoutes.value) {
      if (r.path === route.path) {
        if (scrollPaneRef.value) {
          scrollPaneRef.value.moveToTarget(r)
        }
        break
      }
    }
  })
}

// 添加当前路由到标签页
const addCurrentRouteToTags = () => {
  if (route.path && route.meta && route.meta.title) {
    // 排除一些不需要标签页的路由
    if (['/login', '/auth-redirect', '/404', '/401', '/lock'].includes(route.path)) {
      return
    }

    // 检查是否已经存在该标签页
    const exists = tagsViewStore.visitedViews.some(v => v.path === route.path)
    if (!exists) {
      tagsViewStore.addView(route)
    }
  }
}

onMounted(() => {
  // 添加当前路由到标签页
  addCurrentRouteToTags()

  // 添加滚动监听
  const scrollWrapper = scrollPaneRef.value?.wrapRef
  if (scrollWrapper) {
    scrollWrapper.addEventListener('scroll', updateArrowState)
  }
  // 初始更新箭头状态
  updateArrowState()
  // 添加全局点击事件监听，用于关闭右键菜单
  document.addEventListener('click', closeMenu)
})

// 监听路由变化，添加新路由到标签页
watch(() => route.path, () => {
  addCurrentRouteToTags()
  moveToCurrentTag()
})

// 监听 visible 变化，更新 selectedDropdownTag
watch(visible, (val) => {
  if (val) {
    selectedDropdownTag.value = selectedTag.value
  }
})

// 监听路由变化，更新 selectedDropdownTag
watch(() => route.path, () => {
  selectedDropdownTag.value = visitedViews.value.find(v => v.path === route.path) || {}
})
</script>

<style lang="scss" scoped>
$tags-bar-height: 34px;

.tags-view-container {
  height: $tags-bar-height;
  width: 100%;
  background: var(--tags-bg, #fff);
  border-bottom: 1px solid var(--tags-item-border, #d8dce5);
  display: flex;
  align-items: center;
  overflow: hidden;

  $btn-width: 28px;
  $btn-color: #71717a;
  $btn-hover-bg: #f0f2f5;
  $btn-hover-color: #303133;
  $btn-disabled-color: #c0c4cc;
  $divider: 1px solid var(--tags-item-border, #d8dce5);

  .tags-nav-btn {
    flex-shrink: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    width: $btn-width;
    height: $tags-bar-height;
    cursor: pointer;
    color: $btn-color;
    font-size: 13px;
    user-select: none;
    transition: background 0.15s, color 0.15s;

    &:hover:not(.disabled) {
      background: $btn-hover-bg;
      color: $btn-hover-color;
    }

    &.disabled {
      color: $btn-disabled-color;
      cursor: not-allowed;
    }

    &--left  { border-right: $divider; }
    &--right { border-left: $divider; }
  }

  .tags-view-wrapper {
    flex: 1;
    min-width: 0;
    height: 100%;

    .tags-view-item {
      display: inline-flex;
      align-items: center;
      position: relative;
      cursor: pointer;
      height: 26px;
      line-height: 26px;
      border: 1px solid var(--tags-item-border, #d8dce5);
      color: var(--tags-item-text, #495060);
      background: var(--tags-item-bg, #fff);
      padding: 0 8px;
      font-size: 12px;
      margin-left: 5px;
      border-radius: 3px;
      text-decoration: none;
      vertical-align: middle;
      padding-top: 2px !important;

      &:first-of-type { margin-left: 6px; }
      &:last-of-type  { margin-right: 15px; }
    }
  }

  &:not(.tags-view-container--chrome) .tags-view-wrapper .tags-view-item.active {
    background-color: #42b983;
    color: #fff;
    border-color: #42b983;

    &::before {
      content: '';
      background: #fff;
      display: inline-block;
      width: 8px;
      height: 8px;
      border-radius: 50%;
      position: relative;
      margin-right: 5px;
    }
  }

  &:not(.tags-view-container--chrome) .tags-view-wrapper .tags-view-item.active.has-icon::before {
    content: none !important;
  }

  .tags-action-dropdown {
    flex-shrink: 0;
    display: flex;
    align-items: center;
  }

  .tags-action-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: $btn-width;
    height: $tags-bar-height;
    cursor: pointer;
    color: $btn-color;
    font-size: 13px;
    border-left: $divider;
    user-select: none;
    transition: background 0.15s, color 0.15s;

    &:hover {
      background: $btn-hover-bg;
      color: $btn-hover-color;
    }
  }

  .tags-refresh-btn {
    width: 60px;
  }

  .contextmenu {
    margin: 0;
    background: var(--el-bg-color-overlay, #fff);
    z-index: 3000;
    position: fixed;
    list-style-type: none;
    padding: 5px 0;
    border-radius: 4px;
    font-size: 12px;
    font-weight: 400;
    color: var(--tags-item-text, #333);
    box-shadow: 2px 2px 3px 0 rgba(0, 0, 0, .3);
    border: 1px solid var(--el-border-color-light, #e4e7ed);

    li {
      margin: 0;
      padding: 7px 16px;
      cursor: pointer;

      &:hover {
        background: var(--tags-item-hover, #eee);
      }
    }
  }

  &.tags-view-container--chrome {
    --chrome-strip-bg: #ffffff;
    --chrome-strip-border: var(--el-border-color-lighter, #e4e7ed);
    --chrome-tab-active-bg: var(--el-color-primary-light-9);
    --chrome-tab-text: var(--el-text-color-regular, #606266);
    --chrome-tab-text-active: var(--el-color-primary);
    --chrome-wing-r: 10px;

    overflow: visible;
    background: var(--chrome-strip-bg);
    border-bottom: 1px solid var(--chrome-strip-border);
    align-items: flex-end;

    .tags-nav-btn {
      align-self: stretch;
      height: auto;
      min-height: $tags-bar-height;
      border-color: var(--chrome-strip-border);
    }

    .tags-action-btn {
      border-color: var(--chrome-strip-border);
    }

    .tags-view-wrapper {
      .tags-view-item {
        display: inline-flex !important;
        align-items: center;
        justify-content: center;
        position: relative;
        z-index: 1;
        height: 30px;
        min-height: 30px;
        margin: 0 0 -1px;
        padding: 0 12px;
        font-size: 13px;
        font-weight: 400;
        line-height: 1.2;
        border: none !important;
        border-radius: 0;
        background: transparent !important;
        color: var(--chrome-tab-text);
        padding-top: 0 !important;
        box-shadow: none !important;
        transition: background 0.12s ease, color 0.12s ease, border-radius 0.12s ease;

        &::before,
        &::after {
          content: '' !important;
          display: block !important;
          position: absolute;
          bottom: 0;
          width: var(--chrome-wing-r);
          height: var(--chrome-wing-r);
          margin: 0 !important;
          pointer-events: none;
          background: transparent !important;
          border-radius: 0 !important;
          transition: box-shadow 0.12s ease;
        }

        &::before {
          left: calc(-1 * var(--chrome-wing-r));
          border-bottom-right-radius: var(--chrome-wing-r) !important;
          box-shadow: none;
        }

        &::after {
          right: calc(-1 * var(--chrome-wing-r));
          border-bottom-left-radius: var(--chrome-wing-r) !important;
          box-shadow: none;
        }

        &:first-of-type {
          margin-left: 6px;
        }

        &:last-of-type {
          margin-right: 10px;
        }

        &:not(.active) + .tags-view-item:not(.active) {
          border-left: 1px solid var(--el-border-color-lighter, #e4e7ed);
          padding-left: 11px;
        }

        &:hover:not(.active) {
          background: var(--el-fill-color-light, #f5f7fa) !important;
          border-radius: 6px 6px 0 0;
          color: var(--el-text-color-primary, #303133);
        }

        &.active {
          height: 31px;
          min-height: 31px;
          padding: 0 14px;
          color: var(--chrome-tab-text-active) !important;
          font-weight: 500;
          background: var(--chrome-tab-active-bg) !important;
          border: none !important;
          border-radius: var(--chrome-wing-r) var(--chrome-wing-r) 0 0;
          box-shadow: 0 1px 4px rgba(0, 0, 0, 0.06);

          &::before {
            box-shadow: calc(var(--chrome-wing-r) * 0.5) calc(var(--chrome-wing-r) * 0.5) 0 calc(var(--chrome-wing-r) * 0.5) var(--chrome-tab-active-bg);
          }

          &::after {
            box-shadow: calc(var(--chrome-wing-r) * -0.5) calc(var(--chrome-wing-r) * 0.5) 0 calc(var(--chrome-wing-r) * 0.5) var(--chrome-tab-active-bg);
          }
        }
      }
    }
  }
}
</style>

<style lang="scss">
.tags-view-wrapper {
  .tags-view-item {
    .tags-close-btn {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      width: 16px;
      height: 16px;
      margin-left: 4px;
      border-radius: 50%;
      transition: all 0.3s cubic-bezier(0.645, 0.045, 0.355, 1);
      cursor: pointer;

      .el-icon-close {
        width: 1em;
        height: 1em;
        vertical-align: 0;
        line-height: 1;
        display: inline-flex;
        align-items: center;
        justify-content: center;
      }

      &:hover {
        background-color: var(--tags-close-hover, #b4bccc);

        .el-icon-close {
          color: #fff;
        }
      }
    }
  }
}

</style>
