## 分布式id调研


源码
[go-id-alloc](https://github.com/owenliang/go-id-alloc)

![](http://ww1.sinaimg.cn/large/006dizvAly1g1xng0yh1uj30gc0duwen.jpg)

以外卖订单请求为例；
假设分布式id服务是部署在三台机器(A,B,C)上;他们访问的是同一个数据库;

tag | max_id |step
----|-------|-----
waimai| 0 | 1000


  1. id服务A第一次接受请求，去读取数据库；
   此时数据库tag == "waimai"的记录
    的max_id=0;  更新max_id= max_id +step(1000) 更新后max_id为1000；
   服务A将(1~1000)号码段存入内存中; 
   并按顺序读取返回给外卖订单请求,即第一次读取id=1返回；并更新内存中的号码段为 2-1000;

  2. id服务A第二次接受请求，此时内存中有 2-1000这个号码段；读取id=2返回给外卖订单请求；并更新内存中的号码段为 3-1000;

  3. id服务B第一次接受请求，此时内存中没有号码段，则去请求数据库，
    此时数据库max_id = 1000; 更新max_id= max_id +step(1000);更新后max_id为2000;
    服务B将(1001-2000)号码段存入内存中；
    并按顺序读取返回给外卖订单请求,即第一次读取id=1001返回；并更新内存中的号码段为 1002-2000;

  3. id服务C第一次请求同上，读取 2001-3000存入内存中

  4. id服务A号码段消耗完的时候；就会再去读取数据库，以此类推
  








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


