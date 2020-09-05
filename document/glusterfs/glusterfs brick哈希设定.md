
## glusterfs brick哈希范围设定

### 调试环境搭建

```
rm -rf /glusterfs/dht/*
mkdir /glusterfs/dht/brick1
mkdir /glusterfs/dht/brick2
mkdir /glusterfs/dht/brick3
rm -rf /var/log/glusterfs/bricks/*
gluster volume create dht3 10.211.55.3:/glusterfs/dht/brick1 10.211.55.3:/glusterfs/dht/brick2 10.211.55.3:/glusterfs/dht/brick3 force
glusterd --log-level TRACE
gluster volume set dht3  diagnostics.client-log-level TRACE
gluster volume set dht3 diagnostics.brick-log-level TRACE
gluster volume start dht3 
mount -t glusterfs -o acl 10.211.55.3:/dht3 /mnt/dht3 
```

### 基本线索

- glusterfs哈希计算是在cluster/dht中，这个xlator是被加载到glusterfs(客户端上),重点查看客户端的日志
- 客户端日志查看发现信息
```
[2020-08-31 01:29:57.098107] T [MSGID: 0] [client-rpc-fops_v2.c:2633:client4_0_lookup_cbk] 0-stack-trace: stack-address: 0x7fdf040011e8, dht3-client-0 returned 0
[2020-08-31 01:29:57.098139] D [MSGID: 0] [dht-common.c:1382:dht_lookup_dir_cbk] 0-dht3-dht: /: lookup on dht3-client-0 returned with op_ret = 0, op_errno = 0
[2020-08-31 01:29:57.098149] T [MSGID: 0] [dht-layout.c:347:dht_layout_merge] 0-dht3-dht: Missing disk layout on dht3-client-0. err = -1
[2020-08-31 01:29:57.098162] D [MSGID: 0] [dht-common.c:1478:dht_lookup_dir_cbk] 0-dht3-dht: /: mds xattr trusted.glusterfs.dht.mds is not present on dht3-client-0(gfid =
 00000000-0000-0000-0000-000000000001)
[2020-08-31 01:29:57.098180] I [MSGID: 109063] [dht-layout.c:641:dht_layout_normalize] 0-dht3-dht: Found anomalies [{path=/}, {gfid=00000000-0000-0000-0000-000000000001},
 {holes=1}, {overlaps=0}]
[2020-08-31 01:29:57.098187] D [MSGID: 0] [dht-common.c:1325:dht_needs_selfheal] 0-dht3-dht: fixing assignment on /
[2020-08-31 01:29:57.098207] D [MSGID: 0] [dht-selfheal.c:1763:dht_selfheal_layout_new_directory] 0-dht3-dht: chunk size = 0xffffffff / 41949 = 102385.451262
[2020-08-31 01:29:57.098214] T [MSGID: 0] [dht-hashfn.c:95:dht_hash_compute] 0-dht3-dht: trying regex for /
[2020-08-31 01:29:57.098228] D [MSGID: 0] [dht-selfheal.c:1799:dht_selfheal_layout_new_directory] 0-dht3-dht: assigning range size 0x55555555 to dht3-client-1
[2020-08-31 01:29:57.098234] T [MSGID: 0] [dht-selfheal.c:1800:dht_selfheal_layout_new_directory] 0-dht3-dht: gave fix: 0x0 - 0x55555554, with commit-hash 0x1 on dht3-cli
ent-1 for /
[2020-08-31 01:29:57.098240] D [MSGID: 0] [dht-selfheal.c:1799:dht_selfheal_layout_new_directory] 0-dht3-dht: assigning range size 0x55555555 to dht3-client-2
[2020-08-31 01:29:57.098245] T [MSGID: 0] [dht-selfheal.c:1800:dht_selfheal_layout_new_directory] 0-dht3-dht: gave fix: 0x55555555 - 0xaaaaaaa9, with commit-hash 0x1 on d
ht3-client-2 for /
[2020-08-31 01:29:57.098251] D [MSGID: 0] [dht-selfheal.c:1799:dht_selfheal_layout_new_directory] 0-dht3-dht: assigning range size 0x55555555 to dht3-client-0
[2020-08-31 01:29:57.098256] T [MSGID: 0] [dht-selfheal.c:1800:dht_selfheal_layout_new_directory] 0-dht3-dht: gave fix: 0xaaaaaaaa - 0xfffffffe, with commit-hash 0x1 on d
ht3-client-0 for /
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
(gdb) r
(gdb) br glusterfs_process_volfp
(gdb) br glusterfs_graph_construct
//出现fork进程时候，设置进入子进程
[Detaching after fork from child process 31104]
(gdb) set follow-fork-mode child 
(gdb) set detach-on-fork off
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
