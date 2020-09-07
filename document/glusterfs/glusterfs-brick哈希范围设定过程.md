
## glusterfs brick哈希范围设定过程

### 调试环境

- 测试卷的部署
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
    subvolumes dht_debug-client-0 dht_debug-client-1 dht_debug-client-2
end-volume

```

### 查找基本的线索

- glusterfs哈希计算是在cluster/dht中，这个xlator是被加载到glusterfs(客户端上),重点查看客户端的日志,发现有对应的函数为"/“目录的设置对应的哈希范围
```
//

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

```
- 任何节点第一次挂载时候会初始化每个brick根目录("/")哈希范围,会在每个b rick目录扩展属性里面设置哈希范围，brick的哈希函数是在volume dht_debug-dht 这个xlator中设定，同时请求服务端会把brick上设置trusted.glusterfs.dht 等于哈希范围，这个相当于持久化到brick上的文件属性上。以后直接读取brick上的这个扩展属性，卷的brick的哈希范围分配是通过dht_selfheal_layout_new_directory函数来设置，这个函数是如何被触发和调用的等，下面将会介绍



### 调试客户端
- gdb调试信息

```
[root@glusterfs4 ~]# gdb /usr/local/sbin/glusterfs
(gdb) set args  --acl --process-name fuse --volfile-server=172.25.78.14 --volfile-id=dht_debug /mnt/dht_debug
(gdb) br create_fuse_mount 
Breakpoint 1 at 0x4072e3: file glusterfsd.c, line 557.
(gdb) r
Breakpoint 1, create_fuse_mount (ctx=0x63c010) at glusterfsd.c:557
557         int ret = 0;
Missing separate debuginfos, use: debuginfo-install glibc-2.17-260.el7.x86_64 keyutils-libs-1.5.8-3.el7.x86_64 krb5-libs-1.15.1-46.el7.x86_64 libcom_err-1.42.9-17.el7.x86_64 libselinux-2.5-15.el7.x86_64 libtirpc-0.2.4-0.16.el7.x86_64 libuuid-2.23.2-63.el7.x86_64 openssl-libs-1.0.2k-19.el7.x86_64 pcre-8.32-17.el7.x86_64 zlib-1.2.7-18.el7.x86_64
(gdb) n
558         cmd_args_t *cmd_args = NULL;
(gdb) 
606         ret = xlator_init(master);
(gdb) 
Detaching after fork from child process 275182.
607         if (ret) {
(gdb) set follow-fork-mode child 
(gdb) set detach-on-fork off
(gdb) set print pretty on
(gdb) br dht_lookup
(gdb) c
Reading symbols from /usr/local/lib/glusterfs/2020.09.07/rpc-transport/socket.so...done.
Reading symbols from /usr/local/lib/glusterfs/2020.09.07/xlator/protocol/client.so...done.
Reading symbols from /usr/local/lib/glusterfs/2020.09.07/xlator/cluster/distribute.so...done.
Reading symbols from /usr/local/lib/glusterfs/2020.09.07/xlator/features/utime.so...done.
Reading symbols from /usr/local/lib/glusterfs/2020.09.07/xlator/performance/write-behind.so...done.
Reading symbols from /usr/local/lib/glusterfs/2020.09.07/xlator/performance/readdir-ahead.so...done.
Reading symbols from /usr/local/lib/glusterfs/2020.09.07/xlator/performance/open-behind.so...done.
Reading symbols from /usr/local/lib/glusterfs/2020.09.07/xlator/performance/quick-read.so...done.
Reading symbols from /usr/local/lib/glusterfs/2020.09.07/xlator/performance/md-cache.so...done.
Reading symbols from /usr/local/lib/glusterfs/2020.09.07/xlator/performance/io-threads.so...done.
Reading symbols from /usr/local/lib/glusterfs/2020.09.07/xlator/debug/io-stats.so...done.
Reading symbols from /usr/local/lib/glusterfs/2020.09.07/xlator/system/posix-acl.so...done.
Reading symbols from /usr/local/lib/glusterfs/2020.09.07/xlator/meta.so...done.

Breakpoint 2, dht_lookup (frame=0x2aaac8001bf8, this=0x2aaab8010e30, loc=0x2aaac4002950, xattr_req=0x2aaac4000b38) at dht-common.c:3492
3492        xlator_t *hashed_subvol = NULL;
(gdb) br dht_do_fresh_lookup
Breakpoint 3 at 0x2aaab706bf11: file dht-common.c, line 3275.
(gdb)  br dht_set_dir_xattr_req
Breakpoint 4 at 0x2aaab706bd70: file dht-common.c, line 3230.
(gdb)  br dht_lookup_cbk
Breakpoint 5 at 0x2aaab706ab5d: file dht-common.c, line 3048.
(gdb)  br dht_lookup_directory
Breakpoint 6 at 0x2aaab705fecc: file dht-common.c, line 1586.
(gdb)  br dht_lookup_dir_cbk
Breakpoint 7 at 0x2aaab705eb6e: file dht-common.c, line 1356.
(gdb)  br dht_selfheal_directory
 Breakpoint 8 at 0x2aaab7046f64: file dht-selfheal.c, line 1929.
(gdb)  br dht_selfheal_dir_getafix
Breakpoint 9 at 0x2aaab7046aec: file dht-selfheal.c, line 1815.
(gdb)  br dht_selfheal_layout_new_directory
Breakpoint 10 at 0x2aaab70464a6: file dht-selfheal.c, line 1719.
(gdb) p this->name
$1 = 0x2aaab8001a30 "dht_debug-dht"
(gdb) p this->type
$2 = 0x2aaab8010790 "cluster/distribute"
(gdb) c
Continuing.

Breakpoint 3, dht_do_fresh_lookup (frame=0x2aaac8001bf8, this=0x2aaab8010e30, loc=0x2aaac4002950) at dht-common.c:3275
3275        int ret = -1;
(gdb) p this->name
$3 = 0x2aaab8001a30 "dht_debug-dht"
(gdb) p this->type
$4 = 0x2aaab8010790 "cluster/distribute"
(gdb) c
Continuing.

Breakpoint 4, dht_set_dir_xattr_req (this=0x2aaab8010e30, loc=0x2aaac4002950, xattr_req=0x2aaac4000b38) at dht-common.c:3230
3230        int ret = -EINVAL;
(gdb) p this->name
$5 = 0x2aaab8001a30 "dht_debug-dht"
(gdb) p this->type
$6 = 0x2aaab8010790 "cluster/distribute"
(gdb) c
Continuing.
[Switching to Thread 0x2aaab6d7b700 (LWP 275599)]

Breakpoint 5, dht_lookup_cbk (frame=0x2aaac8001bf8, cookie=0x2aaab8007790, this=0x2aaab8010e30, op_ret=0, op_errno=0, inode=0x2aaab8000a38, stbuf=0x2aaab6d798e0, 
    xattr=0x2aaabc00b698, postparent=0x2aaab6d79840) at dht-common.c:3048
(gdb) p this->name
$7 = 0x2aaab8001a30 "dht_debug-dht"
(gdb) p this->type
$8 = 0x2aaab8010790 "cluster/distribute"
(gdb) c
Continuing.

Breakpoint 6, dht_lookup_directory (frame=0x2aaac8001bf8, this=0x2aaab8010e30, loc=0x2aaac8002318) at dht-common.c:1586
1586        int call_cnt = 0;
(gdb) p this->type
$9 = 0x2aaab8010790 "cluster/distribute"
(gdb) p this->name
$10 = 0x2aaab8001a30 "dht_debug-dht"
(gdb) c
Continuing.
[Switching to Thread 0x2aaab6b7a700 (LWP 275598)]

Breakpoint 7, dht_lookup_dir_cbk (frame=0x2aaac8001bf8, cookie=0x2aaab8007790, this=0x2aaab8010e30, op_ret=0, op_errno=0, inode=0x2aaab8000a38, stbuf=0x2aaab6b788e0, 
    xattr=0x2aaab8057988, postparent=0x2aaab6b78840) at dht-common.c:1356
1356        dht_local_t *local = NULL;
(gdb) p this->name
$11 = 0x2aaab8001a30 "dht_debug-dht"
(gdb) p this->type
$12 = 0x2aaab8010790 "cluster/distribute"
(gdb) c
Continuing.

Breakpoint 7, dht_lookup_dir_cbk (frame=0x2aaac8001bf8, cookie=0x2aaab800de00, this=0x2aaab8010e30, op_ret=0, op_errno=0, inode=0x2aaab8000a38, stbuf=0x2aaab6b788e0, 
    xattr=0x2aaab8030af8, postparent=0x2aaab6b78840) at dht-common.c:1356
1356        dht_local_t *local = NULL;
(gdb) p this->type
$13 = 0x2aaab8010790 "cluster/distribute"
(gdb) p this->name
$14 = 0x2aaab8001a30 "dht_debug-dht"
(gdb) c
Continuing.
[Switching to Thread 0x2aaab6d7b700 (LWP 275599)]

Breakpoint 7, dht_lookup_dir_cbk (frame=0x2aaac8001bf8, cookie=0x2aaab800add0, this=0x2aaab8010e30, op_ret=0, op_errno=0, inode=0x2aaab8000a38, stbuf=0x2aaab6d798e0, 
    xattr=0x2aaabc0088b8, postparent=0x2aaab6d79840) at dht-common.c:1356
1356        dht_local_t *local = NULL;
(gdb) p this->name
$15 = 0x2aaab8001a30 "dht_debug-dht"
(gdb) p this->type
$16 = 0x2aaab8010790 "cluster/distribute"
(gdb) c
Continuing.

Breakpoint 8, dht_selfheal_directory (frame=0x2aaac8001bf8, dir_cbk=0x2aaab7059f37 <dht_lookup_selfheal_cbk>, loc=0x2aaac8002318, layout=0x2aaabc003970)
    at dht-selfheal.c:1929
1929        dht_local_t *local = NULL;
(gdb) p *loc->path
$17 = 47 '/'
(gdb) c
Continuing.

Breakpoint 9, dht_selfheal_dir_getafix (frame=0x2aaac8001bf8, loc=0x2aaac8002318, layout=0x2aaabc003970) at dht-selfheal.c:1815
1815        dht_local_t *local = NULL;
(gdb) c
Continuing.

Breakpoint 10, dht_selfheal_layout_new_directory (frame=0x2aaac8001bf8, loc=0x2aaac8002318, layout=0x2aaabc003970) at dht-selfheal.c:1719
1719        xlator_t *this = NULL;
(gdb) bt
#0  dht_selfheal_layout_new_directory (frame=0x2aaac8001bf8, loc=0x2aaac8002318, layout=0x2aaabc003970) at dht-selfheal.c:1719
#1  0x00002aaab7046b64 in dht_selfheal_dir_getafix (frame=0x2aaac8001bf8, loc=0x2aaac8002318, layout=0x2aaabc003970) at dht-selfheal.c:1832
#2  0x00002aaab7047515 in dht_selfheal_directory (frame=0x2aaac8001bf8, dir_cbk=0x2aaab7059f37 <dht_lookup_selfheal_cbk>, loc=0x2aaac8002318, layout=0x2aaabc003970)
    at dht-selfheal.c:2015
#3  0x00002aaab705f846 in dht_lookup_dir_cbk (frame=0x2aaac8001bf8, cookie=0x2aaab800add0, this=0x2aaab8010e30, op_ret=0, op_errno=0, inode=0x2aaab8000a38, 
    stbuf=0x2aaab6d798e0, xattr=0x2aaabc0088b8, postparent=0x2aaab6d79840) at dht-common.c:1577
