/*
Navicat MySQL Data Transfer

Source Server         : localhost_3306
Source Server Version : 50719
Source Host           : localhost:3306
Source Database       : student

Target Server Type    : MYSQL
Target Server Version : 50719
File Encoding         : 65001

Date: 2020-04-25 22:09:49
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for lesson
-- ----------------------------
DROP TABLE IF EXISTS `lesson`;
CREATE TABLE `lesson` (
  `Stunum` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL,
  `lesson` varchar(255) NOT NULL,
  `class` varchar(255) NOT NULL,
  `time` varchar(255) NOT NULL,
  `place` varchar(255) NOT NULL,
  `Selectlesson` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`Stunum`,`lesson`,`class`),
  KEY `idx_Stunum_lesson_class` (`Stunum`,`lesson`,`class`) USING BTREE,
  KEY `lesson` (`lesson`),
  KEY `Stunum` (`Stunum`),
  CONSTRAINT `lesson` FOREIGN KEY (`lesson`) REFERENCES `selectlesson` (`lesson`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
