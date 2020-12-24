## Glusterfs性能调优


| 作者 | 时间 |QQ技术交流群 |
| ------ | ------ |------ |
| perrynzhou@gmail.com |2020/12/01 |中国开源存储技术交流群(672152841) |

###  volume信息

```
$ gluster volume info
 
Volume Name: dht-vol
Type: Distribute
Volume ID: dd150400-ef24-4b7a-bb11-93092e7c4100
Status: Started
Snapshot Count: 0
Number of Bricks: 3
Transport-type: tcp
Bricks:
Brick1: 172.168.56.40:/dht-vol-pool/brick
Brick2: 172.168.56.41:/dht-vol-pool/brick
Brick3: 172.168.56.42:/dht-vol-pool/brick
Options Reconfigured:
storage.fips-mode-rchecksum: on
transport.address-family: inet
nfs.disable: on
```

### 参数说明

```
//查看某个卷的当前参数的设置
gluster volume get dht-vol  all |grep event -A7

// 查看参数的说明和默认值
gluster volume set help |grep event -A7
```
### 参数调优

```
// 打开metadata-cache,打开这个选项可以提高在mount端操作文件、目录元数据的性能，这个cache的是有一个过期时间，默认是10分钟，如下命令是打开客户端的元数据cache的命令
gluster volume set dht-vol group metadata-cache


// 增加cache的inode的数量，默认是20000,采用lru的淘汰策略进行过期inode
gluster volume set dht-vol network.inode-lru-limit 50000

// cluster.lookup-optimize 选项，在处理查找卷中不存在的条目时会有性能损失。因为DHT会试图在所有子卷中查找文件，所以这种查找代价很高，并且通常会减慢文件的创建速度。 这尤其会影响小文件的性能，其中大量文件被快速连续地添加/创建。 查找卷中不存在的条目的查找扇出行为可以通过在一个均衡过的卷中不进行相同的执行进行优化
gluster volume set dht-vol cluster.lookup-optimize on

// 目录预读的优化
gluster volume get dht-vol performance.readdir-ahead on

// 设置performance.readdir-ahead的内存，默认是10mb，可以适当调大，比如设置为128MB
gluster volume set monitoring_vol performance.rda-cache-limit 60mb

// 目录并行读的优化
gluster volume set dht-vol performance.parallel-readdir on

// 指定客户端网络请求的同时处理的个数，默认是2，这个参数不要超过cpu core的个数
gluster volume set dht-vol client.event-threads  32

// 指定服务端网络请求的同时处理的个数，默认是2，这个参数值不要超过cpu core的个数
gluster volume set dht-vol server.event-threads  32


// glusterfs开启IO缓存的功能
gluster volume set dht-vol  performance.io-cache  on
// 数据读取的cache的内存大小，按照业务特性和机器配置来设定这个值
gluster volume set dht-vol performance.cache-size 16GB 

// 设定缓存文件的最大尺寸，默认是0
gluster volume set dht-vol performance.cache-max-file-size 256MB

// 设置缓存文件的最小尺寸，默认是0
gluster volume set dht-vol performance.cache-min-file-size 1MB


// 在dht上生效的，是指在查找时候，如果在hash所在节点上没有找到相应文件的话，去所有节点上查找一遍。
gluster volume set dht-vol lookup-unhashed off


// 当执行IO操作时候会在客户端把IO入一个内部队列后，返回操作结果给客户端；等内部队列积累的数量达到一定aggregate-size后统一进行通过网络发到后端存储或者经过下一个xlator的处理，这个是异步处理
gluster volume set dht-vol write-behind on

// 设置write-behind开启后，内部队列积累的数据量上线，默认是128KB，这个值视情况而定
gluster volume set dht-vol aggregate-size 8mb 

// 在write-behind开启后，设置flush-behind开启后，用户数据写入到内部队列后直接返回给操作结果给用户
gluster volume set dht-vol flush-behind on


// 这个选项仅仅是针对EC卷生效，并行读取EC卷数据的线程，因为EC数做分片的，所以提供整个参数的值可以提高读取数据的并发度
gluster volume set dht-vol  performance.client-iothreads on

// 设置实际做IO操作线程的数量，建议不超过cpu core的数量
gluster volume set dht-vol  performance.io-thread-count 16


// 默认值是1,设置每次读取数据选择subvolume的策略，1是根据文件的gfid选择子卷，2是根据客户端mount的pid和gfid选择子卷，3是根据最少请求读取子卷，4是选择网络延迟最小策略选择子卷
// 0 = first readable child of AFR, starting from 1st child.
// 1 = hash by GFID of file (all clients use same subvolume).
// 2 = hash by GFID of file and client PID.
// 3 = brick having the least outstanding read requests.
// 4 = brick having the least network ping latency.
gluster volume set dht-vol   cluster.read-hash-mode 1
```
