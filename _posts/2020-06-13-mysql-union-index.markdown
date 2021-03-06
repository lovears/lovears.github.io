---
layout:     post 
title:      "【MySQL】MySQL联合索引"
subtitle:   "MySQL联合索引使用规则"
date:       2020-06-13
author:     "老回路"
header-img: "img/post-bg-unix-linux.jpg"
tags:
    - MySQL
    - 数据库
    - 索引
---
# MySql联合索引
> 对数据库方面的知识有点匮乏或者说缺乏系统性的学习，平时的工作中虽然有用到了数据库，也仅限于基本的数据库的增删改查，具体的实现机制或者原理了解的比较浅显或者甚至比较混乱，所以最近也在看一些数据库相关的文章或者书籍。不过在看了多篇文章后，发现这些文章居然有些地方是相互矛盾的，这就让人很疑惑。和朋友讨论，也找了一些有参考价值的文章或者评论，但是最近还是觉得自己亲自实践一下比较好。

### 1.前期准备：
- 创建一张表，并且在(c1,c2,c3)上创建一个联合索引`IDX_U1`

```sql
CREATE TABLE `i_test` (
  `ID` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `c1` int(11) DEFAULT NULL COMMENT 'c1',
  `c2` int(11) DEFAULT NULL COMMENT 'c2',
  `c3` int(11) DEFAULT NULL COMMENT 'c3',
  `c4` int(11) DEFAULT NULL COMMENT 'c4',
  `c5` int(11) DEFAULT NULL COMMENT 'c5',
  `c6` int(11) DEFAULT NULL COMMENT 'c6',
  PRIMARY KEY (`ID`),
  KEY `IDX_U1` (`c1`,`c2`,`c3`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='i_test'
```  

### 2. 测试用例

#### 2.1 可以分以下几种情况：
- 查询所有列，索引按顺序使用，在索引列顺序的前中后插入不在索引中的列
- 查询所有列，索引乱序使用，在索引条件的前中后插入不在索引中的列
- 查询所有列，顺序使用部份索引，主要有(c1c2,c1c3,c2c3)三种情况
- 查询所有列，乱序使用部份索引，主要有(c2c1,c3c1,c3c2)三种情况
- 查询所有列，单独使用索引中的一列
- 查询索引中的列，索引按顺序使用，在索引列顺序的前中后插入不在索引中的列
- 查询索引中的列，索引乱序使用，在索引条件的前中后插入不在索引中的列
- 查询索引中的列，顺序使用部份索引，主要有(c1c2,c1c3,c2c3)三种情况，在索引列顺序的前中后插入不在索引中的列
- 查询索引中的列，乱序使用部份索引，主要有(c2c1,c3c1,c3c2)三种情况，在索引列顺序的前中后插入不在索引中的列
- 查询索引中的列，单独使用索引中的一列，在索引列顺序的前中后插入不在索引中的列
- 查询索引中的列，使用非索引列查询


### 3.先说结论
- **联合索引的使用顺序和where条件中的顺序并没有很强的关系，因为sql执行器会优化sql来匹配索引顺序；在使用中，为了减少执行器优化开销，应该尽量按照索引顺序使用**
- **联合索引是按照顺序最左匹配的；即：如果跳过某个索引列，其后的索引列都不会使用索引。**
- **如果查询列为索引列的话，会使用到覆盖索引**，比如： `select c1,c2,c3 from i_test where c2='' and c3='' ` 就会使用到覆盖索引，而且会减少回表查询的操作

### 4 相关概念
#### 4.1 覆盖索引
> 如果一个索引包含（或者说覆盖了）所有满足查询所需要的数据，那么就称这类索引为覆盖索引

#### 4.2 Extra字段解释
>using index ：使用覆盖索引的时候就会出现
>using where：在查找使用索引的情况下，需要回表去查询所需的数据
>using index condition：查找使用了索引，但是需要回表查询数据
>using index & using where：查找使用了索引，但是需要的数据都在索引列中能找到，所以不需要回表查询数据


