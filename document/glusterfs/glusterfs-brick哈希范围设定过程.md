
## glusterfs brick哈希范围设定过程

### 调试环境搭建

- 基本配置
```
rm -rf /glusterfs/*
mkdir /glusterfs/data1/brick
mkdir /glusterfs/data2/brick
mkdir /glusterfs/data3/brick
rm -rf /var/log/glusterfs/bricks/*
gluster volume create dht_debug  172.25.78.14:/glusterfs/data1/brick 172.25.78.14:/glusterfs/data2/brick 172.25.78.14:/glusterfs/data3/brick force
glusterd --log-level TRACE

gluster volume set dht_debug  diagnostics.client-log-level TRACE
gluster volume set dht_debug diagnostics.brick-log-level TRACE
gluster volume start dht_debug 
mount -t glusterfs -o acl 172.25.78.14:/dht_debug /mnt/dht_debug 
```
- glusterd 配置和基本目录说明
```
//glusterd的配置信息，比如peer信息
[root@glusterfs4 ~]# stat /var/lib/glusterd/
//glusterfs的日志目录
[root@glusterfs4 ~]# stat /var/log/glusterfs/
[root@glusterfs4 ~]# cat  /usr/lib/systemd/system/glusterd.service
[Unit]
Description=GlusterFS, a clustered file-system server
Documentation=man:glusterd(8)
Requires=
After=network.target 
Before=network-online.target

[Service]
Type=forking
PIDFile=/var/run/glusterd.pid
LimitNOFILE=65536
Environment="LOG_LEVEL=INFO"
EnvironmentFile=-/etc/sysconfig/glusterd
ExecStart=/usr/local/sbin/glusterd -p /var/run/glusterd.pid  --log-level TRACE  $GLUSTERD_OPTIONS
KillMode=process
SuccessExitStatus=15

[Install]
WantedBy=multi-user.target
```
### 基本线索

- glusterfs哈希计算是在cluster/dht中，这个xlator是被加载到glusterfs(客户端上),重点查看客户端的日志
```
//任何节点第一次挂载时候会初始化每个brick根目录("/")哈希范围,会在每个b rick目录扩展属性里面设置哈希范围，第一次以后直接读取brick上的这个扩展属性,

[client-rpc-fops_v2.c:2633:client4_0_lookup_cbk] 0-stack-trace: stack-address: 0x2b1fa00011e8, dht_debug-client-1 returned 0 
[dht-common.c:1382:dht_lookup_dir_cbk] 0-dht_debug-dht: /: lookup on dht_debug-client-1 returned with op_ret = 0, op_errno = 0 
[dht-layout.c:347:dht_layout_merge] 0-dht_debug-dht: Missing disk layout on dht_debug-client-1. err = -1 
[dht-layout.c:641:dht_layout_normalize] 0-dht_debug-dht: Found anomalies [{path=/}, {gfid=00000000-0000-0000-0000-000000000001}, {holes=1}, {overlaps=0}] 
[dht-common.c:1325:dht_needs_selfheal] 0-dht_debug-dht: fixing assignment on / 
[dht-selfheal.c:1763:dht_selfheal_layout_new_directory] 0-dht_debug-dht: chunk size = 0xffffffff / 22886544 = 187.663428 
[dht-hashfn.c:95:dht_hash_compute] 0-dht_debug-dht: trying regex for / 
[dht-selfheal.c:1799:dht_selfheal_layout_new_directory] 0-dht_debug-dht: assigning range size 0x55555555 to dht_debug-client-1 
[dht-selfheal.c:1800:dht_selfheal_layout_new_directory] 0-dht_debug-dht: gave fix: 0x0 - 0x55555554, with commit-hash 0x1 on dht_debug-client-1 for / 
[dht-selfheal.c:1799:dht_selfheal_layout_new_directory] 0-dht_debug-dht: assigning range size 0x55555555 to dht_debug-client-2 
[dht-selfheal.c:1800:dht_selfheal_layout_new_directory] 0-dht_debug-dht: gave fix: 0x55555555 - 0xaaaaaaa9, with commit-hash 0x1 on dht_debug-client-2 for / 
[dht-selfheal.c:1799:dht_selfheal_layout_new_directory] 0-dht_debug-dht: assigning range size 0x55555555 to dht_debug-client-0 
[dht-selfheal.c:1800:dht_selfheal_layout_new_directory] 0-dht_debug-dht: gave fix: 0xaaaaaaaa - 0xfffffffe, with commit-hash 0x1 on dht_debug-client-0 for / 
[dht-selfheal.c:1025:dht_selfheal_dir_setattr] 0-stack-trace: stack-address: 0x2b1fa00011e8, winding from dht_debug-dht to dht_debug-client-0 
[dht-selfheal.c:1025:dht_selfheal_dir_setattr] 0-stack-trace: stack-address: 0x2b1fa00011e8, winding from dht_debug-dht to dht_debug-client-1 
[dht-selfheal.c:1025:dht_selfheal_dir_setattr] 0-stack-trace: stack-address: 0x2b1fa00011e8, winding from dht_debug-dht to dht_debug-client-2 
```