#4  0x00002aaab6deba7f in client4_0_lookup_cbk (req=0x2aaabc005398, iov=0x2aaabc0053c8, count=1, myframe=0x2aaabc001908) at client-rpc-fops_v2.c:2632
#5  0x00002aaaab0244a6 in rpc_clnt_handle_reply (clnt=0x2aaab8052ba0, pollin=0x2aaabc008b80) at rpc-clnt.c:768
#6  0x00002aaaab0249cf in rpc_clnt_notify (trans=0x2aaab8052df0, mydata=0x2aaab8052bd0, event=RPC_TRANSPORT_MSG_RECEIVED, data=0x2aaabc008b80) at rpc-clnt.c:935
#7  0x00002aaaab0209cf in rpc_transport_notify (this=0x2aaab8052df0, event=RPC_TRANSPORT_MSG_RECEIVED, data=0x2aaabc008b80) at rpc-transport.c:520
#8  0x00002aaab60d3c7a in socket_event_poll_in_async (xl=0x2aaab800add0, async=0x2aaabc008c98) at socket.c:2502
#9  0x00002aaab60cc27c in gf_async (async=0x2aaabc008c98, xl=0x2aaab800add0, cbk=0x2aaab60d3c23 <socket_event_poll_in_async>)
    at ../../../../libglusterfs/src/glusterfs/async.h:189