### 5.实践验证
#### 5.1 查询所有列，索引按顺序使用，在索引列顺序的前中后插入不在索引中的列
```bash
mysql> explain select * from i_test where c1='' and c2='' and c3='';
+----+-------------+--------+------------+------+---------------+--------+---------+-------------------+------+----------+-------+
| id | select_type | table  | partitions | type | possible_keys | key    | key_len | ref               | rows | filtered | Extra |
+----+-------------+--------+------------+------+---------------+--------+---------+-------------------+------+----------+-------+
|  1 | SIMPLE      | i_test | NULL       | ref  | IDX_U1        | IDX_U1 | 15      | const,const,const |    1 |   100.00 | NULL  |
+----+-------------+--------+------------+------+---------------+--------+---------+-------------------+------+----------+-------+
1 row in set, 1 warning (1.89 sec)

mysql> explain select * from i_test where c1='' and c2='' and c3='' and c4='';
+----+-------------+--------+------------+------+---------------+--------+---------+-------------------+------+----------+-------------+
| id | select_type | table  | partitions | type | possible_keys | key    | key_len | ref               | rows | filtered | Extra       |
+----+-------------+--------+------------+------+---------------+--------+---------+-------------------+------+----------+-------------+
|  1 | SIMPLE      | i_test | NULL       | ref  | IDX_U1        | IDX_U1 | 15      | const,const,const |    1 |   100.00 | Using where |
+----+-------------+--------+------------+------+---------------+--------+---------+-------------------+------+----------+-------------+
1 row in set, 1 warning (0.00 sec)

mysql> explain select * from i_test where c4='' and c1='' and c2='' and c3='';
+----+-------------+--------+------------+------+---------------+--------+---------+-------------------+------+----------+-------------+
| id | select_type | table  | partitions | type | possible_keys | key    | key_len | ref               | rows | filtered | Extra       |
+----+-------------+--------+------------+------+---------------+--------+---------+-------------------+------+----------+-------------+
|  1 | SIMPLE      | i_test | NULL       | ref  | IDX_U1        | IDX_U1 | 15      | const,const,const |    1 |   100.00 | Using where |
+----+-------------+--------+------------+------+---------------+--------+---------+-------------------+------+----------+-------------+
1 row in set, 1 warning (0.00 sec)

mysql> explain select * from i_test where c1='' and c4='' and c2='' and c3='';
+----+-------------+--------+------------+------+---------------+--------+---------+-------------------+------+----------+-------------+
| id | select_type | table  | partitions | type | possible_keys | key    | key_len | ref               | rows | filtered | Extra       |
+----+-------------+--------+------------+------+---------------+--------+---------+-------------------+------+----------+-------------+
|  1 | SIMPLE      | i_test | NULL       | ref  | IDX_U1        | IDX_U1 | 15      | const,const,const |    1 |   100.00 | Using where |
+----+-------------+--------+------------+------+---------------+--------+---------+-------------------+------+----------+-------------+
1 row in set, 1 warning (0.00 sec)
```
##### 5.1.1结论
可以看到，以上几种情况都使用到了索引；可以说，联合索引按顺序使用，不论顺序中是否使用了其他条件，都会使用到索引。

