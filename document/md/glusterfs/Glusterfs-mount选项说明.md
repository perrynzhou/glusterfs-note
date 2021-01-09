###  Gluster挂载选项说明

| 作者                 | 时间       | QQ技术交流群                      |
| -------------------- | ---------- | --------------------------------- |
| perrynzhou@gmail.com | 2020/12/01 | 中国开源存储技术交流群(672152841) |

####  Mount命令

```
  mount -t glusterfs -o dump-fuse=filename  backup-volfile-servers=volfile_server2:volfile_server3,transport-type tcp,log-level=WARNING,reader-thread-count=2,logfile=/var/log/gluster.log server1:/test-volume /mnt/glusterfs
```

#### Mount参数说明
- backup-volfile-servers
  
  - backup-volfile-servers参数提供一组volfile server列表，当第一个volfile server挂了,glusterfs server会指定从backup-volfile-servers列表中执行可用的volfile server给客户端使用，直到客户端挂载成功
  
- log-level

  - log-level 参数说明客户端日志的级别，有效的日志级别分别有TRACE, DEBUG, WARNING, ERROR, CRITICAL INFO and NONE。

- log-file

  - log-file参数指定日志存储的文件

- transport-type

  - transport-type 参数指定fuse客户端和glusterfsd通信的协议，目前支持tcp,rdma已经废弃。

- dump-fuse

  - dum-fuse参数是指定一个文件，用于dump fuse在glusterfs client和linux kernel之间的流量信息

    ```
    # mount -t glusterfs -o dump-fuse=filename hostname:/volname mount-path
    ```

- ro

  - ro 指定mount文件系统以只读模式挂载

    ```
    # mount -t glusterfs -o ro,dump-fuse=filename hostname:/volname mount-path
    ```
  
- acl

  - acl启用posox access contro list功能

- background-qlen = n

  - background-qlen指定fuse请求处理之前最大的请求队列，默认是64

- reader-thread-count=n

  - reader-thread-count指定fuse的读线程数，默认是1；增大这个线程可以获取比较好的读性能

- lru-limit

  - 采用lru方式限制客户端缓存的最大inodes数量，强烈建议这个参数要大于2000或者更大；默认是131072

    ```
    # mount -o lru-limit=131072 -t glusterfs hostname:/volname /mnt/mountdir
    ```

    