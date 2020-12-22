#  gluster 源码的阅读笔记

| 作者 | 时间 |QQ技术交流群 |
| ------ | ------ |------ |
| perrynzhou@gmail.com |2020/12/01 |中国开源存储技术交流群(672152841) |

## 目标

- minio原理和源代码分析
- glusterfs的架构设计,梳理glusterfs中的哈希卷、副本卷、EC卷的读写过程
- glusterfs块设备实现项目gluster-block原理和源码分析
- libfuse实现原理

## 文章中的图片无法显示问题解决

```
1.打开https://www.ipaddress.com/网址
2.查询 raw.githubusercontent.com 域名对应的ip
3.修改C:\Windows\System32\drivers\etc\hosts文件追加如下内容
  199.232.68.133 githubusercontent.com
  199.232.68.133 raw.githubusercontent.com
4.刷新windows的dns，即可访问文章中的图片
```

## libfuse

- [libfuse3-10源码编译](./document/libfuse/2020-12-06-libfuse-3.10源码编译.md)


## glusterfs 

### glusterfs 

- [Glusterfs多副本服务端数据丢失演练](./document/glusterfs/Glusterfs多副本服务端数据丢失演练.md)
- [too many files引起glusterfsd crash](./document/glusterfs/glusterfsd出现crash的分析和总结.md)
- [glusterfs安装及创建卷使用](./document/glusterfs/glusterfs安装及创建卷使用.md)
- [glusterfs添加节点成功但状态异常](./document/glusterfs/glusterfs添加节点错误.md)
- [glusterfs opencas IO加速方案](./document/opencas/OpenCAS缓存加速方案.md)
- [glusterfs-nfs-ganesha初探](./document/glusterfs/2020-11-24-基于glusterfs的nfs-ganesha方案.md)
- [nfs-ganesha容器部署](./document/glusterfs/2020-12-01-nfs-ganesha容器部署.md)

### glusterfs 
- [glustefs 101](./document/glusterfs101-courses)
- [glustefs调试](./document/glusterfs/glusterfs调试.md)
- [glusterfsd启动过程](./document/glusterfs/glusterfsd启动过程.md)
- [glusterfs客户端挂载init流程](./document/glusterfs/glusterfs客户端挂载init流程.md)
- [gluster-create-volume处理过程](./document/glusterfs/gluster-create-volume处理过程.md)
- [glusterfs架构和基本概念](./document/glusterfs/glusterfs架构和基本概念.md)
- [glusterfs-brick哈希范围设定过程](./document/glusterfs/glusterfs-brick哈希范围设定过程.md)
- [glusterfs客户端写数据分析](./document/glusterfs/glusterfs客户端写数据分析.md)
- [glusterfs问题诊断](./document/glusterfs/glusterfs问题诊断.md)
- [glusterfs-fuse实现(持续更新)](./document/glusterfs/glusterfs-fuse实现.md)
- [cluster.read-hash-mode工作原理](./document/glusterfs/cluster.read-hash-mode工作原理.md)
- [多副本情况下mount挂载目录如何选择可用的副本目录](./document/glusterfs/多副本情况下mount挂载目录如何选择可用的副本目录.md)
- [gfapi如何工作的](./document/glusterfs/2020-11-04-gfapi如何工作的.md)
- [gluste-block安装](./document/glusterfs/gluste-block介绍.md) 
- [nfs-ganesha源码安装](./document/glusterfs/2020-12-20-nfs-ganesha源码安装.md)


### glusterfs贡献的pr

- [glusterfs read-hash-mode的bug](https://github.com/gluster/glusterfs/commit/268faabed00995537394c04ac168c018167fbe27)


### glusterfs代码提交流程
- [glustefs代码提交流程](./document/glusterfs/glusterfs代码提交流程.md)






