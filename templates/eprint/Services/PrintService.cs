using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.IO;
using System.Net;
using System.Net.Http;
using System.Text;
using System.Threading.Tasks;
using eprint.Models;

namespace eprint.Services
{
    /// <summary>
    /// 任务文件下载与 SumatraPDF 静默打印。
    /// 打印参数完全来自服务端任务，客户端不自行改写纸张 / 彩色 / 双面。
    /// </summary>
    public static class PrintService
    {
        private const int PrintTimeoutMs = 180000;

        private static readonly HttpClient DownloadClient = new HttpClient
        {
            Timeout = TimeSpan.FromSeconds(120)
        };

        /// <summary>
        /// 下载任务 PDF 到指定目录，返回本地文件路径。
        /// fileUrl 支持绝对地址和 /storage 开头的相对地址；下载失败抛 PrintException。
        /// </summary>
        public static async Task<string> DownloadAsync(TaskDto task, string directory)
        {
            if (task == null)
            {
                throw new PrintException("FILE_BROKEN", "打印任务为空");
            }

            Uri uri = ResolveFileUri(task.FileUrl);
            Directory.CreateDirectory(directory);
            string path = Path.Combine(directory, SafeFileName(task.TaskNo) + ".pdf");

            try
            {
                using (HttpResponseMessage response = await DownloadClient
                    .GetAsync(uri, HttpCompletionOption.ResponseHeadersRead)
                    .ConfigureAwait(false))
                {
                    if (!response.IsSuccessStatusCode)
                    {
                        string code = response.StatusCode == HttpStatusCode.NotFound ? "FILE_BROKEN" : "UNKNOWN";
                        throw new PrintException(code, "下载文件失败：HTTP " + (int)response.StatusCode);
                    }

                    using (Stream input = await response.Content.ReadAsStreamAsync().ConfigureAwait(false))
                    using (FileStream output = new FileStream(path, FileMode.Create, FileAccess.Write, FileShare.None))
                    {
                        await input.CopyToAsync(output).ConfigureAwait(false);
                    }
                }
            }
            catch (PrintException)
            {
                TryDelete(path);
                throw;
            }
            catch (Exception ex)
            {
                TryDelete(path);
                throw new PrintException("UNKNOWN", "下载文件失败：" + ex.Message, ex);
            }

            if (!File.Exists(path) || new FileInfo(path).Length == 0)
            {
                TryDelete(path);
                throw new PrintException("FILE_BROKEN", "文件内容为空");
            }

            return path;
        }

        /// <summary>
        /// 用 SumatraPDF 把 PDF 静默打印到任务指定的打印队列。
        /// 退出码 2/3 视为文件异常，4/5 视为打印机异常，其它按未知错误处理。
        /// </summary>
        public static async Task PrintAsync(TaskDto task, string pdfPath)
        {
            if (task == null || string.IsNullOrWhiteSpace(task.SystemName))
            {
                throw new PrintException("PRINTER_OFFLINE", "任务未指定打印机");
            }
            if (string.IsNullOrWhiteSpace(pdfPath) || !File.Exists(pdfPath))
            {
                throw new PrintException("FILE_BROKEN", "待打印文件不存在");
            }

            string engine = ResolveEnginePath();
            string appDataDir = Path.Combine(Path.GetTempPath(), "eprint", "sumatra");
            Directory.CreateDirectory(appDataDir);

            var startInfo = new ProcessStartInfo
            {
                FileName = engine,
                Arguments = BuildArguments(task, BuildSettings(task), pdfPath, appDataDir),
                UseShellExecute = false,
                CreateNoWindow = true,
                RedirectStandardOutput = true,
                RedirectStandardError = true,
                WorkingDirectory = Path.GetDirectoryName(engine) ?? AppDomain.CurrentDomain.BaseDirectory
            };

            Process process;
            try
            {
                process = Process.Start(startInfo);
            }
            catch (Exception ex)
            {
                throw new PrintException("UNKNOWN", "启动打印引擎失败：" + ex.Message, ex);
            }

            if (process == null)
            {
                throw new PrintException("UNKNOWN", "打印引擎未能启动");
            }

            using (process)
            {
                Task<string> outputTask = process.StandardOutput.ReadToEndAsync();
                Task<string> errorTask = process.StandardError.ReadToEndAsync();

                bool exited = await Task.Run(() => process.WaitForExit(PrintTimeoutMs)).ConfigureAwait(false);
                if (!exited)
                {
                    try
                    {
                        process.Kill();
                    }
                    catch
                    {
                        // 进程已退出或无法结束时不阻断回执
                    }

                    throw new PrintException("UNKNOWN", "打印超时（" + (PrintTimeoutMs / 1000) + " 秒）");
                }

                string output = await outputTask.ConfigureAwait(false);
                string error = await errorTask.ConfigureAwait(false);
                if (process.ExitCode == 0)
                {
                    return;
                }

                string reason = DescribeExitCode(process.ExitCode);
                string detail = CleanOutput(error);
                if (detail.Length == 0)
                {
                    detail = CleanOutput(output);
                }
                if (detail.Length > 0)
                {
                    if (detail.Length > 160)
                    {
                        detail = detail.Substring(0, 160);
                    }
                    reason += "：" + detail;
                }

                throw new PrintException(MapExitCode(process.ExitCode), reason);
            }
        }

