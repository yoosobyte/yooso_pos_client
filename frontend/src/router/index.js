import { createWebHistory, createRouter } from 'vue-router'
import Layout from '@/layout'
import {Printer} from "@element-plus/icons-vue";

export const constantRoutes = [
    {
        path: "/:pathMatch(.*)*",
        component: () => import('@/views/error/404'),
        hidden: true
    },
    {
        path: '/401',
        component: () => import('@/views/error/401'),
        hidden: true
    },
    {
        path: '/lock',
        component: () => import('@/views/lock'),
        hidden: true,
        meta: { title: '锁定屏幕' }
    },
    // 首页
    {
        path: '/',
        component: Layout,
        redirect: '/index',
        children: [
            {
                path: '/index',
                component: () => import('@/views/index.vue'),
                name: 'Index',
                meta: { title: '工作台', icon: 'server' }
            }
        ]
    },
    {
        path: '/cashier',
        component: Layout,
        redirect: '/sales/cashier',
        children: [
            {
                path: '/sales/cashier',
                component: () => import('@/views/sales/cashier.vue'),
                name: 'Cashier',
                meta: {
                    title: '收银台',
                    icon: 'printer'
                }
            }
        ]
    },
    // 商品管理
    {
        path: '/goods',
        component: Layout,
        redirect: '/goods/list',
        meta: { title: '商品管理', icon: 'completedInventory' },
        children: [
            {
                path: 'list',
                component: () => import('@/views/goods/index.vue'),
                name: 'Goods',
                meta: { title: '商品列表', icon: 'product2' }
            },
            {
                path: 'category',
                component: () => import('@/views/goods/category.vue'),
                name: 'GoodsCategory',
                meta: { title: '商品分类', icon: 'tree' }
            }
        ]
    },
    // 库存管理
    {
        path: '/stock',
        component: Layout,
        redirect: '/stock/overview',
        meta: { title: '库存管理', icon: 'inventory' },
        children: [
            {
                path: 'overview',
                component: () => import('@/views/stock/index.vue'),
                name: 'Stock',
                meta: { title: '库存总览', icon: 'check' }
            },
            {
                path: 'log',
                component: () => import('@/views/stock/log.vue'),
                name: 'StockLog',
                meta: { title: '出入库明细', icon: 'log' }
            }
        ]
    },
    // 采购管理
    {
        path: '/purchase',
        component: Layout,
        redirect: '/purchase/input',
        meta: { title: '采购管理', icon: 'flowRecord' },
        children: [
            {
                path: 'input',
                component: () => import('@/views/purchase/index.vue'),
                name: 'Purchase',
                meta: { title: '采购入库', icon: 'delivery' }
            },
            {
                path: 'record',
                component: () => import('@/views/purchase/record.vue'),
                name: 'PurchaseRecord',
                meta: { title: '采购记录', icon: 'plan' }
            },
            {
                path: 'supplier',
                component: () => import('@/views/purchase/supplier.vue'),
                name: 'Supplier',
                meta: { title: '供应商管理', icon: 'factory' }
            }
        ]
    },
    // 销售管理
    {
        path: '/sales',
        component: Layout,
        redirect: '/sales/cashier',
        meta: { title: '销售管理', icon: 'wallet' },
        children: [
            {
                path: 'record',
                component: () => import('@/views/sales/record.vue'),
                name: 'SalesRecord',
                meta: { title: '销售订单', icon: 'erpDict' }
            },
            {
                path: 'customer',
                component: () => import('@/views/sales/customer.vue'),
                name: 'Customer',
                meta: { title: '客户管理', icon: 'customer' }
            }
        ]
    },
    // 报表分析
    {
        path: '/report',
        component: Layout,
        redirect: '/report/sales',
        meta: { title: '报表分析', icon: 'usecount' },
        children: [
            {
                path: 'sales',
                component: () => import('@/views/report/sales.vue'),
                name: 'SalesReport',
                meta: { title: '销售报表', icon: 'excel' }
            },
            {
                path: 'purchase',
                component: () => import('@/views/report/purchase.vue'),
                name: 'PurchaseReport',
                meta: { title: '采购报表', icon: 'build' }
            }
        ]
    },
    {
        path: '/website',
        component: Layout,
        redirect: '/website/index',
        children: [
            {
                path: '/website/index',
                component: () => import('@/views/website/index.vue'),
                name: 'Website',
                meta: { title: '悦速官网', icon: 'promotion' }
            }
        ]
    },
]

export const dynamicRoutes = []

const router = createRouter({
    history: createWebHistory(),
    routes: constantRoutes,
    scrollBehavior(to, from, savedPosition) {
        if (savedPosition) {
            return savedPosition
        }
        return { top: 0 }
    }
})

export default router
