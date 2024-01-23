# cronin 服务器定时任务

#### 介绍
linux、windows服务器定时任务管理平台。
cronin 采用网页的形式对定时任务管理设置；
支持http请求、cmd/shell脚本、grpc请求、sql脚本，周期循环任务或单次脚本任务。
任务请求记录日志，并支持错误告警通知。




#### 安装教程
1.  创建配置文件。
2.  编译服务 或者直接下载发行版
    GOOS=linux go build -ldflags "-X main.version=v0.4.1 -X main.isBuildResource=true" -o cronin.xx ./main.go
3.  linux服务端运行服务
    ./cronin.xx

#### 使用说明

1.  入口页地址： http://127.0.0.1:9003/
2.  设置自己的任务
![image](./work/set.png)
![image](./work/list.png)


