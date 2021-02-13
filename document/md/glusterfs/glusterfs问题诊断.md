

# glusterfs问题诊断方法

| 作者 | 时间 |QQ技术交流群 |
| ------ | ------ |------ |
| perrynzhou@gmail.com |2020/12/01 |中国开源存储技术交流群(672152841) |


### 检查硬盘是否有故障

```
smartctl -H /dev/sdm1
```
### 查看进程文件描述符

```
# ps -ef|grep 169810
root     107496 175632  0 16:21 pts/1    00:00:00 grep --color=auto 169810
root     169810      1  7 Oct10 ?        03:17:25 /usr/sbin/glusterfsd -s 192.168.12.132 --volfile-id sharing_vol.192.168.12.132.data11-brick_sharing_vol -p /var/run/gluster/vols/sharing_vol/192.168.12.132-data11-brick_sharing_vol.pid -S /var/run/gluster/4db36e75931c2470.socket --brick-name /data11/brick_sharing_vol -l /var/log/glusterfs/bricks/data11-brick_sharing_vol.log --xlator-option *-posix.glusterd-uuid=e4abe33a-6b84-4b55-becf-c6354afa0926 --process-name brick --brick-port 49158 --xlator-option sharing_vol-server.listen-port=49158

# ls /proc/169810/fd |wc -l
849

# grep open  /proc/169810/limits
Max open files            1048576              1048576              files    

```
### glusterfs调试诊断

```
gdb --args /usr/local/sbin/gluster  --acl --process-name fuse --volfile-server=192.168.15.153 --volfile-id=rep3_vol /mnt/rep3_vol

or

gdb /usr/local/sbin/gluster
set args --acl --process-name fuse --volfile-server=192.168.15.153 --volfile-id=rep3_vol /mnt/rep3_vol

// create_fuse_mount函数中调用xlator_init函数返回后设置如下两个GDB的选型
set follow-fork-mode child
set detach-on-fork off
```

### EC卷命令行使用说明

```
//如果想要建立一个4+2（6)的EC卷，可以用如下命令来做

gluster volume create disperse 6 redundancy 2 node1:/brick node2:/brick node3:/brick node4:/brick node5:/brick node6:/brick

//或者按照如下方式建立EC卷

gluster volume create disperse-data 4 redundancy 2 node1:/brick node2:/brick node3:/brick node4:/brick node5:/brick node6:/brick
```
### 针对进程的资源消耗

```
yum install perf 
perf top -p {进程号}

//收集信息，会生成perf.data文件
perf record -e cpu-clock -g -p {进程号}
perf record -p `pgrep -d ',' glusterfsd`

// 分析perf.data
perf report -i perf.data

// call-graph
sudo perf record --call-graph dwarf -p {进程PID}
sudo perf script | FlameGraph/stackcollapse-perf.pl | FlameGraph/flamegraph.pl > process.svg
```
### 查看进程D状态
- 进程D状态，一般是进程等待IO，处于D状态的进程是无法kill，只能reboot机器才能解决，如何查看进程处于D状态，按照如下方法
```
$ ps -eo ppid,pid,user,stat,pcpu,comm,wchan:32


//这个命令可以把D状态的进程的内核栈信息trace到/var/log/messages中
$ echo w > /proc/sysrq-trigger
```


### 查看glusterfs卷相关状态  

```
gluster volume status volume_name
Lists status information for each brick in the volume.

gluster volume status volume_name detail
Lists more detailed status information for each brick in the volume.

gluster volume status volume_name clients
Lists the clients connected to the volume.

gluster volume status volume_name mem
Lists the memory usage and memory pool details for each brick in the volume.

gluster volume status volume_name inode
Lists the inode tables of the volume.

gluster volume status volume_name fd
Lists the open file descriptor tables of the volume.

gluster volume status volume_name callpool
Lists the pending calls for the volume.
```
### glusterfs设置进程的调试级别

```
glusterd --log-level TRACE

gluster volume set dht_debug  diagnostics.client-log-level TRACE
gluster volume set dht_debug diagnostics.brick-log-level TRACE
```
### glusterfs二进制调试方法

