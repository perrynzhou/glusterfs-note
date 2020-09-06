
## glusterfs brick哈希范围设定

### 调试环境搭建

```
rm -rf /glusterfs/dht/*
mkdir /glusterfs/dht/brick1
mkdir /glusterfs/dht/brick2
mkdir /glusterfs/dht/brick3
rm -rf /var/log/glusterfs/bricks/*
gluster volume create dht_test  10.211.55.3:/glusterfs/dht/brick1 10.211.55.3:/glusterfs/dht/brick2 10.211.55.3:/glusterfs/dht/brick3 force
glusterd --log-level TRACE

gluster volume set dht_test  diagnostics.client-log-level TRACE
gluster volume set dht_test diagnostics.brick-log-level TRACE
gluster volume start dht_test 
mount -t glusterfs -o acl 10.211.55.3:/dht3 /mnt/dht3 
```

### 基本线索

- glusterfs哈希计算是在cluster/dht中，这个xlator是被加载到glusterfs(客户端上),重点查看客户端的日志
- 客户端日志查看发现信息
```
//任何节点第一次挂载时候会初始化每个brick根目录("/")哈希范围,会在每个b rick目录扩展属性里面设置哈希范围，第一次以后直接读取brick上的这个扩展属性

[2020-09-06 05:43:54.069230] D [MSGID: 0] [dht-common.c:1478:dht_lookup_dir_cbk] 0-dht_test-dht: /: mds xattr trusted.glusterfs.dht.mds is not present on dht_test-client-0(gfid = 00000000-0000-0000-0000-000000000001) 
[2020-09-06 05:43:54.069239] T [MSGID: 0] [client-rpc-fops_v2.c:2633:client4_0_lookup_cbk] 0-stack-trace: stack-address: 0x7fdcdc0011e8, dht_test-client-2 returned 0 
[2020-09-06 05:43:54.069247] D [MSGID: 0] [dht-common.c:1382:dht_lookup_dir_cbk] 0-dht_test-dht: /: lookup on dht_test-client-2 returned with op_ret = 0, op_errno = 0 
[2020-09-06 05:43:54.069255] T [MSGID: 0] [dht-layout.c:347:dht_layout_merge] 0-dht_test-dht: Missing disk layout on dht_test-client-2. err = -1 
[2020-09-06 05:43:54.069263] D [MSGID: 0] [dht-common.c:1478:dht_lookup_dir_cbk] 0-dht_test-dht: /: mds xattr trusted.glusterfs.dht.mds is not present on dht_test-client-2(gfid = 00000000-0000-0000-0000-000000000001) 
[2020-09-06 05:43:54.069280] I [MSGID: 109063] [dht-layout.c:641:dht_layout_normalize] 0-dht_test-dht: Found anomalies [{path=/}, {gfid=00000000-0000-0000-0000-000000000001}, {holes=1}, {overlaps=0}] 
[2020-09-06 05:43:54.069286] D [MSGID: 0] [dht-common.c:1325:dht_needs_selfheal] 0-dht_test-dht: fixing assignment on / 
[2020-09-06 05:43:54.069300] D [MSGID: 0] [dht-selfheal.c:1763:dht_selfheal_layout_new_directory] 0-dht_test-dht: chunk size = 0xffffffff / 41949 = 102385.451262 
[2020-09-06 05:43:54.069305] T [MSGID: 0] [dht-hashfn.c:95:dht_hash_compute] 0-dht_test-dht: trying regex for / 
[2020-09-06 05:43:54.069315] D [MSGID: 0] [dht-selfheal.c:1799:dht_selfheal_layout_new_directory] 0-dht_test-dht: assigning range size 0x55555555 to dht_test-client-1 
[2020-09-06 05:43:54.069320] T [MSGID: 0] [dht-selfheal.c:1800:dht_selfheal_layout_new_directory] 0-dht_test-dht: gave fix: 0x0 - 0x55555554, with commit-hash 0x1 on dht_test-client-1 for / 
[2020-09-06 05:43:54.069325] D [MSGID: 0] [dht-selfheal.c:1799:dht_selfheal_layout_new_directory] 0-dht_test-dht: assigning range size 0x55555555 to dht_test-client-2 
[2020-09-06 05:43:54.069330] T [MSGID: 0] [dht-selfheal.c:1800:dht_selfheal_layout_new_directory] 0-dht_test-dht: gave fix: 0x55555555 - 0xaaaaaaa9, with commit-hash 0x1 on dht_test-client-2 for / 
[2020-09-06 05:43:54.069334] D [MSGID: 0] [dht-selfheal.c:1799:dht_selfheal_layout_new_directory] 0-dht_test-dht: assigning range size 0x55555555 to dht_test-client-0 
[2020-09-06 05:43:54.069339] T [MSGID: 0] [dht-selfheal.c:1800:dht_selfheal_layout_new_directory] 0-dht_test-dht: gave fix: 0xaaaaaaaa - 0xfffffffe, with commit-hash 0x1 on dht_test-client-0 for / 
[2020-09-06 05:43:54.069348] T [MSGID: 0] [dht-selfheal.c:1025:dht_selfheal_dir_setattr] 0-stack-trace: stack-address: 0x7fdcdc0011e8, winding from dht_test-dht to dht_test-client-0 
[2020-09-06 05:43:54.072671] T [MSGID: 0] [dht-layout.c:347:dht_layout_merge] 0-dht_test-dht: Missing disk layout on dht_test-client-0. err = -1 
[2020-09-06 05:43:54.072677] T [MSGID: 0] [dht-selfheal.c:900:dht_selfheal_dir_xattr] 0-dht_test-dht: 3 subvolumes missing xattr for / 
[2020-09-06 05:43:54.072689] D [MSGID: 109036] [dht-common.c:11213:dht_log_new_layout_for_dir_selfheal] 0-dht_test-dht: Setting layout of / with [Subvol_name: dht_test-client-0, Err: -1 , Start: 0xaaaaaaaa, Stop: 0xffffffff, Hash: 0x1], [Subvol_name: dht_test-client-1, Err: -1 , Start: 0x0, Stop: 0x55555554, Hash: 0x1], [Subvol_name: dht_test-client-2, Err: -1 , Start: 0x55555555, Stop: 0xaaaaaaa9, Hash: 0x1],  
[2020-09-06 05:43:54.072704] T [MSGID: 0] [dht-selfheal.c:769:dht_selfheal_dir_xattr_persubvol] 0-dht_test-dht: setting hash range 0xaaaaaaaa - 0xffffffff (type 0) on subvolume dht_test-client-0 for / 
[2020-09-06 05:43:54.072711] T [MSGID: 0] [dht-selfheal.c:795:dht_selfheal_dir_xattr_persubvol] 0-stack-trace: stack-address: 0x7fdcdc0011e8, winding from dht_test-dht to dht_test-client-0 
[2020-09-06 05:43:54.072727] D [MSGID: 101016] [glusterfs3.h:781:dict_to_xdr] 0-dict: key 'trusted.glusterfs.dht' would not be sent on wire in the future [Invalid argument]
[2020-09-06 05:43:54.072780] T [MSGID: 0] [dht-selfheal.c:769:dht_selfheal_dir_xattr_persubvol] 0-dht_test-dht: setting hash range 0x0 - 0x55555554 (type 0) on subvolume dht_test-client-1 for / 
[2020-09-06 05:43:54.072786] T [MSGID: 0] [dht-selfheal.c:795:dht_selfheal_dir_xattr_persubvol] 0-stack-trace: stack-address: 0x7fdcdc0011e8, winding from dht_test-dht to dht_test-client-1 
[2020-09-06 05:43:54.072833] T [MSGID: 0] [dht-selfheal.c:769:dht_selfheal_dir_xattr_persubvol] 0-dht_test-dht: setting hash range 0x55555555 - 0xaaaaaaa9 (type 0) on subvolume dht_test-client-2 for / 
[2020-09-06 05:43:54.072840] T [MSGID: 0] [dht-selfheal.c:795:dht_selfheal_dir_xattr_persubvol] 0-stack-trace: stack-address: 0x7fdcdc0011e8, winding from dht_test-dht to dht_test-client-2 

```
- 每个birck对应的protocol/client的实例
```
//brick1对应的protocol/client的xlator
volume dht3-client-0
    type protocol/client
    option remote-host 10.211.55.3
    option remote-subvolume /glusterfs/dht/brick1
end-volume

//brick2对应的protocol/client的xlator
volume dht3-client-1
    type protocol/client
    option remote-host 10.211.55.3
    option remote-subvolume /glusterfs/dht/brick2
end-volume

//brick3对应的protocol/client的xlator
volume dht3-client-2
    type protocol/client
    option remote-host 10.211.55.3
    option remote-subvolume /glusterfs/dht/brick3
end-volume
```


