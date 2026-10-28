<!-- 会员详情抽屉：基础资料 + 会员卡 + 收货地址 + 第三方绑定（只读） -->
<template>
  <ElDrawer v-model="drawerVisible" title="会员详情" size="640px" destroy-on-close>
    <div v-loading="loading">
      <div class="flex items-center mb-5">
        <ElAvatar :size="64" :src="member?.avatar || defaultAvatar" />
        <div class="ml-4">
          <p class="text-lg font-medium">{{ member?.nickname || '-' }}</p>
          <p class="mt-1 text-sm text-g-600">{{ member?.mobile || '-' }}</p>
        </div>
      </div>

      <ElDescriptions :column="2" border>
        <ElDescriptionsItem label="真实姓名">{{ member?.name || '-' }}</ElDescriptionsItem>
        <ElDescriptionsItem label="性别">{{ SEX_TEXT[member?.sex ?? 0] ?? '未知' }}</ElDescriptionsItem>
        <ElDescriptionsItem label="生日">{{ formatDate(member?.birthday) }}</ElDescriptionsItem>
        <ElDescriptionsItem label="状态">
          <ElTag :type="member?.status === 1 ? 'success' : 'info'" effect="plain">
            {{ member?.status === 1 ? '启用' : '禁用' }}
          </ElTag>
        </ElDescriptionsItem>
        <ElDescriptionsItem label="用户标记" :span="2">{{ member?.mark || '-' }}</ElDescriptionsItem>
        <ElDescriptionsItem label="注册时间" :span="2">{{ member?.createdAt || '-' }}</ElDescriptionsItem>
      </ElDescriptions>

      <ElDivider content-position="left">会员卡</ElDivider>
      <ElDescriptions v-if="card" :column="2" border>
        <ElDescriptionsItem label="卡号">{{ card.cardNo }}</ElDescriptionsItem>
        <ElDescriptionsItem label="卡名称">{{ card.cardName }}</ElDescriptionsItem>
        <ElDescriptionsItem label="卡等级">{{ card.cardLevel }}</ElDescriptionsItem>
        <ElDescriptionsItem label="折扣">
          {{ card.discountPercent === 100 ? '无折扣' : `${card.discountPercent / 10} 折` }}
        </ElDescriptionsItem>
        <ElDescriptionsItem label="有效期">
          {{ card.expireTime ? card.expireTime : '永久有效' }}
        </ElDescriptionsItem>
        <ElDescriptionsItem label="可用状态">
          <ElTag :type="card.enabled ? 'success' : 'info'" effect="plain">
            {{ card.enabled ? '可用' : '不可用' }}
          </ElTag>
        </ElDescriptionsItem>
      </ElDescriptions>
      <ElEmpty v-else description="该会员还没有会员卡" :image-size="72" />

      <ElTabs v-model="activeTab" class="mt-4">
        <ElTabPane label="收货地址" name="address">
          <ElTable :data="addresses" size="small" border>
            <ElTableColumn prop="name" label="收件人" width="100" />
            <ElTableColumn prop="mobile" label="手机号" width="130" />
            <ElTableColumn prop="detailAddress" label="详细地址" show-overflow-tooltip />
            <ElTableColumn label="默认" width="70">
              <template #default="{ row }">
                <ElTag v-if="row.defaultStatus === 1" type="success" effect="plain">默认</ElTag>
                <span v-else>-</span>
              </template>
            </ElTableColumn>
          </ElTable>
        </ElTabPane>
        <ElTabPane label="第三方绑定" name="social">
          <ElTable :data="bindings" size="small" border>
            <ElTableColumn prop="socialType" label="渠道" width="140" />
            <ElTableColumn prop="createdAt" label="绑定时间" />
          </ElTable>
        </ElTabPane>
      </ElTabs>
    </div>
  </ElDrawer>
</template>

<script setup lang="ts">
  import defaultAvatar from '@imgs/user/avatar.webp'
  import { fetchGetMemberDetail, fetchGetMemberAddresses, fetchGetMemberSocialBindings } from '@/api/member/user'
  import { fetchGetMemberCard } from '@/api/member/card'

  defineOptions({ name: 'MemberDetail' })

  interface Props {
    visible: boolean
    memberId: string
  }
  interface Emits {
    (e: 'update:visible', value: boolean): void
    (e: 'updated'): void
  }

  const props = defineProps<Props>()
  const emit = defineEmits<Emits>()

  const SEX_TEXT: Record<number, string> = { 0: '未知', 1: '男', 2: '女' }

  const drawerVisible = computed({
    get: () => props.visible,
    set: (value) => emit('update:visible', value)
  })

  const loading = ref(false)
  const activeTab = ref('address')
  const member = ref<Api.Member.MemberItem | null>(null)
  const card = ref<Api.Member.CardItem | null>(null)
  const addresses = ref<Api.Member.AddressItem[]>([])
  const bindings = ref<Api.Member.SocialBindingItem[]>([])

  // 日期只展示年月日，接口返回的是 datetime
  const formatDate = (value?: string): string => (value ? value.slice(0, 10) : '-')

  // 会员卡未发放时后端返回业务错误码，这里吞掉错误按「未发卡」展示
  const loadCard = async (memberId: string): Promise<void> => {
    card.value = null
    try {
      card.value = await fetchGetMemberCard(memberId)
    } catch {
      card.value = null
    }
  }

  const loadDetail = async (): Promise<void> => {
    if (!props.memberId) return

    loading.value = true
    try {
      const [detail, addressList, bindingList] = await Promise.all([
        fetchGetMemberDetail(props.memberId),
        fetchGetMemberAddresses(props.memberId),
        fetchGetMemberSocialBindings(props.memberId)
      ])
      member.value = detail
      addresses.value = addressList
      bindings.value = bindingList
      await loadCard(props.memberId)
    } catch {
      // 错误提示已由请求层统一弹出
    } finally {
      loading.value = false
    }
  }

  watch(
    () => [props.visible, props.memberId] as const,
    ([visible]) => {
      if (!visible) return
      loadDetail()
    }
  )
</script>
