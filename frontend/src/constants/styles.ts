import type { CSSProperties } from 'vue'

/**
 * 取件码高亮样式
 */
export const PICKUP_CODE_STYLE: CSSProperties = {
  color: '#18a058',
  fontSize: '14px',
  fontWeight: 'bold'
}

/**
 * 根据百分比获取阈值颜色
 * 用于滞留天数、使用率等场景
 * @param value 当前值
 * @param thresholds 阈值配置 [高阈值, 中阈值]，默认 [70, 30]
 * @returns 对应的颜色值
 */
export function getThresholdColor(
  value: number,
  thresholds: [number, number] = [70, 30]
): string {
  const [high, medium] = thresholds
  if (value >= high) return '#f5222d'  // 红色 - 高危
  if (value >= medium) return '#fa8c16' // 橙色 - 警告
  return '#52c41a'                      // 绿色 - 正常
}

/**
 * 获取滞留天数颜色
 * @param days 滞留天数
 */
export function getOverdueDaysColor(days: number): string {
  return getThresholdColor(days, [7, 3])
}

/**
 * 获取使用率颜色
 * @param rate 使用率百分比
 */
export function getUsageRateColor(rate: number): string {
  return getThresholdColor(rate, [90, 70])
}
