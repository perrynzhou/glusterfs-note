

## glusterfsd出现crash的分析和总结

- glusterfs版本

```
# gluster --version
glusterfs 7.2
```
- glusterfs volume 信息

```
# gluster volume info sharing_vol 
 
Volume Name: sharing_vol
Type: Distributed-Replicate
Volume ID: 9dbb8cdc-68b9-40f5-8c35-8eb7c80c6bed
Status: Started
Snapshot Count: 0
Number of Bricks: 12 x 3 = 36
Transport-type: tcp
Bricks:
Brick1: 130.114.10.129:/data1/brick_sharing_vol
Brick2: 130.114.10.132:/data1/brick_sharing_vol
Brick3: 130.114.10.133:/data1/brick_sharing_vol
Brick4: 130.114.10.129:/data2/brick_sharing_vol
Brick5: 130.114.10.132:/data2/brick_sharing_vol
Brick6: 130.114.10.133:/data2/brick_sharing_vol
Brick7: 130.114.10.129:/data3/brick_sharing_vol
Brick8: 130.114.10.132:/data3/brick_sharing_vol
Brick9: 130.114.10.133:/data3/brick_sharing_vol
Brick10: 130.114.10.129:/data4/brick_sharing_vol
Brick11: 130.114.10.132:/data4/brick_sharing_vol
Brick12: 130.114.10.133:/data4/brick_sharing_vol
Brick13: 130.114.10.129:/data5/brick_sharing_vol
Brick14: 130.114.10.132:/data5/brick_sharing_vol
Brick15: 130.114.10.133:/data5/brick_sharing_vol
Brick16: 130.114.10.129:/data6/brick_sharing_vol
Brick17: 130.114.10.132:/data6/brick_sharing_vol
Brick18: 130.114.10.133:/data6/brick_sharing_vol
Brick19: 130.114.10.129:/data7/brick_sharing_vol
Brick20: 130.114.10.132:/data7/brick_sharing_vol
Brick21: 130.114.10.133:/data7/brick_sharing_vol
Brick22: 130.114.10.129:/data8/brick_sharing_vol
Brick23: 130.114.10.132:/data8/brick_sharing_vol
Brick24: 130.114.10.133:/data8/brick_sharing_vol
Brick25: 130.114.10.129:/data9/brick_sharing_vol
Brick26: 130.114.10.132:/data9/brick_sharing_vol
Brick27: 130.114.10.133:/data9/brick_sharing_vol
Brick28: 130.114.10.129:/data10/brick_sharing_vol
Brick29: 130.114.10.132:/data10/brick_sharing_vol
Brick30: 130.114.10.133:/data10/brick_sharing_vol
Brick31: 130.114.10.129:/data11/brick_sharing_vol
Brick32: 130.114.10.132:/data11/brick_sharing_vol
Brick33: 130.114.10.133:/data11/brick_sharing_vol
Brick34: 130.114.10.129:/data12/brick_sharing_vol
Brick35: 130.114.10.132:/data12/brick_sharing_vol
Brick36: 130.114.10.133:/data12/brick_sharing_vol
Options Reconfigured:
storage.fips-mode-rchecksum: on
nfs.disable: on
```
- 某个卷的后端服务批量的crash,日志如下

