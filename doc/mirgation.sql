CREATE DATABASE IF NOT EXISTS fileserver DEFAULT CHARSET utf8 COLLATE utf8_general_ci;

use fileserver;

CREATE TABLE `file_metas`(
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `file_name` VARCHAR(256) NOT NULL DEFAULT '' COMMENT '文件名',
  `file_sha1` CHAR(40) NOT NULL DEFAULT '' COMMENT '文件has',
  `file_path` VARCHAR(256) NOT NULL DEFAULT '' COMMENT '文件路径',
  `file_size` BIGINT(30) NOT NULL DEFAULT 0 COMMENT '文件大小',
  `created_at` datetime default NOW() COMMENT '创建日期',
  `updated_at` datetime default NOW() on update current_timestamp() COMMENT '更新时间',
  `status` int(11) NOT NULL DEFAULT '0' COMMENT '状态(可用/禁用/删除)',
  `username` varchar(255) NOT NULL DEFAULT '用户名',
  `ext1` int(11) DEFAULT '0' COMMENT '备用字段1',
  `ext2` text COMMENT '备用字段2',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_fils_hash` (`file_sha1`),
  key `idx_status` (`status`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 COMMENT='文件元数据';

CREATE TABLE `user`(
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `user_name` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '用户名',
    `user_pwd` VARCHAR(100) NOT NULL DEFAULT '' COMMENT '密码',
    `token` CHAR(128) NOT NULL DEFAULT '' COMMENT '登录token',
    `created_at` datetime default NOW() COMMENT '创建时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户信息';


CREATE TABLE `tbl_user_file`(
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `user_name` varchar(64) NOT NULL,
    `file_sha1` varchar(64) NOT NULL DEFAULT '' COMMENT '文件hash',
    `file_name` varchar(256) NOT NULL DEFAULT '' COMMENT '文件名',
    `file_size` bigint  NOT NULL DEFAULT 0 COMMENT '文件大小',
    `uploaded_at` datetime default CURRENT_TIMESTAMP COMMENT '上传时间',
    `updated_at` datetime DEFAULT CURRENT_TIMESTAMP
        ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `status` int(11) NOT NULL DEFAULT '0' COMMENT '文件状态（0删除，1正常，2禁用)',
    PRIMARY KEY (`id`),
    unique key `id_user_file` (`user_name`, `file_sha1`),
    KEY `id_status` (`status`),
    KEY `id_user_id` (`user_name`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8;

