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
    modelValue: Api.Pay.WalletTransactionSearchParams
  }
  interface Emits {
    (e: 'update:modelValue', value: Api.Pay.WalletTransactionSearchParams): void
    (e: 'search', params: Api.Pay.WalletTransactionSearchParams): void
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

  // 与 pay_wallet_transaction.type 一致
  const typeOptions = [
    { label: '收入', value: 1 },
    { label: '支出', value: 2 }
  ]

  // 业务类型与 pay_enums.BizType* 一致
  const bizOptions = [
    { label: '钱包充值', value: 'wallet_recharge' },
    { label: '打印消费', value: 'print_consume' },
    { label: '打印退款', value: 'print_refund' },
    { label: '后台调账', value: 'manual_adjust' }
  ]

  const formItems = computed(() => [
    {
      label: '会员 ID',
      key: 'memberId',
      type: 'input',
      placeholder: '填写后可看余额，留空查全部流水',
      clearable: true
    },
    {
      label: '收支',
      key: 'type',
      type: 'select',
      props: { placeholder: '请选择', options: typeOptions, clearable: true }
    },
    {
      label: '业务类型',
      key: 'bizType',
      type: 'select',
      props: { placeholder: '请选择', options: bizOptions, clearable: true }
    }
  ])

  function handleReset() {
    emit('reset')
  }

  async function handleSearch(params: Api.Pay.WalletTransactionSearchParams) {
    await searchBarRef.value.validate()
    emit('search', params)
  }
</script>
