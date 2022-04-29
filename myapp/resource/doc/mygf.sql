/*
Navicat MySQL Data Transfer

Source Server         : 192.168.212.117
Source Server Version : 50505
Source Host           : 192.168.212.117:3306
Source Database       : myapp

Target Server Type    : MYSQL
Target Server Version : 50505
File Encoding         : 65001

Date: 2022-04-22 14:30:22
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for `desktop`
-- ----------------------------
DROP TABLE IF EXISTS `desktop`;
CREATE TABLE `desktop` (
  `uuid` char(32) NOT NULL COMMENT '桌面uuid',
  `vm_uuid` char(32) NOT NULL COMMENT '虚拟化平台上虚机的uuid',
  `display_name` varchar(255) NOT NULL COMMENT '桌面的显示名称',
  `gpu_attach_status` enum('preAttached','attached','unattached') NOT NULL DEFAULT 'unattached' COMMENT '桌面挂载gpu状态，''preAttached''表示预挂载，即关联了GPU规格的关机态虚机,''attached''表示已经挂载,''unattached''表示未挂载',
  `enabled` enum('enabled','disabled') NOT NULL DEFAULT 'enabled' COMMENT '桌面的启用状态，enabled表示启用，disabled表示桌面已禁用',
  `node_uuid` char(32) NOT NULL DEFAULT '' COMMENT '物理机uuid',
  `node_name` varchar(255) NOT NULL DEFAULT '' COMMENT '物理机名称',
  `is_default` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否是默认桌面，False表示不是，True表示是。默认False。',
  `desc` varchar(255) NOT NULL DEFAULT '' COMMENT '桌面的描述信息',
  `created_at` timestamp NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT NULL COMMENT '最后修改时间',
  PRIMARY KEY (`uuid`),
  KEY `ix_desktop_vm_uuid` (`vm_uuid`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of desktop
-- ----------------------------
INSERT INTO `desktop` VALUES ('111', '222', 'sss', 'unattached', 'enabled', 'sss', 'aaa', '0', '', '2022-02-16 11:12:09', '2022-04-06 17:26:15');
INSERT INTO `desktop` VALUES ('6c2fi10fk80cj3246aye17s2003afi79', 'string', 'string', 'unattached', 'enabled', 'string', 'string', '0', '333', '2022-04-06 17:28:23', '2022-04-06 17:41:45');

-- ----------------------------
-- Table structure for `user`
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
  `uuid` char(32) NOT NULL COMMENT 'uuid',
  `user_name` varchar(64) NOT NULL COMMENT '登录名',
  `display_name` varchar(255) NOT NULL COMMENT '姓名',
  `password` varchar(1024) NOT NULL COMMENT '密码',
  `email` varchar(255) NOT NULL DEFAULT '' COMMENT '邮箱',
  `phone` char(15) NOT NULL DEFAULT '' COMMENT '电话',
  `enabled` enum('enabled','disabled') NOT NULL DEFAULT 'enabled' COMMENT '用户的启用状态，1表示enabled：启用，0表示disabled：禁用',
  `desc` varchar(255) NOT NULL DEFAULT '' COMMENT '描述信息',
  `created_at` timestamp NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT NULL COMMENT '最后修改时间',
  PRIMARY KEY (`uuid`),
  UNIQUE KEY `user_name` (`user_name`),
  KEY `ix_user_display_name` (`display_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of user
-- ----------------------------
INSERT INTO `user` VALUES ('6c2fi104x40cj8vq9a047j0200bht99c', 'admin', 'admin', '77e89dd878d0357d961e9566b70c0291', '', '', 'enabled', 'admin user', '2022-04-13 13:57:28', '2022-04-13 13:57:28');
