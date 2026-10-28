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
    modelValue: Api.Pay.RefundSearchParams
  }
  interface Emits {
    (e: 'update:modelValue', value: Api.Pay.RefundSearchParams): void
    (e: 'search', params: Api.Pay.RefundSearchParams): void
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

  const statusOptions = [
    { label: '待退款', value: 0 },
    { label: '退款中', value: 1 },
    { label: '退款成功', value: 2 },
    { label: '退款失败', value: 3 }
  ]

  const formItems = computed(() => [
    {
      label: '支付单号',
      key: 'paymentNo',
      type: 'input',
      placeholder: '请输入支付单号',
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
      label: '状态',
      key: 'status',
      type: 'select',
      props: { placeholder: '请选择状态', options: statusOptions }
    }
  ])

  function handleReset() {
    emit('reset')
  }

  async function handleSearch(params: Api.Pay.RefundSearchParams) {
    await searchBarRef.value.validate()
    emit('search', params)
  }
</script>
