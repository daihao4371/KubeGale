-- MySQL dump 10.13  Distrib 5.7.31, for macos10.14 (x86_64)
--
-- Host: localhost    Database: cmdb
-- ------------------------------------------------------
-- Server version	5.7.31

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `sys_users`
--

DROP TABLE IF EXISTS `sys_users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `sys_users` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `uuid` varchar(191) DEFAULT NULL COMMENT '用户UUID',
  `username` varchar(191) DEFAULT NULL COMMENT '用户登录名',
  `password` varchar(191) DEFAULT NULL COMMENT '用户登录密码',
  `nick_name` varchar(191) DEFAULT '系统用户' COMMENT '用户昵称',
  `header_img` varchar(191) DEFAULT 'https://pic.cnblogs.com/avatar/2399534/20220419203643.png' COMMENT '用户头像',
  `authority_id` bigint(20) unsigned DEFAULT '888' COMMENT '用户角色ID',
  `phone` varchar(191) DEFAULT NULL COMMENT '用户手机号',
  `email` varchar(191) DEFAULT NULL COMMENT '用户邮箱',
  `enable` bigint(20) DEFAULT '1' COMMENT '用户是否被冻结 1正常 2冻结',
  `origin_setting` text COMMENT '配置',
  PRIMARY KEY (`id`),
  KEY `idx_sys_users_deleted_at` (`deleted_at`),
  KEY `idx_sys_users_uuid` (`uuid`),
  KEY `idx_sys_users_username` (`username`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_users`
--

LOCK TABLES `sys_users` WRITE;
/*!40000 ALTER TABLE `sys_users` DISABLE KEYS */;
INSERT INTO `sys_users` VALUES (1,'2025-03-20 22:12:42.429','2025-04-11 15:43:14.206',NULL,'39378593-7867-4b66-8601-b55b09769809','admin','$2a$10$uwFq1KAOHqzHVDTp6DBCsedsV.mx/pPovpKIOKt56DESfGjljcto6','花海','',889,'17899998888','9999@qq.com',1,NULL),
                               (2,'2025-04-09 11:06:29.619','2025-04-09 21:21:47.448','','19f0ac64-d07c-4123-b83a-d0451eba7c48','zhangsan','$2a$10$99231svZ78GrLaQIZhOm3ucDda7HnI.Ou.OpplBJttuuNRAOgHlu6','张三','',888,'19018182828','sbka@qq.com',1,NULL),
                               (3,'2025-04-09 11:12:47.771','2025-04-09 11:12:47.771','','3b167e14-165e-4c61-b3a8-a5c47fe400a8','lisi1','$2a$10$b317lWtkxoqjYKKbCENT5uD/MDmfGmDWu82.WlFbjR6voR7GraB0q','李四','https://pic.cnblogs.com/avatar/2399534/20220419203643.png',9528,'17811112222','qwe@qq.com',1,NULL),
                               (4,'2025-04-09 13:41:31.648','2025-04-09 15:56:05.557','','10065896-c3dc-4396-ac50-72487b649ce2','wangwu','$2a$10$6px3hWOpBYVIJZUnHXg3TujBM3NU2dIvyOsJPCs6BhYxuWQ95cMOW','王五','',9528,'13388991122','sbka@qq.com',2,NULL),
                               (5,'2025-04-09 15:49:15.842','2025-04-09 16:56:43.685','','bd5ef14b-37b1-4f07-810e-389e79d16663','zhaoliu','$2a$10$rxLDVKdochD74UMphX9w.eWBcVvTDoJZtb1ZbayWAkqHQ.rscSRzG','赵六','',9528,'13411887722','123adbj@qq.com',2,NULL),
                               (6,'2025-04-09 20:16:42.826','2025-04-12 22:39:33.564',NULL,'9b6b28aa-b142-4566-9c61-535e150177cd','devops','$2a$10$LA73LunLwW8XOaZQ9foc2OYqrwttHASlcrIn1SemY.yAu2mNyZnDG','运维开发','',892,'13388991122','12315dbj@qq.com',1,NULL),
                               (8,'2025-04-10 10:54:58.090','2025-04-10 12:11:36.159',NULL,'bb0972ca-34b6-42c3-af50-f04475b93c62','test2','$2a$10$C7TRYdtij9btY3.aD4/JCeqxPFIkkcGseRgGPsrM.6D4OQ.yqagu.','王五','',888,'13766667777','wangwu@qq.com',1,NULL);
/*!40000 ALTER TABLE `sys_users` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2025-04-16  9:32:28
