declare namespace Api {
  namespace Member {
    /** 会员档案（对齐 yudao member_user） */
    interface MemberItem {
      id: string
      mobile: string
      nickname: string
      avatar: string
      /** 真实姓名 */
      name: string
      /** 0 未知 / 1 男 / 2 女 */
      sex: number
      birthday: string | null
      areaId: number
      /** 用户标记 */
      mark: string
      registerIp: string
      loginIp: string
      loginDate: string | null
      /** 0 禁用 / 1 启用 */
      status: number
      creator?: string
      createdAt?: string
      updatedAt?: string
    }

    /** 会员列表 */
    interface MemberList {
      list: MemberItem[]
      total: number
    }

    /** 会员搜索参数 */
    interface MemberSearchParams {
      nickname?: string
      mobile?: string
      status?: number
      pageNo?: number
      pageSize?: number
    }

    /** 会员新增 / 编辑参数 */
    interface MemberSaveParams {
      nickname: string
      mobile: string
      avatar?: string
      name?: string
      sex: number
      /** 格式：YYYY-MM-DD HH:mm:ss */
      birthday?: string
      areaId?: number
      mark?: string
      status: number
    }

    /** 会员卡：卡自带等级与折扣，与会员等级体系无关 */
    interface CardItem {
      id: string
      memberUserId: string
      /** 卡号，服务端生成 */
      cardNo: string
      /** 卡等级编码：normal / silver / gold / black */
      cardLevel: string
      /** 卡名称，展示用 */
      cardName: string
      /** 折扣百分比：100 无折扣 / 95 表示 95 折 */
      discountPercent: number
      /** 卡类型：level 折扣卡（预留 stored 储值卡） */
      cardType: string
      /** 开卡来源：register / admin / purchase */
      source: string
      /** 0 停用 / 1 正常 / 2 已过期 */
      status: number
      startTime: string | null
      expireTime: string | null
      remark: string
      /** 派生字段：当前是否可用于折扣 */
      enabled: boolean
      createdAt?: string
      updatedAt?: string
    }

    /** 会员卡列表 */
    interface CardList {
      list: CardItem[]
      total: number
    }

    /** 会员卡搜索参数 */
    interface CardSearchParams {
      cardNo?: string
      status?: number
      pageNo?: number
      pageSize?: number
    }

    /** 会员卡发放 / 修改参数 */
    interface CardSaveParams {
      cardLevel: string
      cardName: string
      discountPercent: number
      status: number
      startTime?: string | null
      expireTime?: string | null
      remark?: string
    }

    /** 第三方登录渠道配置（对齐 yudao system_social_client） */
    interface SocialClientItem {
      id: string
      name: string
      /** wechat_mini / wechat_mp / wechat_open */
      socialType: string
      clientId: string
      /** 列表不返回密钥 */
      agentId: string
      /** 0 禁用 / 1 启用 */
      status: number
      createdAt?: string
      updatedAt?: string
    }

    interface SocialClientList {
      list: SocialClientItem[]
      total: number
    }

    interface SocialClientSearchParams {
      keyword?: string
      status?: number
      pageNo?: number
      pageSize?: number
    }

    interface SocialClientSaveParams {
      name: string
      socialType: string
      clientId: string
      /** 编辑留空表示不修改 */
      clientSecret?: string
      agentId?: string
      status: number
    }

    /** 会员第三方账号绑定 */
    interface SocialBindingItem {
      id: string
      /** wechat_mini / wechat_mp / wechat_open */
      socialType: string
      nickname: string
      avatar: string
      createdAt: string
    }

    /** 收货地址 */
    interface AddressItem {
      id: string
      userId: string
      name: string
      mobile: string
      areaId: number
      detailAddress: string
      /** 0 否 / 1 是 */
      defaultStatus: number
    }
  }
}
