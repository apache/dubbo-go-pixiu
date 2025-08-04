/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

/*pixiu_admin Database table design demo example*/

CREATE DATABASE IF NOT EXISTS `pixiu` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE `pixiu`;

-- ----------------------------
-- Table structure for pixiu_user
-- ----------------------------
DROP TABLE IF EXISTS `pixiu_user`;
CREATE TABLE `pixiu_user` (
                              `id` INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
                              `username` VARCHAR(255) NOT NULL,
                              `password` VARCHAR(255) NOT NULL,
                              `role` INT(5) NOT NULL DEFAULT 1,
                              `enabled` TINYINT(1) NOT NULL DEFAULT 0 COMMENT 'delete or not',
                              `date_created` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
                              `date_updated` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'update time'
);

-- ----------------------------
-- Table structure for pixiu_role
-- ----------------------------
DROP TABLE IF EXISTS `pixiu_role`;
CREATE TABLE `pixiu_role` (
                              `id` INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
                              `role_name` VARCHAR(255) NOT NULL,
                              `description` VARCHAR(255) NOT NULL,
                              `date_created` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
                              `date_updated` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'update time'
);

-- ----------------------------
-- Table structure for pixiu_user_role
-- ----------------------------
DROP TABLE IF EXISTS `pixiu_user_role`;
CREATE TABLE `pixiu_user_role` (
                                   `id` INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
                                   `user_id` INT NOT NULL,
                                   `role_id` INT NOT NULL,
                                   `date_created` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
                                   `date_updated` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'update time'
);
