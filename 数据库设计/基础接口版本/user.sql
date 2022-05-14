/*
 Navicat MySQL Data Transfer

 Source Server         : Tiktok
 Source Server Type    : MySQL
 Source Server Version : 80027
 Source Host           : 81.70.17.190:3306
 Source Schema         : Tiktok

 Target Server Type    : MySQL
 Target Server Version : 80027
 File Encoding         : 65001

 Date: 14/05/2022 12:46:12
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user`  (
  `user_id` bigint NOT NULL AUTO_INCREMENT COMMENT '用户id，唯一标识',
  `name` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '用户名，唯一，登陆时使用，也是用户昵称',
  `password` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '登录密码',
  PRIMARY KEY (`user_id`) USING BTREE,
  UNIQUE INDEX `index_name`(`name`) USING BTREE COMMENT '加速查询用户名，用于登录与注册查询用户名是否存在'
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
