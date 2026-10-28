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
    modelValue: Api.Pay.ChannelConfigSearchParams
  }
  interface Emits {
    (e: 'update:modelValue', value: Api.Pay.ChannelConfigSearchParams): void
    (e: 'search', params: Api.Pay.ChannelConfigSearchParams): void
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
    { label: '微信支付', value: 'wechat' },
    { label: '支付宝', value: 'alipay' }
  ]

  const statusOptions = [
    { label: '启用', value: 1 },
    { label: '停用', value: 0 }
  ]

  const formItems = computed(() => [
    {
      label: '配置编码',
      key: 'code',
      type: 'input',
      placeholder: '如 wechat_mini',
      clearable: true
    },
    {
      label: '渠道',
      key: 'channel',
      type: 'select',
      props: {
        placeholder: '请选择渠道',
        options: channelOptions,
        clearable: true
      }
    },
    {
      label: '状态',
      key: 'status',
      type: 'select',
      props: {
        placeholder: '请选择状态',
        options: statusOptions
      }
    }
  ])

  function handleReset() {
    emit('reset')
  }

  async function handleSearch(params: Api.Pay.ChannelConfigSearchParams) {
    await searchBarRef.value.validate()
    emit('search', params)
  }
</script>