        /// <summary>
        /// 查找 SumatraPDF 可执行文件：配置路径 → tools 目录 → 程序目录 → Program Files。
        /// 找不到时抛 PrintException，提示检查部署文件。
        /// </summary>
        public static string ResolveEnginePath()
        {
            string baseDir = AppDomain.CurrentDomain.BaseDirectory;
            var candidates = new List<string>();

            if (!string.IsNullOrWhiteSpace(AppConfig.SumatraPdfPath))
            {
                string configured = AppConfig.SumatraPdfPath.Trim();
                candidates.Add(Path.IsPathRooted(configured) ? configured : Path.Combine(baseDir, configured));
            }

            candidates.Add(Path.Combine(baseDir, "tools", "SumatraPDF.exe"));
            candidates.Add(Path.Combine(baseDir, "tools", "SumatraPDF-3.6.1-64.exe"));
            candidates.Add(Path.Combine(baseDir, "tools", "SumatraPDF-3.6.1.exe"));
            candidates.Add(Path.Combine(baseDir, "SumatraPDF.exe"));

            string programFiles = Environment.GetFolderPath(Environment.SpecialFolder.ProgramFiles);
            if (!string.IsNullOrWhiteSpace(programFiles))
            {
                candidates.Add(Path.Combine(programFiles, "SumatraPDF", "SumatraPDF.exe"));
            }
            string programFilesX86 = Environment.GetFolderPath(Environment.SpecialFolder.ProgramFilesX86);
            if (!string.IsNullOrWhiteSpace(programFilesX86))
            {
                candidates.Add(Path.Combine(programFilesX86, "SumatraPDF", "SumatraPDF.exe"));
            }

            foreach (string candidate in candidates)
            {
                if (!string.IsNullOrWhiteSpace(candidate) && File.Exists(candidate))
                {
                    return Path.GetFullPath(candidate);
                }
            }

            throw new PrintException("UNKNOWN",
                "未找到 SumatraPDF 打印引擎，请把 SumatraPDF.exe 放到 tools 目录，或配置 SumatraPdfPath");
        }

        /// <summary>
        /// 按服务端任务拼 SumatraPDF 的 -print-settings：
        /// 页码范围、纸张、彩色、单双面、方向与份数，顺序不影响解析。
        /// </summary>
        private static string BuildSettings(TaskDto task)
        {
            var settings = new List<string>();

            string pageRange = NormalizePageRange(task.PageRange);
            if (pageRange.Length > 0)
            {
                settings.Add(pageRange);
            }

            settings.Add(string.Equals(task.PaperType, "A3", StringComparison.OrdinalIgnoreCase)
                ? "paper=A3"
                : "paper=A4");
            settings.Add(string.Equals(task.ColorType, "color", StringComparison.OrdinalIgnoreCase)
                ? "color"
                : "monochrome");
            settings.Add(string.Equals(task.PrintMode, "double", StringComparison.OrdinalIgnoreCase)
                ? "duplexlong"
                : "simplex");
            settings.Add(string.Equals(task.Orientation, "landscape", StringComparison.OrdinalIgnoreCase)
                ? "landscape"
                : "portrait");

            if (task.Copies > 1)
            {
                settings.Add(task.Copies + "x");
            }

            return string.Join(",", settings);
        }

        /// <summary>
        /// 只保留数字、逗号和短横线，避免把非法字符带进 SumatraPDF 参数。
        /// 返回空串表示忽略页码范围，按整份文件打印。
        /// </summary>
        private static string NormalizePageRange(string pageRange)
        {
            if (string.IsNullOrWhiteSpace(pageRange))
            {
                return string.Empty;
            }

            var builder = new StringBuilder();
            foreach (char value in pageRange.Trim())
            {
                if (char.IsDigit(value) || value == ',' || value == '-')
                {
                    builder.Append(value);
                }
            }

            return builder.ToString().Trim(',');
        }

        /// <summary>
        /// 拼 SumatraPDF 命令行，所有动态值都按 Windows 参数规则加引号。
        /// </summary>
        private static string BuildArguments(TaskDto task, string settings, string pdfPath, string appDataDir)
        {
            var parts = new List<string>
            {
                "-print-to",
                Quote(task.SystemName),
                "-print-settings",
                Quote(settings),
                "-silent",
                "-appdata",
                Quote(appDataDir),
                Quote(pdfPath)
            };

            return string.Join(" ", parts);
        }

