import request from '@/utils/http'

// 分页查询支付流水
export function fetchGetPaymentList(params: Api.Pay.PaymentSearchParams) {
  return request.get<Api.Pay.PaymentList>({
    url: '/api/v1/pay/payments',
    params
  })
}

// 会员钱包人工调账：正数加款、负数扣款
export function fetchAdjustWallet(memberId: string, data: Api.Pay.WalletAdjustParams) {
  return request.post<{ balance: number }>({
    url: `/api/v1/pay/wallets/${memberId}/adjust`,
    params: data
  })
}
