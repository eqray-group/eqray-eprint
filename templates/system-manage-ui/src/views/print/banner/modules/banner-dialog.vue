<!-- 轮播图新增 / 编辑：门店传 0 表示全平台可见 -->
<template>
  <ElDialog
    v-model="dialogVisible"
    :title="bannerId ? '编辑轮播图' : '新增轮播图'"
    width="600px"
    align-center
    :close-on-click-modal="false"
  >
    <ElForm ref="formRef" :model="form" :rules="rules" label-width="110px">
      <ElFormItem label="所属门店">
        <ElSelect v-model="form.printShopId" class="w-full" filterable>
          <ElOption label="全平台（所有门店可见）" value="0" />
          <ElOption v-for="shop in shops" :key="shop.id" :label="shop.name" :value="shop.id" />
        </ElSelect>
      </ElFormItem>
      <ElFormItem label="标题">
        <ElInput v-model="form.title" maxlength="64" placeholder="展示用，可留空" />
      </ElFormItem>
      <ElFormItem label="图片地址" prop="imageUrl">
        <ElInput v-model="form.imageUrl" maxlength="255" placeholder="/storage/... 或完整 URL" />
        <span class="mt-1 text-xs text-g-500">建议先用「基础设施 / 文件配置」上传，再填返回的 url</span>
      </ElFormItem>
      <ElFormItem label="跳转类型">
        <ElRadioGroup v-model="form.linkType">
          <ElRadio value="none">不跳转</ElRadio>
          <ElRadio value="page">站内页</ElRadio>
          <ElRadio value="url">外部链接</ElRadio>
        </ElRadioGroup>
      </ElFormItem>
      <ElFormItem v-if="form.linkType !== 'none'" label="跳转地址" prop="linkUrl">
        <ElInput v-model="form.linkUrl" maxlength="255" placeholder="站内页填路径，外链填完整 URL" />
      </ElFormItem>
      <ElFormItem label="排序">
        <ElInputNumber v-model="form.sort" :min="0" />
        <span class="ml-2 text-xs text-g-500">数字小的靠前</span>
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
  import {
    fetchCreatePrintBanner,
    fetchGetPrintBanner,
    fetchUpdatePrintBanner
  } from '@/api/print/banner'
  import { fetchGetPrintShopList } from '@/api/print/shop'

  defineOptions({ name: 'PrintBannerDialog' })

  interface Props {
    visible: boolean
    bannerId: string | null
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

  const form = reactive<Api.Print.BannerSaveParams>({
    printShopId: '0',
    title: '',
    imageUrl: '',
    linkType: 'none',
    linkUrl: '',
    sort: 0,
    status: 1
  })

  const rules: FormRules = {
    imageUrl: [{ required: true, message: '请填写图片地址', trigger: 'blur' }],
    linkUrl: [{ required: true, message: '请填写跳转地址', trigger: 'blur' }]
  }

  watch(
    () => props.visible,
    async (visible) => {
      if (!visible) return

      Object.assign(form, {
        printShopId: '0',
        title: '',
        imageUrl: '',
        linkType: 'none',
        linkUrl: '',
        sort: 0,
        status: 1
      })

      if (shops.value.length === 0) {
        try {
          const res = await fetchGetPrintShopList({ pageNo: 1, pageSize: 200 })
          shops.value = res.list
        } catch {
          shops.value = []
        }
      }

      if (props.bannerId) {
        try {
          const detail = await fetchGetPrintBanner(props.bannerId)
          Object.assign(form, {
            printShopId: detail.printShopId,
            title: detail.title,
            imageUrl: detail.imageUrl,
            linkType: detail.linkType || 'none',
            linkUrl: detail.linkUrl,
            sort: detail.sort,
            status: detail.status
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
      const payload: Api.Print.BannerSaveParams = { ...form }
      if (props.bannerId) {
        await fetchUpdatePrintBanner(props.bannerId, payload)
      } else {
        await fetchCreatePrintBanner(payload)
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
