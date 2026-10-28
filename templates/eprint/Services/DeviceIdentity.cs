using System;
using System.Net.NetworkInformation;
using System.Security.Cryptography;
using System.Text;

namespace eprint.Services
{
    /// <summary>
    /// 生成稳定的设备机器码，作为同一台电脑重复登录时复用设备记录的依据。
    /// </summary>
    public static class DeviceIdentity
    {
        private static string _cachedCode;

        /// <summary>
        /// 用「计算机名 + 第一块启用网卡的 MAC」计算 SHA-256，取前 24 位作为设备编码。
        /// 返回形如 pc-xxxxxxxxxxxxxxxxxxxxxxxx 的字符串；取不到 MAC 时只用计算机名。
        /// </summary>
        public static string GetDeviceCode()
        {
            if (!string.IsNullOrEmpty(_cachedCode))
            {
                return _cachedCode;
            }

            string raw = Environment.MachineName + "|" + TryGetMacAddress();

            using (SHA256 sha = SHA256.Create())
            {
                byte[] hash = sha.ComputeHash(Encoding.UTF8.GetBytes(raw));
                var builder = new StringBuilder(hash.Length * 2);
                foreach (byte value in hash)
                {
                    builder.Append(value.ToString("x2"));
                }

                _cachedCode = "pc-" + builder.ToString(0, 24);
            }

            return _cachedCode;
        }

        /// <summary>
        /// 读取第一块处于启用状态、非回环网卡的物理地址；读取失败返回空串。
        /// </summary>
        private static string TryGetMacAddress()
        {
            try
            {
                foreach (NetworkInterface nic in NetworkInterface.GetAllNetworkInterfaces())
                {
                    if (nic.OperationalStatus != OperationalStatus.Up ||
                        nic.NetworkInterfaceType == NetworkInterfaceType.Loopback)
                    {
                        continue;
                    }

                    PhysicalAddress address = nic.GetPhysicalAddress();
                    if (address != null && address.GetAddressBytes().Length > 0)
                    {
                        return address.ToString();
                    }
                }
            }
            catch
            {
                // 网卡信息读取失败时退化为仅用计算机名，不影响登录与登记
            }

            return string.Empty;
        }
    }
}
