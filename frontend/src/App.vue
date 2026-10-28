<template>
  <router-view />
</template>

<script setup>
import { onMounted, nextTick } from 'vue'
import useSettingsStore from '@/store/modules/settings'
import { handleThemeStyle } from '@/utils/theme'
import {Call} from "@wailsio/runtime";

onMounted(() => {
  nextTick(() => {
    console.log('Wails v3 runtime 已就绪')
    console.log('Wails v3 Call:', typeof Call)
    handleThemeStyle(useSettingsStore().theme)
  })
})
</script>