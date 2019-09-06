## Redis
Redis是一个开源（BSD许可）的内存型数据结构存储；常用作数据库，缓存和消息代理。  

支持的数据结构：
- String
- Hash
- List
- Set
- Sorted Set(带有范围查询的排序集)

除了以上五种数据结构，Redis还支持 
- bitmap
- hyperloglog
- geospatial（半径查询和流的地理空间索引）

这几个都没用到过，以后有用到在补充

Redis具有内置复制，Lua脚本，LRU驱逐，事务和不同级别的磁盘持久性，并通过Redis Sentinel提供高可用性并使用Redis Cluster自动分区。

以上是Redis的基本介绍；具体可以参考[官网](https://redis.io/topics/introduction)
### Redis 有什么优缺点？
#### 优点
- 速度快： 因为Redis的数据存在内存中，内存的读取速度更快；也支持数据持久化到磁盘中。
    + 持久化： 可以通过每隔一段时间将数据集转储到磁盘或记录每一条命令到日志来持久化。
- 支持丰富的数据类型： String，List,Hash,Set, SotedSet
- 支持事务： 所有的操作都是原子性的，也支持对几个操作合并后的原子执行（MULTI, EXEC, DISCARD and WATCH）。
- 丰富的特性 

### Redis 缓存失效策略