```
  /*------------------open too many files,happens many times----------------------*/
  [2020-10-10 09:22:50.117882] W [socket.c:3126:socket_server_event_handler] 0-tcp.sharing_vol-server: accept on 11 failed (Too many open files)
  [2020-10-10 09:22:50.117890] W [socket.c:3126:socket_server_event_handler] 0-tcp.sharing_vol-server: accept on 11 failed (Too many open files)
  [2020-10-10 09:22:50.117894] W [socket.c:3126:socket_server_event_handler] 0-tcp.sharing_vol-server: accept on 11 failed (Too many open files)
  [2020-10-10 09:22:50.117897] W [socket.c:3126:socket_server_event_handler] 0-tcp.sharing_vol-server: accept on 11 failed (Too many open files)
  [2020-10-10 09:22:50.117901] W [socket.c:3126:socket_server_event_handler] 0-tcp.sharing_vol-server: accept on 11 failed (Too many open files)
  
  /*-------------------------------open file and opendir failed many times-------------*/
  [2020-10-10 09:22:56.499536] E [MSGID: 115070] [server-rpc-fops_v2.c:1502:server4_open_cbk] 0-sharing_vol-server: 5756661176: OPEN /public_data/speech_wakeup/datasets/PD2001/multi-noisy/data/panhaiquan/XVXV/20200726_A/drive_closewindow/eval/panhaiquan_XVXV_20200726_A_drive_closewindow_eval_69.wav (ab21069a-8344-4122-9f82-9a658b76f995), client: CTX_ID:d1b7e7b6-30be-4572-9509-a2c8d7d49544-GRAPH_ID:0-PID:44983-HOST:ai-vtraining-prd-10-193-85-11.v-bj-4.vivo.lan-PC_NAME:sharing_vol-client-31-RECON_NO:-0, error-xlator: sharing_vol-posix [Too many open files]
  [2020-10-10 09:22:56.499904] E [MSGID: 113039] [posix-inode-fd-ops.c:1523:posix_open] 0-sharing_vol-posix: open on gfid-handle /data11/brick_sharing_vol/.glusterfs/10/78/10784702-07b8-432f-937e-9f6a037bfa6c (path: /public_data/speech_wakeup/datasets/PD2001/multi-noisy/data/panhaiquan/XVXV/20200726_A/drive_closewindow/eval/panhaiquan_XVXV_20200726_A_drive_closewindow_eval_70.wav), flags: 0 [Too many open files]
  [2020-10-10 09:22:56.499933] E [MSGID: 115070] [server-rpc-fops_v2.c:1502:server4_open_cbk] 0-sharing_vol-server: 5756661178: OPEN /public_data/speech_wakeup/datasets/PD2001/multi-noisy/data/panhaiquan/XVXV/20200726_A/drive_closewindow/eval/panhaiquan_XVXV_20200726_A_drive_closewindow_eval_70.wav (10784702-07b8-432f-937e-9f6a037bfa6c), client: CTX_ID:d1b7e7b6-30be-4572-9509-a2c8d7d49544-GRAPH_ID:0-PID:44983-HOST:ai-vtraining-prd-10-193-85-11.v-bj-4.vivo.lan-PC_NAME:sharing_vol-client-31-RECON_NO:-0, error-xlator: sharing_vol-posix [Too many open files]
  [2020-10-10 09:22:56.547027] E [MSGID: 113015] [posix-inode-fd-ops.c:1235:posix_opendir] 0-sharing_vol-posix: opendir failed on gfid-handle: /data11/brick_sharing_vol/11101488/ly-data/imagenet/train/n04285008 (path: /11101488/ly-data/imagenet/train/n04285008) [Too many open files]
  
  
  /*-----------------------glusterfsd check heal---------------------------*/
  [2020-10-10 09:22:56.896357] E [MSGID: 113099] [posix-helpers.c:448:_posix_xattr_get_set] 0-sharing_vol-posix: Opening file /data11/brick_sharing_vol/11101479/workspace/youtube-8m/feature_extractor/download_img/414925.jpg failed [Too many open files]
  [2020-10-10 09:22:56.907259] W [MSGID: 113075] [posix-helpers.c:2111:posix_fs_health_check] 0-sharing_vol-posix: open_for_write() on /data11/brick_sharing_vol/.glusterfs/health_check returned [Too many open files]
  [2020-10-10 09:22:56.907304] M [MSGID: 113075] [posix-helpers.c:2185:posix_health_check_thread_proc] 0-sharing_vol-posix: health-check failed, going down 
  [2020-10-10 09:22:56.913949] M [MSGID: 113075] [posix-helpers.c:2203:posix_health_check_thread_proc] 0-sharing_vol-posix: still alive! -> SIGTERM 
  ```

  

- 后端glusterfsd进程的Crash的本质原因

  - 由于节点的glusterfsd的进程打开太多文件描述符，超过了操作系统的限制，首先在accept系统调用阶段出现too many files

  - 然后用户请求不断的开发文件和目录，同时也出现了这个错误

  - glusterfsd自身机制需要检查当前进程所对应的文件系统、磁盘是否是正常状态，单独启动一个thread,调用posix_health_check_thread_proc方法打开".glusterfs/health_check"文件，这个时候open这个文件时候出现了这个错误，glusterfsd的健康检查是无法区分是磁盘坏了、还是文件系统出现了问题，但是总体表现是打开这个文件失败，然后glusterfsd误认为磁盘坏了或者文件系统坏了就直接退出了。
  - 不排除当前glusterfsd在处理IO流程的时候，出现了资源泄露，目前在和官方讨论这个问题

- 解决方法
	- 需要找到too many files的本质原因，然后解决掉即可
	- glusterfsd heal check的检查机制有些粗暴，后续可以考虑更加优雅的方法来判断具体是那种情况