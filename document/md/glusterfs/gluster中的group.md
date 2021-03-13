### Gluster中的group

| 作者 | 时间 |QQ群 |
| ------ | ------ |------ |
| perrynzhou@gmail.com |2020/12/01 |672152841 |

- gluster中的group概念，gluster中的group是一个或者多个xlator的功能集合，比如设置gluster volume set vol1 db-workload,针对db-workload会有一个或者多个xlator的功能组合对应这个db-workload。当设置以后glusterd会去按照/var/lib/glusterd/groups/db-workload文件中定义的一组xlator对应的参数进行设置完成这个db-workload的操作
- gluster中的group定义如下几种group
```

$ cd /var/lib/glusterd/groups/
$ ls 
db-workload  distributed-virt  gluster-block  metadata-cache  nl-cache  samba  virt

db-workload：适合数据库负载类型的应用
distributed-virt：
gluster-block：运行gluster块存储的设置的group
metadata-cache：用于提供客户端缓存元数据功能
nl-cache：用于提高优化用户文件/目录创建性能
samba：针对windows下使用samba下需要设置的group
virt：针对跑在虚拟机上面的需要设置的group
```
- 每种group定义的xlator的组合的功能如下
```
// db-workload  group定义的xlator功能集合
[root@node /var/lib/glusterd/groups 16:44:55]$ cat db-workload 
performance.open-behind=on
performance.write-behind=off
performance.stat-prefetch=off
performance.quick-read=off
performance.strict-o-direct=on
performance.read-ahead=off
performance.io-cache=off
performance.readdir-ahead=off
performance.client-io-threads=on
server.event-threads=4
client.event-threads=4
performance.read-after-open=yes

// distributed-virt  group定义的xlator功能集合
[root@node /var/lib/glusterd/groups 16:44:58]$ cat distributed-virt 
performance.quick-read=off
performance.read-ahead=off
performance.io-cache=off
performance.low-prio-threads=32
network.remote-dio=enable
features.shard=on
user.cifs=off
client.event-threads=4
server.event-threads=4
performance.client-io-threads=on

// gluster-block   group定义的xlator功能集合
[root@node /var/lib/glusterd/groups 16:45:00]$ cat gluster-block 
performance.quick-read=off
performance.read-ahead=off
performance.io-cache=off
performance.stat-prefetch=off
performance.open-behind=off
performance.readdir-ahead=off
performance.strict-o-direct=on
performance.client-io-threads=on
performance.io-thread-count=32
performance.high-prio-threads=32
performance.normal-prio-threads=32
performance.low-prio-threads=32
performance.least-prio-threads=4
client.event-threads=8
server.event-threads=8
network.remote-dio=disable
cluster.eager-lock=enable
cluster.quorum-type=auto
cluster.data-self-heal-algorithm=full
cluster.locking-scheme=granular
cluster.shd-max-threads=8
cluster.shd-wait-qlength=10000
features.shard=on
features.shard-block-size=64MB
user.cifs=off
server.allow-insecure=on
cluster.choose-local=off

// metadata-cache  group定义的xlator功能集合
[root@node /var/lib/glusterd/groups 16:45:02]$ cat metadata-cache 
features.cache-invalidation=on
features.cache-invalidation-timeout=600
performance.stat-prefetch=on
performance.cache-invalidation=on
performance.md-cache-timeout=600
network.inode-lru-limit=200000
[root@node /var/lib/glusterd/groups 16:45:06]$ cat nl-cache 
features.cache-invalidation=on
features.cache-invalidation-timeout=600
performance.nl-cache=on
performance.nl-cache-timeout=600
network.inode-lru-limit=200000

// samba group定义的xlator功能集合
[root@node /var/lib/glusterd/groups 16:45:08]$ cat samba 
features.cache-invalidation=on
features.cache-invalidation-timeout=600
performance.cache-samba-metadata=on
performance.stat-prefetch=on
performance.cache-invalidation=on
performance.md-cache-timeout=600
network.inode-lru-limit=200000
performance.nl-cache=on
performance.nl-cache-timeout=600
performance.readdir-ahead=on
performance.parallel-readdir=on

// virt group定义的xlator功能集合
[root@node /var/lib/glusterd/groups 16:45:12]$ cat virt 
performance.quick-read=off
performance.read-ahead=off
performance.io-cache=off
performance.low-prio-threads=32
network.remote-dio=disable
performance.strict-o-direct=on
cluster.eager-lock=enable
cluster.quorum-type=auto
cluster.server-quorum-type=server
cluster.data-self-heal-algorithm=full
cluster.locking-scheme=granular
cluster.shd-max-threads=8
cluster.shd-wait-qlength=10000
features.shard=on
user.cifs=off
cluster.choose-local=off
client.event-threads=4
server.event-threads=4
performance.client-io-threads=on
network.ping-timeout=20
server.tcp-user-timeout=20
server.keepalive-time=10
server.keepalive-interval=2
server.keepalive-count=5
```
