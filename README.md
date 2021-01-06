#  gluster 源码的阅读笔记

| 作者 | 时间 |QQ技术交流群 |
| ------ | ------ |------ |
| perrynzhou@gmail.com |2020/12/01 |中国开源存储技术交流群(672152841) |



### glusterfs/pNFS/NFS/fuse public document
- [glustefs](./document/pdf/glusterfs)
- [pNFS](./document/pdf/pNFS)
- [NFS](./document/pdf/NFS)
- [fuse](./document/pdf/fuse)


### fuse 原理分析
- [libfuse3-10源码编译](./document/md/fuse/2020-12-06-libfuse-3.10源码编译.md)


### glusterfs 原理分析
- [glusterfs源码安装](./document/md/glusterfs/glusterfs源码安装.md)
- [glusterfsd启动过程](./document/md/glusterfs/glusterfsd启动过程.md)
- [glusterfs客户端挂载init流程](./document/md/glusterfs/glusterfs客户端挂载init流程.md)
- [gluster-create-volume处理过程](./document/md/glusterfs/gluster-create-volume处理过程.md)
- [glusterfs架构和基本概念](./document/md/glusterfs/glusterfs架构和基本概念.md)
- [glusterfs-brick哈希范围设定过程](./document/md/glusterfs/glusterfs-brick哈希范围设定过程.md)
- [glusterfs客户端写数据分析](./document/md/glusterfs/glusterfs客户端写数据分析.md)
- [glusterfs诊断](./document/md/glusterfs/glusterfs诊断.md)
- [glusterfs-fuse实现(持续更新)](./document/md/glusterfs/glusterfs-fuse实现.md)
- [cluster.read-hash-mode工作原理](./document/md/glusterfs/cluster.read-hash-mode工作原理.md)
- [多副本情况下mount挂载目录如何选择可用的副本目录](./document/md/glusterfs/多副本情况下mount挂载目录如何选择可用的副本目录.md)
- [gfapi如何工作的](./document/md/glusterfs/2020-11-04-gfapi如何工作的.md)
- [gluste-block安装](./document/md/glusterfs/gluste-block介绍.md) 
- [glusterfs目录创建深入分析](./document/md/glusterfs/2020-12-25-glusterfs目录创建深入分析.md) 
- [event-threads设定后都做了什么](./document/md/glusterfs/event-threads设定后都做了什么.md) 

### glusterfs 运维

- [Glusterfs多副本服务端数据丢失演练](./document/md/glusterfs/Glusterfs多副本服务端数据丢失演练.md)
- [too many files引起glusterfsd crash](./document/md/glusterfs/glusterfsd出现crash的分析和总结.md)
- [glusterfs安装及创建卷使用](./document/md/glusterfs/glusterfs安装及创建卷使用.md)
- [glusterfs添加节点成功但状态异常](./document/md/glusterfs/glusterfs添加节点错误.md)
- [glusterfs opencas IO加速方案](./document/md/glusterfs/OpenCAS缓存加速方案.md)
- [2020-12-24-glusterfs性能调优.md](./document/md/glusterfs/2020-12-24-glusterfs性能调优.md)
- [2020-12-31-Glusterfs-mount选项说明](./document/md/glusterfs/2020-12-31-Glusterfs-mount选项说明.md)

### glusterfs贡献的pr

- [glusterfs read-hash-mode的bug](https://github.com/gluster/md/glusterfs/commit/268faabed00995537394c04ac168c018167fbe27)


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