        /// <summary>
        /// 把文件地址转成绝对 URI：绝对地址直接使用，/storage 相对地址按服务端根地址补全。
        /// </summary>
        private static Uri ResolveFileUri(string fileUrl)
        {
            if (string.IsNullOrWhiteSpace(fileUrl))
            {
                throw new PrintException("FILE_BROKEN", "任务缺少文件地址");
            }

            Uri absolute;
            if (Uri.TryCreate(fileUrl, UriKind.Absolute, out absolute))
            {
                return absolute;
            }

            string baseUrl = AppSession.Api == null ? AppConfig.ApiBaseUrl : AppSession.Api.BaseUrl;
            Uri baseUri;
            if (!Uri.TryCreate(baseUrl.TrimEnd('/') + "/", UriKind.Absolute, out baseUri) ||
                !Uri.TryCreate(baseUri, fileUrl, out absolute))
            {
                throw new PrintException("FILE_BROKEN", "文件地址不合法");
            }

            return absolute;
        }

        /// <summary>
        /// 把 SumatraPDF 退出码映射成服务端约定的失败错误码。
        /// 1 是实际打印失败（打印机不可用 / 驱动拒绝），按可重试的打印机异常处理。
        /// </summary>
        private static string MapExitCode(int exitCode)
        {
            switch (exitCode)
            {
                case 1:
                case 4:
                case 5:
                    return "PRINTER_OFFLINE";
                case 2:
                case 3:
                    return "FILE_BROKEN";
                default:
                    return "UNKNOWN";
            }
        }

        /// <summary>
        /// 把 SumatraPDF 退出码翻译成可读原因，避免把调试日志直接写进任务失败原因。
        /// </summary>
        private static string DescribeExitCode(int exitCode)
        {
            switch (exitCode)
            {
                case 1:
                    return "打印失败，打印机不可用或驱动拒绝任务";
                case 2:
                    return "无法打开 PDF 文件";
                case 3:
                    return "PDF 文件不允许打印";
                case 4:
                    return "打印机不存在";
                case 5:
                    return "打印机驱动或设备异常";
                case 6:
                    return "打印被策略限制";
                default:
                    return "打印引擎返回未知错误（退出码 " + exitCode + "）";
            }
        }

        /// <summary>
        /// 过滤 SumatraPDF 的调试日志，只保留真正的错误文本。
        /// </summary>
        private static string CleanOutput(string value)
        {
            if (string.IsNullOrWhiteSpace(value))
            {
                return string.Empty;
            }

            var kept = new List<string>();
            foreach (string line in value.Split(new[] { '\r', '\n' }, StringSplitOptions.RemoveEmptyEntries))
            {
                string text = line.Trim();
                if (text.Length == 0)
                {
                    continue;
                }
                if (text.StartsWith("ParseFlags:", StringComparison.Ordinal) ||
                    text.StartsWith("Starting:", StringComparison.Ordinal) ||
                    text.StartsWith("SaveSettings", StringComparison.Ordinal) ||
                    text.StartsWith("LoadSettings", StringComparison.Ordinal))
                {
                    continue;
                }

                kept.Add(text);
            }

            return string.Join(" ", kept);
        }

        /// <summary>
        /// 任务号只保留安全字符，防止生成非法临时文件名。
        /// </summary>
        private static string SafeFileName(string value)
        {
            var builder = new StringBuilder();
            foreach (char c in value ?? string.Empty)
            {
                if (char.IsLetterOrDigit(c) || c == '-' || c == '_')
                {
                    builder.Append(c);
                }
            }

            return builder.Length == 0 ? "task" : builder.ToString();
        }

        /// <summary>
        /// 按 Windows CreateProcess 规则给参数加双引号并转义内部引号 / 结尾反斜杠。
        /// </summary>
        private static string Quote(string value)
        {
            if (string.IsNullOrEmpty(value))
            {
                return "\"\"";
            }

            var builder = new StringBuilder("\"");
            int backslashes = 0;
            foreach (char c in value)
            {
                if (c == '\\')
                {
                    backslashes++;
                    continue;
                }

                if (c == '"')
                {
                    builder.Append('\\', backslashes * 2 + 1);
                    builder.Append('"');
                    backslashes = 0;
                    continue;
                }

                builder.Append('\\', backslashes);
                backslashes = 0;
                builder.Append(c);
            }

            builder.Append('\\', backslashes * 2);
            builder.Append('"');
            return builder.ToString();
        }

        /// <summary>
        /// 删除临时文件，失败时忽略，避免影响打印结果回执。
        /// </summary>
        private static void TryDelete(string path)
        {
            try
            {
                if (File.Exists(path))
                {
                    File.Delete(path);
                }
            }
            catch
            {
                // 临时文件清理失败不影响主流程
            }
        }
    }

    /// <summary>
    /// 打印异常：Code 为服务端任务失败错误码（FILE_BROKEN / PRINTER_OFFLINE / UNKNOWN 等）。
    /// </summary>
    public sealed class PrintException : Exception
    {
        public string Code { get; private set; }

        public PrintException(string code, string message)
            : base(message)
        {
            Code = code;
        }

        public PrintException(string code, string message, Exception inner)
            : base(message, inner)
        {
            Code = code;
        }
    }
}