#### 5.2 查询所有列，索引乱序使用，在索引条件的前中后插入不在索引中的列
```bash
mysql> explain select * from i_test where c3='' and c2='' and c1='';
+----+-------------+--------+------------+------+---------------+--------+---------+-------------------+------+----------+-------+
| id | select_type | table  | partitions | type | possible_keys | key    | key_len | ref               | rows | filtered | Extra |
+----+-------------+--------+------------+------+---------------+--------+---------+-------------------+------+----------+-------+
|  1 | SIMPLE      | i_test | NULL       | ref  | IDX_U1        | IDX_U1 | 15      | const,const,const |    1 |   100.00 | NULL  |
+----+-------------+--------+------------+------+---------------+--------+---------+-------------------+------+----------+-------+
1 row in set, 1 warning (0.00 sec)

mysql> explain select * from i_test where c4='' and c3='' and c2='' and c1='';
+----+-------------+--------+------------+------+---------------+--------+---------+-------------------+------+----------+-------------+
| id | select_type | table  | partitions | type | possible_keys | key    | key_len | ref               | rows | filtered | Extra       |
+----+-------------+--------+------------+------+---------------+--------+---------+-------------------+------+----------+-------------+
|  1 | SIMPLE      | i_test | NULL       | ref  | IDX_U1        | IDX_U1 | 15      | const,const,const |    1 |   100.00 | Using where |
+----+-------------+--------+------------+------+---------------+--------+---------+-------------------+------+----------+-------------+
1 row in set, 1 warning (0.00 sec)

mysql> explain select * from i_test where c3='' and c4='' and c2='' and c1='';
+----+-------------+--------+------------+------+---------------+--------+---------+-------------------+------+----------+-------------+
| id | select_type | table  | partitions | type | possible_keys | key    | key_len | ref               | rows | filtered | Extra       |
+----+-------------+--------+------------+------+---------------+--------+---------+-------------------+------+----------+-------------+
|  1 | SIMPLE      | i_test | NULL       | ref  | IDX_U1        | IDX_U1 | 15      | const,const,const |    1 |   100.00 | Using where |
+----+-------------+--------+------------+------+---------------+--------+---------+-------------------+------+----------+-------------+
1 row in set, 1 warning (0.00 sec)

mysql> explain select * from i_test where c3='' and c2='' and c1='' and c4='';
+----+-------------+--------+------------+------+---------------+--------+---------+-------------------+------+----------+-------------+
| id | select_type | table  | partitions | type | possible_keys | key    | key_len | ref               | rows | filtered | Extra       |
+----+-------------+--------+------------+------+---------------+--------+---------+-------------------+------+----------+-------------+
|  1 | SIMPLE      | i_test | NULL       | ref  | IDX_U1        | IDX_U1 | 15      | const,const,const |    1 |   100.00 | Using where |
+----+-------------+--------+------------+------+---------------+--------+---------+-------------------+------+----------+-------------+
1 row in set, 1 warning (0.00 sec)
```
##### 5.2.1 结论
乱序使用所有索引列，也会使用到索引。所以可以得出结论：**使用索引和`where`条件中的索引列顺序无关，只是和使用到的索引列的顺序有关。**应该是MySQL的执行器优化了语句的结果；因为执行器对SQL的语句优化也是有开销的，所以实际使用中`where`中的语句最好还是按照索引的顺序写会比较好。

