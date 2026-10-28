<!-- 渠道配置弹窗：编辑时回显密钥，留空表示不修改 -->
<template>
  <ElDialog
    v-model="dialogVisible"
    :title="isEdit ? '编辑支付渠道' : '新增支付渠道'"
    width="720px"
    align-center
    :close-on-click-modal="false"
  >
    <ElForm ref="formRef" :model="form" :rules="rules" label-width="120px">
      <ElFormItem label="配置编码" prop="code">
        <ElInput v-model="form.code" placeholder="如 wechat_mini / wechat_h5 / alipay_wap" maxlength="32" />
        <span class="ml-3 text-xs text-g-600">编码唯一，决定调用哪套商户配置</span>
      </ElFormItem>
      <ElFormItem label="配置名称" prop="name">
        <ElInput v-model="form.name" placeholder="如 微信小程序支付" maxlength="64" />
      </ElFormItem>
      <ElFormItem label="渠道" prop="channel">
        <ElSelect v-model="form.channel" class="w-full" @change="handleChannelChange">
          <ElOption label="微信支付" value="wechat" />
          <ElOption label="支付宝" value="alipay" />
        </ElSelect>
      </ElFormItem>
      <ElFormItem label="支付场景" prop="scene">
        <ElSelect v-model="form.scene" class="w-full">
          <ElOption v-for="item in sceneOptions" :key="item.value" :label="item.label" :value="item.value" />
        </ElSelect>
      </ElFormItem>
      <ElFormItem label="AppID" prop="appId">
        <ElInput v-model="form.appId" placeholder="微信 appid / 支付宝 appid" maxlength="64" />
      </ElFormItem>
      <ElFormItem label="商户号">
        <ElInput v-model="form.mchId" placeholder="微信 mchid / 支付宝 PID" maxlength="64" />
      </ElFormItem>
      <ElFormItem v-if="form.channel === 'wechat'" label="证书序列号">
        <ElInput v-model="form.serialNo" placeholder="APIv3 证书序列号" maxlength="64" />
      </ElFormItem>
      <ElFormItem label="商户私钥">
        <ElInput
          v-model="form.privateKey"
          type="textarea"
          :rows="3"
          :placeholder="isEdit ? '留空表示不修改' : '粘贴商户私钥内容'"
        />
      </ElFormItem>
      <ElFormItem v-if="form.channel === 'wechat'" label="APIv3 密钥">
        <ElInput v-model="form.apiKey" :placeholder="isEdit ? '留空表示不修改' : '32 位 APIv3 密钥'" maxlength="256" />
      </ElFormItem>
      <ElFormItem label="平台公钥">
        <ElInput
          v-model="form.platformPublicKey"
          type="textarea"
          :rows="3"
          :placeholder="isEdit ? '留空表示不修改' : '微信平台证书 / 支付宝公钥'"
        />
      </ElFormItem>
      <ElFormItem v-if="form.channel === 'alipay'" label="网关地址">
        <ElInput v-model="form.gateway" placeholder="https://openapi.alipay.com/gateway.do" maxlength="128" />
      </ElFormItem>
      <ElFormItem label="回调地址">
        <ElInput v-model="form.notifyUrl" placeholder="支付结果回调，空则用默认路由" maxlength="255" />
      </ElFormItem>
      <ElFormItem label="回跳地址">
        <ElInput v-model="form.returnUrl" placeholder="H5 支付完成后的跳回页" maxlength="255" />
      </ElFormItem>
      <ElFormItem label="扩展参数">
        <ElInput v-model="form.extra" type="textarea" :rows="2" placeholder="JSON，可留空" />
      </ElFormItem>
      <ElFormItem label="状态">
        <ElSwitch v-model="form.status" :active-value="1" :inactive-value="0" />
      </ElFormItem>
      <ElFormItem label="排序">
        <ElInputNumber v-model="form.sort" :min="0" :max="999" />
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
    fetchCreateChannelConfig,
    fetchGetChannelConfig,
    fetchUpdateChannelConfig
  } from '@/api/pay/channel-config'

  defineOptions({ name: 'PayChannelConfigDialog' })

  interface Props {
    visible: boolean
    /** 编辑时的配置 ID，空表示新增 */
    configId?: string | null
  }
  interface Emits {
    (e: 'update:visible', value: boolean): void
    (e: 'success'): void
  }

  const props = withDefaults(defineProps<Props>(), { configId: null })
  const emit = defineEmits<Emits>()

  const dialogVisible = computed({
    get: () => props.visible,
    set: (value) => emit('update:visible', value)
  })

  const isEdit = computed(() => !!props.configId)

  const formRef = ref<FormInstance>()
  const submitting = ref(false)

  const form = reactive<Api.Pay.ChannelConfigSaveParams>({
    code: '',
    name: '',
    channel: 'wechat',
    scene: 'mini',
    appId: '',
    mchId: '',
    serialNo: '',
    privateKey: '',
    apiKey: '',
    platformPublicKey: '',
    gateway: '',
    notifyUrl: '',
    returnUrl: '',
    extra: '',
    status: 1,
    sort: 0,
    remark: ''
  })

  // 场景随渠道变化：微信 mini/jsapi/h5，支付宝只支持手机网站 wap
  const sceneOptions = computed(() =>
    form.channel === 'wechat'
      ? [
          { label: '小程序 mini', value: 'mini' },
          { label: '公众号 jsapi', value: 'jsapi' },
          { label: '手机浏览器 h5', value: 'h5' }
        ]
      : [{ label: '手机网站 wap', value: 'wap' }]
  )

  const rules: FormRules = {
    code: [{ required: true, message: '请输入配置编码', trigger: 'blur' }],
    name: [{ required: true, message: '请输入配置名称', trigger: 'blur' }],
    channel: [{ required: true, message: '请选择渠道', trigger: 'change' }],
    scene: [{ required: true, message: '请选择支付场景', trigger: 'change' }]
  }

  // 渠道切换后场景要重置，避免出现 wechat + wap 这类无效组合
  const handleChannelChange = (): void => {
    form.scene = form.channel === 'wechat' ? 'mini' : 'wap'
  }

  // 打开弹窗：新增清空表单，编辑拉详情回显（含密钥）
  watch(
    () => [props.visible, props.configId] as const,
    async ([visible]) => {
      if (!visible) return

      Object.assign(form, {
        code: '',
        name: '',
        channel: 'wechat',
        scene: 'mini',
        appId: '',
        mchId: '',
        serialNo: '',
        privateKey: '',
        apiKey: '',
        platformPublicKey: '',
        gateway: '',
        notifyUrl: '',
        returnUrl: '',
        extra: '',
        status: 1,
        sort: 0,
        remark: ''
      })

      if (props.configId) {
        try {
          const detail = await fetchGetChannelConfig(props.configId)
          Object.assign(form, {
            code: detail.code,
            name: detail.name,
            channel: detail.channel,
            scene: detail.scene,
            appId: detail.appId,
            mchId: detail.mchId,
            serialNo: detail.serialNo,
            privateKey: detail.privateKey ?? '',
            apiKey: detail.apiKey ?? '',
            platformPublicKey: detail.platformPublicKey ?? '',
            gateway: detail.gateway,
            notifyUrl: detail.notifyUrl,
            returnUrl: detail.returnUrl,
            extra: detail.extra ?? '',
            status: detail.status,
            sort: detail.sort,
            remark: detail.remark
          })
        } catch {
          dialogVisible.value = false
        }
      }

      nextTick(() => formRef.value?.clearValidate())
    },
    { immediate: true }
  )

  // 提交：编辑时密钥留空表示保持原值
  const handleSubmit = async (): Promise<void> => {
    if (!formRef.value) return

    const valid = await formRef.value.validate().catch(() => false)
    if (!valid) return

    submitting.value = true
    try {
      const payload: Api.Pay.ChannelConfigSaveParams = {
        ...form,
        code: form.code.trim(),
        name: form.name.trim(),
        privateKey: form.privateKey?.trim() ?? '',
        apiKey: form.apiKey?.trim() ?? '',
        platformPublicKey: form.platformPublicKey?.trim() ?? '',
        extra: form.extra?.trim() ?? ''
      }

      if (props.configId) {
        await fetchUpdateChannelConfig(props.configId, payload)
      } else {
        await fetchCreateChannelConfig(payload)
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
