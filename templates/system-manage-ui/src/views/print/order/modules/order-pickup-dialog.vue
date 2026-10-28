<!-- 取件核销：门店 + 6 位取件码，仅待取件订单可核销 -->
<template>
  <ElDialog v-model="dialogVisible" title="取件核销" width="460px" align-center :close-on-click-modal="false">
    <ElAlert
      type="info"
      :closable="false"
      show-icon
      title="核对用户取件码后核销，核销后订单转「已完成」且不可撤销"
      class="mb-4"
    />

    <ElForm label-width="90px">
      <ElFormItem label="门店" required>
        <ElSelect v-model="form.printShopId" class="w-full" filterable placeholder="请选择门店">
          <ElOption v-for="shop in shops" :key="shop.id" :label="shop.name" :value="shop.id" />
        </ElSelect>
      </ElFormItem>
      <ElFormItem label="取件码" required>
        <ElInput v-model="form.code" maxlength="16" placeholder="6 位数字" @keyup.enter="handleSubmit" />
      </ElFormItem>
    </ElForm>

    <template #footer>
      <ElButton @click="dialogVisible = false">取消</ElButton>
      <ElButton type="primary" :loading="submitting" @click="handleSubmit">核销</ElButton>
    </template>
  </ElDialog>
</template>

<script setup lang="ts">
  import { ElMessage } from 'element-plus'
  import { fetchPickupPrintOrder } from '@/api/print/order'
  import { fetchGetPrintShopList } from '@/api/print/shop'

  defineOptions({ name: 'PrintOrderPickupDialog' })

  interface Props {
    visible: boolean
  }
  interface Emits {
    (e: 'update:visible', value: boolean): void
    (e: 'success'): void
  }

  const props = defineProps<Props>()
  const emit = defineEmits<Emits>()

  const dialogVisible = computed({
    get: () => props.visible,
    set: (value) => emit('update:visible', value)
  })

  const submitting = ref(false)
  const shops = ref<Api.Print.ShopItem[]>([])
  const form = reactive<Api.Print.PickupParams>({ printShopId: '', code: '' })

  watch(
    () => props.visible,
    async (visible) => {
      if (!visible) return

      form.code = ''

      if (shops.value.length === 0) {
        try {
          const res = await fetchGetPrintShopList({ pageNo: 1, pageSize: 200 })
          shops.value = res.list
        } catch {
          shops.value = []
        }
      }
      form.printShopId = form.printShopId || shops.value[0]?.id || ''
    }
  )

  const handleSubmit = async (): Promise<void> => {
    if (!form.printShopId || !form.code.trim()) {
      ElMessage.warning('请选择门店并输入取件码')
      return
    }

    submitting.value = true
    try {
      const order = await fetchPickupPrintOrder({
        printShopId: form.printShopId,
        code: form.code.trim()
      })
      ElMessage.success(`核销成功：${order.orderNo}`)
      dialogVisible.value = false
      emit('success')
    } catch {
      // 错误提示已由请求层统一弹出
    } finally {
      submitting.value = false
    }
  }
</script>
