# cronin 内部定时任务

#### 介绍
linux服务器下crontab的平替工具。
cronin 以通过网页的形式进行定时任务的设置和管理，并展示任务近期执行结果日志。建议通过html <iframe>元素嵌入到已有的管理系统中，做为原系统的一部分；也可以独立页面进行管理。




#### 安装教程
1.  创建配置文件。
2.  编译服务 或者直接下载发行版
    GOOS=linux go build -ldflags "-X main.version=v0.4.1 -X main.isBuildResource=true" -o cronin.v0.4.1 ./main.go
3.  linux服务端运行服务
    ./cronin.xxx

#### 使用说明

1.  入口页地址： http://127.0.0.1:9003/view/cron/list
2.  设置自己的任务
![image](./work/set.png)
![image](./work/list.png)


