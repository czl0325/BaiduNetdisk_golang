# BaiduNetdisk_golang
采用golang从0开始打造百度网盘服务端，包含文件上传下载，不同用户文件隔离，文件的秒传，断点续传，分块上传，离线下载，分布式云存储。


### mysql数据库

* 文件表

```mysql
create table `tbl_file` (
    `id` int(11) NOT NULL auto_increment,
    `file_sha1` char(40) NOT NULL default '' COMMENT '文件hash',
    `file_name` varchar(256) NOT NULL default '' COMMENT '文件名',
    `file_size` bigint(20) default '0' COMMENT '文件大小',
    `file_addr` varchar(1024) NOT NULL default '' COMMENT '文件存储路径',
    `create_time` datetime default NOW() COMMENT '创建日期',
    `update_time` datetime default NOW() COMMENT '更新日期',
    `status` tinyint default '1' COMMENT '状态(可用/禁用/已删除)',
    `ext1` int(11) default '0' COMMENT '备用字段1',
    `ext2` text COMMENT '备用字段2',
    primary key (`id`),
    UNIQUE KEY `idx_file_hash` (`file_sha1`),
    key `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

* 用户表

```mysql
create table `tbl_user` (
    `id` int(11) NOT NULL auto_increment,
    `user_name` varchar(40) NOT NULL default '' COMMENT '用户名',
    `user_pwd` varchar(256) NOT NULL default '' COMMENT '用户密码',
    `email` varchar(60) default '' COMMENT '邮箱',
    `phone` varchar(40) default '' COMMENT '电话号码',
    `email_validated` tinyint(1) default 0 COMMENT '邮箱是否已验证',
    `phone_validated` tinyint(1) default 0 COMMENT '手机是否已验证',
    `create_time` datetime default NOW() COMMENT '创建日期',
    `update_time` datetime default NOW() COMMENT '更新日期',
    `status` tinyint default '1' COMMENT '状态(可用/禁用/已删除)',
    `profile` text COMMENT '用户属性',
    primary key (`id`),
    UNIQUE KEY `idx_phone` (`phone`),
    key `idx_status` (`status`)
) ENGINE=InnoDB AUTO_INCREMENT=100 DEFAULT CHARSET=utf8mb4;
```