#10 0x00002aaab60d3e08 in socket_event_poll_in (this=0x2aaab8052df0, notify_handled=true) at socket.c:2543
#11 0x00002aaab60d4ccf in socket_event_handler (fd=14, idx=4, gen=1, data=0x2aaab8052df0, poll_in=1, poll_out=0, poll_err=0, event_thread_died=0 '\000')
    at socket.c:2934
#12 0x00002aaaaad77d18 in event_dispatch_epoll_handler (event_pool=0x672ce0, event=0x2aaab6d7a040) at event-epoll.c:640
#13 0x00002aaaaad7825b in event_dispatch_epoll_worker (data=0x6d1710) at event-epoll.c:751
#14 0x00002aaaac1d4dd5 in start_thread () from /lib64/libpthread.so.0
#15 0x00002aaaac949ead in clone () from /lib64/libc.so.6
(gdb) info break
Num     Type           Disp Enb Address            What
1       breakpoint     keep y   <MULTIPLE>         
        breakpoint already hit 1 time
1.1                         y     0x00000000004072e3 in create_fuse_mount at glusterfsd.c:557 inf 1
1.2                         y     0x00000000004072e3 in create_fuse_mount at glusterfsd.c:557 inf 2
2       breakpoint     keep y   0x00002aaab706df41 in dht_lookup at dht-common.c:3492 inf 2
        breakpoint already hit 1 time
