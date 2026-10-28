declare namespace Api {
  namespace Print {
    /** 打印门店 */
    interface ShopItem {
      id: string
      name: string
      address: string
      longitude: number
      latitude: number
      avatar: string
      /** 营业时段 JSON；空表示全天营业 */
      businessHours: string | null
      /** 0 停业 / 1 营业 */
      status: number
      /** 0 待审核 / 1 通过 / 2 拒绝 */
      auditStatus: number
      auditRemark: string
      remark: string
      createdAt?: string
    }

    interface ShopSearchParams {
      name?: string
      status?: number
      auditStatus?: number
      pageNo?: number
      pageSize?: number
    }

    interface ShopSaveParams {
      name: string
      address?: string
      longitude?: number
      latitude?: number
      avatar?: string
      businessHours?: string
      status: number
      remark?: string
    }

    /** 门店成员 */
    interface StaffItem {
      id: string
      userId: string
      shopId: string
      /** owner 店主 / staff 店员 */
      role: string
      status: number
      createdAt?: string
    }

    /** 打印设备 */
    interface DeviceItem {
      id: string
      printShopId: string
      code: string
      name: string
      clientVersion: string
      /** 0 禁用 / 1 启用 */
      status: number
      /** 0 离线 / 1 在线 */
      onlineStatus: number
      ip: string
      lastHeartbeatAt: string | null
      lastLoginUserId: string
      lastErrorInfo: string
      remark: string
      createdAt?: string
    }

    interface DeviceSearchParams {
      printShopId?: string
      keyword?: string
      status?: number
      onlineStatus?: number
      pageNo?: number
      pageSize?: number
    }

    interface DeviceSaveParams {
      code?: string
      printShopId: string
      name?: string
      status: number
      remark?: string
    }

    /** 打印机（Windows 打印队列） */
    interface PrinterItem {
      id: string
      printShopId: string
      deviceId: string
      name: string
      /** Windows 队列名，客户端静默打印按此指定 */
      systemName: string
      model: string
      /** 逗号分隔：A4,A3 */
      paperTypes: string
      colorSupport: number
      duplexSupport: number
      status: number
      remark: string
      createdAt?: string
    }

    interface PrinterSearchParams {
      printShopId?: string
      deviceId?: string
      status?: number
      pageNo?: number
      pageSize?: number
    }

    interface PrinterSaveParams {
      deviceId: string
      name: string
      systemName: string
      model?: string
      paperTypes: string
      colorSupport: number
      duplexSupport: number
      status: number
      remark?: string
    }

    /** 价格规则：门店 × 纸张 × 色彩 × 单双面 唯一 */
    interface PriceItem {
      id: string
      printShopId: string
      paperType: string
      colorType: string
      printMode: string
      /** 单价（分/页） */
      price: number
      status: number
      createdAt?: string
    }

    interface PriceSearchParams {
      printShopId?: string
      pageNo?: number
      pageSize?: number
    }

    interface PriceSaveParams {
      printShopId: string
      paperType: string
      colorType: string
      printMode: string
      price: number
      status: number
    }

    /** 可售组合体检提示 */
    interface CapabilityWarning {
      paperType: string
      colorType: string
      printMode: string
      message: string
    }

    /** 打印订单 */
    interface OrderItem {
      id: string
      orderNo: string
      printShopId: string
      memberUserId: string
      deviceId: string
      pickupCode: string
      pickupExpireTime: string | null
      itemCount: number
      totalPages: number
      totalSheets: number
      /** 金额单位：分 */
      totalAmount: number
      discountPercent: number
      discountAmount: number
      payAmount: number
      refundAmount: number
      /** 0 未支付 / 1 已支付 / 2 退款中 / 3 已退款 */
      payStatus: number
      payType: string
      /** 0 待支付 / 1 待打印 / 2 打印中 / 3 待取件 / 4 已完成 / 5 已取消 / 6 退款中 / 7 已退款 / 8 超期关闭 */
      orderStatus: number
      remark: string
      createdAt?: string
    }

    interface OrderSearchParams {
      printShopId?: string
      orderNo?: string
      orderStatus?: number
      payStatus?: number
      pageNo?: number
      pageSize?: number
    }

    /** 打印文件（会员上传，后台只读） */
    interface FileItem {
      id: string
      memberUserId: string
      infraFileId: string
      /** 可打印 PDF 的 infra_file ID；与源文件相同表示无需转换 */
      pdfFileId: string
      fileName: string
      /** pdf / docx / jpg / jpeg / png */
      fileExt: string
      fileSize: number
      /** 可打印页数，转换完成前为 0 */
      pageCount: number
      /** 0 待处理 / 1 转换中 / 2 可打印 / 3 转换失败 */
      convertStatus: number
      convertError: string
      expireTime: string | null
      createdAt?: string
    }

    interface FileSearchParams {
      memberId?: string
      fileName?: string
      convertStatus?: number
      pageNo?: number
      pageSize?: number
    }

    /** 控制台打印概览，金额单位：分 */
    interface StatisticsOverview {
      todayOrders: number
      todayPaidOrders: number
      todayPaidAmount: number
      pendingPickup: number
      printingOrders: number
      failedTasks: number
      onlineDevices: number
    }

    /** 轮播图：printShopId = 0 表示全平台 */
    interface BannerItem {
      id: string
      printShopId: string
      title: string
      imageUrl: string
      /** none 不跳转 / page 站内页 / url 外链 */
      linkType: string
      linkUrl: string
      sort: number
      status: number
      createdAt?: string
    }

    interface BannerSearchParams {
      printShopId?: string
      pageNo?: number
      pageSize?: number
    }

    interface BannerSaveParams {
      printShopId: string
      title?: string
      imageUrl: string
      linkType: string
      linkUrl?: string
      sort: number
      status: number
    }

    /** 订单明细，同时是一条打印任务 */
    interface OrderTaskItem {
      id: string
      orderId: string
      taskNo: string
      fileId: string
      fileName: string
      filePageCount: number
      paperType: string
      colorType: string
      printMode: string
      orientation: string
      pageRange: string
      pageCount: number
      copies: number
      /** 印张数 = 页数 × 份数 */
      sheetCount: number
      /** 行金额（分） */
      amount: number
      /** 0 待打印 / 1 已分配 / 2 打印中 / 3 已完成 / 4 失败 / 5 已取消 */
      status: number
      deviceId: string
      printerId: string
      retryCount: number
      errorCode: string
      failReason: string
      printedPages: number
      startedAt: string | null
      finishedAt: string | null
    }

    /** 订单详情：订单 + 明细 */
    interface OrderDetail {
      order: OrderItem
      items: OrderTaskItem[]
    }

    interface PickupParams {
      printShopId: string
      code: string
    }
  }
}
