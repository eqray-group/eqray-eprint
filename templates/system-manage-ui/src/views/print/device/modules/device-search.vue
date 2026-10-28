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
  import { fetchGetPrintShopList } from '@/api/print/shop'

  interface Props {
    modelValue: Api.Print.DeviceSearchParams
  }
  interface Emits {
    (e: 'update:modelValue', value: Api.Print.DeviceSearchParams): void
    (e: 'search', params: Api.Print.DeviceSearchParams): void
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
  const shops = ref<Api.Print.ShopItem[]>([])

  // 门店下拉一次拉全量，搜索与新增共用，避免每个组件各拉一遍
  onMounted(async () => {
    try {
      const res = await fetchGetPrintShopList({ pageNo: 1, pageSize: 200 })
      shops.value = res.list
    } catch {
      shops.value = []
    }
  })

  const formItems = computed(() => [
    { label: '设备名称', key: 'keyword', type: 'input', placeholder: '名称 / 编码', clearable: true },
    {
      label: '所属门店',
      key: 'printShopId',
      type: 'select',
      props: {
        placeholder: '请选择门店',
        clearable: true,
        filterable: true,
        options: shops.value.map((shop) => ({ label: shop.name, value: shop.id }))
      }
    },
    {
      label: '在线状态',
      key: 'onlineStatus',
      type: 'select',
      props: {
        placeholder: '请选择',
        clearable: true,
        options: [
          { label: '在线', value: 1 },
          { label: '离线', value: 0 }
        ]
      }
    }
  ])

  function handleReset() {
    emit('reset')
  }

  async function handleSearch(params: Api.Print.DeviceSearchParams) {
    await searchBarRef.value.validate()
    emit('search', params)
  }
</script>
