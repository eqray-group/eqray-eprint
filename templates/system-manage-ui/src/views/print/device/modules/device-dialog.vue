<!-- 设备新增 / 编辑：设备编码是客户端机器指纹，创建后不可改 -->
<template>
  <ElDialog
    v-model="dialogVisible"
    :title="deviceId ? '编辑设备' : '新增设备'"
    width="560px"
    align-center
    :close-on-click-modal="false"
  >
    <ElForm ref="formRef" :model="form" :rules="rules" label-width="110px">
      <ElFormItem label="所属门店" prop="printShopId">
        <ElSelect v-model="form.printShopId" class="w-full" filterable placeholder="请选择门店">
          <ElOption v-for="shop in shops" :key="shop.id" :label="shop.name" :value="shop.id" />
        </ElSelect>
      </ElFormItem>
      <ElFormItem label="设备编码" prop="code">
        <ElInput v-model="form.code" :disabled="!!deviceId" placeholder="客户端机器指纹，全局唯一" maxlength="128" />
      </ElFormItem>
      <ElFormItem label="设备名称">
        <ElInput v-model="form.name" placeholder="留空自动生成" maxlength="64" />
      </ElFormItem>
      <ElFormItem label="状态">
        <ElSwitch
          v-model="form.status"
          :active-value="1"
          :inactive-value="0"
          active-text="启用"
          inactive-text="禁用"
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
  import { fetchCreatePrintDevice, fetchGetPrintDevice, fetchUpdatePrintDevice } from '@/api/print/device'
  import { fetchGetPrintShopList } from '@/api/print/shop'

  defineOptions({ name: 'PrintDeviceDialog' })

  interface Props {
    visible: boolean
    deviceId: string | null
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
  const shops = ref<Api.Print.ShopItem[]>([])

  const form = reactive<Api.Print.DeviceSaveParams>({
    code: '',
    printShopId: '',
    name: '',
    status: 1,
    remark: ''
  })

  const rules: FormRules = {
    printShopId: [{ required: true, message: '请选择所属门店', trigger: 'change' }],
    code: [{ required: true, message: '请输入设备编码', trigger: 'blur' }]
  }

  watch(
    () => props.visible,
    async (visible) => {
      if (!visible) return

      Object.assign(form, { code: '', printShopId: '', name: '', status: 1, remark: '' })

      if (shops.value.length === 0) {
        try {
          const res = await fetchGetPrintShopList({ pageNo: 1, pageSize: 200 })
          shops.value = res.list
        } catch {
          shops.value = []
        }
      }

      if (props.deviceId) {
        try {
          const detail = await fetchGetPrintDevice(props.deviceId)
          Object.assign(form, {
            code: detail.code,
            printShopId: detail.printShopId,
            name: detail.name,
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
      if (props.deviceId) {
        await fetchUpdatePrintDevice(props.deviceId, { ...form })
      } else {
        await fetchCreatePrintDevice({ ...form })
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
