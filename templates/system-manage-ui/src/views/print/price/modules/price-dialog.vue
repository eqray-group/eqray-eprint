<!-- 价格规则：门店 × 纸张 × 色彩 × 单双面 唯一；改价不能改组合 -->
<template>
  <ElDialog
    v-model="dialogVisible"
    :title="priceId ? '修改价格' : '新增价格'"
    width="520px"
    align-center
    :close-on-click-modal="false"
  >
    <ElForm ref="formRef" :model="form" :rules="rules" label-width="120px">
      <ElFormItem v-if="!priceId" label="纸张">
        <ElSelect v-model="form.paperType" class="w-full">
          <ElOption label="A4" value="A4" />
          <ElOption label="A3" value="A3" />
        </ElSelect>
      </ElFormItem>
      <ElFormItem v-if="!priceId" label="色彩">
        <ElSelect v-model="form.colorType" class="w-full">
          <ElOption label="黑白" value="bw" />
          <ElOption label="彩色" value="color" />
        </ElSelect>
      </ElFormItem>
      <ElFormItem v-if="!priceId" label="单双面">
        <ElSelect v-model="form.printMode" class="w-full">
          <ElOption label="单面" value="single" />
          <ElOption label="双面" value="double" />
        </ElSelect>
      </ElFormItem>
      <ElFormItem label="单价（分/页）" prop="price">
        <ElInputNumber v-model="form.price" :min="1" :step="10" class="w-full" />
        <span class="ml-2 text-xs text-g-500">展示给用户时按元换算</span>
      </ElFormItem>
      <ElFormItem label="状态">
        <ElSwitch
          v-model="form.status"
          :active-value="1"
          :inactive-value="0"
          active-text="启用"
          inactive-text="停用"
        />
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
  import { fetchCreatePrintPrice, fetchUpdatePrintPrice } from '@/api/print/price'

  defineOptions({ name: 'PrintPriceDialog' })

  interface Props {
    visible: boolean
    /** 为空表示新增 */
    priceId: string | null
    /** 新增时必须带门店 */
    shopId: string
    /** 编辑时传入原始价格与状态，避免再查一次详情 */
    origin?: Api.Print.PriceItem | null
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

  const form = reactive<Api.Print.PriceSaveParams>({
    printShopId: '',
    paperType: 'A4',
    colorType: 'bw',
    printMode: 'single',
    price: 10,
    status: 1
  })

  const rules: FormRules = {
    price: [{ required: true, message: '请输入单价', trigger: 'blur' }]
  }

  watch(
    () => props.visible,
    (visible) => {
      if (!visible) return

      if (props.priceId && props.origin) {
        Object.assign(form, {
          printShopId: props.origin.printShopId,
          paperType: props.origin.paperType,
          colorType: props.origin.colorType,
          printMode: props.origin.printMode,
          price: props.origin.price,
          status: props.origin.status
        })
      } else {
        Object.assign(form, {
          printShopId: props.shopId,
          paperType: 'A4',
          colorType: 'bw',
          printMode: 'single',
          price: 10,
          status: 1
        })
      }

      nextTick(() => formRef.value?.clearValidate())
    }
  )

  const handleSubmit = async (): Promise<void> => {
    if (!formRef.value) return

    const valid = await formRef.value.validate().catch(() => false)
    if (!valid) return

    submitting.value = true
    try {
      if (props.priceId) {
        await fetchUpdatePrintPrice(props.priceId, form.price, form.status)
      } else {
        await fetchCreatePrintPrice({ ...form, printShopId: props.shopId })
      }
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
