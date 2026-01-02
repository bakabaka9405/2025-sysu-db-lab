import { computed, type ComputedRef } from 'vue'

/**
 * 分页配置 Props 类型
 */
export interface PaginationProps {
  page: number
  pageSize: number
  total: number
  onChange: (page: number) => void
  onSizeChange: (size: number) => void
}

/**
 * NaiveUI DataTable 分页配置类型
 */
export interface TablePaginationConfig {
  page: number
  pageSize: number
  itemCount: number
  showSizePicker: boolean
  pageSizes: number[]
  onUpdatePage: (page: number) => void
  onUpdatePageSize: (size: number) => void
  prefix: () => string
}

/**
 * 创建表格分页配置的 composable
 * @param getPagination 获取分页 props 的函数
 * @returns 分页配置计算属性
 */
export function usePaginationConfig(
  getPagination: () => PaginationProps | false | undefined
): ComputedRef<TablePaginationConfig | false> {
  return computed(() => {
    const pagination = getPagination()
    if (!pagination) return false

    const { page, pageSize, total, onChange, onSizeChange } = pagination
    return {
      page,
      pageSize,
      itemCount: total,
      showSizePicker: true,
      pageSizes: [10, 20, 50],
      onUpdatePage: onChange,
      onUpdatePageSize: onSizeChange,
      prefix: () => `共 ${total} 条`
    }
  })
}
