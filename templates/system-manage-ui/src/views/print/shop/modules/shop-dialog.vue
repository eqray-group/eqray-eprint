<!-- 门店新增 / 编辑：只维护基础信息，审核与营业状态走独立入口 -->
<template>
  <ElDialog
    v-model="dialogVisible"
    :title="shopId ? '编辑门店' : '新增门店'"
    width="640px"
    align-center
    :close-on-click-modal="false"
  >
    <ElForm ref="formRef" :model="form" :rules="rules" label-width="110px">
      <ElFormItem label="门店名称" prop="name">
        <ElInput v-model="form.name" maxlength="64" placeholder="如 大学城打印店" />
      </ElFormItem>
      <ElFormItem label="详细地址" prop="address">
        <ElInput v-model="form.address" maxlength="255" />
      </ElFormItem>
      <ElFormItem label="门店位置">
        <div class="w-full">
          <div class="flex flex-wrap gap-2 items-center">
            <ElButton plain @click="locationVisible = true">地图选点</ElButton>
            <span class="text-xs text-g-500">{{ coordText }}</span>
            <ElButton v-if="hasPoint" link type="danger" @click="clearPoint">清除</ElButton>
          </div>
          <div class="mt-1 text-xs text-g-500">
            在高德地图上选点，自动回填经纬度与详细地址；坐标用于会员端按距离排序
          </div>
        </div>
      </ElFormItem>
      <ElFormItem label="营业状态">
        <div>
          <ElSwitch
            v-model="form.status"
            :active-value="1"
            :inactive-value="0"
            active-text="营业"
            inactive-text="停业"
          />
          <div class="mt-1 text-xs text-g-500">审核通过且营业中的门店才接单</div>
        </div>
      </ElFormItem>
      <ElFormItem label="营业时间">
        <div class="w-full">
          <div v-for="(slot, index) in hours" :key="index" class="mb-2 flex items-center gap-2">
            <ElTimeSelect
              v-model="slot.start"
              start="00:00"
              end="23:30"
              step="00:30"
              placeholder="开始"
              class="w-32"
            />
            <span class="text-g-500">至</span>
            <ElTimeSelect
              v-model="slot.end"
              start="00:30"
              end="23:59"
              step="00:30"
              placeholder="结束"
              class="w-32"
            />
            <ElButton link type="danger" @click="removeHour(index)">删除</ElButton>
          </div>
          <div class="flex items-center gap-2">
            <ElButton plain size="small" @click="addHour">添加时段</ElButton>
            <span class="text-xs text-g-500">
              {{ hours.length ? '一天有多个营业时段时可添加多条' : '不填表示全天营业' }}
            </span>
          </div>
        </div>
      </ElFormItem>
      <ElFormItem label="备注">
        <ElInput v-model="form.remark" type="textarea" :rows="2" maxlength="255" />
      </ElFormItem>
    </ElForm>

    <template #footer>
      <ElButton @click="dialogVisible = false">取消</ElButton>
      <ElButton type="primary" :loading="submitting" @click="handleSubmit">保存</ElButton>
    </template>
  </ElDialog>

  <ShopLocationDialog
    v-model:visible="locationVisible"
    :longitude="form.longitude"
    :latitude="form.latitude"
    :address="form.address"
    @confirm="handleLocationConfirm"
  />
</template>

