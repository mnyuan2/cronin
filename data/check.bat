@echo off
:: 解决中文输出乱码问题
chcp 65001 >nul

setlocal enabledelayedexpansion

:: ========================================
:: 服务健康监控脚本 (Windows 版本)
:: 当状态码非200时发送HTTP POST通知
:: 使用方法：以管理员打开cmd，执行任务创建语法↓，每小时执行一次：
::   schtasks /create /tn "CroninCheck" /tr "绝对路径前缀\incron\data\check.bat" /sc hourly /mo 1
:: ========================================

:: 配置监控的URL
set "health_url=http://127.0.0.1:9013/health"

:: 使用curl
for /f "tokens=*" %%i in ('curl -s -o check.txt -w "%%{http_code}" "%health_url%"') do set "status_code=%%i"
:: 获取响应体
for /f "delims=" %%b in (check.txt) do set "response=%%b"
if "%status_code%" equ "000" (
    set "response=服务未启动"
)

:: 显示当前状态码
echo [%date% %time%] 健康检查状态码: !status_code!  响应内容: %response%

:: 检查状态码是否为200
if "%status_code%" equ "200" (
    echo [%date% %time%] 服务正常
) else (
    echo [%date% %time%] 服务异常

    :: 配置通知的Webhook URL（请替换为您的实际通知地址）
    set "webhook_url=https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=xxxxx"
    set "wechat_body={\"msgtype\":\"text\",\"text\":{\"content\":\"服务健康检查失败：%status_code%.%response%\"}}"
    :: 检查微信发送是否成功
    for /f "delims=" %%d in ('curl -s -X POST -H "Content-Type: application/json" -d "%wechat_body%" "%webhook_url%"') do (
        set "push_response=%%d"
    )
    echo 消息结果:!push_response!
)
