/** 数据库 **/
CREATE DATABASE blog CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

/** 博客表 **/
CREATE TABLE `blogs` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `content` text COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_blog_created` (`created_at`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

/** session表 **/
CREATE TABLE sessions (
  token CHAR(43) PRIMARY KEY, 
  data BLOB NOT NULL,
  expiry TIMESTAMP(6) NOT NULL 
);
CREATE INDEX sessions_expiry_idx ON sessions (expiry);

/** 用户表 **/
CREATE TABLE users (
  id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT, 
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL, 
  hashed_password CHAR(60) NOT NULL, 
  created DATETIME NOT NULL
);
ALTER TABLE users ADD CONSTRAINT users_uc_email UNIQUE (email);