#### 5.3 查询所有列，顺序使用部份索引，主要有(c1c2,c1c3,c2c3)三种情况
```bash

mysql> explain select * from i_test where c1='' and c2='';
+----+-------------+--------+------------+------+---------------+--------+---------+-------------+------+----------+-------+
| id | select_type | table  | partitions | type | possible_keys | key    | key_len | ref         | rows | filtered | Extra |
+----+-------------+--------+------------+------+---------------+--------+---------+-------------+------+----------+-------+
|  1 | SIMPLE      | i_test | NULL       | ref  | IDX_U1        | IDX_U1 | 10      | const,const |    1 |   100.00 | NULL  |
+----+-------------+--------+------------+------+---------------+--------+---------+-------------+------+----------+-------+
1 row in set, 1 warning (0.00 sec)

mysql> explain select * from i_test where c1='' and c3='';
+----+-------------+--------+------------+------+---------------+--------+---------+-------+------+----------+-----------------------+
| id | select_type | table  | partitions | type | possible_keys | key    | key_len | ref   | rows | filtered | Extra                 |
+----+-------------+--------+------------+------+---------------+--------+---------+-------+------+----------+-----------------------+
|  1 | SIMPLE      | i_test | NULL       | ref  | IDX_U1        | IDX_U1 | 5       | const |    1 |   100.00 | Using index condition |
+----+-------------+--------+------------+------+---------------+--------+---------+-------+------+----------+-----------------------+
1 row in set, 1 warning (0.00 sec)

mysql> explain select * from i_test where c2='' and c3='';
+----+-------------+--------+------------+------+---------------+------+---------+------+------+----------+-------------+
| id | select_type | table  | partitions | type | possible_keys | key  | key_len | ref  | rows | filtered | Extra       |
+----+-------------+--------+------------+------+---------------+------+---------+------+------+----------+-------------+
|  1 | SIMPLE      | i_test | NULL       | ALL  | NULL          | NULL | NULL    | NULL |    1 |   100.00 | Using where |
+----+-------------+--------+------------+------+---------------+------+---------+------+------+----------+-------------+
1 row in set, 1 warning (0.00 sec)
```
##### 5.3.1 结论
在使用索引中的部份索引列的时候，必须保证索引的顺序使用即`c1,c2` 会用到 `c1,c2`两个列上的索引，`c1,c3`因为跳过了`c2`所以不会使用到`c3`的索引；`c2,c3`跳过了`c1`,不会使用索引。所以可以得出一个结论：**索引是按照顺序最左匹配的；即：如果跳过某个索引列，其后的索引列都不会使用索引。**

#### 5.4 查询所有列，乱序使用部份索引，主要有(c2c1,c3c1,c3c2)三种情况
基于 `5.2.1` 的结论可以推断得知，此处的查询结果应该和`5.3`中的结果一致，执行过程就不贴出来了。

##### 5.4.1 结论
参见 `5.3.1 结论`

#### 5.5 查询所有列，单独使用索引中的一列
基于 `5.3.1` 的结论可以推断得知，只有在使用`c1`查询的时候会使用到索引，其他时候都不会使用到。执行过程就不贴出来了。

#### 5.6 以下几种情况一起讨论
##### 5.6.1 查询索引中的列，索引按顺序使用，在索引列顺序的前中后插入不在索引中的列
此种情况同 `5.1`,会使用到覆盖索引
```bash
mysql> explain select c1,c2,c3 from i_test  where c1='' and  c2='' and c3='';
+----+-------------+--------+------------+------+---------------+--------+---------+-------------------+------+----------+-------------+
| id | select_type | table  | partitions | type | possible_keys | key    | key_len | ref               | rows | filtered | Extra       |
+----+-------------+--------+------------+------+---------------+--------+---------+-------------------+------+----------+-------------+
|  1 | SIMPLE      | i_test | NULL       | ref  | IDX_U1        | IDX_U1 | 15      | const,const,const |    1 |   100.00 | Using index |
+----+-------------+--------+------------+------+---------------+--------+---------+-------------------+------+----------+-------------+
1 row in set, 1 warning (0.00 sec)

mysql> explain select c1,c2,c3 from i_test  where c3='' and  c2='' and c1='';
+----+-------------+--------+------------+------+---------------+--------+---------+-------------------+------+----------+-------------+
| id | select_type | table  | partitions | type | possible_keys | key    | key_len | ref               | rows | filtered | Extra       |
+----+-------------+--------+------------+------+---------------+--------+---------+-------------------+------+----------+-------------+
|  1 | SIMPLE      | i_test | NULL       | ref  | IDX_U1        | IDX_U1 | 15      | const,const,const |    1 |   100.00 | Using index |
+----+-------------+--------+------------+------+---------------+--------+---------+-------------------+------+----------+-------------+
1 row in set, 1 warning (0.00 sec)
```
###### 5.6.1.1结论
可以看到，`Extra`是`using index`表示使用到了覆盖索引，

##### 5.6.2 查询索引中的列，索引乱序使用，在索引条件的前中后插入不在索引中的列
基于 5.2.1 的结论可以推断得知，此种情况同 `5.6.1`