3       breakpoint     keep y   0x00002aaab706bf11 in dht_do_fresh_lookup at dht-common.c:3275 inf 2
        breakpoint already hit 1 time
4       breakpoint     keep y   0x00002aaab706bd70 in dht_set_dir_xattr_req at dht-common.c:3230 inf 2
        breakpoint already hit 1 time
5       breakpoint     keep y   0x00002aaab706ab5d in dht_lookup_cbk at dht-common.c:3048 inf 2
        breakpoint already hit 1 time
6       breakpoint     keep y   0x00002aaab705fecc in dht_lookup_directory at dht-common.c:1586 inf 2
        breakpoint already hit 1 time
7       breakpoint     keep y   0x00002aaab705eb6e in dht_lookup_dir_cbk at dht-common.c:1356 inf 2
        breakpoint already hit 3 times
8       breakpoint     keep y   0x00002aaab7046f64 in dht_selfheal_directory at dht-selfheal.c:1929 inf 2
        breakpoint already hit 1 time
9       breakpoint     keep y   0x00002aaab7046aec in dht_selfheal_dir_getafix at dht-selfheal.c:1815 inf 2
        breakpoint already hit 1 time
10      breakpoint     keep y   0x00002aaab70464a6 in dht_selfheal_layout_new_directory at dht-selfheal.c:1719 inf 2
        breakpoint already hit 1 time
(gdb) 
```

- brick哈希设定核心函数，经过gdb调试信息分析发现dht_selfheal_layout_new_directory的调用链如下

```
//这里忽略了glusterfs的启动过程，直接在cluster/distribute的dht xlator上设置断点，如果想要查看客户端的启动过程，可以添加如下的断点信息

mgmt_getspec_cbk  
glusterfs_volfile_fetch_one
glusterfs_volfile_fetch 
glusterfs_mgmt_init
glusterfs_volumes_init
glusterfs_process_volfp
glusterfs_graph_construct


//gdb glustefs客户端，调用顺序依次是从上到下
create_fuse_mount
dht_lookup
dht_do_fresh_lookup
dht_set_dir_xattr_req
dht_lookup_cbk
dht_lookup_directory
dht_lookup_dir_cbk
dht_selfheal_directory
dht_selfheal_dir_getafix
//brick的哈希范围设置
dht_selfheal_layout_new_directory
```
- 函数功能和简单说明
  - dht_lookup
  - dht_do_fresh_lookup
  - dht_set_dir_xattr_req
  - dht_lookup_cbk
  - dht_lookup_directory
  - dht_lookup_dir_cbk
  - dht_selfheal_dir_getafix
  - dht_selfheal_layout_new_directory
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

### gdb调试脚本

```
//断点信息
br create_fuse_mount
br dht_lookup
br dht_do_fresh_lookup
br dht_set_dir_xattr_req
br dht_lookup_cbk
br dht_lookup_directory
br dht_lookup_dir_cbk
br dht_selfheal_directory
br dht_selfheal_dir_getafix
br dht_selfheal_layout_new_directory
 

//追踪子进程
set follow-fork-mode child 
set detach-on-fork off
set print pretty on


//测试卷的创建和删除

rm -rf /var/log/glusterfs/bricks/*
rm -rf /var/log/glusterfs/mnt*


rm -rf   /glusterfs/data3/brick
rm -rf   /glusterfs/data2/brick
rm -rf   /glusterfs/data1/brick

mkdir  /glusterfs/data3/brick
mkdir  /glusterfs/data2/brick
mkdir  /glusterfs/data1/brick

gluster volume create dht_debug  172.25.78.14:/glusterfs/data1/brick 172.25.78.14:/glusterfs/data2/brick 172.25.78.14:/glusterfs/data3/brick force
gluster volume set dht_debug  diagnostics.client-log-level TRACE
gluster volume set dht_debug diagnostics.brick-log-level TRACE
gluster volume start dht_debug


gluster volume stop dht_debug
gluster volume delete dht_debug

//调试命令和调试参数设置
gdb /usr/local/sbin/glusterfs
set args  --acl --process-name fuse --volfile-server=172.25.78.14 --volfile-id=dht_debug /mnt/dht_debug


```
### 相关日志
- [1.glusterd服务器](../../document/logs/glusterd.log)
- [2.glusterfs客户端日志](../../document/logs/mnt-dht_debug.log)
- [3.glusterfsd brick1日志](../../document/logs/glusterfs-data1-brick.log)
- [4.glusterfsd brick2日志](../../document/logs/glusterfs-data2-brick.log)
- [5.glusterfsd brick3日志](../../document/logs/glusterfs-data3-brick.log)


