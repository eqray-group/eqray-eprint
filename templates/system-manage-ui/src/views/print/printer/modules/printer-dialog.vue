<!-- 打印机新增 / 编辑：能力字段直接决定会员端能选哪些参数组合 -->
<template>
  <ElDialog
    v-model="dialogVisible"
    :title="printerId ? '编辑打印机' : '新增打印机'"
    width="600px"
    align-center
    :close-on-click-modal="false"
  >
    <ElAlert
      type="info"
      :closable="false"
      show-icon
      title="能力配置决定能否接单：彩色 / 双面未勾选时，对应组合不会出现在会员端"
      class="mb-4"
    />

    <ElForm ref="formRef" :model="form" :rules="rules" label-width="120px">
      <ElFormItem label="所属设备" prop="deviceId">
        <ElSelect v-model="form.deviceId" class="w-full" filterable placeholder="请选择设备">
          <ElOption
            v-for="device in devices"
            :key="device.id"
            :label="`${device.name}（${device.code}）`"
            :value="device.id"
          />
        </ElSelect>
      </ElFormItem>
      <ElFormItem label="打印机名称" prop="name">
        <ElInput v-model="form.name" placeholder="如 A4 黑白 1 号机" maxlength="64" />
      </ElFormItem>
      <ElFormItem label="队列名称" prop="systemName">
        <ElInput v-model="form.systemName" placeholder="Windows 打印队列名，客户端按此静默打印" maxlength="128" />
      </ElFormItem>
      <ElFormItem label="型号">
        <ElInput v-model="form.model" maxlength="64" />
      </ElFormItem>
      <ElFormItem label="支持纸张" prop="paperTypes">
        <ElCheckboxGroup v-model="paperSelection">
          <ElCheckbox label="A4" value="A4" />
          <ElCheckbox label="A3" value="A3" />
        </ElCheckboxGroup>
      </ElFormItem>
      <ElFormItem label="支持彩色">
        <ElSwitch v-model="form.colorSupport" :active-value="1" :inactive-value="0" />
      </ElFormItem>
      <ElFormItem label="支持双面">
        <ElSwitch v-model="form.duplexSupport" :active-value="1" :inactive-value="0" />
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
  import {
    fetchCreatePrintPrinter,
    fetchGetPrintPrinter,
    fetchUpdatePrintPrinter
  } from '@/api/print/printer'
  import { fetchGetPrintDeviceList } from '@/api/print/device'

  defineOptions({ name: 'PrintPrinterDialog' })

  interface Props {
    visible: boolean
    printerId: string | null
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
  const devices = ref<Api.Print.DeviceItem[]>([])
  const paperSelection = ref<string[]>(['A4'])

  const form = reactive<Api.Print.PrinterSaveParams>({
    deviceId: '',
    name: '',
    systemName: '',
    model: '',
    paperTypes: 'A4',
    colorSupport: 0,
    duplexSupport: 0,
    status: 1,
    remark: ''
  })

  const rules: FormRules = {
    deviceId: [{ required: true, message: '请选择所属设备', trigger: 'change' }],
    name: [{ required: true, message: '请输入打印机名称', trigger: 'blur' }],
    systemName: [{ required: true, message: '请输入 Windows 打印队列名称', trigger: 'blur' }]
  }

  watch(
    () => props.visible,
    async (visible) => {
      if (!visible) return

      Object.assign(form, {
        deviceId: '',
        name: '',
        systemName: '',
        model: '',
        paperTypes: 'A4',
        colorSupport: 0,
        duplexSupport: 0,
        status: 1,
        remark: ''
      })
      paperSelection.value = ['A4']

      if (devices.value.length === 0) {
        try {
          const res = await fetchGetPrintDeviceList({ pageNo: 1, pageSize: 200 })
          devices.value = res.list
        } catch {
          devices.value = []
        }
      }

      if (props.printerId) {
        try {
          const detail = await fetchGetPrintPrinter(props.printerId)
          Object.assign(form, {
            deviceId: detail.deviceId,
            name: detail.name,
            systemName: detail.systemName,
            model: detail.model,
            paperTypes: detail.paperTypes,
            colorSupport: detail.colorSupport,
            duplexSupport: detail.duplexSupport,
            status: detail.status,
            remark: detail.remark
          })
          paperSelection.value = detail.paperTypes ? detail.paperTypes.split(',').filter(Boolean) : []
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
    if (paperSelection.value.length === 0) {
      ElMessage.warning('请至少选择一种纸张')
      return
    }

    submitting.value = true
    try {
      const payload: Api.Print.PrinterSaveParams = {
        ...form,
        paperTypes: paperSelection.value.join(',')
      }
      if (props.printerId) {
        await fetchUpdatePrintPrinter(props.printerId, payload)
      } else {
        await fetchCreatePrintPrinter(payload)
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