- 每个birck对应的protocol/client的实例
```
//brick1对应的protocol/client的xlator
volume dht_debug-client-0
    type protocol/client
    option remote-host 172.25.78.14
    option remote-subvolume /glusterfs/data1/brick
    option transport-type socket
    option transport.address-family inet
end-volume

volume dht_debug-client-1
    type protocol/client
    option remote-host 172.25.78.14
    option remote-subvolume /glusterfs/data2/brick
    option transport-type socket
    option transport.address-family inet	
end-volume

volume dht_debug-client-2
    type protocol/client
    option remote-host 172.25.78.14
    option remote-subvolume /glusterfs/data3/brick
    option transport-type socket
    option transport.address-family inet
end-volume

//brick的哈希范围设定在是dht_debug-dht这个xlator去做的
volume dht_debug-dht
    type cluster/distribute
    option lock-migration off
    option force-migration off
    subvolumes dht_debug-client-0 dht_debug-client-1 dht_debug-client-2
end-volume

```


### 调试客户端

#### brick哈希设定核心函数

```
client4_0_lookup_cbk
dht_lookup_dir_cbk
dht_selfheal_directory
dht_selfheal_dir_getafix
dht_selfheal_layout_new_directory
```

#### 调试

```
gdb /usr/local/sbin/glusterfs
(gdb) set args  --acl --process-name fuse --volfile-server=172.25.78.14 --volfile-id=dht_debug /mnt/dht_debug
(gdb) br main
(gdb) br create_fuse_mount 
//出现fork进程时候，设置进入子进程
[Detaching after fork from child process 31104]
(gdb) set follow-fork-mode child 
(gdb) set detach-on-fork off
(gdb) set print pretty on
(gdb) r
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
Breakpoint 1, dht_selfheal_layout_new_directory (frame=0x2aaac8001bf8, loc=0x2aaac8002318, layout=0x2aaab8057e10) at dht-selfheal.c:1719
1719        xlator_t *this = NULL;
(gdb) p this->name
$2 = 0x2aaab8001a30 "dht_debug-dht"
(gdb) bt
#0  dht_selfheal_dir_setattr (frame=0x2aaac8001bf8, loc=0x2aaac8002318, stbuf=0x2aaac80023a8, valid=-1, layout=0x2aaab8057e10) at dht-selfheal.c:980
#1  0x00002aaab6e445b7 in dht_selfheal_dir_mkdir (frame=0x2aaac8001bf8, loc=0x2aaac8002318, layout=0x2aaab8057e10, force=0) at dht-selfheal.c:1369
#2  0x00002aaab6e4656c in dht_selfheal_directory (frame=0x2aaac8001bf8, dir_cbk=0x2aaab6e58f37 <dht_lookup_selfheal_cbk>, loc=0x2aaac8002318, layout=0x2aaab8057e10)
    at dht-selfheal.c:2022
#3  0x00002aaab6e5e846 in dht_lookup_dir_cbk (frame=0x2aaac8001bf8, cookie=0x2aaab800de00, this=0x2aaab8010e30, op_ret=0, op_errno=0, inode=0x2aaab8057988, 
    stbuf=0x2aaab6b788e0, xattr=0x2aaab8057268, postparent=0x2aaab6b78840) at dht-common.c:1577
#4  0x00002aaab6beaa7f in client4_0_lookup_cbk (req=0x2aaab805b318, iov=0x2aaab805b348, count=1, myframe=0x2aaab804e6d8) at client-rpc-fops_v2.c:2632
#5  0x00002aaaab0244a6 in rpc_clnt_handle_reply (clnt=0x2aaab804e980, pollin=0x2aaab8002e60) at rpc-clnt.c:768
#6  0x00002aaaab0249cf in rpc_clnt_notify (trans=0x2aaab804ec80, mydata=0x2aaab804e9b0, event=RPC_TRANSPORT_MSG_RECEIVED, data=0x2aaab8002e60) at rpc-clnt.c:935
#7  0x00002aaaab0209cf in rpc_transport_notify (this=0x2aaab804ec80, event=RPC_TRANSPORT_MSG_RECEIVED, data=0x2aaab8002e60) at rpc-transport.c:520
#8  0x00002aaab60d3c7a in socket_event_poll_in_async (xl=0x2aaab800de00, async=0x2aaab8002f78) at socket.c:2502
#9  0x00002aaab60cc27c in gf_async (async=0x2aaab8002f78, xl=0x2aaab800de00, cbk=0x2aaab60d3c23 <socket_event_poll_in_async>)
    at ../../../../libglusterfs/src/glusterfs/async.h:189
#10 0x00002aaab60d3e08 in socket_event_poll_in (this=0x2aaab804ec80, notify_handled=true) at socket.c:2543
#11 0x00002aaab60d4ccf in socket_event_handler (fd=14, idx=2, gen=4, data=0x2aaab804ec80, poll_in=1, poll_out=0, poll_err=0, event_thread_died=0 '\000')
    at socket.c:2934
#12 0x00002aaaaad77d18 in event_dispatch_epoll_handler (event_pool=0x672ce0, event=0x2aaab6b79040) at event-epoll.c:640
#13 0x00002aaaaad7825b in event_dispatch_epoll_worker (data=0x6d1420) at event-epoll.c:751
#14 0x00002aaaac1d4dd5 in start_thread () from /lib64/libpthread.so.0
```
### getfattr获取brick的哈希范围

