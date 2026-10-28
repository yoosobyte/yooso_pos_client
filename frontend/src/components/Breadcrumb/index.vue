<template>
  <el-breadcrumb class="app-breadcrumb" separator="/">
    <transition-group name="breadcrumb">
      <el-breadcrumb-item v-for="(item, index) in levelList" :key="item.path">
        <span v-if="index === levelList.length - 1" class="no-redirect">{{ item.meta?.title }}</span>
        <a v-else @click.prevent="handleLink(item)">{{ item.meta?.title }}</a>
      </el-breadcrumb-item>
    </transition-group>
  </el-breadcrumb>
</template>

<script setup>
import {useRoute, useRouter} from 'vue-router'
import { defineComponent, ref,watchEffect  } from 'vue'
const route = useRoute()
const router = useRouter()
const levelList = ref([])

function getBreadcrumb() {
  // 直接从 route.matched 获取匹配的路由
  let matched = route.matched.filter(item => item.meta && item.meta.title)

  // 如果匹配结果为空，从当前路由获取
  if (matched.length === 0 && route.meta?.title) {
    matched = [{ path: route.path, meta: route.meta }]
  }

  levelList.value = matched
}

function handleLink(item) {
  const { redirect, path } = item
  if (redirect) {
    router.push(redirect)
    return
  }
  router.push(path)
}

watchEffect(() => {
  if (route.path.startsWith('/redirect/')) {
    return
  }
  getBreadcrumb()
})
getBreadcrumb()
</script>

<style lang='scss' scoped>
.app-breadcrumb.el-breadcrumb {
  display: inline-block;
  font-size: 14px;
  line-height: 50px;

  .no-redirect {
    color: #97a8be;
    cursor: text;
  }
}
</style>
