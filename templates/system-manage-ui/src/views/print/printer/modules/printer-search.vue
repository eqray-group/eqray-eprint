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
  import { fetchGetPrintDeviceList } from '@/api/print/device'

  interface Props {
    modelValue: Api.Print.PrinterSearchParams
  }
  interface Emits {
    (e: 'update:modelValue', value: Api.Print.PrinterSearchParams): void
    (e: 'search', params: Api.Print.PrinterSearchParams): void
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
  const devices = ref<Api.Print.DeviceItem[]>([])

  onMounted(async () => {
    try {
      const res = await fetchGetPrintDeviceList({ pageNo: 1, pageSize: 200 })
      devices.value = res.list
    } catch {
      devices.value = []
    }
  })

  const formItems = computed(() => [
    {
      label: '所属设备',
      key: 'deviceId',
      type: 'select',
      props: {
        placeholder: '请选择设备',
        clearable: true,
        filterable: true,
        options: devices.value.map((device) => ({
          label: `${device.name}（${device.code}）`,
          value: device.id
        }))
      }
    },
    {
      label: '状态',
      key: 'status',
      type: 'select',
      props: {
        placeholder: '请选择',
        clearable: true,
        options: [
          { label: '启用', value: 1 },
          { label: '停用', value: 0 }
        ]
      }
    }
  ])

  function handleReset() {
    emit('reset')
  }

  async function handleSearch(params: Api.Print.PrinterSearchParams) {
    await searchBarRef.value.validate()
    emit('search', params)
  }
</script>
