import router from './router'
import NProgress from 'nprogress'
import 'nprogress/nprogress.css'

NProgress.configure({ showSpinner: false })

router.beforeEach((to, from, next) => {
  NProgress.start()
  // 锁屏拦截：如果已锁定且目标不是锁屏页，强制跳转锁屏
  const isLock = JSON.parse(localStorage.getItem('screen-lock') || 'false')
  if (isLock && to.path !== '/lock') {
    next('/lock')
    return
  }
  next()
})

router.afterEach(() => {
  NProgress.done()
})
