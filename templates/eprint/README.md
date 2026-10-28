# 易云印 Windows 云打印客户端

配合本仓库的 Goravel 后端使用，完成「登录 → 选择门店 → 登记设备 → 上报打印机能力 → 轮询任务 → 下载 PDF → 静默打印 → 回执」。

## 环境要求

- Windows 10 / 11
- .NET Framework 4.8
- MSBuild（Visual Studio 2022 生成工具即可）

## 构建与运行

```powershell
# 在 Visual Studio Developer PowerShell 中，于仓库根目录执行
msbuild "templates\eprint\eprint.csproj" /t:Rebuild /p:Configuration=Debug
```

生成结果：`templates/eprint/bin/Debug/eprint.exe`。

服务端地址以 `eprint.exe.config` 的 `ApiBaseUrl` 作为首次启动的默认值（默认 `http://127.0.0.1:11111`），
之后一律以独立配置文件为准，详见下一节。

## 独立配置文件与自动登录

客户端配置不再依赖 `exe.config`，而是保存在当前用户目录：

```
%LOCALAPPDATA%\yiyunyin\client.json
```

| 字段 | 说明 |
| --- | --- |
| `apiBaseUrl` | 服务端地址。登录页右上角「设置」里修改后立即写盘，下次启动直接使用 |
| `autoLogin` | 是否自动登录（对应登录页「自动登录（下次启动免输入）」勾选状态） |
| `username` | 上次登录账号，登录页自动预填 |
| `password` | DPAPI 加密后的密码，仅勾选自动登录时保存，登录页自动预填 |
| `refreshToken` | DPAPI 加密后的刷新令牌，自动登录的唯一凭据 |

行为说明：

- 勾选「自动登录」登录成功后，客户端会保存密码与刷新令牌；下次启动先用刷新令牌调
  `POST /api/v1/system/auth/refresh` 续期，成功后直接进入主界面，**不再弹验证码**。
- 刷新令牌默认有效期 7 天，服务端每次续期都会轮换旧令牌，客户端会立刻回写新令牌；
  只要在有效期内启动过客户端，自动登录就能一直延续。
- 令牌过期或被后台注销时，客户端会清掉本地令牌并回到登录页，提示「自动登录已失效」；
  网络不通等临时故障不会清令牌，下次启动继续尝试。
- 未勾选自动登录时只记住服务端地址与账号，密码与令牌会被清空。
- 点「退出登录」会清掉本地刷新令牌，下次启动必须重新登录（账号密码会预填，补验证码即可）。
- `password` / `refreshToken` 用 Windows DPAPI 按当前用户加密，换用户或拷到别的机器都无法解密，
  会退化成需要重新登录；需要重置配置时直接删除该文件即可。

## 使用前提

1. 后端已启动（`go run .`），并已执行 `go run . artisan migrate`。
2. 登录账号已在后台绑定打印门店（打印管理 → 门店 → 员工），否则登录后会提示「未绑定任何打印门店」。
3. 以管理员身份或对打印机有访问权限的账号运行客户端，否则读取驱动能力可能不完整。

## 当前功能

- 账号密码 + 图形验证码登录（后台账号）。
- 登录后按「机器名 + 首块网卡 MAC」生成稳定设备码（`pc-` 前缀），登记设备并刷新心跳。
- 读取 Windows 已安装打印机，自动探测 A4 / A3、彩色、双面能力。
- 表格中勾选需要上报的打印机，能力可手动调整后点「上报选中」。
- 上报为全量覆盖：新增的创建、已有的更新、**未勾选的由服务端停用（不删除）**。
- 每 5 秒拉取一次待打印任务；任务里的 `fileUrl` 可以是绝对地址，也可以是 `/storage/...` 相对地址。
- PDF 下载到本地临时目录后，调用随包附带的 `tools/SumatraPDF.exe` 静默打印。
- 打印前回执 `2`（打印中），成功后回执 `3`；失败回执 `4` 并带上错误码与原因。
- 登录页右上角齿轮可修改服务端地址，保存后按新地址重新获取验证码。
- 服务器地址、账号与自动登录凭据保存在独立配置文件，改一次下次直接生效。
- 勾选自动登录后，下次启动用刷新令牌免验证码进入主界面；退出登录会清掉本地令牌。

## 服务端接口

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | `/api/v1/system/auth/captcha` | 图形验证码 |
| POST | `/api/v1/system/auth/login` | 账号密码登录 |
| POST | `/api/v1/system/auth/refresh` | 用刷新令牌续期（自动登录） |
| POST | `/api/v1/system/auth/logout` | 退出登录 |
| GET | `/api/v1/print/device/shops` | 当前账号绑定的门店 |
| POST | `/api/v1/print/device/devices/register` | 登记 / 复用设备 |
| GET | `/api/v1/print/device/devices/{id}/printers` | 设备已上报的打印机 |
| POST | `/api/v1/print/device/devices/{id}/printers/sync` | 全量上报打印机与能力 |
| GET | `/api/v1/print/device/devices/{id}/tasks?limit=5` | 拉取待打印任务（同时刷新心跳） |
| POST | `/api/v1/print/device/devices/{id}/tasks/report` | 打印状态回执 |

设备接口使用后台访问令牌；服务端会校验当前账号是否为该门店员工。

## 打印机能力约定

- 纸张：`paperTypes` 为逗号分隔的 `A4` / `A3`，例如 `A4,A3`。
- 彩色：`colorSupport` 为 `0` 或 `1`。
- 双面：`duplexSupport` 为 `0` 或 `1`。
- 驱动探测结果只作为默认值，最终以界面上报的勾选为准。

## 自动打印流程

1. 每 5 秒调用任务接口，按 `taskNo` 去重，避免同一任务重复打印。
2. 任务必须是 `PDF`，`fileUrl` 相对地址按登录页配置的服务端地址补全。
3. 下载到临时目录后，先回执 `2`，再调用 `tools/SumatraPDF.exe`：
   - 静默打印，不弹出阅读器窗口；
   - 纸张、彩色、双面参数由服务端任务下发，客户端不自行改写。
4. SumatraPDF 退出码为 `0` 时回执 `3`；否则回执 `4`，错误码取 `PRINT_FAILED` /
   `PRINTER_OFFLINE` / `FILE_BROKEN` 等，原因写入服务端便于排查。
5. 临时 PDF 在打印结束后删除。

## 发布注意

- `tools/SumatraPDF.exe` 已随项目打包，构建时复制到输出目录 `tools/`；
  不要只拷贝 `eprint.exe`，否则打印时会提示找不到引擎。
- 首次运行建议以管理员身份启动，确保能读取所有打印机驱动能力。
- 客户端不持久化令牌，重启后需要重新登录。
