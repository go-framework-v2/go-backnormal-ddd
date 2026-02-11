CREATE TABLE `biz_app` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `bundle_id` varchar(255) NOT NULL COMMENT '包名，唯一区分应用，新品和马甲不是一个应用',
  `app_name` varchar(255) NOT NULL COMMENT '应用名',
  `project_id` int NOT NULL COMMENT '项目ID，新品和马甲共用',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_biz_app_bundle` (`bundle_id`),
  KEY `idx_biz_ap_project_id` (`project_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='app应用';


INSERT INTO `biz_app` (`bundle_id`, `app_name`, `project_id`)
VALUES
	('com.happyhorse.alert', '每日打卡小闹钟', '53'),
	('com.happyhorse.billjc', '安居水电燃气管家', '59'),
	('com.happyhorse.bluecad', 'CAD蓝图-专业看图', '5'),
	('com.happyhorse.music', '全民畅听-音乐播放器', '1'),
	('com.happyhorse.planet', '环球影视大全解说', '10'),
	('com.happyhorse.budgetpro', '小财神记账', '30'),
	('com.happyhorse.onenet', 'WiFi一键连网', '13'),
	('com.happyhorse.snow', '雪峰高考指南', '20'),
    ('com.happyhorse.nfc.hh', 'NFC智能读卡器', '15'),
    ('com.happyhorse.remotee', 'E家遥控器', '6'),
    ('com.happyhorse.redapple', '红菓短剧', '67'),
    ('com.happyhorse.soda', 'Soda音乐播放器', '1'),
    ('com.happyhorse.stepcash', '每日走路赚多多', '47'),
    ('com.happyhorse.zhi', '小智超级计算助手', '14'),
    ('com.happyhorse.train', '火车票实时速查', '50'),
    ('com.nato.tax', '个人所得查税退税', '62'),
    ('com.happyhorse.cool', '酷驹音乐', '82'),
    ('com.nato.taxx', '个人所得税查', '69'),
    ('com.nato.zhsb', '智慧社保云', '66'),
    ('com.happyhorse.nineone', '9I浏览器-极速版', '4'),
    ('com.happyhorse.camss', '专业数数相机', '8'),
    ('com.happyhorse.bus', '免费公交即刻查询', '16');