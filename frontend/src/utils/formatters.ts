/**
 * 格式化日期时间
 * @param dateStr 日期字符串
 * @returns 格式化后的日期时间字符串
 */
export function formatDateTime(dateStr?: string | null): string {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString('zh-CN')
}

/**
 * 格式化金额
 * @param amount 金额数值
 * @param precision 小数位数，默认2
 * @returns 格式化后的金额字符串
 */
export function formatAmount(amount?: number | null, precision: number = 2): string {
  if (amount === undefined || amount === null) return '-'
  return `¥${amount.toFixed(precision)}`
}

/**
 * 格式化重量
 * @param weight 重量数值
 * @returns 格式化后的重量字符串
 */
export function formatWeight(weight?: number | null): string {
  if (weight === undefined || weight === null) return '-'
  return `${weight} kg`
}
