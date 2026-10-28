<!-- 发起退款：余额支付立即到账，微信/支付宝按渠道处理 -->
<template>
  <ElDialog v-model="dialogVisible" title="发起退款" width="560px" align-center :close-on-click-modal="false">
    <ElAlert
      type="warning"
      :closable="false"
      show-icon
      title="退款不可撤销，金额不得超过该支付单剩余可退金额"
      class="mb-4"
    />

    <ElForm ref="formRef" :model="form" :rules="rules" label-width="100px">
      <ElFormItem label="支付单号" prop="paymentNo">
        <ElInput v-model="form.paymentNo" placeholder="请输入要退款的支付单号" maxlength="32" />
      </ElFormItem>
      <ElFormItem label="退款金额" prop="amount">
        <ElInputNumber v-model="form.amount" :min="1" :step="100" class="w-full" />
        <span class="ml-3 text-xs text-g-600">单位：分</span>
      </ElFormItem>
      <ElFormItem label="退款原因" prop="reason">
        <ElInput v-model="form.reason" type="textarea" :rows="2" maxlength="255" placeholder="必填，会记录在退款单上" />
      </ElFormItem>
    </ElForm>

    <template #footer>
      <ElButton @click="dialogVisible = false">取消</ElButton>
      <ElButton type="primary" :loading="submitting" @click="handleSubmit">确认退款</ElButton>
    </template>
  </ElDialog>
</template>

<script setup lang="ts">
  import { ElMessage } from 'element-plus'
  import type { FormInstance, FormRules } from 'element-plus'
  import { fetchCreateRefund } from '@/api/pay/refund'

  defineOptions({ name: 'PayRefundDialog' })

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

  const formRef = ref<FormInstance>()
  const submitting = ref(false)

  const form = reactive({ paymentNo: '', amount: 0, reason: '' })

  const rules: FormRules = {
    paymentNo: [{ required: true, message: '请输入支付单号', trigger: 'blur' }],
    amount: [{ required: true, message: '请输入退款金额', trigger: 'blur' }],
    reason: [{ required: true, message: '请填写退款原因', trigger: 'blur' }]
  }

  watch(
    () => props.visible,
    (visible) => {
      if (!visible) return

      Object.assign(form, { paymentNo: '', amount: 0, reason: '' })
      nextTick(() => formRef.value?.clearValidate())
    }
  )

  const handleSubmit = async (): Promise<void> => {
    if (!formRef.value) return

    const valid = await formRef.value.validate().catch(() => false)
    if (!valid) return

    submitting.value = true
    try {
      const refund = await fetchCreateRefund({
        paymentNo: form.paymentNo.trim(),
        amount: form.amount,
        reason: form.reason.trim()
      })
      ElMessage.success(refund.status === 2 ? '退款成功' : '退款已提交，等待渠道处理')
      dialogVisible.value = false
      emit('success')
    } catch {
      // 错误提示已由请求层统一弹出
    } finally {
      submitting.value = false
    }
  }
</script>
