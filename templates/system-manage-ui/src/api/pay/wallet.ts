import request from '@/utils/http'

// 查询会员钱包余额与累计收支
export function fetchGetWallet(memberId: string) {
  return request.get<Api.Pay.WalletItem>({
    url: `/api/v1/pay/wallets/${memberId}`
  })
}

// 分页查询钱包流水
export function fetchGetWalletTransactions(params: Api.Pay.WalletTransactionSearchParams) {
  return request.get<{ list: Api.Pay.WalletTransactionItem[]; total: number }>({
    url: '/api/v1/pay/wallet-transactions',
    params
  })
}
