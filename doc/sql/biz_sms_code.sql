CREATE TABLE `biz_sms_code` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `app_id` bigint NOT NULL,
  `mobile` varchar(20) NOT NULL,
  `code` varchar(10) NOT NULL,
  `expired_at` datetime NOT NULL COMMENT '过期时间',
  `is_used` tinyint NOT NULL DEFAULT '0' COMMENT '0-未使用, 1-已使用',
  `scene` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '场景: login-登录,bind-绑定,reset-重置密码',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_mobile_app` (`mobile`,`app_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='短信验证码';