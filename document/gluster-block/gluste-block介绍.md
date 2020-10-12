## gluste-block 介绍


## 服务器节点
- glusterfs安装

```
yum install centos-release-gluster
yum install glusterfs-fuse.x86_64 glusterfs-server.x86_64 glusterfs-cli.x86_64 glusterfs-api.x86_64 glusterfs-libs.x86_64 glusterfs-rdma.x86_64 glusterfs-api-devel.x86_64  glusterfs-client-xlators.x86_64 glusterfs-extra-xlators.x86_64 glusterfs-geo-replication.x86_64 glusterfs-client-xlators.x86_64 glusterfs-api.x86_64 -y
systemctl start glusterd
```

- gluster-block RPM安装

```
yum install install gluster-block
systemctl start gluster-block
```

- 创建块设备
```
//创建卷以后，在创建块设备
gluster-block create {卷名称/块设备名称} 172.25.78.12 300GiB --json-pretty
```
- gluster-block源码安装
  
```
# git clone https://github.com/gluster/gluster-block.git
# cd gluster-block/

# dnf install gcc autoconf automake make file libtool libuuid-devel json-c-devel glusterfs-api-devel glusterfs-server tcmu-runner targetcli

On Fedora27 and Centos7 [Which use legacy glibc RPC], pass '--enable-tirpc=no' flag at configure time
# ./autogen.sh  
# CFLAGS="-ggdb3 -O0" CLIBS="-lgfapi" ./configure --enable-tirpc=no && make -j install
```

## 客户端节点

```
yum install install iscsi-initiator-utils
//会发现在客户客户节点出现了一个虚拟盘
iscsiadm -m discovery -t st -p 172.25.78.12 -l

```