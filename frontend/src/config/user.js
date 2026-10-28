import { reactive } from 'vue'
import defaultAvatar from '@/assets/logo/logo.png'
import { getStoreName, setStoreName, getSetting, setSetting } from '@/api'

// 统一用户配置，所有模块从此处读取
export const userConfig = reactive({
  storeName: '',
  avatar: defaultAvatar,
  secondaryPassword: ''
})

// 从数据库加载所有配置
export async function loadUserConfig() {
  const nameRes = await getStoreName()
  if (nameRes && nameRes.code === 0 && nameRes.data) {
    userConfig.storeName = nameRes.data
  }
  const avatar = await getSetting('avatar')
  if (avatar) userConfig.avatar = avatar
  const pwd = await getSetting('secondary_password')
  if (pwd) userConfig.secondaryPassword = pwd
}

// 保存指定配置项
export async function saveUserSetting(key, value) {
  const res = await setSetting(key, value)
  if (res && res.code === 0) {
    if (key === 'store_name') userConfig.storeName = value
    else if (key === 'avatar') userConfig.avatar = value
    else if (key === 'secondary_password') userConfig.secondaryPassword = value
    return true
  }
  return false
}
