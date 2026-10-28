<template>
  <div class="navbar" :class="'nav' + settingsStore.navType">
    <hamburger id="hamburger-container" :is-active="appStore.sidebar.opened" class="hamburger-container" @toggleClick="toggleSideBar" />
    <breadcrumb v-if="settingsStore.navType == 1" id="breadcrumb-container" class="breadcrumb-container" />
    <top-nav v-if="settingsStore.navType == 2" id="topmenu-container" class="topmenu-container" />
    <template v-if="settingsStore.navType == 3">
      <logo v-show="settingsStore.sidebarLogo" :collapse="false"></logo>
      <top-bar id="topbar-container" class="topbar-container" />
    </template>

    <div class="right-menu">
      <template v-if="appStore.device !== 'mobile'">
        <header-search id="header-search" class="right-menu-item" />

        <div class="right-menu-item hover-effect theme-switch-wrapper" @click="toggleTheme">
          <svg-icon v-if="settingsStore.isDark" icon-class="sunny" />
          <svg-icon v-if="!settingsStore.isDark" icon-class="moon" />
        </div>

      </template>

      <el-dropdown @command="handleCommand" class="avatar-container right-menu-item hover-effect" trigger="hover">
        <div class="avatar-wrapper">
          <img :src="userStore.avatar || defaultAvatar" class="user-avatar" />
          <span class="user-storename">{{ userStore.storeName || '' }}</span>
        </div>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="shopInfo">
              <el-icon><Shop /></el-icon>
              <span>店铺信息</span>
            </el-dropdown-item>
            <el-dropdown-item command="setLayout" v-if="settingsStore.showSettings">
              <el-icon><Setting /></el-icon>
              <span>布局设置</span>
            </el-dropdown-item>
            <el-dropdown-item command="lockScreen">
              <el-icon><Lock /></el-icon>
              <span>锁定屏幕</span>
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>

    <!-- 系统设置弹窗 -->
    <el-dialog :title="title" v-model="shopDialogVisible" width="420px" draggable >
      <el-form ref="shopFormRef" :model="shopForm" :rules="shopRules" label-position="top">
        <el-form-item label="店铺名称" prop="storeName">
          <el-input v-model="shopForm.storeName" placeholder="请输入店铺名称" />
        </el-form-item>
<el-form-item label="锁屏密码">
          <el-input v-model="shopForm.oldPassword" placeholder="输入旧密码（默认密码:123456）" show-password />
          <el-input v-model="shopForm.secondaryPassword" placeholder="输入新密码（留空不修改）" show-password style="margin-top: 8px;" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="shopDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitShopForm">保存</el-button>
      </template>
    </el-dialog>

    <UpdateDialog ref="updateDialogRef" />
  </div>
</div>
</template>

<script setup>
import { ElMessageBox, ElMessage } from 'element-plus'
import UpdateDialog from '@/components/UpdateDialog/index.vue'
import Breadcrumb from '@/components/Breadcrumb'
import TopNav from './TopNav'
import TopBar from './TopBar'
import Logo from './Sidebar/Logo'
import Hamburger from '@/components/Hamburger'
import HeaderSearch from '@/components/HeaderSearch'
import useAppStore from '@/store/modules/app'
import useLockStore from '@/store/modules/lock'
import useSettingsStore from '@/store/modules/settings'
import { onMounted, onBeforeUnmount, ref, reactive } from 'vue'
import { useRoute, useRouter } from "vue-router";
import { userConfig, loadUserConfig, saveUserSetting } from '@/config/user';
import { checkVersion } from '@/api'
import defaultAvatar from '@/assets/logo/logo.png'
import {Close, Setting, Shop} from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const appStore = useAppStore()
const lockStore = useLockStore()
const settingsStore = useSettingsStore()

const userStore = userConfig

// 店铺信息弹窗
const shopDialogVisible = ref(false)
const shopFormRef = ref(null)
const title = ref('店铺信息')
const shopForm = reactive({ storeName: '', avatar: '', secondaryPassword: '', oldPassword: '' })
const shopRules = {
  storeName: [{ required: true, message: '请输入店铺名称', trigger: 'blur' }]
}

function openShopDialog() {
  shopForm.storeName = userStore.storeName
  shopForm.avatar = userStore.avatar !== defaultAvatar ? userStore.avatar : ''
  shopForm.secondaryPassword = ''
  shopForm.oldPassword = ''
  shopDialogVisible.value = true
}

async function submitShopForm() {
  const valid = await shopFormRef.value.validate().catch(() => false)
  if (!valid) return
  // 如果填了新密码，必须验证旧密码
  if (shopForm.secondaryPassword) {
    if (!shopForm.oldPassword) {
      ElMessage.error('请输入旧密码')
      return
    }
    if (shopForm.oldPassword !== userStore.secondaryPassword) {
      ElMessage.error('旧密码不正确')
      return
    }
  }
  let ok = true
  ok = (await saveUserSetting('store_name', shopForm.storeName.trim())) && ok
  ok = (await saveUserSetting('avatar', shopForm.avatar.trim())) && ok
  if (shopForm.secondaryPassword) {
    ok = (await saveUserSetting('secondary_password', shopForm.secondaryPassword.trim())) && ok
  }
  if (ok) {
    if (shopForm.avatar.trim()) userStore.avatar = shopForm.avatar.trim()
    else userStore.avatar = defaultAvatar
    const changed = shopForm.secondaryPassword ? `，新密码：${shopForm.secondaryPassword}` : ''
    ElMessage.success('保存成功' + changed)
    shopDialogVisible.value = false
  } else {
    ElMessage.error('保存失败')
  }
}