#### 5.7 查询索引中的列，顺序使用部份索引，主要有(c1c2,c1c3,c2c3)三种情况，在索引列顺序的前中后插入不在索引中的列
```bash
mysql> explain select c1,c2,c3 from i_test where c1='' and c2='';
+----+-------------+--------+------------+------+---------------+--------+---------+-------------+------+----------+-------------+
| id | select_type | table  | partitions | type | possible_keys | key    | key_len | ref         | rows | filtered | Extra       |
+----+-------------+--------+------------+------+---------------+--------+---------+-------------+------+----------+-------------+
|  1 | SIMPLE      | i_test | NULL       | ref  | IDX_U1        | IDX_U1 | 10      | const,const |    1 |   100.00 | Using index |
+----+-------------+--------+------------+------+---------------+--------+---------+-------------+------+----------+-------------+
1 row in set, 1 warning (0.00 sec)

mysql> explain select c1,c2,c3 from i_test where c1='' and c3='';
+----+-------------+--------+------------+------+---------------+--------+---------+-------+------+----------+--------------------------+
| id | select_type | table  | partitions | type | possible_keys | key    | key_len | ref   | rows | filtered | Extra                    |
+----+-------------+--------+------------+------+---------------+--------+---------+-------+------+----------+--------------------------+
|  1 | SIMPLE      | i_test | NULL       | ref  | IDX_U1        | IDX_U1 | 5       | const |    1 |   100.00 | Using where; Using index |
+----+-------------+--------+------------+------+---------------+--------+---------+-------+------+----------+--------------------------+
1 row in set, 1 warning (0.00 sec)

mysql> explain select c1,c2,c3 from i_test where c2='' and c3='';
+----+-------------+--------+------------+-------+---------------+--------+---------+------+------+----------+--------------------------+
| id | select_type | table  | partitions | type  | possible_keys | key    | key_len | ref  | rows | filtered | Extra                    |
+----+-------------+--------+------------+-------+---------------+--------+---------+------+------+----------+--------------------------+
|  1 | SIMPLE      | i_test | NULL       | index | NULL          | IDX_U1 | 15      | NULL |    1 |   100.00 | Using where; Using index |
+----+-------------+--------+------------+-------+---------------+--------+---------+------+------+----------+--------------------------+
1 row in set, 1 warning (0.00 sec)
```
##### 5.7.1 结论
可以看到，以上几种情况都使用了索引，和`5.3`的结论不一样；这里我们就需要注意`Extra`字段了，可以看到，在`c2c3`的组合中，使用到了`using where using index`表示使用到了覆盖索引。

#### 5.8 查询索引中的列，乱序使用部份索引，主要有(c2c1,c3c1,c3c2)三种情况，在索引列顺序的前中后插入不在索引中的列
基于 5.2.1 的结论可以推断得知，此种情况同 `5.7`,会使用到覆盖索引

#### 5.9 查询索引中的列，单独使用索引中的一列，在索引列顺序的前中后插入不在索引中的列
此种情况同 `5.7`,会使用到覆盖索引

#### 5.10查询索引中的列，使用非索引列查询
```bash
mysql> explain select c1,c2,c3 from i_test where c4='';
+----+-------------+--------+------------+------+---------------+------+---------+------+------+----------+-------------+
| id | select_type | table  | partitions | type | possible_keys | key  | key_len | ref  | rows | filtered | Extra       |
+----+-------------+--------+------------+------+---------------+------+---------+------+------+----------+-------------+
|  1 | SIMPLE      | i_test | NULL       | ALL  | NULL          | NULL | NULL    | NULL |    1 |   100.00 | Using where |
+----+-------------+--------+------------+------+---------------+------+---------+------+------+----------+-------------+
1 row in set, 1 warning (0.00 sec)
```
##### 5.10.1
这种情况，虽然查询的是索引列中的数据，但是没有使用到索引；`Extra`字段为`using where`表示使用到索引的情况下，会去回表查询。

