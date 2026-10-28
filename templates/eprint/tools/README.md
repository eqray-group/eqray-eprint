# SumatraPDF 打印引擎

`SumatraPDF.exe` 是随客户端分发的便携版 PDF 打印引擎，用于把后端任务 PDF 静默打印到指定 Windows 打印队列。

| 项 | 值 |
| --- | --- |
| 版本 | 3.6.1（64 位） |
| 下载地址 | https://www.sumatrapdfreader.org/dl/rel/3.6.1/SumatraPDF-3.6.1-64.zip |
| 官方页面 | https://www.sumatrapdfreader.org/download-free-pdf-viewer |
| 许可证 | GPLv3，全文见 `COPYING-GPLv3.txt` |
| SHA-256 | `719F689B34F47BE8CA105CE8484948474DAFDE0E106BAB599E4A89326070C3D0` |

客户端启动打印时会追加 `-appdata <临时目录>`，因此不会在 `tools` 目录下生成
`SumatraPDF-settings.txt`。如需升级，下载新的 64 位便携版并替换本文件即可；若升级后
文件名变化，可同步修改 `App.config` 的 `SumatraPdfPath`，或保留 `SumatraPDF.exe` 这个名称。

SumatraPDF 以 GPLv3 分发，源码可从官方仓库获取：https://github.com/sumatrapdfreader/sumatrapdf
