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
    modelValue: Api.Print.ShopSearchParams
  }
  interface Emits {
    (e: 'update:modelValue', value: Api.Print.ShopSearchParams): void
    (e: 'search', params: Api.Print.ShopSearchParams): void
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
    { label: '营业', value: 1 },
    { label: '停业', value: 0 }
  ]

  // 审核状态与 print_shop.audit_status 一致
  const auditOptions = [
    { label: '待审核', value: 0 },
    { label: '已通过', value: 1 },
    { label: '已拒绝', value: 2 }
  ]

  const formItems = computed(() => [
    { label: '门店名称', key: 'name', type: 'input', placeholder: '请输入门店名称', clearable: true },
    {
      label: '营业状态',
      key: 'status',
      type: 'select',
      props: { placeholder: '请选择', options: statusOptions, clearable: true }
    },
    {
      label: '审核状态',
      key: 'auditStatus',
      type: 'select',
      props: { placeholder: '请选择', options: auditOptions, clearable: true }
    }
  ])

  function handleReset() {
    emit('reset')
  }

  async function handleSearch(params: Api.Print.ShopSearchParams) {
    await searchBarRef.value.validate()
    emit('search', params)
  }
</script>
