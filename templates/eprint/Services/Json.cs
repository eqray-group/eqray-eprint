using System.Web.Script.Serialization;

namespace eprint.Services
{
    /// <summary>
    /// 轻量 JSON 序列化封装，基于 .NET Framework 自带的 JavaScriptSerializer，
    /// 避免客户端额外依赖第三方 JSON 库。
    /// </summary>
    public static class Json
    {
        private static readonly JavaScriptSerializer Serializer = new JavaScriptSerializer
        {
            MaxJsonLength = int.MaxValue
        };

        /// <summary>
        /// 把对象序列化为 JSON 字符串。
        /// </summary>
        public static string Serialize(object value)
        {
            return Serializer.Serialize(value);
        }

        /// <summary>
        /// 把 JSON 字符串反序列化为指定类型；字段名不区分大小写。
        /// </summary>
        public static T Deserialize<T>(string json)
        {
            return Serializer.Deserialize<T>(json);
        }
    }
}
