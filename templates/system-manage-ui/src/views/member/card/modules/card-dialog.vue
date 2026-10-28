<!-- 会员卡编辑弹窗：卡号由服务端生成，这里只维护等级、折扣与有效期 -->
<template>
  <ElDialog
    v-model="dialogVisible"
    title="编辑会员卡"
    width="560px"
    align-center
    :close-on-click-modal="false"
  >
    <ElForm ref="formRef" :model="form" :rules="rules" label-width="100px">
      <ElFormItem label="卡号">
        <ElInput :model-value="cardData?.cardNo ?? '-'" disabled />
      </ElFormItem>
      <ElFormItem label="卡等级" prop="cardLevel">
        <ElSelect v-model="form.cardLevel" class="w-full">
          <ElOption label="普通卡 normal" value="normal" />
          <ElOption label="银卡 silver" value="silver" />
          <ElOption label="金卡 gold" value="gold" />
          <ElOption label="黑金卡 black" value="black" />
        </ElSelect>
      </ElFormItem>
      <ElFormItem label="卡名称" prop="cardName">
        <ElInput v-model="form.cardName" placeholder="如：金卡会员" maxlength="64" />
      </ElFormItem>
      <ElFormItem label="折扣百分比" prop="discountPercent">
        <ElInputNumber v-model="form.discountPercent" :min="1" :max="100" :step="1" />
        <span class="ml-3 text-xs text-g-600">100 表示无折扣，95 表示 95 折</span>
      </ElFormItem>
      <ElFormItem label="状态" prop="status">
        <ElSelect v-model="form.status" class="w-full">
          <ElOption label="正常" :value="1" />
          <ElOption label="停用" :value="0" />
          <ElOption label="已过期" :value="2" />
        </ElSelect>
      </ElFormItem>
      <ElFormItem label="到期时间">
        <ElDatePicker
          v-model="form.expireTime"
          type="datetime"
          value-format="YYYY-MM-DD HH:mm:ss"
          placeholder="留空表示永久有效"
          class="w-full"
        />
      </ElFormItem>
      <ElFormItem label="备注">
        <ElInput v-model="form.remark" type="textarea" :rows="2" maxlength="255" />
      </ElFormItem>
    </ElForm>

    <template #footer>
      <ElButton @click="dialogVisible = false">取消</ElButton>
      <ElButton type="primary" :loading="submitting" @click="handleSubmit">保存</ElButton>
    </template>
  </ElDialog>
</template>

<script setup lang="ts">
  import { ElMessage } from 'element-plus'
  import type { FormInstance, FormRules } from 'element-plus'
  import { fetchSaveMemberCard } from '@/api/member/card'

  defineOptions({ name: 'MemberCardDialog' })

  interface Props {
    visible: boolean
    cardData?: Api.Member.CardItem | null
  }
  interface Emits {
    (e: 'update:visible', value: boolean): void
    (e: 'success'): void
  }

  const props = withDefaults(defineProps<Props>(), { cardData: null })
  const emit = defineEmits<Emits>()

  const dialogVisible = computed({
    get: () => props.visible,
    set: (value) => emit('update:visible', value)
  })

  const formRef = ref<FormInstance>()
  const submitting = ref(false)

  const form = reactive({
    cardLevel: 'normal',
    cardName: '普通会员卡',
    discountPercent: 100,
    status: 1,
    expireTime: '' as string | null,
    remark: ''
  })

  // 折扣必须是 1-100 的整数百分比
  const rules: FormRules = {
    cardLevel: [{ required: true, message: '请选择卡等级', trigger: 'change' }],
    cardName: [{ required: true, message: '请输入卡名称', trigger: 'blur' }],
    discountPercent: [{ required: true, message: '请填写折扣百分比', trigger: 'blur' }]
  }

  // 每次打开按当前卡数据回填
  watch(
    () => [props.visible, props.cardData] as const,
    ([visible]) => {
      if (!visible) return

      const data = props.cardData
      Object.assign(form, {
        cardLevel: data?.cardLevel ?? 'normal',
        cardName: data?.cardName ?? '普通会员卡',
        discountPercent: data?.discountPercent ?? 100,
        status: data?.status ?? 1,
        expireTime: data?.expireTime ?? '',
        remark: data?.remark ?? ''
      })

      nextTick(() => formRef.value?.clearValidate())
    },
    { immediate: true }
  )

  // 提交：一个会员一张卡，后端按会员 ID upsert
  const handleSubmit = async (): Promise<void> => {
    if (!formRef.value || !props.cardData) return

    const valid = await formRef.value.validate().catch(() => false)
    if (!valid) return

    submitting.value = true
    try {
      await fetchSaveMemberCard(props.cardData.memberUserId, {
        cardLevel: form.cardLevel,
        cardName: form.cardName.trim(),
        discountPercent: form.discountPercent,
        status: form.status,
        expireTime: form.expireTime || null,
        remark: form.remark.trim()
      })
      ElMessage.success('保存成功')
      dialogVisible.value = false
      emit('success')
    } catch {
      // 错误提示已由请求层统一弹出
    } finally {
      submitting.value = false
    }
  }
</script>
