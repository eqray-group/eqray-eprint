<!-- 钱包人工调账：正数加款、负数扣款，直接改余额并写流水 -->
<template>
  <ElDialog v-model="dialogVisible" title="钱包调账" width="480px" align-center :close-on-click-modal="false">
    <ElForm ref="formRef" :model="form" :rules="rules" label-width="90px">
      <ElFormItem label="会员 ID">
        <ElInput :model-value="memberId" disabled />
      </ElFormItem>
      <ElFormItem label="调账金额" prop="amount">
        <ElInputNumber v-model="form.amount" :min="-1000000" :max="1000000" :step="100" class="w-full" />
        <span class="ml-3 text-xs text-g-600">单位：分，正数加款、负数扣款</span>
      </ElFormItem>
      <ElFormItem label="调账原因" prop="reason">
        <ElInput v-model="form.reason" type="textarea" :rows="2" maxlength="200" placeholder="必填，会写入流水备注" />
      </ElFormItem>
    </ElForm>

    <template #footer>
      <ElButton @click="dialogVisible = false">取消</ElButton>
      <ElButton type="primary" :loading="submitting" @click="handleSubmit">确认调账</ElButton>
    </template>
  </ElDialog>
</template>

<script setup lang="ts">
  import { ElMessage } from 'element-plus'
  import type { FormInstance, FormRules } from 'element-plus'
  import { fetchAdjustWallet } from '@/api/pay/payment'

  defineOptions({ name: 'PayWalletAdjustDialog' })

  interface Props {
    visible: boolean
    memberId?: string
  }
  interface Emits {
    (e: 'update:visible', value: boolean): void
    (e: 'success'): void
  }

  const props = withDefaults(defineProps<Props>(), { memberId: '' })
  const emit = defineEmits<Emits>()

  const dialogVisible = computed({
    get: () => props.visible,
    set: (value) => emit('update:visible', value)
  })

  const formRef = ref<FormInstance>()
  const submitting = ref(false)

  const form = reactive({ amount: 0, reason: '' })

  const rules: FormRules = {
    amount: [
      {
        required: true,
        validator: (_rule, value: number, callback) => {
          if (!value) {
            callback(new Error('调账金额不能为 0'))
            return
          }
          callback()
        },
        trigger: 'blur'
      }
    ],
    reason: [{ required: true, message: '请填写调账原因', trigger: 'blur' }]
  }

  watch(
    () => props.visible,
    (visible) => {
      if (!visible) return

      Object.assign(form, { amount: 0, reason: '' })
      nextTick(() => formRef.value?.clearValidate())
    }
  )

  const handleSubmit = async (): Promise<void> => {
    if (!formRef.value || !props.memberId) return

    const valid = await formRef.value.validate().catch(() => false)
    if (!valid) return

    submitting.value = true
    try {
      const res = await fetchAdjustWallet(props.memberId, {
        amount: form.amount,
        reason: form.reason.trim()
      })
      ElMessage.success(`调账成功，当前余额 ${(res.balance / 100).toFixed(2)} 元`)
      dialogVisible.value = false
      emit('success')
    } catch {
      // 错误提示已由请求层统一弹出
    } finally {
      submitting.value = false
    }
  }
</script>
