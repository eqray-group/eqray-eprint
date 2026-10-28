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
    modelValue: Api.Member.CardSearchParams
  }
  interface Emits {
    (e: 'update:modelValue', value: Api.Member.CardSearchParams): void
    (e: 'search', params: Api.Member.CardSearchParams): void
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

  // 卡状态与 member_card.status 一致：0 停用 / 1 正常 / 2 已过期
  const statusOptions = [
    { label: '正常', value: 1 },
    { label: '停用', value: 0 },
    { label: '已过期', value: 2 }
  ]

  const formItems = computed(() => [
    {
      label: '卡号',
      key: 'cardNo',
      type: 'input',
      placeholder: '请输入卡号',
      clearable: true
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

  async function handleSearch(params: Api.Member.CardSearchParams) {
    await searchBarRef.value.validate()
    emit('search', params)
  }
</script>
