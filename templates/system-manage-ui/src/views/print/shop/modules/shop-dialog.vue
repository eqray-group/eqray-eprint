<!-- 门店新增 / 编辑：只维护基础信息，审核与营业状态走独立入口 -->
<template>
  <ElDialog
    v-model="dialogVisible"
    :title="shopId ? '编辑门店' : '新增门店'"
    width="640px"
    align-center
    :close-on-click-modal="false"
  >
    <ElForm ref="formRef" :model="form" :rules="rules" label-width="110px">
      <ElFormItem label="门店名称" prop="name">
        <ElInput v-model="form.name" maxlength="64" placeholder="如 大学城打印店" />
      </ElFormItem>
      <ElFormItem label="详细地址" prop="address">
        <ElInput v-model="form.address" maxlength="255" />
      </ElFormItem>
      <ElFormItem label="经纬度">
        <div class="flex gap-2">
          <ElInputNumber v-model="form.longitude" :precision="7" :step="0.0000001" placeholder="经度" />
          <ElInputNumber v-model="form.latitude" :precision="7" :step="0.0000001" placeholder="纬度" />
        </div>
        <span class="ml-2 text-xs text-g-500">用于会员端按距离排序，可留 0</span>
      </ElFormItem>
      <ElFormItem label="营业状态">
        <ElSwitch
          v-model="form.status"
          :active-value="1"
          :inactive-value="0"
          active-text="营业"
          inactive-text="停业"
        />
      </ElFormItem>
      <ElFormItem label="营业时间">
        <ElInput v-model="form.businessHours" placeholder='JSON，如 [{"start":"09:00","end":"21:00"}]，留空表示全天' />
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
  import { fetchCreatePrintShop, fetchGetPrintShop, fetchUpdatePrintShop } from '@/api/print/shop'

  defineOptions({ name: 'PrintShopDialog' })

  interface Props {
    visible: boolean
    /** 为空表示新增 */
    shopId: string | null
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

  const form = reactive<Api.Print.ShopSaveParams>({
    name: '',
    address: '',
    longitude: 0,
    latitude: 0,
    businessHours: '',
    status: 0,
    remark: ''
  })

  const rules: FormRules = {
    name: [{ required: true, message: '请输入门店名称', trigger: 'blur' }]
  }

  // 每次打开重置表单；编辑时拉详情回填，避免列表数据不全
  watch(
    () => props.visible,
    async (visible) => {
      if (!visible) return

      Object.assign(form, {
        name: '',
        address: '',
        longitude: 0,
        latitude: 0,
        businessHours: '',
        status: 0,
        remark: ''
      })

      if (props.shopId) {
        try {
          const detail = await fetchGetPrintShop(props.shopId)
          Object.assign(form, {
            name: detail.name,
            address: detail.address,
            longitude: detail.longitude,
            latitude: detail.latitude,
            businessHours: detail.businessHours ?? '',
            status: detail.status,
            remark: detail.remark
          })
        } catch {
          dialogVisible.value = false
          return
        }
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
      if (props.shopId) {
        await fetchUpdatePrintShop(props.shopId, { ...form })
      } else {
        await fetchCreatePrintShop({ ...form })
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
