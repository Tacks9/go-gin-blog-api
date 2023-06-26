-- 表 blog_tag 结构
-- [博客-标签表]
CREATE TABLE IF NOT EXISTS `blog_tag` (
	`id` 		      INT(11) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增ID',
	`name` 	      VARCHAR(255) NOT NULL DEFAULT ''  COMMENT '标签名称',
	`created_on`  INT(10) UNSIGNED      DEFAULT '0' COMMENT '创建时间',
	`created_by` 	VARCHAR(255) NOT NULL DEFAULT ''  COMMENT '创建人',
	`modified_on` INT(10) UNSIGNED      DEFAULT '0' COMMENT '修改时间',
	`modified_by` VARCHAR(255) NOT NULL DEFAULT ''  COMMENT '修改人',
	`deleted_on`  INT(10) UNSIGNED      DEFAULT '0' COMMENT '删除时间',
  `state`       TINYINT(1) UNSIGNED   DEFAULT '1' COMMENT '状态 0为禁用、1为启用',
	PRIMARY KEY (`id`) USING BTREE
)
ENGINE=InnoDB
CHARSET='utf8mb4'
COLLATE='utf8mb4_unicode_ci'
COMMENT='博客-标签表'
;


-- 表 blog_article 结构
-- [博客-文章表]
CREATE TABLE IF NOT EXISTS `blog_article` (
	`id` 		      INT(11) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `tag_id`      INT(10) UNSIGNED      DEFAULT '0' COMMENT '标签ID',
	`title` 	    VARCHAR(255) NOT NULL DEFAULT ''  COMMENT '文章标题',
	`desc` 	      VARCHAR(255) NOT NULL DEFAULT ''  COMMENT '简述',
  `content`     text NOT NULL                     COMMENT '正文',
	`created_on`  INT(10) UNSIGNED      DEFAULT '0' COMMENT '创建时间',
	`created_by` 	VARCHAR(255) NOT NULL DEFAULT ''  COMMENT '创建人',
	`modified_on` INT(10) UNSIGNED      DEFAULT '0' COMMENT '修改时间',
	`modified_by` VARCHAR(255) NOT NULL DEFAULT ''  COMMENT '修改人',
	`deleted_on`  INT(10) UNSIGNED      DEFAULT '0' COMMENT '删除时间',
  `state`       TINYINT(1) UNSIGNED   DEFAULT '1' COMMENT '状态 0为禁用、1为启用',
	PRIMARY KEY (`id`) USING BTREE
)
ENGINE=InnoDB
CHARSET='utf8mb4'
COLLATE='utf8mb4_unicode_ci'
COMMENT='博客-文章表'
;

 

-- 表 blog_auth 结构
-- [博客-用户表]
CREATE TABLE IF NOT EXISTS `blog_auth` (
	`id` 		      INT(11) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增ID',
	`username` 	    VARCHAR(50) NOT NULL DEFAULT ''  COMMENT '账号',
	`password` 	    VARCHAR(50) NOT NULL DEFAULT ''  COMMENT '密码',
	PRIMARY KEY (`id`) USING BTREE
)
ENGINE=InnoDB
CHARSET='utf8mb4'
COLLATE='utf8mb4_unicode_ci'
COMMENT='博客-用户表'
;


-- 插入默认数据
INSERT INTO `blog`.`blog_auth` (`id`, `username`, `password`) VALUES (null, 'admin', '123456');

