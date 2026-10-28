<!-- 订单详情：订单汇总 + 明细（明细即打印任务，含失败原因与重试次数） -->
<template>
  <ElDialog v-model="dialogVisible" title="订单详情" width="900px" align-center>
    <div v-loading="loading">
      <ElDescriptions v-if="order" :column="3" border size="small" class="mb-4">
        <ElDescriptionsItem label="订单号">{{ order.orderNo }}</ElDescriptionsItem>
        <ElDescriptionsItem label="取件码">{{ order.pickupCode || '-' }}</ElDescriptionsItem>
        <ElDescriptionsItem label="实付">{{ money(order.payAmount) }} 元</ElDescriptionsItem>
        <ElDescriptionsItem label="文件数">{{ order.itemCount }}</ElDescriptionsItem>
        <ElDescriptionsItem label="总页数">{{ order.totalPages }}</ElDescriptionsItem>
        <ElDescriptionsItem label="印张数">{{ order.totalSheets }}</ElDescriptionsItem>
        <ElDescriptionsItem label="原价">{{ money(order.totalAmount) }} 元</ElDescriptionsItem>
        <ElDescriptionsItem label="折扣">
          {{ order.discountPercent === 100 ? '无折扣' : `${order.discountPercent / 10} 折` }}
        </ElDescriptionsItem>
        <ElDescriptionsItem label="已退">{{ money(order.refundAmount) }} 元</ElDescriptionsItem>
        <ElDescriptionsItem label="订单状态">{{ orderStatusText(order.orderStatus) }}</ElDescriptionsItem>
        <ElDescriptionsItem label="支付状态">{{ payStatusText(order.payStatus) }}</ElDescriptionsItem>
        <ElDescriptionsItem label="下单时间">{{ order.createdAt ?? '-' }}</ElDescriptionsItem>
        <ElDescriptionsItem label="备注" :span="3">{{ order.remark || '-' }}</ElDescriptionsItem>
      </ElDescriptions>

      <ElTable :data="items" border size="small" max-height="360">
        <ElTableColumn prop="fileName" label="文件" min-width="180" show-overflow-tooltip />
        <ElTableColumn label="参数" min-width="150">
          <template #default="{ row }">
            {{ row.paperType }} {{ row.colorType === 'color' ? '彩色' : '黑白' }}
            {{ row.printMode === 'double' ? '双面' : '单面' }}
          </template>
        </ElTableColumn>
        <ElTableColumn label="页码" width="110">
          <template #default="{ row }">{{ row.pageRange || '全部' }}</template>
        </ElTableColumn>
        <ElTableColumn prop="pageCount" label="页数" width="80" />
        <ElTableColumn prop="copies" label="份数" width="80" />
        <ElTableColumn prop="sheetCount" label="印张" width="80" />
        <ElTableColumn label="金额" width="100">
          <template #default="{ row }">{{ money(row.amount) }} 元</template>
        </ElTableColumn>
        <ElTableColumn label="状态" width="100">
          <template #default="{ row }">
            <ElTag :type="taskStatusType(row.status)" effect="plain">
              {{ taskStatusText(row.status) }}
            </ElTag>
          </template>
        </ElTableColumn>
        <ElTableColumn label="重试" width="70">
          <template #default="{ row }">{{ row.retryCount }}</template>
        </ElTableColumn>
        <ElTableColumn label="失败原因" min-width="160" show-overflow-tooltip>
          <template #default="{ row }">{{ row.failReason || '-' }}</template>
        </ElTableColumn>
      </ElTable>
    </div>

    <template #footer>
      <ElButton @click="dialogVisible = false">关闭</ElButton>
    </template>
  </ElDialog>
</template>

<script setup lang="ts">
  import { fetchGetPrintOrder } from '@/api/print/order'

  defineOptions({ name: 'PrintOrderDetailDialog' })

  interface Props {
    visible: boolean
    orderId: string | null
  }
  interface Emits {
    (e: 'update:visible', value: boolean): void
  }

  const props = defineProps<Props>()
  const emit = defineEmits<Emits>()

  const dialogVisible = computed({
    get: () => props.visible,
    set: (value) => emit('update:visible', value)
  })

  const loading = ref(false)
  const order = ref<Api.Print.OrderItem | null>(null)
  const items = ref<Api.Print.OrderTaskItem[]>([])

  const money = (value: number): string => (value / 100).toFixed(2)

  const ORDER_TEXT: Record<number, string> = {
    0: '待支付',
    1: '待打印',
    2: '打印中',
    3: '待取件',
    4: '已完成',
    5: '已取消',
    6: '退款中',
    7: '已退款',
    8: '超期关闭'
  }
  const PAY_TEXT: Record<number, string> = { 0: '未支付', 1: '已支付', 2: '退款中', 3: '已退款' }
  const TASK_TEXT: Record<number, string> = {
    0: '待打印',
    1: '已分配',
    2: '打印中',
    3: '已完成',
    4: '失败',
    5: '已取消'
  }
  const TASK_TYPE: Record<number, 'warning' | 'success' | 'danger' | 'info' | 'primary'> = {
    0: 'warning',
    1: 'warning',
    2: 'primary',
    3: 'success',
    4: 'danger',
    5: 'info'
  }

  const orderStatusText = (value: number): string => ORDER_TEXT[value] ?? '未知'
  const payStatusText = (value: number): string => PAY_TEXT[value] ?? '-'
  const taskStatusText = (value: number): string => TASK_TEXT[value] ?? '未知'
  const taskStatusType = (value: number): 'warning' | 'success' | 'danger' | 'info' | 'primary' =>
    TASK_TYPE[value] ?? 'info'

  watch(
    () => props.visible,
    async (visible) => {
      if (!visible || !props.orderId) return

      loading.value = true
      order.value = null
      items.value = []
      try {
        const detail = await fetchGetPrintOrder(props.orderId)
        order.value = detail.order
        items.value = detail.items
      } catch {
        dialogVisible.value = false
      } finally {
        loading.value = false
      }
    }
  )
</script>
