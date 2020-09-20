

# glusterfs问题诊断方法
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
[root@CentOS1 ~]$ getfattr -n glusterfs.gfid.string  /mnt/rep_test/test1
getfattr: Removing leading '/' from absolute path names
# file: mnt/rep_test/test1
glusterfs.gfid.string="b85f1ece-7d38-41c6-873d-79a4b14f99f4"
```
### 查看文件的分布式的的信息

```
# getfattr -n trusted.glusterfs.pathinfo -e text  /data/glusterfs_speech_04_v6/11085164/espnet/vivo_input2/espnet_aishell/data/train/wav.scp 
getfattr: Removing leading '/' from absolute path names
# file: data/glusterfs_speech_04_v6/11085164/espnet/vivo_input2/espnet_aishell/data/train/wav.scp
trusted.glusterfs.pathinfo="(<DISTRIBUTE:speech_v6_rep_vol-dht> (<REPLICATE:speech_v6_rep_vol-replicate-1> <POSIX(/speech_v6/data2/brick):ai-storage-center-prd-10-194-39-15.v-bj-4.vivo.lan:/speech_v6/data2/brick/11085164/espnet/vivo_input2/espnet_aishell/data/train/wav.scp> <POSIX(/speech_v6/data2/brick):ai-storage-center-prd-10-194-39-7.v-bj-4.vivo.lan:/speech_v6/data2/brick/11085164/espnet/vivo_input2/espnet_aishell/data/train/wav.scp> <POSIX(/speech_v6/data2/brick):ai-storage-center-prd-10-194-39-6.v-bj-4.vivo.lan:/speech_v6/data2/brick/11085164/espnet/vivo_input2/espnet_aishell/data/train/wav.scp>))"
```