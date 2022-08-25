-- 任务表创建
CREATE TABLE `cron_config`  (
    `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
    `name` varchar(255) default '' COMMENT '任务名称',
    `spec` varchar(32) default '' COMMENT '执行时间，表达式',
    `protocol` tinyint(2) default 0 COMMENT '协议:1.http、2.grpc、3.系统命令行',
    `command` varchar(255) default '' COMMENT '命令',
    `remark` varchar(255) default '' COMMENT '备注',
    `status` tinyint(2) DEFAULT 1 COMMENT '状态：1.停止、2.启用',
    `create_dt` datetime NULL COMMENT '创建时间',
    `update_dt` datetime NULL COMMENT '更新时间',
    PRIMARY KEY (`id`)
);

-- 任务日志表创建
CREATE TABLE `cron_log` (
    `id` INT ( 11 ) NOT NULL AUTO_INCREMENT COMMENT '主键',
    `conf_id` INT ( 11 ) NOT NULL COMMENT '任务编号',
    `create_dt` datetime NULL COMMENT '创建时间',
    `status` TINYINT ( 2 ) NULL COMMENT '状态：1.错误、2.成功',
    `body` text NULL COMMENT '日志内容',
    `snap` text NULL COMMENT '任务快照',
    PRIMARY KEY ( `id` ),
    INDEX ( `conf_id` )
);