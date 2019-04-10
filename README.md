
## 初始化数据库

```sql
create database id_alloc_db;

use id_alloc_db;

CREATE TABLE `segments` (
 `biz_tag` varchar(32) NOT NULL,
 `max_id` bigint NOT NULL,
 `step` bigint NOT NULL,
 `description` varchar(1024) DEFAULT '' NOT NULL,
 `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
 PRIMARY KEY (`biz_tag`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO segments(`biz_tag`, `max_id`, `step`, `description`) VALUES('test', 0, 100000, "test业务ID池");
```