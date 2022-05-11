/*
 Navicat MySQL Data Transfer

 Source Server         : LocalMySQL
 Source Server Type    : MySQL
 Source Server Version : 80029
 Source Host           : localhost:3306
 Source Schema         : titoktest

 Target Server Type    : MySQL
 Target Server Version : 80029
 File Encoding         : 65001

 Date: 10/05/2022 22:52:58
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for comment
-- ----------------------------
DROP TABLE IF EXISTS `comment`;
CREATE TABLE `comment`  (
  `comment_id` int NOT NULL COMMENT '评论id，唯一，标识评论',
  `user_id` int NOT NULL COMMENT '用户id，唯一，标识用户',
  `video_id` int NOT NULL COMMENT '视频id，唯一，标识视频',
  `content_place` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '评论内容存放的地址',
  `create_date` datetime NULL DEFAULT NULL COMMENT '评论创建的时间',
  PRIMARY KEY (`comment_id`, `user_id`, `video_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of comment
-- ----------------------------

SET FOREIGN_KEY_CHECKS = 1;