### 调试客户端

#### 目录哈希范围分配核心函数

```
dht_selfheal_new_directory
dht_selfheal_layout_new_directory
dht_mkdir_cbk
client4_0_mkdir_cbk
```

#### 调试

```
 gdb /usr/local/sbin/glusterfs 
(gdb)  set args --acl --process-name fuse --volfile-server=10.211.55.3 --volfile-id=dht3 /mnt/dht-31
(gdb) br main
(gdb) br create_fuse_mount 
//出现fork进程时候，设置进入子进程
[Detaching after fork from child process 31104]
(gdb) set follow-fork-mode child 
(gdb) set detach-on-fork off
(gdb) r
(gdb) br glusterfs_process_volfp
(gdb) br glusterfs_graph_construct
//当进入glusterfs_graph_construct时候，出现如下加载链接库时候设置每个xlator的函数
Reading symbols from /usr/local/lib/glusterfs/2020.08.30/xlator/protocol/client.so...done.
Reading symbols from /usr/local/lib/glusterfs/2020.08.30/xlator/cluster/distribute.so...done.
Reading symbols from /usr/local/lib/glusterfs/2020.08.30/xlator/features/utime.so...done.
Reading symbols from /usr/local/lib/glusterfs/2020.08.30/xlator/performance/write-behind.so...done.
Reading symbols from /usr/local/lib/glusterfs/2020.08.30/xlator/performance/readdir-ahead.so...done.
Reading symbols from /usr/local/lib/glusterfs/2020.08.30/xlator/performance/open-behind.so...done.
Reading symbols from /usr/local/lib/glusterfs/2020.08.30/xlator/performance/quick-read.so...done.
Reading symbols from /usr/local/lib/glusterfs/2020.08.30/xlator/performance/md-cache.so...done.
Reading symbols from /usr/local/lib/glusterfs/2020.08.30/xlator/performance/io-threads.so...done.
Reading symbols from /usr/local/lib/glusterfs/2020.08.30/xlator/debug/io-stats.so...done.
(gdb) br dht_layout_search
(gdb) br dht_hash_compute
(gdb) br dht_selfheal_layout_new_directory
```
### 命令行获取brick的哈希范围