```
$ gdb /usr/local/sbin/glusterfs 
$ set args --acl --process-name fuse --volfile-server=10.193.189.153 --volfile-id=rep3_vol /mnt/rep3_vol
$ br main
```

### glusterfs 开启debug.trace模块

- 开启设置
```
// 可以查看请求的是哪一个FOP
//客户端启用
gluster volume set volname debug.trace client
gluster volume set volname debug.log-file 1
// 在服务端启用
gluster volume set volname debug.trace posix
gluster volume set volname debug.log-file 1

```

- debug.trace选项
```
static volgen_brick_xlator_t server_graph_table[] = {
    {brick_graph_add_server, NULL},
    {brick_graph_add_io_stats, "NULL"},
    {brick_graph_add_sdfs, "sdfs"},
    {brick_graph_add_namespace, "namespace"},
    {brick_graph_add_cdc, NULL},
    {brick_graph_add_quota, "quota"},
    {brick_graph_add_index, "index"},
    {brick_graph_add_barrier, NULL},
    {brick_graph_add_marker, "marker"},
    {brick_graph_add_selinux, "selinux"},
    {brick_graph_add_fdl, "fdl"},
    {brick_graph_add_iot, "io-threads"},
    {brick_graph_add_upcall, "upcall"},
    {brick_graph_add_leases, "leases"},
    {brick_graph_add_pump, NULL},
    {brick_graph_add_ro, NULL},
    {brick_graph_add_worm, NULL},
    {brick_graph_add_locks, "locks"},
    {brick_graph_add_acl, "acl"},
    {brick_graph_add_bitrot_stub, "bitrot-stub"},
    {brick_graph_add_changelog, "changelog"},
    {brick_graph_add_trash, "trash"},
    {brick_graph_add_arbiter, "arbiter"},
    {brick_graph_add_posix, "posix"},
};
```
### perf查看gluterfs相关进程函数

```
// 列举出当前可以采集的指标集合
perf list

//采集进程112547 CPU 时间消耗分析
perf record -e cpu-clock -g -p 112547

//分析采集到的数据
perf report -i perf.data
```

### glusterd的service模式配置

```
glusterd --log-level TRACE
```
```
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


[root@glusterfs4 ~]# systemctl daemon-reload
```
###  glusterfs客户端进程的statedump

  
- 生成statedump信息
```
//针对glusterfd/glusterfsd/glusterfs进程启动一个statedump
kill -SIGUSR1 {glusterd/glusterfsd/glusterfs-process-pod}

//statedump保存路径
/var/run/gluster/
```

### 显示file的gfid挂载方式

```
$ mount -t glusterfs -o aux-gfid-mount vm1:test /mnt/testvol

$ getfattr -n glusterfs.gfid.string  /mnt/rep_test/test1
getfattr: Removing leading '/' from absolute path names
# file: mnt/rep_test/test1
glusterfs.gfid.string="b85f1ece-7d38-41c6-873d-79a4b14f99f4"
```
### 查看文件的分布式的的信息

```
# getfattr -n trusted.glusterfs.pathinfo -e text  /data/glusterfs_speech_04_v6/11085164/espnet/hello_input2/espnet_aishell/data/train/wav.scp 
getfattr: Removing leading '/' from absolute path names
# file: data/glusterfs_speech_04_v6/11085164/espnet/hello_input2/espnet_aishell/data/train/wav.scp
trusted.glusterfs.pathinfo="(<DISTRIBUTE:speech_v6_rep_vol-dht> (<REPLICATE:speech_v6_rep_vol-replicate-1> <POSIX(/speech_v6/data2/brick):node.hello.lan:/speech_v6/data2/brick/11085164/espnet/hello_input2/espnet_aishell/data/train/wav.scp> <POSIX(/speech_v6/data2/brick):test-node:/speech_v6/data2/brick/11085164/espnet/hello_input2/espnet_aishell/data/train/wav.scp> <POSIX(/speech_v6/data2/brick):ai-storage-center-prd-10-194-39-6.v-bj-4.hello.lan:/speech_v6/data2/brick/11085164/espnet/hello_input2/espnet_aishell/data/train/wav.scp>))"
```