```
[root@glusterfs4 bricks]# getfattr  -d -m . -e hex /glusterfs/data3/brick
getfattr: Removing leading '/' from absolute path names
# file: glusterfs/data3/brick
trusted.gfid=0x00000000000000000000000000000001
trusted.glusterfs.dht=0x000000010000000055555555aaaaaaa9
trusted.glusterfs.mdata=0x010000000000000000000000005f561f1c000000002ed98f53000000005f56130100000000325bd66d000000005f5621b8000000000944b075
trusted.glusterfs.volume-id=0x33555f3d1cd541cea2e8a6fd6657a703

[root@glusterfs4 bricks]# getfattr  -d -m . -e hex /glusterfs/data1/brick
getfattr: Removing leading '/' from absolute path names
# file: glusterfs/data1/brick
trusted.gfid=0x00000000000000000000000000000001
trusted.glusterfs.dht=0x0000000100000000aaaaaaaaffffffff
trusted.glusterfs.volume-id=0x33555f3d1cd541cea2e8a6fd6657a703

[root@glusterfs4 bricks]# getfattr  -d -m . -e hex /glusterfs/data2/brick
getfattr: Removing leading '/' from absolute path names
# file: glusterfs/data2/brick
trusted.gfid=0x00000000000000000000000000000001
trusted.glusterfs.dht=0x00000001000000000000000055555554
trusted.glusterfs.mdata=0x010000000000000000000000005f561f1c000000002ed98f53000000005f56130100000000325bd66d000000005f5621b8000000000944b075
trusted.glusterfs.volume-id=0x33555f3d1cd541cea2e8a6fd6657a703
```

### 相关日志
- [1.glusterd服务器](../../document/logs/glusterd.log)
- [2.glusterfs客户端日志](../../document/logs/mnt-dht_debug.log)
- [3.glusterfsd brick1日志](../../document/logs/glusterfs-data1-brick.log)
- [4.glusterfsd brick2日志](../../document/logs/glusterfs-data2-brick.log)
- [5.glusterfsd brick3日志](../../document/logs/glusterfs-data3-brick.log)


