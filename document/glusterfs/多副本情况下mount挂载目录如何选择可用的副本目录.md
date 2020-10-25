# 多副本情况下mount挂载目录如何选择可用的副本目录
| author | update |
| ------ | ------ |
| perrynzhou@gmail.com | 2020/10/22 |

## 场景
- 一个多副本（>=3)副本集群，如果一组副本对应的brick全部宕机或者磁盘损坏，这时候glusterfs mount子目录，如果恰好选择这组副本中副本恰好在磁盘损坏或者机器宕机的那个节点，这时候glusterfs就会有问题，比如数据目录找不到了。在进行挂载时候，glusterfs客户端是如何选择可用副本的，如果挂载还是选择已经宕机的副本组，那永远就挂载失败。
- 一组副本所在节点宕机，这时候进行gluster volume replace-brick 操作进行替换已经损坏的brick,然后进行gluster volume rebalance操作，如果数量量非常大，在迁移过程中进行mont子目录，如果选择的子目录是正在同步的birck,同时brick的子目录还没有从其他的birck同步过来，这时候挂载会失败

## volume信息
```
[root@CentOS73 /rep_vol/data1/brick/public]$ gluster volume info
 
Volume Name: rep_vol
Type: Replicate
Volume ID: 2f54d945-3c83-494b-89de-6cea2ef3dd7d
Status: Started
Snapshot Count: 0
Number of Bricks: 1 x 3 = 3
Transport-type: tcp
Bricks:
Brick1: 10.211.55.9:/rep_vol/data1/brick
Brick2: 10.211.55.10:/rep_vol/data1/brick
Brick3: 10.211.55.11:/rep_vol/data1/brick
Options Reconfigured:
cluster.read-hash-mode: 0
storage.fips-mode-rchecksum: on
transport.address-family: inet
nfs.disable: on
performance.client-io-threads: off
[root@CentOS73 /rep_vol/data1/brick/public]$ gluster volume status
Status of volume: rep_vol
Gluster process                             TCP Port  RDMA Port  Online  Pid
------------------------------------------------------------------------------
Brick 10.211.55.9:/rep_vol/data1/brick      49152     0          Y       2573 
Brick 10.211.55.10:/rep_vol/data1/brick     49152     0          Y       2344 
Brick 10.211.55.11:/rep_vol/data1/brick     49152     0          Y       2332 
Self-heal Daemon on localhost               N/A       N/A        Y       2349 
Self-heal Daemon on centos71.shared         N/A       N/A        Y       2590 
Self-heal Daemon on CentOS72                N/A       N/A        Y       2361 
 
Task Status of Volume rep_vol
------------------------------------------------------------------------------
There are no active volume tasks
 
```
## 现象模拟
- 模式其中brick1中的数据正在同步，但是没有同步完成
```
//模拟数据没有同步
[root@CentOS71 /rep_vol/data1/brick/public]$ pwd
/rep_vol/data1/brick/public
[root@CentOS71 /rep_vol/data1/brick/public]$ ls -l
total 0
drwxr-xr-x. 2 root root 37 Oct 22 08:24 1024
drwxr-xr-x. 2 root root 22 Oct 22 08:24 1998
drwxr-xr-x. 2 root root 22 Oct 22 08:24 2020
drwxr-xr-x. 2 root root  6 Oct 22 08:24 209
[root@CentOS71 /rep_vol/data1/brick/public]$ rm -rf 2020
[root@CentOS71 /rep_vol/data1/brick/public]$ ls -l
total 0
drwxr-xr-x. 2 root root 37 Oct 22 08:24 1024
drwxr-xr-x. 2 root root 22 Oct 22 08:24 1998
drwxr-xr-x. 2 root root  6 Oct 22 08:24 209
```
- 进行挂载，提示挂载失
```
[root@CentOS74 ~]$ mount -t glusterfs 10.211.55.9:rep_vol/public/2020  /mnt/public/2020/
WARNING: getfattr not found, certain checks will be skipped..
```
- mount日志
```

[2020-10-22 00:30:27.712019] W [MSGID: 114043] [client-handshake.c:727:client_setvolume_cbk] 0-rep_vol-client-0: failed to set the volume [{errno=2}, {error=No such file or directory}] 
[2020-10-22 00:30:27.712047] E [MSGID: 114044] [client-handshake.c:757:client_setvolume_cbk] 0-rep_vol-client-0: SETVOLUME on remote-host failed [{remote-error=subdirectory for mount "/public/2020" is not found}, {errno=2}, {error=No such file or directory}] 
[2020-10-22 00:30:27.712056] I [MSGID: 114049] [client-handshake.c:865:client_setvolume_cbk] 0-rep_vol-client-0: sending AUTH_FAILED event [] 
[2020-10-22 00:30:27.712076] E [fuse-bridge.c:6495:notify] 0-fuse: Server authenication failed. Shutting down.
[2020-10-22 00:30:27.712112] I [fuse-bridge.c:7074:fini] 0-fuse: Unmounting '/mnt/public/2020'.
[2020-10-22 00:30:27.713960] I [MSGID: 114057] [client-handshake.c:1128:select_server_supported_programs] 0-rep_vol-client-1: Using Program [{Program-name=GlusterFS 4.x v1}, {Num=1298437}, {Version=400}] 
[2020-10-22 00:30:27.714050] I [rpc-clnt.c:1967:rpc_clnt_reconfig] 0-rep_vol-client-2: changing port to 49152 (from 0)
[2020-10-22 00:30:27.714065] I [socket.c:849:__socket_shutdown] 0-rep_vol-client-2: intentional socket shutdown(13)
[2020-10-22 00:30:27.716317] I [MSGID: 114046] [client-handshake.c:857:client_setvolume_cbk] 0-rep_vol-client-1: Connected, attached to remote volume [{conn-name=rep_vol-client-1}, {remote_subvol=/rep_vol/data1/brick}] 
[2020-10-22 00:30:27.716675] I [fuse-bridge.c:7079:fini] 0-fuse: Closing fuse connection to '/mnt/public/2020'.
[2020-10-22 00:30:27.716794] I [MSGID: 108005] [afr-common.c:5995:__afr_handle_child_up_event] 0-rep_vol-replicate-0: Subvolume 'rep_vol-client-1' came back up; going online. 
[2020-10-22 00:30:27.717027] I [MSGID: 114057] [client-handshake.c:1128:select_server_supported_programs] 0-rep_vol-client-2: Using Program [{Program-name=GlusterFS 4.x v1}, {Num=1298437}, {Version=400}] 
[2020-10-22 00:30:27.717053] W [glusterfsd.c:1439:cleanup_and_exit] (-->/lib64/libpthread.so.0(+0x7ea5) [0x7f753522bea5] -->/usr/local/sbin/glusterfs(glusterfs_sigwaiter+0xe4) [0x40b799] -->/usr/local/sbin/glusterfs(cleanup_and_exit+0x77) [0x409a2f] ) 0-: received signum (15), shutting down 
[2020-10-22 00:30:27.717139] I [timer.c:86:gf_timer_call_cancel] (-->/usr/local/lib/libgfrpc.so.0(+0x1904e) [0x7f75363d104e] -->/usr/local/lib/libgfrpc.so.0(+0x18c1c) [0x7f75363d0c1c] -->/usr/local/lib/libglusterfs.so.0(gf_timer_call_cancel+0xcd) [0x7f753663bc80] ) 0-timer: ctx cleanup started 
[2020-10-22 00:30:27.717163] E [timer.c:34:gf_timer_call_after] (-->/usr/local/lib/libgfrpc.so.0(+0x19109) [0x7f75363d1109] -->/usr/local/lib/libgfrpc.so.0(+0x18b3b) [0x7f75363d0b3b] -->/usr/local/lib/libglusterfs.so.0(gf_timer_call_after+0x9e) [0x7f753663b9ff] ) 0-timer: Either ctx is NULL or ctx cleanup started [Invalid argument]
[2020-10-22 00:30:27.717186] W [rpc-clnt-ping.c:61:__rpc_clnt_rearm_ping_timer] 0-rep_vol-client-2: unable to setup ping timer
[2020-10-22 00:30:27.717194] W [rpc-clnt-ping.c:219:rpc_clnt_ping_cbk] 0-rep_vol-client-2: failed to set the ping timer

```
## 客户端调试