<script setup lang="ts">
  import { ElMessage, ElTimeSelect } from 'element-plus'
  import type { FormInstance, FormRules } from 'element-plus'
  import { fetchCreatePrintShop, fetchGetPrintShop, fetchUpdatePrintShop } from '@/api/print/shop'
  import ShopLocationDialog from './shop-location-dialog.vue'

  defineOptions({ name: 'PrintShopDialog' })

  interface Props {
    visible: boolean
    /** 为空表示新增 */
    shopId: string | null
  }
  interface Emits {
    (e: 'update:visible', value: boolean): void
    (e: 'success'): void
  }

  const props = defineProps<Props>()
  const emit = defineEmits<Emits>()

  const dialogVisible = computed({
    get: () => props.visible,
    set: (value) => emit('update:visible', value)
  })

  const formRef = ref<FormInstance>()
  const submitting = ref(false)

  /** 营业时段，空数组表示全天营业；提交前序列化成 businessHours JSON */
  interface BusinessHourSlot {
    start: string
    end: string
  }

  const form = reactive<Api.Print.ShopSaveParams>({
    name: '',
    address: '',
    longitude: 0,
    latitude: 0,
    businessHours: '',
    // 新增门店默认营业，避免保存后忘记切状态
    status: 1,
    remark: ''
  })

  const hours = ref<BusinessHourSlot[]>([])

  const rules: FormRules = {
    name: [{ required: true, message: '请输入门店名称', trigger: 'blur' }]
  }

  /**
   * 解析详情里的营业时间 JSON
   * @param raw 后端返回的 businessHours，形如 [{"start":"09:00","end":"21:00"}]
   * @returns 为空或格式不合法时返回空数组（按全天营业处理），不阻断编辑
   */
  const parseBusinessHours = (raw?: string | null): BusinessHourSlot[] => {
    if (!raw) return []

    try {
      const parsed = JSON.parse(raw)
      if (!Array.isArray(parsed)) return []

      return parsed
        .filter((item) => typeof item?.start === 'string' && typeof item?.end === 'string')
        .map((item) => ({ start: item.start, end: item.end }))
    } catch {
      return []
    }
  }

  /**
   * 把时段列表序列化成后端存储格式
   * @returns 没有完整时段时返回空串，后端会存 NULL 表示全天营业
   */
  const buildBusinessHours = (): string => {
    const slots = hours.value
      .filter((slot) => slot.start && slot.end)
      .map((slot) => ({ start: slot.start, end: slot.end }))

    return slots.length ? JSON.stringify(slots) : ''
  }

  /**
   * 添加一个默认时段 09:00-21:00
   * 参数：无；返回 void
   */
  const addHour = (): void => {
    hours.value.push({ start: '09:00', end: '21:00' })
  }

  /**
   * 删除指定时段
   * @param index 时段下标
   */
  const removeHour = (index: number): void => {
    hours.value.splice(index, 1)
  }

  /** 地图选点弹窗是否可见 */
  const locationVisible = ref(false)

  /** 是否已选坐标：经纬度都为 0 视为未选，保存后门店不参与按距离排序 */
  const hasPoint = computed(() => Number(form.longitude) > 0 || Number(form.latitude) > 0)

  /** 坐标展示文案，保留 6 位小数便于核对 */
  const coordText = computed(() =>
    hasPoint.value
      ? `经纬度：${Number(form.longitude).toFixed(6)}, ${Number(form.latitude).toFixed(6)}`
      : '尚未选择坐标'
  )

  /**
   * 清空坐标
   * 参数：无；返回 void，只清经纬度，保留手填的详细地址
   */
  const clearPoint = (): void => {
    form.longitude = 0
    form.latitude = 0
  }

  /**
   * 地图选点确认后回填表单
   * @param payload 高德选点回传的经纬度与逆地理编码地址
   */
  const handleLocationConfirm = (payload: {
    longitude: number
    latitude: number
    address: string
  }): void => {
    form.longitude = payload.longitude
    form.latitude = payload.latitude
    if (payload.address) {
      form.address = payload.address
    }
  }

  // 每次打开重置表单；编辑时拉详情回填，避免列表数据不全
  watch(
    () => props.visible,
    async (visible) => {
      if (!visible) return

      Object.assign(form, {
        name: '',
        address: '',
        longitude: 0,
        latitude: 0,
        businessHours: '',
        // 新增门店默认营业
        status: 1,
        remark: ''
      })
      hours.value = []

      if (props.shopId) {
        try {
          const detail = await fetchGetPrintShop(props.shopId)
          Object.assign(form, {
            name: detail.name,
            address: detail.address,
            longitude: detail.longitude,
            latitude: detail.latitude,
            businessHours: detail.businessHours ?? '',
            status: detail.status,
            remark: detail.remark
          })
          hours.value = parseBusinessHours(detail.businessHours)
        } catch {
          dialogVisible.value = false
          return
        }
      }

      nextTick(() => formRef.value?.clearValidate())
    }
  )

  const handleSubmit = async (): Promise<void> => {
    if (!formRef.value) return

    const valid = await formRef.value.validate().catch(() => false)
    if (!valid) return

    // 时段必须成对填写，且结束时间晚于开始时间（暂不支持跨夜营业）
    const incomplete = hours.value.find(
      (slot) => (!slot.start && slot.end) || (slot.start && !slot.end)
    )
    if (incomplete) {
      ElMessage.warning('请把营业时段的开始和结束时间都填上')
      return
    }
    const invalid = hours.value.find((slot) => slot.start && slot.end && slot.start >= slot.end)
    if (invalid) {
      ElMessage.warning(`营业时段 ${invalid.start}-${invalid.end} 不合法，结束时间要晚于开始时间`)
      return
    }

    const payload = { ...form, businessHours: buildBusinessHours() }

    submitting.value = true
    try {
      if (props.shopId) {
        await fetchUpdatePrintShop(props.shopId, payload)
      } else {
        await fetchCreatePrintShop(payload)
      }
      ElMessage.success('保存成功')
      dialogVisible.value = false
      emit('success')
    } catch {
      // 错误提示已由请求层统一弹出
    } finally {
      submitting.value = false
    }
  }
</script>
