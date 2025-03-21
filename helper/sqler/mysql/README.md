# mysql

## docker install mysql

```bash
docker search mysql

docker pull mysql@latest
```

`docker run --name mysql-container -e MYSQL_ROOT_PASSWORD=my-secret-pw -d mysql:latest`

1. --name mysql-container：指定容器名称为 mysql-container（您可以根据需要更改它）。
2. -e MYSQL_ROOT_PASSWORD=my-secret-pw：设置 MySQL root 用户的密码为 my-secret-pw（请根据需要更改密码）。
3. -d：以后台模式运行容器。
4. mysql:latest：使用最新版本的 MySQL 镜像。

`docker run --name mysql -e MYSQL_ROOT_PASSWORD=123456 -p 3306:3306 -d mysql:latest`

## mysql create demo script

```SQL
docker exec -it mysql mysql -u root -p123456
```

```SQL
create database IF NOT EXISTS demo;
use demo;

# create table user
CREATE TABLE IF NOT EXISTS `t_user` (
　　`id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT 'Primary key, user id',

　　`user_name` VARCHAR(64) NOT NULL COMMENT 'User Name',
　　`password` VARCHAR(32) NOT NULL COMMENT 'Encrypted login password; Convention all caps; MD5(MD5(source password) + password_salt)',
　　`password_salt` VARCHAR(32) NOT NULL COMMENT '32 CHAR UUID; Password encrypted salt; SELECT UPPER(REPLACE(UUID(),"-",""))',
　　`email_not_verified` VARCHAR(320) NOT NULL DEFAULT '' COMMENT 'Unverified Email',
　　`user_email` VARCHAR(320) NOT NULL COMMENT 'Verified Email can be used for login; First Set UPPER(REPLACE(UUID(),"-",""))',

　　`update_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'update time',
　　`create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
　　`delete_time` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Soft delete timestamp',

　　UNIQUE KEY `unique_user_name` (`user_name`, `delete_time`),
　　UNIQUE KEY `unique_user_email` (`user_email`, `delete_time`)
) ENGINE = InnoDB DEFAULT CHARSET = UTF8MB4 COMMENT = 'user table';
```

```SQL
INSERT INTO `t_user` (`user_name`, `password`, `password_salt`, `email_not_verified`, `user_email`, `create_time`, `update_time`)
VALUES
('john_doe', 'e99a18c428cb38d5f260853678922e03', 'A1B2C3D4E5F6789A0B1C2D3E4F5A6789', '', 'john@example.com', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('jane_smith', '72b302bf297a228a75730123efef7c41', 'C2D3E4F5A6789B0A1B2C3D4E5F6789A0', '', 'jane@example.com', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('bob_jones', '06a21793bfcad1d9e9f6b4f3497430c3', 'D1E2F3A4B5C6789A0B1C2D3E4F5A6789', '', 'bob@example.com', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('alice_williams', 'c20ad4d76fe97759aa27a0c99bff6710', 'F1A2B3C4D5E6789A0B1C2D3E4F5A678A', '', 'alice@example.com', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('charlie_brown', 'dbab3f3a743e77720448f93950770e6c', 'E2F3A4B5C6789B0A1B2C3D4E5F6789B0', '', 'charlie@example.com', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('eve_chen', 'adbcdca9f5db9f92ee77644b82787016', 'F3A4B5C6D7E8F9A0B1C2D3E4F5A678B', '', 'eve@example.com', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('grace_lee', '6d74a1c03e1d68cbf4e7b089849b8b3c', 'F4A5B6C7D8E9F0A1B2C3D4E5F6A789C', '', 'grace@example.com', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('frank_taylor', 'e2cb82f2a5ccdb9439e95eaf92732b6b', 'G1A2B3C4D5E6789A0B1C2D3E4F5A678D', '', 'frank@example.com', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('hannah_moore', 'd5d7a66d9e6c26fd83720a9233b1188f', 'H2A3B4C5D6E7F8A0B1C2D3E4F5A789E', '', 'hannah@example.com', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('ivan_johnson', 'd3b3a5293f2a56c12b9f22f99c6b601a', 'I1A2B3C4D5E6789A0B1C2D3E4F5A678F', '', 'ivan@example.com', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
```
