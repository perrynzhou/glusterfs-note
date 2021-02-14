#  gluster/fuse/nfs-ganesha 笔记

| 作者 | 时间 |QQ技术交流群 |
| ------ | ------ |------ |
| perrynzhou@gmail.com |2020/12/01 |672152841 |


### gluster blogs

- [why-brick-multiplexing](https://gluster.home.blog/2019/05/06/why-brick-multiplexing/)

### glusterfs/pNFS/NFS/fuse/HPC
- [glustefs](./document/pdf/glusterfs)
- [pNFS](./document/pdf/pNFS)
- [NFS](./document/pdf/NFS)
- [fuse](./document/pdf/fuse)
- [hpc](./document/pdf/hpc)


### fuse 分析
- [libfuse3-10源码编译](./document/md/fuse/libfuse-3.10源码编译.md)

### nfs-ganesha分析
- [nfs-ganesha容器部署](./document/md/nfs-ganesha/nfs-ganesha容器部署.md)
- [nfs-ganesha源码安装](./document/md/nfs-ganesha/nfs-ganesha源码安装.md)
- [基于glusterfs的nfs-ganesha方案](./document/md/nfs-ganesha/基于glusterfs的nfs-ganesha方案.md)


### glusterfs commit pr

- [fixed read-hash-mode选择读取模式问题](https://review.gluster.org/#/c/glusterfs/+/25062/)
- [fixed posix_disk_space_check_thread_proc函数检测磁盘剩余空间间隔时间可配置方式](https://github.com/perrynzhou/glusterfs/commit/3256de978f29801b7d29af56c4cc7587ec421cc9)


### glusterfs 分析
- [glusterfs源码安装](./document/md/glusterfs/glusterfs源码安装.md)
- [glusterfsd启动过程](./document/md/glusterfs/glusterfsd启动过程.md)
- [glusterfs客户端挂载init流程](./document/md/glusterfs/glusterfs客户端挂载init流程.md)
- [gluster-create-volume处理过程](./document/md/glusterfs/gluster-create-volume处理过程.md)
- [glusterfs架构和基本概念](./document/md/glusterfs/glusterfs架构和基本概念.md)
- [glusterfs-brick哈希范围设定过程](./document/md/glusterfs/glusterfs-brick哈希范围设定过程.md)
- [glusterfs客户端写数据分析](./document/md/glusterfs/glusterfs客户端写数据分析.md)
- [glusterfs问题诊断](./document/md/glusterfs/glusterfs问题诊断.md)
- [glusterfs-fuse实现(持续更新)](./document/md/glusterfs/glusterfs-fuse实现.md)
- [cluster.read-hash-mode工作原理](./document/md/glusterfs/cluster.read-hash-mode工作原理.md)
- [多副本情况下mount挂载目录如何选择可用的副本目录](./document/md/glusterfs/多副本情况下mount挂载目录如何选择可用的副本目录.md)
- [gfapi如何工作的](./document/md/glusterfs/gfapi如何工作的.md)
- [gluste-block安装](./document/md/glusterfs/gluste-block介绍.md) 
- [glusterfs写入一个文件深入分析](./document/md/glusterfs/glusterfs写入一个文件深入分析.md) 
- [event-threads设定后都做了什么](./document/md/glusterfs/event-threads设定后都做了什么.md) 
- [perf分析glusterfs写操作-持续更新](./document/md/glusterfs/perf分析glusterfs写操作.md) 
- [io-thread线程工作方式](./document/md/glusterfs/io-thread线程工作方式.md) 
- [gluster中的group](./document/md/glusterfs/gluster中的group.md) 
- [Gluster如何限制brick的预留空间](./document/md/glusterfs/Gluster如何限制brick的预留空间.md) 
- [fixes-storoge.reserve检测磁盘预留空闲时间间隔问题](./document/md/glusterfs/fixes-storoge.reserve检测磁盘预留空闲时间间隔问题.md) 
- [quota介绍](./document/md/glusterfs/quota介绍.md)

### glusterfs 运维

- [Glusterfs多副本服务端数据丢失演练](./document/md/glusterfs/Glusterfs多副本服务端数据丢失演练.md)
- [too many files引起glusterfsd crash](./document/md/glusterfs/glusterfsd出现crash的分析和总结.md)
- [glusterfs安装及创建卷使用](./document/md/glusterfs/glusterfs安装及创建卷使用.md)
- [glusterfs添加节点成功但状态异常](./document/md/glusterfs/glusterfs添加节点错误.md)
- [glusterfs opencas IO加速方案](./document/md/glusterfs/OpenCAS缓存加速方案.md)
- [glusterfs性能调优](./document/md/glusterfs/glusterfs性能调优.md)
- [Glusterfs-mount选项说明](./document/md/glusterfs/Glusterfs-mount选项说明.md)
- [glusterfs替换掉brick](./document/md/glusterfs/glusterfs替换掉brick.md)
- [Glusterfs扩缩容方法](./document/md/glusterfs/Glusterfs扩缩容方法.md)
- [根据gfid定位到文件具体路径](./document/md/glusterfs/根据gfid定位到文件具体路径.md)
- [posix_disk_space_check_thread_proc函数检测磁盘剩余空间间隔时间可配置方式](./document/md/glusterfs/修改posix_disk_space_check_thread_proc函数检测磁盘剩余空间间隔时间可配置方式.md)



### glusterfs代码提交流程
- [glustefs代码提交流程](./document/md/glusterfs/glusterfs代码提交流程.md)



### 文章中的图片无法显示问题解决

```
1.打开https://www.ipaddress.com/网址
2.查询 raw.githubusercontent.com 域名对应的ip
3.修改C:\Windows\System32\drivers\etc\hosts文件追加如下内容
  199.232.68.133 githubusercontent.com
  199.232.68.133 raw.githubusercontent.com
4.刷新windows的dns，即可访问文章中的图片
```





