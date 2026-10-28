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
    modelValue: Api.Print.BannerSearchParams
  }
  interface Emits {
    (e: 'update:modelValue', value: Api.Print.BannerSearchParams): void
    (e: 'search', params: Api.Print.BannerSearchParams): void
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

  onMounted(async () => {
    try {
      const res = await fetchGetPrintShopList({ pageNo: 1, pageSize: 200 })
      shops.value = res.list
    } catch {
      shops.value = []
    }
  })

  // 门店留空表示查全部门店（含全平台图）
  const formItems = computed(() => [
    {
      label: '门店',
      key: 'printShopId',
      type: 'select',
      props: {
        placeholder: '不选则查全部',
        clearable: true,
        filterable: true,
        options: shops.value.map((shop) => ({ label: shop.name, value: shop.id }))
      }
    }
  ])

  function handleReset() {
    emit('reset')
  }

  async function handleSearch(params: Api.Print.BannerSearchParams) {
    await searchBarRef.value.validate()
    emit('search', params)
  }
</script>
