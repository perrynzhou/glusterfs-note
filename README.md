#  glusterfs-note 阅读笔记

## 目标

- glusterfs的架构设计
- 梳理glusterfs中的哈希卷、副本卷、EC卷的读写过程


## glusterfs源码分析
- [1.glustefs调试](./document/glusterfs/glusterfs调试.md)
- [2.glusterfsd启动过程](./document/glusterfs/glusterfsd启动过程.md)
- [3.glusterfs客户端挂载init流程](./document/glusterfs/glusterfs客户端挂载init流程.md)
- [4.gluster-create-volume处理过程](./document/glusterfs/gluster-create-volume处理过程.md)
- [5.glusterfs架构和基本概念](./document/glusterfs/glusterfs架构和基本概念.md)
- [6.glusterfs-brick哈希范围设定过程](./document/glusterfs/glusterfs-brick哈希范围设定过程.md)
- [7.glusterfs-write调用链分析](./document/glusterfs/glusterfs-write调用链分析.md)
- [8.glusterfs问题诊断和调试方法](./document/glusterfs/glusterfs问题诊断和调试方法.md)
- [9.glusterfs-fuse实现(持续更新)](./document/glusterfs/glusterfs-fuse实现.md)
- [10.cluster.read-hash-mode工作原理(持续更新)](./document/glusterfs/cluster.read-hash-mode工作原理.md)

## gluster-block使用
- [1.gluste-block安装](./document/gluster-block/gluste-block介绍.md)
## gluster-block源码分析



## glusterfs官方issue

- 源码分析相关
  - [doubt for dht_selfheal_layout_new_directory and trusted.glusterfs.mdata](https://github.com/gluster/glusterfs/issues/1467)

- 性能相关
  - [performance bottleneck about glusterfs](https://github.com/gluster/glusterfs/issues/1462)
- 使用相关 
  - [{features.shard}:sharding-mount glusterfs volume, files larger than 64Mb only show 64Mb](https://github.com/gluster/glusterfs/issues/1384)
  - [{features.shard}:Copying large files (with shard on) fails](https://github.com/gluster/glusterfs/issues/1474)