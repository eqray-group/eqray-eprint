import request from '@/utils/http'

// 分页查询会员卡
export function fetchGetMemberCardList(params: Api.Member.CardSearchParams) {
  return request.get<Api.Member.CardList>({
    url: '/api/v1/member/member-cards',
    params
  })
}

// 查询指定会员的会员卡
export function fetchGetMemberCard(memberId: string) {
  return request.get<Api.Member.CardItem>({
    url: `/api/v1/member/users/${memberId}/card`
  })
}

// 发放或修改会员卡（一个会员一张卡，走后端 upsert）
export function fetchSaveMemberCard(memberId: string, data: Api.Member.CardSaveParams) {
  return request.put<{ id: string }>({
    url: `/api/v1/member/users/${memberId}/card`,
    params: data
  })
}

// 停用会员卡
export function fetchDisableMemberCard(memberId: string) {
  return request.put<null>({
    url: `/api/v1/member/users/${memberId}/card/disable`
  })
}
