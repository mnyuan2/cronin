#!/bin/bash

# ========================================
# 服务健康监控脚本 (Linux Shell 版本)
# 当状态码非200时发送HTTP POST通知到企业微信
# 使用方法：
#   crontab -e
#      0 * * * * bash /path/to/check.sh > /dev/check.log 2>&1 &
# ========================================

# 配置监控的URL
health_url="http://127.0.0.1:9013/health"

# 使用curl获取状态码和响应体
# -s: 静默模式
# -w "%{http_code}": 输出HTTP状态码
# -o response_body.tmp: 将响应体输出到临时文件
http_code=$(curl -s -o response_body.tmp -w "%{http_code}" "$health_url")
response_body=$(cat response_body.tmp) # 读取临时文件内容到变量
# rm -f response_body.tmp # 删除临时文件
if [ "$http_code" -eq "000" ]; then
    response_body="服务未启动"
fi


# 显示当前状态码和响应 (调试时可取消注释)
# echo "[$(date '+%Y-%m-%d %H:%M:%S')] 健康检查状态码: $http_code 响应内容: $response_body"

# 检查状态码是否为200
if [ "$http_code" -eq 200 ]; then
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] 服务正常"
else
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] 服务异常"

    # 配置通知的Webhook URL（请替换为您的实际通知地址）
    webhook_url="https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=xxxxx"
    # 构造JSON消息体，注意转义双引号
    wechat_body="{\"msgtype\":\"text\",\"text\":{\"content\":\"服务健康检查失败：状态码 ${http_code}，响应体: ${response_body}\"}}"

    # 发送POST请求到Webhook
    push_response=$(curl -s -X POST -H "Content-Type: application/json" -d "$wechat_body" "$webhook_url")
    echo "消息结果: $push_response"
fi