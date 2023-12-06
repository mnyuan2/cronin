-- 任务设置
CREATE TABLE `cron_setting`  (
    `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
    `scene` varchar(255) NULL COMMENT '使用场景',
    `title` varchar(255) NULL COMMENT '名称',
    `env` varchar(255) NULL COMMENT '环境:system.系统信息、其它.业务环境信息',
    `content` text NULL COMMENT '内容',
    `create_dt` datetime(0) NULL COMMENT '创建时间',
    `update_dt` datetime(0) NULL COMMENT '更新时间',
    `status` tinyint(2) NULL DEFAULT 1 COMMENT '状态:枚举由业务定义',
    PRIMARY KEY (`id`),
    INDEX `env`(`env`, `key`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- 默认数据
INSERT INTO `cron_setting`(`key`, `title`, `env`, `create_dt`) VALUES ('env', '默认环境', 'system', '2023-11-29 00:00:00');

-- 任务表创建
CREATE TABLE `cron_config`  (
    `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
    `entry_id` int(11) NULL DEFAULT 0 COMMENT '执行队列id，status启用时有效',
    `name` varchar(255) default '' COMMENT '任务名称',
    `spec` varchar(32) default '' COMMENT '执行时间，表达式',
    `protocol` tinyint(2) default 0 COMMENT '协议:1.http、2.grpc、3.系统命令行',
    `command` json null COMMENT '命令',
    `remark` varchar(255) default '' COMMENT '备注',
    `status` tinyint(2) DEFAULT 1 COMMENT '状态：1.停止、2.启用',
    `type` tinyint(2) NULL DEFAULT 1 COMMENT '类型：1.周期任务（默认）、2.单次任务',
    `create_dt` datetime NULL COMMENT '创建时间',
    `update_dt` datetime NULL COMMENT '更新时间',
    PRIMARY KEY (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 任务日志表创建
CREATE TABLE `cron_log` (
    `id` INT ( 11 ) NOT NULL AUTO_INCREMENT COMMENT '主键',
    `conf_id` INT ( 11 ) NOT NULL COMMENT '任务编号',
    `create_dt` datetime NULL COMMENT '创建时间',
    `duration` double(10, 3) NULL default 0 COMMENT '耗时/秒',
    `status` TINYINT ( 2 ) NULL COMMENT '状态：1.错误、2.成功',
    `body` text NULL COMMENT '日志内容',
    `snap` text NULL COMMENT '任务快照',
    PRIMARY KEY ( `id` ),
    INDEX ( `conf_id` )
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;