/*
Navicat MySQL Data Transfer

Source Server         : localhost_3306
Source Server Version : 50719
Source Host           : localhost:3306
Source Database       : student

Target Server Type    : MYSQL
Target Server Version : 50719
File Encoding         : 65001

Date: 2020-04-25 22:10:11
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for student
-- ----------------------------
DROP TABLE IF EXISTS `student`;
CREATE TABLE `student` (
  `name` varchar(255) NOT NULL,
  `stunum` varchar(255) NOT NULL,
  `Lesson` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`stunum`),
  KEY `idx_stunum` (`stunum`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