function toggleSideBar() {
  appStore.toggleSideBar()
}

function handleCommand(command) {
  switch (command) {
    case "shopInfo":
      openShopDialog()
      break
    case "setLayout":
      setLayout()
      break
    case "lockScreen":
      lockScreen()
      break
    case "logout":
      logout()
      break
    default:
      break
  }
}

function logout() {
  ElMessageBox.confirm('确定注销并退出系统吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    // 直接跳转，不调用登出接口
    location.href = '/'
  }).catch(() => { })
}

const emits = defineEmits(['setLayout'])
function setLayout() {
  emits('setLayout')
}

function lockScreen() {
  const currentPath = route.fullPath
  lockStore.lockScreen(currentPath)
  router.replace('/lock')
}

const updateDialogRef = ref(null)

function onShowUpdateDialog(e) {
  updateDialogRef.value?.open(e.detail)
}

onMounted(async () => {
  await loadUserConfig()
  window.addEventListener('show-update-dialog', onShowUpdateDialog)
  const info = await checkVersion()
  if (info && info.isNew) {
    updateDialogRef.value?.open(info)
  }
})

onBeforeUnmount(() => {
  window.removeEventListener('show-update-dialog', onShowUpdateDialog)
})

async function toggleTheme(event) {
  const x = event?.clientX || window.innerWidth / 2
  const y = event?.clientY || window.innerHeight / 2
  const wasDark = settingsStore.isDark

  if (!document.startViewTransition) {
    settingsStore.toggleTheme()
    return
  }

  const transition = document.startViewTransition(() => {
    settingsStore.toggleTheme()
  })
  await transition.ready

  const endRadius = Math.hypot(
    Math.max(x, window.innerWidth - x),
    Math.max(y, window.innerHeight - y)
  )
  document.documentElement.animate(
    {
      clipPath: !wasDark
        ? [`circle(${endRadius}px at ${x}px ${y}px)`, `circle(0px at ${x}px ${y}px)`]
        : [`circle(0px at ${x}px ${y}px)`, `circle(${endRadius}px at ${x}px ${y}px)`]
    }, {
      duration: 500,
      easing: 'cubic-bezier(0.4, 0, 0.2, 1)',
      pseudoElement: !wasDark ? '::view-transition-old(root)' : '::view-transition-new(root)'
    }
  )
  await transition.finished
}
</script>

<style lang='scss' scoped>
.navbar.nav3 {
  .hamburger-container {
    display: none !important;
  }
}

.navbar {
  height: 50px;
  overflow: hidden;
  position: relative;
  background: var(--navbar-bg);
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
  display: flex;
  align-items: center;
  // padding: 0 8px;
  box-sizing: border-box;

  .hamburger-container {
    line-height: 46px;
    height: 100%;
    cursor: pointer;
    transition: background 0.3s;
    -webkit-tap-highlight-color: transparent;
    display: flex;
    align-items: center;
    flex-shrink: 0;
    margin-right: 8px;

    &:hover {
      background: rgba(0, 0, 0, 0.025);
    }
  }

  .breadcrumb-container {
    flex-shrink: 0;
  }

  .topmenu-container {
    position: absolute;
    left: 50px;
  }

  .topbar-container {
    flex: 1;
    min-width: 0;
    display: flex;
    align-items: center;
    overflow: hidden;
    margin-left: 8px;
  }

  .right-menu {
    height: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    margin-left: auto;

    &:focus {
      outline: none;
    }

    .right-menu-item {
      display: flex;
      align-items: center;
      padding: 0 8px;
      height: 100%;
      font-size: 18px;
      color: #5a5e66;

      &.hover-effect {
        cursor: pointer;
        transition: background 0.3s;

        &:hover {
          background: rgba(0, 0, 0, 0.025);
        }
      }

      &.theme-switch-wrapper {
        display: flex;
        align-items: center;
        margin-right: 4px;

        svg {
          transition: transform 0.3s;

          &:hover {
            transform: scale(1.15);
          }
        }
      }
    }

    .avatar-container {
      padding-right: 8px;

      .avatar-wrapper {
        display: flex;
        align-items: center;
        gap: 6px;

        .user-avatar {
          cursor: pointer;
          width: 30px;
          height: 30px;
          border-radius: 6px;
          user-select: none;
        }

        .user-storename{
          font-size: 14px;
          font-weight: 540;
          color: var(--el-text-color-primary);
          user-select: none;
        }

        i {
          cursor: pointer;
          position: absolute;
          right: -20px;
          top: 25px;
          font-size: 12px;
        }
      }
    }
  }
}
.dialog-header {
  padding: 0;
  display: flex;
  justify-content: flex-start;
  align-items: flex-start;
  background-color: red;
}

</style>