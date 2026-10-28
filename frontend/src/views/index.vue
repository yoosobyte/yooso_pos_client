<template>
  <div class="dashboard-container">
    <!-- 统计卡片 -->
    <el-row :gutter="16">
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card clickable" @click="$router.push('/goods')">
          <div class="stat-icon" style="background: #409eff22; color: #409eff;">
            <svg-icon icon-class="product2" />
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ dashboard.totalProducts }}</div>
            <div class="stat-label">商品总数</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card clickable" @click="$router.push('/stock')">
          <div class="stat-icon" style="background: #f56c6c22; color: #f56c6c;">
            <svg-icon icon-class="wtip" />
          </div>
          <div class="stat-info">
            <div class="stat-value" :class="dashboard.lowStockCount>0?'danger':'info'">{{ dashboard.lowStockCount>0?dashboard.lowStockCount:'正常' }}</div>
            <div class="stat-label">低库存预警</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card clickable" @click="$router.push('/sales/record')">
          <div class="stat-icon" style="background: #409eff22; color: #409eff;">
            <svg-icon icon-class="plan" />
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ dashboard.todaySalesCount }}单</div>
            <div class="stat-label">累计订单</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card clickable" @click="$router.push('/report/sales')">
          <div class="stat-icon" style="background: #67c23a22; color: #67c23a;">
            <svg-icon icon-class="money2" />
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ fmtMoney(dashboard.todaySalesAmount) }}</div>
            <div class="stat-label">今日销售额</div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 快捷入口 + 低库存 -->
    <el-row :gutter="16" style="margin-top: 16px;">
      <el-col :span="8">
        <el-card shadow="hover" header="快捷操作" class="show-card">
          <div class="quick-actions">
            <el-button type="primary" size="large" @click="$router.push('/sales/cashier')">
              <svg-icon icon-class="printer" /> 收银开单
            </el-button>
            <el-button type="success" size="large" @click="$router.push('/purchase')">
              <svg-icon icon-class="delivery" /> 采购入库
            </el-button>
            <el-button type="warning" size="large" @click="$router.push('/goods')">
              <svg-icon icon-class="product2" /> 商品管理
            </el-button>
            <el-button type="info" size="large" @click="$router.push('/stock')">
              <svg-icon icon-class="inventory" /> 库存查询
            </el-button>
          </div>
        </el-card>
      </el-col>
      <el-col :span="16">
        <el-card shadow="hover" header="低库存预警商品" class="show-card">
          <el-table
              :data="lowStockList"
              stripe
              border
              size="small"
              min-height="300"
          >
            <el-table-column prop="name" label="商品名称" min-width="140" align="center" />
            <el-table-column label="规格" min-width="120" align="center">
              <template #default="{ row }">{{ row.spec || '-' }}</template>
            </el-table-column>
            <el-table-column prop="stockQty" label="当前库存" width="110" align="center">
              <template #default="{ row }">
                {{ row.stockQty }} {{ row.unit?.name || '' }}
              </template>
            </el-table-column>
            <el-table-column prop="minStock" label="最低库存" width="110" align="center">
              <template #default="{ row }">
                {{ row.minStock }} {{ row.unit?.name || '' }}
              </template>
            </el-table-column>
            <el-table-column label="状态" width="100" align="center">
              <template #default="{ row }">
                <el-tag v-if="row.stockQty <= 0" type="danger">缺货</el-tag>
                <el-tag v-else type="warning">低于预警</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="100" align="center">
              <template #default="{ row }">
                <el-button type="primary" size="small" @click="$router.push('/purchase')">
                  <svg-icon icon-class="purchase" style="margin-right: 5px;"/> 补货
                </el-button>
              </template>
            </el-table-column>

            <!-- 空状态用插槽 -->
            <template #empty>
              <el-empty description="库存健康，无需补货" />
            </template>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getDashboard, getLowStockProducts } from '@/api'
import { WarningFilled } from '@element-plus/icons-vue'

const dashboard = ref({
  totalProducts: 0,
  totalCustomers: 0,
  totalSuppliers: 0,
  todaySalesCount: 0,
  todaySalesAmount: 0,
  monthSalesAmount: 0,
  lowStockCount: 0
})
const lowStockList = ref([])

function fmtMoney(val) {
  const num = parseFloat(val)
  if (isNaN(num)) return '0'
  return num.toFixed(2).replace(/\.?0+$/, '')
}

onMounted(async () => {
  try {
    const res = await getDashboard()
    if (res && res.code === 0) dashboard.value = res.data
  } catch (e) {
    console.error('获取仪表盘数据失败:', e)
  }

  try {
    const res = await getLowStockProducts()
    if (res && res.code === 0) lowStockList.value = res.data
  } catch (e) {
    console.error('获取低库存产品失败:', e)
  }
})
</script>

<style scoped>
.dashboard-container {
  padding: 20px;
  background-color: #FFFFFF;
  min-height: 100%;
}

html.dark .dashboard-container {
  background-color: var(--el-bg-color);
}

.stat-card {
  cursor: default;
  transition: background-color 0.3s, border-color 0.3s;
}

.clickable {
  cursor: pointer;
}

.stat-card :deep(.el-card__body) {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 16px;
}

html.dark .stat-card {
  background-color: var(--el-bg-color-overlay);
  border-color: var(--el-border-color-darker);
}

.stat-icon {
  width: 44px;
  height: 44px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  font-size: 22px;
  transition: background 0.3s;
}

html.dark .stat-icon {
  backdrop-filter: brightness(1.2);
}

.stat-value {
  font-size: 18px;
  font-weight: 700;
  color: var(--el-text-color-primary);
}

.stat-value.danger {
  color: #f56c6c;
}

html.dark .stat-value.danger {
  color: #f78989;
}


.stat-value.info {
  font-size: 16px;
  color: var(--el-text-color-secondary);
}

html.dark .stat-value.info {
  font-size: 16px;
  color: var(--el-text-color-secondary);
}

.stat-label {
  font-size: 12px;
  color: var(--el-text-color-secondary);
  margin-top: 2px;
}

.quick-actions {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  gap: 12px;
}

.quick-actions .el-button .svg-icon {
  margin-right: 6px;
}

.quick-actions .el-button {
  min-width: 100%;
  margin-left: 0;
}

html.dark :deep(.el-card__header) {
  background-color: var(--el-bg-color-overlay);
  border-bottom-color: var(--el-border-color-darker);
  color: var(--el-text-color-primary);
}

html.dark :deep(.el-table) {
  --el-table-bg-color: var(--el-bg-color-overlay);
  --el-table-tr-bg-color: var(--el-bg-color-overlay);
  --el-table-header-bg-color: var(--el-bg-color);
  --el-table-border-color: var(--el-border-color-darker);
  --el-table-text-color: var(--el-text-color-primary);
  --el-table-header-text-color: var(--el-text-color-primary);
}

html.dark :deep(.el-empty__description) {
  color: var(--el-text-color-secondary);
}
</style>