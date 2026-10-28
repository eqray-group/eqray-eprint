<template>
  <ArtSearchBar
    ref="searchBarRef"
    v-model="formData"
    :items="formItems"
    :rules="rules"
    @reset="handleReset"
    @search="handleSearch"
  >
  </ArtSearchBar>
</template>

<script setup lang="ts">
  interface Props {
    modelValue: Api.Pay.PaymentSearchParams
  }
  interface Emits {
    (e: 'update:modelValue', value: Api.Pay.PaymentSearchParams): void
    (e: 'search', params: Api.Pay.PaymentSearchParams): void
    (e: 'reset'): void
  }
  const props = defineProps<Props>()
  const emit = defineEmits<Emits>()

  const searchBarRef = ref()
  const formData = computed({
    get: () => props.modelValue,
    set: (val) => emit('update:modelValue', val)
  })

  const rules = {}

  const channelOptions = [
    { label: '微信', value: 'wechat' },
    { label: '支付宝', value: 'alipay' },
    { label: '钱包余额', value: 'balance' }
  ]

  // 支付单状态与 pay/payment_service 一致
  const statusOptions = [
    { label: '待支付', value: 0 },
    { label: '支付中', value: 1 },
    { label: '已支付', value: 2 },
    { label: '支付失败', value: 3 },
    { label: '已关闭', value: 4 },
    { label: '已退款', value: 5 }
  ]

  const formItems = computed(() => [
    {
      label: '支付单号',
      key: 'paymentNo',
      type: 'input',
      placeholder: '支持模糊搜索',
      clearable: true
    },
    {
      label: '会员 ID',
      key: 'memberId',
      type: 'input',
      placeholder: '请输入会员 ID',
      clearable: true
    },
    {
      label: '渠道',
      key: 'channel',
      type: 'select',
      props: { placeholder: '请选择渠道', options: channelOptions, clearable: true }
    },
    {
      label: '状态',
      key: 'status',
      type: 'select',
      props: { placeholder: '请选择状态', options: statusOptions }
    }
  ])

  function handleReset() {
    emit('reset')
  }

  async function handleSearch(params: Api.Pay.PaymentSearchParams) {
    await searchBarRef.value.validate()
    emit('search', params)
  }
</script>
