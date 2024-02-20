-- MySQL dump 10.13  Distrib 8.0.34, for macos13 (arm64)
--
-- Host: 127.0.0.1    Database: test
-- ------------------------------------------------------
-- Server version	8.0.35

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `transaction`
--

DROP TABLE IF EXISTS `transaction`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `transaction` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '流水號 主鍵',
  `transaction_id` varchar(255) COLLATE utf8mb4_general_ci NOT NULL COMMENT '交易訂單',
  `from_user_id` bigint DEFAULT NULL COMMENT '來源用戶ID',
  `to_user_id` bigint DEFAULT NULL COMMENT '目的用戶ID',
  `product_name` varchar(256) COLLATE utf8mb4_general_ci NOT NULL COMMENT '產品名稱',
  `product_count` bigint DEFAULT NULL COMMENT '產品數量',
  `amount` decimal(20,2) DEFAULT NULL COMMENT '金額',
  `currency` varchar(32) COLLATE utf8mb4_general_ci NOT NULL COMMENT '幣種',
  `created_at` datetime DEFAULT NULL COMMENT '創建時間',
  `uodate_at` datetime DEFAULT NULL COMMENT '更新時間',
  `status` tinyint(1) DEFAULT '0' COMMENT '交易狀態 0:未完成 1:已完成 2:取消 3:錯誤',
  PRIMARY KEY (`id`),
  UNIQUE KEY `transaction_id` (`transaction_id`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `transaction`
--

LOCK TABLES `transaction` WRITE;
/*!40000 ALTER TABLE `transaction` DISABLE KEYS */;
INSERT INTO `transaction` VALUES (1,'1-0-000000000001',1,2,'ETH',1,2000.00,'TWD','2024-02-19 10:36:48','2024-02-19 10:36:58',1),(2,'2-1-000000000002',2,1,'ETH',1,2000.00,'TWD','2024-02-19 10:36:52','2024-02-19 10:36:58',1),(3,'2-1-000000000003',2,1,'ETH',1,2000.00,'TWD','2024-02-19 10:37:11','2024-02-19 10:37:28',1),(4,'1-0-000000000004',1,2,'ETH',1,2000.00,'TWD','2024-02-19 10:37:14','2024-02-19 10:37:28',1),(5,'2-1-000000000005',2,1,'BTC',1,200.00,'TWD','2024-02-19 10:37:28','2024-02-19 10:37:58',1),(6,'1-0-000000000006',1,2,'BTC',1,200.00,'TWD','2024-02-19 10:37:31','2024-02-19 10:37:58',1);
/*!40000 ALTER TABLE `transaction` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2024-02-19 10:38:33
