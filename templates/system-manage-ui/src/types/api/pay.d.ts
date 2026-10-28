declare namespace Api {
  namespace Pay {
    /** 支付渠道配置（列表不返回密钥） */
    interface ChannelConfigItem {
      id: string
      /** 配置编码：wechat_mini / wechat_mp / wechat_h5 / alipay_wap */
      code: string
      name: string
      /** 渠道大类：wechat / alipay */
      channel: string
      /** 场景：mini / jsapi / h5 / wap */
      scene: string
      appId: string
      mchId: string
      serialNo: string
      gateway: string
      notifyUrl: string
      returnUrl: string
      /** 其他渠道参数 JSON 字符串 */
      extra: string
      /** 0 停用 / 1 启用 */
      status: number
      sort: number
      remark: string
      createdAt?: string
      updatedAt?: string
    }

    /** 渠道配置详情：编辑页回显密钥 */
    interface ChannelConfigDetail extends ChannelConfigItem {
      privateKey: string
      apiKey: string
      platformPublicKey: string
    }

    interface ChannelConfigList {
      list: ChannelConfigItem[]
      total: number
    }

    interface ChannelConfigSearchParams {
      code?: string
      channel?: string
      status?: number
      pageNo?: number
      pageSize?: number
    }

    /** 渠道配置保存参数；三个密钥留空表示不修改 */
    interface ChannelConfigSaveParams {
      code: string
      name: string
      channel: string
      scene: string
      appId?: string
      mchId?: string
      serialNo?: string
      privateKey?: string
      apiKey?: string
      platformPublicKey?: string
      gateway?: string
      notifyUrl?: string
      returnUrl?: string
      extra?: string
      status: number
      sort: number
      remark?: string
    }

    /** 支付单 */
    interface PaymentItem {
      id: string
      paymentNo: string
      /** print_order / wallet_recharge */
      bizType: string
      bizId: string
      memberUserId: string
      /** wechat / alipay / balance */
      channel: string
      channelCode: string
      scene: string
      clientIp: string
      /** 金额单位：分 */
      amount: number
      paidAmount: number
      refundAmount: number
      /** 0 待支付 / 1 支付中 / 2 已支付 / 3 失败 / 4 已关闭 / 5 已退款 */
      status: number
      channelTradeNo: string
      openId: string
      expireTime: string | null
      successTime: string | null
      closeTime: string | null
      remark: string
      createdAt?: string
    }

    interface PaymentList {
      list: PaymentItem[]
      total: number
    }

    interface PaymentSearchParams {
      paymentNo?: string
      memberId?: string
      channel?: string
      status?: number
      pageNo?: number
      pageSize?: number
    }

    /** 退款单 */
    interface RefundItem {
      id: string
      refundNo: string
      paymentId: string
      paymentNo: string
      bizType: string
      bizId: string
      memberUserId: string
      channel: string
      channelCode: string
      amount: number
      reason: string
      /** 0 待退款 / 1 退款中 / 2 成功 / 3 失败 */
      status: number
      channelRefundNo: string
      successTime: string | null
      failReason: string
      operatorId: string
      createdAt?: string
    }

    interface RefundList {
      list: RefundItem[]
      total: number
    }

    interface RefundSearchParams {
      paymentNo?: string
      memberId?: string
      status?: number
      pageNo?: number
      pageSize?: number
    }

    interface RefundCreateParams {
      paymentNo: string
      /** 金额单位：分 */
      amount: number
      reason: string
    }

    /** 会员钱包 */
    interface WalletItem {
      id?: string
      memberUserId: string
      /** 余额（分） */
      balance: number
      totalIncome: number
      totalExpense: number
      /** 0 冻结 / 1 正常 */
      status: number
      updatedAt?: string
    }

    /** 钱包流水 */
    interface WalletTransactionItem {
      id: string
      walletId: string
      memberUserId: string
      bizType: string
      bizId: string
      /** 1 收入 / 2 支出 */
      type: number
      amount: number
      /** 变动后余额（分） */
      balanceAfter: number
      title: string
      remark: string
      createdAt?: string
    }

    interface WalletTransactionSearchParams {
      memberId?: string
      type?: number
      bizType?: string
      pageNo?: number
      pageSize?: number
    }

    /** 钱包人工调账：正数加款、负数扣款 */
    interface WalletAdjustParams {
      amount: number
      reason: string
    }
  }
}