```
[root@CentOS74 /var/log/glusterfs]$  gdb --args /usr/local/sbin/glusterfs --process-name fuse --volfile-server=10.211.55.9 --volfile-id=rep_vol --subdir-mount=/public/2020 /mnt/public/2020
(gdb) br create_fuse_mount 
(gdb) br afr_lookup
(gdb) br afr_discover
(gdb) br afr_discover_do
(gdb) br afr_discover_cbk
(gdb) br afr_lookup_done
(gdb) br afr_read_subvol_decide
(gdb) br afr_read_subvol_select_by_policy
(gdb) br afr_hash_child
(gdb) r
Starting program: /usr/local/sbin/glusterfs --process-name fuse --volfile-server=10.211.55.9 --volfile-id=rep_vol --subdir-mount=/public/2020 /mnt/public/2020
[Thread debugging using libthread_db enabled]
Using host libthread_db library "/lib64/libthread_db.so.1".

Breakpoint 1, create_fuse_mount (ctx=0x63c030) at glusterfsd.c:557
557         int ret = 0;
596         if (cmd_args->fuse_mountopts) {
(gdb) 
606         ret = xlator_init(master);
(gdb) 
[Detaching after fork from child process 3608]
607         if (ret) {
(gdb) set follow-fork-mode child 
(gdb) set detach-on-fork off
(gdb) 
```

