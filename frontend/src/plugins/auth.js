// ========== 权限校验（已简化，全部放行） ==========

/**
 * 字符权限校验
 * @param {Array} value 校验值
 * @returns {Boolean}
 */
export function checkPermi(value) {
  // 直接放行，不做权限校验
  return true
}

/**
 * 角色权限校验
 * @param {Array} value 校验值
 * @returns {Boolean}
 */
export function checkRole(value) {
  // 直接放行，不做权限校验
  return true
}

// ========== 默认导出（兼容若依的 import auth from '@/plugins/auth'） ==========
export default {
  hasPermi: checkPermi,
  hasRole: checkRole,
  hasPermiOr: checkPermi,
  hasRoleOr: checkRole
}