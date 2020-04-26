/*
Navicat MySQL Data Transfer

Source Server         : localhost_3306
Source Server Version : 50719
Source Host           : localhost:3306
Source Database       : student

Target Server Type    : MYSQL
Target Server Version : 50719
File Encoding         : 65001

Date: 2020-04-25 22:10:04
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for selectlesson
-- ----------------------------
DROP TABLE IF EXISTS `selectlesson`;
CREATE TABLE `selectlesson` (
  `lesson` varchar(255) NOT NULL,
  `typee` varchar(255) NOT NULL,
  `teacher` varchar(255) NOT NULL,
  KEY `idx_lesson` (`lesson`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