## 服务端调试
```
[root@CentOS71 /rep_vol/data1/brick/public]$ ps -ef|grep glusterfsd
root      2573     1  0 08:12 ?        00:00:00 /usr/local/sbin/glusterfsd -s 10.211.55.9 --volfile-id rep_vol.10.211.55.9.rep_vol-data1-brick -p /var/run/gluster/vols/rep_vol/10.211.55.9-rep_vol-data1-brick.pid -S /var/run/gluster/b36ab5023b470153.socket --brick-name /rep_vol/data1/brick -l /var/log/glusterfs/bricks/rep_vol-data1-brick.log --xlator-option *-posix.glusterd-uuid=e88c6533-d07a-4450-b1a0-173f9e94cd59 --process-name brick --brick-port 49152 --xlator-option rep_vol-server.listen-port=49152
[root@CentOS71 /rep_vol/data1/brick/public]$ gdb /usr/local/sbin/glusterfsd      
(gdb) attach 2573
(gdb) br server4_lookup_cbk
(gdb) br posix_lookup 
Breakpoint 2 at 0x7f862924edec: file posix-entry-ops.c, line 158.
(gdb) info break
Num     Type           Disp Enb Address            What
1       breakpoint     keep y   0x00007f862262daa2 in server4_lookup_cbk at server-rpc-fops_v2.c:86
2       breakpoint     keep y   0x00007f862924edec in posix_lookup at posix-entry-ops.c:158
(gdb) 
```
## 客户端Final graph:
```
volume rep_vol-client-0
    type protocol/client
    option opversion 80000
    option clnt-lk-version 1
    option volfile-checksum 0
    option volfile-key /rep_vol
    option client-version 2020.10.22
    option process-name fuse
    option process-uuid CTX_ID:572bb692-c4c7-4244-a682-34cfe26dac47-GRAPH_ID:0-PID:3133-HOST:CentOS74-PC_NAME:rep_vol-client-0-RECON_NO:-0
    option fops-version 1298437
    option ping-timeout 42
    option remote-host 10.211.55.9
    option remote-subvolume /rep_vol/data1/brick
    option transport-type socket
    option transport.address-family inet
    option transport.socket.ssl-enabled off
    option transport.tcp-user-timeout 0
    option transport.socket.keepalive-time 20
    option transport.socket.keepalive-interval 2
    option transport.socket.keepalive-count 9
    option strict-locks off
    option send-gids true
end-volume
 
volume rep_vol-client-1
    type protocol/client
    option ping-timeout 42
    option remote-host 10.211.55.10
    option remote-subvolume /rep_vol/data1/brick
    option transport-type socket
    option transport.address-family inet
    option transport.socket.ssl-enabled off
    option transport.tcp-user-timeout 0
    option transport.socket.keepalive-time 20
    option transport.socket.keepalive-interval 2
    option transport.socket.keepalive-count 9
    option strict-locks off
    option send-gids true
end-volume
 
volume rep_vol-client-2
    type protocol/client
    option ping-timeout 42
    option remote-host 10.211.55.11
    option remote-subvolume /rep_vol/data1/brick
    option transport-type socket
    option transport.address-family inet
    option transport.socket.ssl-enabled off
    option transport.tcp-user-timeout 0
    option transport.socket.keepalive-time 20
    option transport.socket.keepalive-interval 2
    option transport.socket.keepalive-count 9
    option strict-locks off
    option send-gids true
end-volume
 
volume rep_vol-replicate-0
    type cluster/replicate
    option afr-pending-xattr rep_vol-client-0,rep_vol-client-1,rep_vol-client-2
    option use-compound-fops off
    subvolumes rep_vol-client-0 rep_vol-client-1 rep_vol-client-2
end-volume
 
volume rep_vol-dht
    type cluster/distribute
    option lock-migration off
    option force-migration off
    subvolumes rep_vol-replicate-0
end-volume
 
volume rep_vol-utime
    type features/utime
    option noatime on
    subvolumes rep_vol-dht
end-volume
 
volume rep_vol-write-behind
    type performance/write-behind
    subvolumes rep_vol-utime
end-volume
 
volume rep_vol-open-behind
    type performance/open-behind
    subvolumes rep_vol-write-behind
end-volume
 
volume rep_vol-quick-read
    type performance/quick-read
    subvolumes rep_vol-open-behind
end-volume
 
volume rep_vol-md-cache
    type performance/md-cache
    subvolumes rep_vol-quick-read
end-volume
 
volume rep_vol
    type debug/io-stats
    option log-level INFO
    option threads 16
    option latency-measurement off
    option count-fop-hits off
    option global-threading off
    subvolumes rep_vol-md-cache
end-volume
 
volume meta-autoload
    type meta
    subvolumes rep_vol
end-volume
```
## 服务端的Final graph:     
```
volume rep_vol-posix
    type storage/posix
    option glusterd-uuid e88c6533-d07a-4450-b1a0-173f9e94cd59
    option directory /rep_vol/data1/brick
    option volume-id 2f54d945-3c83-494b-89de-6cea2ef3dd7d
    option fips-mode-rchecksum on
    option shared-brick-count 1
end-volume
 
volume rep_vol-trash
    type features/trash
    option trash-dir .trashcan
    option brick-path /rep_vol/data1/brick
    option trash-internal-op off
    subvolumes rep_vol-posix
end-volume
 
volume rep_vol-changelog
    type features/changelog
    option changelog-brick /rep_vol/data1/brick
    option changelog-dir /rep_vol/data1/brick/.glusterfs/changelogs
    option changelog-notification off
    option changelog-barrier-timeout 120
    subvolumes rep_vol-trash
end-volume
 
volume rep_vol-bitrot-stub
    type features/bitrot-stub
    option export /rep_vol/data1/brick
    option bitrot disable
    subvolumes rep_vol-changelog
end-volume
 
volume rep_vol-access-control
    type features/access-control
    subvolumes rep_vol-bitrot-stub
end-volume
 
volume rep_vol-locks
    type features/locks
    option enforce-mandatory-lock off
    subvolumes rep_vol-access-control
end-volume
 
volume rep_vol-worm
    type features/worm
    option worm off
    option worm-file-level off
    option worm-files-deletable on
    subvolumes rep_vol-locks
end-volume
 
volume rep_vol-read-only
    type features/read-only
    option read-only off
    subvolumes rep_vol-worm
end-volume
 
volume rep_vol-leases
    type features/leases
    option leases off
    subvolumes rep_vol-read-only
end-volume
 
volume rep_vol-upcall
    type features/upcall
    option cache-invalidation off
    subvolumes rep_vol-leases
end-volume
 
volume rep_vol-io-threads
    type performance/io-threads
    subvolumes rep_vol-upcall
end-volume
 
volume rep_vol-selinux
    type features/selinux
    option selinux on
    subvolumes rep_vol-io-threads
end-volume
 
volume rep_vol-marker
    type features/marker
    option volume-uuid 2f54d945-3c83-494b-89de-6cea2ef3dd7d
    option timestamp-file /var/lib/glusterd/vols/rep_vol/marker.tstamp
    option quota-version 0
    option xtime off
    option gsync-force-xtime off
    option quota off
    option inode-quota off
    subvolumes rep_vol-selinux
end-volume
 
volume rep_vol-barrier
    type features/barrier
    option barrier disable
    option barrier-timeout 120
    subvolumes rep_vol-marker
end-volume
 
volume rep_vol-index
    type features/index
    option index-base /rep_vol/data1/brick/.glusterfs/indices
    option xattrop-dirty-watchlist trusted.afr.dirty
    option xattrop-pending-watchlist trusted.afr.rep_vol-
    subvolumes rep_vol-barrier
end-volume
 
volume rep_vol-quota
    type features/quota
    option volume-uuid rep_vol
    option server-quota off
    option deem-statfs off
    subvolumes rep_vol-index
end-volume
 
volume /rep_vol/data1/brick
    type debug/io-stats
    option auth.addr./rep_vol/data1/brick.allow *
    option auth-path /rep_vol/data1/brick
    option auth.login.f6e53c0a-d557-47e8-a44c-79cd26675394.password 7bf6dbee-343e-4ad3-8238-05ce03ea9b5b
    option auth.login./rep_vol/data1/brick.allow f6e53c0a-d557-47e8-a44c-79cd26675394
    option unique-id /rep_vol/data1/brick
    option volume-id 2f54d945-3c83-494b-89de-6cea2ef3dd7d
    option log-level INFO
    option threads 16
    option latency-measurement off
    option count-fop-hits off
    option global-threading off
    subvolumes rep_vol-quota
end-volume
 
volume rep_vol-server
    type protocol/server
    option transport.socket.listen-port 49152
    option rpc-auth.auth-glusterfs on
    option rpc-auth.auth-unix on
    option rpc-auth.auth-null on
    option rpc-auth-allow-insecure on
    option transport-type tcp
    option transport.address-family inet
    option auth.login./rep_vol/data1/brick.allow f6e53c0a-d557-47e8-a44c-79cd26675394
    option auth.login.f6e53c0a-d557-47e8-a44c-79cd26675394.password 7bf6dbee-343e-4ad3-8238-05ce03ea9b5b
    option auth-path /rep_vol/data1/brick
    option auth.addr./rep_vol/data1/brick.allow *
    option transport.socket.keepalive 1
    option transport.socket.ssl-enabled off
    option transport.socket.keepalive-time 20
    option transport.socket.keepalive-interval 2
    option transport.socket.keepalive-count 9
    option transport.listen-backlog 1024
    subvolumes /rep_vol/data1/brick
end-volume

```