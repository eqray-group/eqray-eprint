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
    modelValue: Api.Print.FileSearchParams
  }
  interface Emits {
    (e: 'update:modelValue', value: Api.Print.FileSearchParams): void
    (e: 'search', params: Api.Print.FileSearchParams): void
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

  // 与 print_file.convert_status 一致
  const statusOptions = [
    { label: '待处理', value: 0 },
    { label: '转换中', value: 1 },
    { label: '可打印', value: 2 },
    { label: '转换失败', value: 3 }
  ]

  const formItems = computed(() => [
    { label: '文件名', key: 'fileName', type: 'input', placeholder: '支持模糊搜索', clearable: true },
    { label: '会员 ID', key: 'memberId', type: 'input', placeholder: '请输入会员 ID', clearable: true },
    {
      label: '转换状态',
      key: 'convertStatus',
      type: 'select',
      props: { placeholder: '请选择', options: statusOptions, clearable: true }
    }
  ])

  function handleReset() {
    emit('reset')
  }

  async function handleSearch(params: Api.Print.FileSearchParams) {
    await searchBarRef.value.validate()
    emit('search', params)
  }
</script>
