/*
 Navicat Premium Data Transfer

 Source Server         : Mysql
 Source Server Type    : MySQL
 Source Server Version : 80029
 Source Host           : 81.70.17.190:3306
 Source Schema         : Tiktok

 Target Server Type    : MySQL
 Target Server Version : 80029
 File Encoding         : 65001

 Date: 13/06/2022 10:16:10
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for comment
-- ----------------------------
DROP TABLE IF EXISTS `comment`;
CREATE TABLE `comment`  (
  `id` bigint NOT NULL COMMENT ' 主键',
  `content` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '评论内容',
  `author_id` bigint NOT NULL COMMENT '评论者 ID',
  `video_id` bigint NOT NULL COMMENT '视频 ID',
  `is_delete` tinyint NOT NULL COMMENT '该评论是否被删除',
  `create_time` timestamp NOT NULL COMMENT '评论时间',
  `update_time` timestamp NOT NULL COMMENT '最后修改时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `AUTHOR`(`author_id` ASC) USING BTREE COMMENT '评论作者',
  INDEX `VIDEO`(`video_id` ASC) USING BTREE COMMENT '视频',
  CONSTRAINT `COMMENTAUTHOR` FOREIGN KEY (`author_id`) REFERENCES `user` (`user_id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `VIDEO` FOREIGN KEY (`video_id`) REFERENCES `video` (`video_id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for favorite
-- ----------------------------
DROP TABLE IF EXISTS `favorite`;
CREATE TABLE `favorite`  (
  `id` bigint NOT NULL COMMENT '主键',
  `video_id` bigint NOT NULL COMMENT '视频 ID',
  `favorite_id` bigint NOT NULL COMMENT '点赞者 ID',
  `create_time` timestamp NOT NULL COMMENT '点赞时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `VIDEOID`(`video_id` ASC) USING BTREE COMMENT '视频ID',
  INDEX `FAVORITEID`(`favorite_id` ASC) USING BTREE COMMENT '点赞者',
  CONSTRAINT `FAVORITEID` FOREIGN KEY (`favorite_id`) REFERENCES `user` (`user_id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `VIDEOID` FOREIGN KEY (`video_id`) REFERENCES `video` (`video_id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for follow
-- ----------------------------
DROP TABLE IF EXISTS `follow`;
CREATE TABLE `follow`  (
  `id` bigint NOT NULL COMMENT '主键',
  `user_id` bigint NOT NULL COMMENT '用户 ID',
  `follower_id` bigint NOT NULL COMMENT '粉丝 ID',
  `create_time` timestamp NOT NULL COMMENT '关注时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `USERID`(`user_id` ASC) USING BTREE COMMENT '被关注者',
  INDEX `FOLLOWERID`(`follower_id` ASC) USING BTREE COMMENT '关注者',
  CONSTRAINT `FOLLOWERID` FOREIGN KEY (`follower_id`) REFERENCES `user` (`user_id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `USERID` FOREIGN KEY (`user_id`) REFERENCES `user` (`user_id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user`  (
  `user_id` bigint NOT NULL COMMENT '用户的 ID',
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户名',
  `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户密码',
  `follow_count` bigint NOT NULL COMMENT '关注数',
  `follower_count` bigint NOT NULL COMMENT '粉丝数',
  `create_time` timestamp NOT NULL COMMENT '创建时间',
  `update_time` timestamp NOT NULL COMMENT '修改时间',
  PRIMARY KEY (`user_id`) USING BTREE,
  UNIQUE INDEX `USERID`(`user_id` ASC) USING BTREE COMMENT '用户 id',
  UNIQUE INDEX `NAME`(`name` ASC) USING BTREE COMMENT '用户名唯一'
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for video
-- ----------------------------
DROP TABLE IF EXISTS `video`;
CREATE TABLE `video`  (
  `video_id` bigint NOT NULL AUTO_INCREMENT COMMENT '视频 ID',
  `author_id` bigint NOT NULL COMMENT '作者 ID',
  `play_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '视频地址',
  `cover_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '封面地址',
  `favorite_count` bigint NOT NULL COMMENT '点赞数',
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '视频标题',
  `create_time` timestamp NOT NULL COMMENT '发布时间',
  `comment_count` bigint NOT NULL COMMENT '评论数',
  PRIMARY KEY (`video_id`) USING BTREE,
  INDEX `CREATETIME`(`create_time` ASC) USING BTREE COMMENT '视频流操作时对时间排序进行搜索',
  INDEX `AUTHOR`(`author_id` ASC) USING BTREE COMMENT '查找发布列表时提高效率',
  CONSTRAINT `AUTHOR` FOREIGN KEY (`author_id`) REFERENCES `user` (`user_id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB AUTO_INCREMENT = 58966961576476673 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