```
[root@CentOS1 ~]$ getfattr  -d -m . -e hex /glusterfs/dht/brick1 
getfattr: Removing leading '/' from absolute path names
# file: glusterfs/dht/brick1
security.selinux=0x756e636f6e66696e65645f753a6f626a6563745f723a64656661756c745f743a733000
trusted.gfid=0x00000000000000000000000000000001
trusted.glusterfs.dht=0x0000000100000000aaaaaaaaffffffff
trusted.glusterfs.mdata=0x010000000000000000000000005f4c52950000000006421021000000005f4c51ec000000000c071c00000000005f4c5295000000000272bfd1
trusted.glusterfs.volume-id=0xafe16957a35147a09a5934c23ba0e09a

[root@CentOS1 ~]$ getfattr  -d -m . -e hex /glusterfs/dht/brick2
getfattr: Removing leading '/' from absolute path names
# file: glusterfs/dht/brick2
security.selinux=0x756e636f6e66696e65645f753a6f626a6563745f723a64656661756c745f743a733000
trusted.gfid=0x00000000000000000000000000000001
trusted.glusterfs.dht=0x00000001000000000000000055555554
trusted.glusterfs.mdata=0x010000000000000000000000005f4c52950000000006421021000000005f4c51ec000000000c071c00000000005f4c5295000000000272bfd1
trusted.glusterfs.volume-id=0xafe16957a35147a09a5934c23ba0e09a

[root@CentOS1 ~]$ getfattr  -d -m . -e hex /glusterfs/dht/brick3
getfattr: Removing leading '/' from absolute path names
# file: glusterfs/dht/brick3
security.selinux=0x756e636f6e66696e65645f753a6f626a6563745f723a64656661756c745f743a733000
trusted.gfid=0x00000000000000000000000000000001
trusted.glusterfs.dht=0x000000010000000055555555aaaaaaa9
trusted.glusterfs.mdata=0x010000000000000000000000005f4c52950000000006421021000000005f4c51ec000000000c071c00000000005f4c5295000000000272bfd1
trusted.glusterfs.volume-id=0xafe16957a35147a09a5934c23ba0e09a
```
