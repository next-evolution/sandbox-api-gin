CREATE TABLE `fx_economic_indicator_data` (
  `id` int(10) unsigned NOT NULL COMMENT '経済指標ID',
  `publication` datetime NOT NULL COMMENT '公表日時',
  `sub_title` varchar(16) DEFAULT NULL COMMENT 'サブタイトル',
  `result_value` varchar(32) NOT NULL COMMENT '結果値',
  `forecast_value` varchar(32) DEFAULT NULL COMMENT '予想値',
  `previous_value` varchar(32) DEFAULT NULL COMMENT '前回値',
  `memo` varchar(255) DEFAULT NULL COMMENT 'メモ',
  PRIMARY KEY (`id`,`publication`),
  CONSTRAINT `economic_indicator_data_fk` FOREIGN KEY (`id`) REFERENCES `fx_economic_indicator` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='経済指標データ情報';
