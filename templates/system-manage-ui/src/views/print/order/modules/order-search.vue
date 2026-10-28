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
    modelValue: Api.Print.OrderSearchParams
  }
  interface Emits {
    (e: 'update:modelValue', value: Api.Print.OrderSearchParams): void
    (e: 'search', params: Api.Print.OrderSearchParams): void
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

  // 状态枚举与 print_order 保持一致
  const statusOptions = [
    { label: '待支付', value: 0 },
    { label: '待打印', value: 1 },
    { label: '打印中', value: 2 },
    { label: '待取件', value: 3 },
    { label: '已完成', value: 4 },
    { label: '已取消', value: 5 },
    { label: '退款中', value: 6 },
    { label: '已退款', value: 7 },
    { label: '超期关闭', value: 8 }
  ]

  const payOptions = [
    { label: '未支付', value: 0 },
    { label: '已支付', value: 1 },
    { label: '退款中', value: 2 },
    { label: '已退款', value: 3 }
  ]

  const formItems = computed(() => [
    { label: '订单号', key: 'orderNo', type: 'input', placeholder: '支持模糊搜索', clearable: true },
    {
      label: '门店',
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
      label: '订单状态',
      key: 'orderStatus',
      type: 'select',
      props: { placeholder: '请选择', options: statusOptions, clearable: true }
    },
    {
      label: '支付状态',
      key: 'payStatus',
      type: 'select',
      props: { placeholder: '请选择', options: payOptions, clearable: true }
    }
  ])

  function handleReset() {
    emit('reset')
  }

  async function handleSearch(params: Api.Print.OrderSearchParams) {
    await searchBarRef.value.validate()
    emit('search', params)
  }
</script>
