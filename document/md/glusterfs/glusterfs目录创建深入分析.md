## glusterfs目录创建深入分析

| 作者                 | 时间       | QQ技术交流群                      |
| -------------------- | ---------- | --------------------------------- |
| perrynzhou@gmail.com | 2020/12/01 | 中国开源存储技术交流群(672152841) |

### 调试卷信息

```
Volume Name: rep_vol
Type: Replicate
Volume ID: 197bdab9-8f4e-438e-8b66-f582ebcb8c1b
Status: Started
Snapshot Count: 0
Number of Bricks: 1 x 3 = 3
Transport-type: tcp
Bricks:
Brick1: 172.16.84.37:/data/rep-vol/brick
Brick2: 172.16.84.41:/data/rep-vol/brick
Brick3: 172.16.84.42:/data/rep-vol/brick
Options Reconfigured:
storage.fips-mode-rchecksum: on
transport.address-family: inet
nfs.disable: on
performance.client-io-threads: off
```
### 客户端断点函数说明

- 设置客户端gdb调试
```
$ gdb glusterfs
(gdb) set args  --acl --process-name fuse --volfile-server=172.16.84.37 --volfile-id=rep-vol  /mnt/rep_vol 
(gdb) br main
// 创建 mount/fuse xlator并且初始化
(gdb) br create_fuse_mount 
Breakpoint 1, create_fuse_mount (ctx=0x2567abaefd0f7200) at glusterfsd.c:556
606         ret = xlator_init(master);
[Detaching after fork from child process 79742]
607         if (ret) {
(gdb) set follow-fork-mode child
(gdb) set detach-on-fork off
// 初始化"mount/fuse" xlator后初始化volume信息
(gdb) br glusterfs_volumes_init
// 注册获取volume spec的函数，volume spec就是从glusterd获取当前集群的配置所有的sub volume
(gdb) br glusterfs_mgmt_init
// 注册rpc连接glusterd服务端和从glusterd服务端断开的函数处理
(gdb) br mgmt_rpc_notify
// 获取volume spec(volume描述信息),整个Final Graph信息
(gdb) br glusterfs_volfile_fetch
// 获取当个sub volume信息
(gdb) br glusterfs_volfile_fetch_one
// 反序列化从glusterd获取到的volume信息，然后依次根据volume file构建volume、预处理volume、激活volume,通过这三步的处理，volume的秒数信息就是挂载日志中的Final Graph信息
(gdb) br mgmt_getspec_cbk
// 从glusterd获取volfile后，需要处理这个volfile
(gdb) br glusterfs_process_volfp
// 根据volfile来构造volume的加载xlator信息
(gdb) br glusterfs_graph_construct
// xlator的graph的预处理
(gdb) br glusterfs_graph_prepare
// 加载和调用所有的sub volume对象的xaltor的函数init函数
(gdb) br glusterfs_graph_activate
[Detaching after fork from child process 79742]
607         if (ret) {
(gdb) set follow-fork-mode child
(gdb) set detach-on-fork off
(gdb) info break
// volume meta-autoload、volume rep_vol-quick-read、volume rep_vol-write-behind、volume rep_vol-open-behind、volume rep_vol-quick-read 在执行mkdir命令时候并没有实际的操作，都是设置为默认的default_mkdir
(gdb) br default_mkdir 
// volume posix-acl-autoload 对应的客户端mkdir命令的执行函数
(gdb) br posix_acl_mkdir 
// volume rep_vol 对应客户端执行mkdir命令的操作函数
(gdb) br io_stats_mkdir 
// volume rep_vol-md-cache 对应客户端执行mkdir命令的操作函数
(gdb) br mdc_mkdir 
// volume rep_vol-utime  对应客户端执行mkdir命令的操作函数
(gdb) br gf_utime_mkdir 
// volume rep_vol-dht 对应客户端执行mkdir命令的操作函数
(gdb) br dht_pt_mkdir 
// volume rep_vol-replicate-0  对应客户端执行mkdir命令的操作函数
(gdb) br afr_mkdir 
(gdb) br afr_mkdir_wind 
// volume rep_vol-client-0 对应client_mkdir函数
(gdb) br  client_mkdir
(gdb) br  client4_0_mkdir
```
-  调试说明:
   -  每个xlator的执行顺序通过graph图的顺序严格执行，这个行为可以在客户端中的Final Graph信息可以看出来。每个客户端执行文件相关的操作，在每个xlator都有对应的操作函数，所有的功能都是通过xlator堆叠出来的。每个xlator要么通过XXX_WIND宏、xxx_wind函数调用下一个xlator。不如一个mkdir命令，每个xlator都有一个xxx_mkdir函数，每个函数都对应这个操作需要处理的逻辑。glusterfs客户端实现是这样，glusterd/glusterfsd设计和实现也是这样.
   -  某些xlator可以缓存一些indo信息，这些inode信息是glusterfs内部定义的。如果是一个open操作，仅仅是在glusterfs客户端层面，执行每个xlator的xxx_open函数，经过层层处理，把此次open的请求的req提交给对应glusterfsd，glusterfsd也是经过层层的处理，最后返回给glusterfs客户端也是glusterfs内部定义的inode信息，紧接着执行echo "test"到glusterfs客户端挂载文件中，调用glusterfs和glusterfsd的所有的xxx_write操作，最终把处理的结果返回给glusterfs客户端，所以glusterfs的一个操作的rpc网络开销是2次，并且操作的inode的信息并不是在glusterfs客户端执行操作的，都需要通过网络请求glusterfsd的

### mkdir命令在glusterfs客户端层面都做了什么


- meta-autoload:执行的是glusterfs/libglusterfs/src/defaults.c,而这个文件是根据glusterfs/libglusterfs/src/default-tmpl.c生成的，在编译器生成的。最终调用的是default_mkdir函数，这个函数所在的xlator是meta-autoload

  - volume 信息

    ```
    volume meta-autoload
        type meta
        subvolumes posix-acl-autoload
    end-volume
    ```

  - meta-autoload的mkdir函数实现

      ```
      // 这个函数是通过glusterfs/libglusterfs/src/default-tmpl.c模板代码生成的函数，用于连接
      int32_t default_mkdir ()
      {
      	   // 这里什么都不做仅仅只调用针对frame进行一些设定，然后调用posix-acl-autoload这个xlaor的posix_acl_mkdir
             (gdb) p this->name
      		$2 = 0x2aaab80225a0 "meta-autoload"
      	   (gdb) p this->children->xlator->name
      		$3 = 0x2aaab80037e0 "posix-acl-autoload"
      	   (gdb) p this->children->xlator->fops->mkdir
      		$4 = (fop_mkdir_t) 0x2aaabc242768 <posix_acl_mkdir>
      	
      }
      ```

 - posix-acl-autoload:当mount时候添加-o acl选项时候会进入这xlator的操作，如果操作是mkdir会执行 posix_acl_mkdir_cbk这个函数
  
   - volume信息
    
       ```
       // 对应 posix_acl_mkdir 函数
       volume posix-acl-autoload
           type system/posix-acl
           subvolumes rep_vol
       end-volume
       ```
    
       
    
   - posix-acl-autoload的mkdir实现
    
       ```
       int posix_acl_mkdir_cbk()
       {
           if (op_ret != 0)
               goto unwind;
       	//获取inode中的ctx，设置文件acl的属性
           posix_acl_ctx_update(inode, this, buf, GF_FOP_MKDIR);
       
       unwind:
       	// 设置回调函数
           STACK_UNWIND_STRICT(mkdir, frame, op_ret, op_errno, inode, buf, preparent,
                               postparent, xdata);
           return 0;
       }
       ```
  
 - rep_vol:如果用户配置了diagnostics.count-fop-hits: on 和 diagnostics.latency-measurement: on，这个是针对gluster volume profile test-volume start和# gluster volume profile *VOLNAME* info
  
   - volume信息
    
       ```
       // io_stats_mkdir函数
       volume rep_vol
           type debug/io-stats
           option log-level INFO
           option threads 16
           option latency-measurement off
        option count-fop-hits off
           option global-threading off
        subvolumes rep_vol-md-cache
       end-volume
       ```
  
   
   - rep_vol的mkdir实现
     
        ```
       int io_stats_mkdir()
       {
           if (loc->path)
               frame->local = gf_strdup(loc->path);
       
           START_FOP_LATENCY(frame);
       
           STACK_WIND(frame, io_stats_mkdir_cbk, FIRST_CHILD(this),
                      FIRST_CHILD(this)->fops->mkdir, loc, mode, umask, xdata);
           return 0;
       }
       ```




### 服务端断点设置





### 客户端xlator加载图

- 客户端xlator加载是从下往上加载，第一个执行的xlator是mount/fuse，第二是执行的是meta-autoload里面关于mkdir的方法，最后一个执行的是protocol/client 中client_mkdir->client4_0_mkdir的函数
```
// 对应client_mkdir函数
volume rep_vol-client-0
    type protocol/client
    option opversion 80000
    option clnt-lk-version 1
    option volfile-checksum 0
    option volfile-key rep_vol
    option client-version 2020.12.16
    option process-name fuse
    option fops-version 1298437
    option ping-timeout 42
    option remote-host 172.16.84.37
    option remote-subvolume /data/rep-vol/brick
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
    option remote-host 172.16.84.41
    option remote-subvolume /data/rep-vol/brick
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
    option remote-host 172.16.84.42
    option remote-subvolume /data/rep-vol/brick
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
 
//对应 afr_mkdir 函数 
volume rep_vol-replicate-0
    type cluster/replicate
    option afr-pending-xattr rep_vol-client-0,rep_vol-client-1,rep_vol-client-2
    option use-compound-fops off
    subvolumes rep_vol-client-0 rep_vol-client-1 rep_vol-client-2
end-volume
 
//对应dht_pt_mkdir这个函数 
volume rep_vol-dht
    type cluster/distribute
    option lock-migration off
    option force-migration off
    subvolumes rep_vol-replicate-0
end-volume

//对应gf_utime_mkdir这个函数, STACK_WIND中下一个xlator执行的函数是next_xl_fn(dht_pt_mkdir)，frame->root->op得到索引， next_xl_fn=get_the_pt_fop(&this->children->xlator->pass_through_fops->stat,frame->root->op)
$56 = (void *) 0x2aaab71920c1 <dht_pt_mkdir>
volume rep_vol-utime
    type features/utime
    option noatime on
    subvolumes rep_vol-dht
end-volume
 
这个volume中没有mkdir实现所以在default_mkdir函数
volume rep_vol-write-behind
    type performance/write-behind
    subvolumes rep_vol-utime
end-volume

// 这个volume中没有mkdir实现所以在default_mkdir函数 
volume rep_vol-open-behind
    type performance/open-behind
    subvolumes rep_vol-write-behind
end-volume
 
// 这个volume中没有mkdir实现所以在default_mkdir函数
volume rep_vol-quick-read
    type performance/quick-read
    subvolumes rep_vol-open-behind
end-volume
 
// 对应mdc_mkdir 函数
volume rep_vol-md-cache
    type performance/md-cache
    option cache-posix-acl true
    subvolumes rep_vol-quick-read
end-volume
 
// io_stats_mkdir函数
volume rep_vol
    type debug/io-stats
    option log-level INFO
    option threads 16
    option latency-measurement off
    option count-fop-hits off
    option global-threading off
    subvolumes rep_vol-md-cache
end-volume

// 对应 posix_acl_mkdir 函数
volume posix-acl-autoload
    type system/posix-acl
    subvolumes rep_vol
end-volume
 
// 对应default_mkdir 这个函数 
volume meta-autoload
    type meta
    subvolumes posix-acl-autoload
end-volume
 
```

### 服务端xlator加载图

- 其中一个brick的信息，加载思路和client一样，加载顺序从下往上加载，每个请求处理，第一经过protocol/server这个xlaor的server4_0_mkdir函数，最后一个执行的是storage/posix 这个xlator的posix_mkdir函数
```
volume rep_vol-posix
    type storage/posix
    option glusterd-uuid 67e6227b-ad22-4092-99f4-bde54f3285d4
    option directory /data/rep_vol/brick
    option volume-id 55d9aec2-df92-4d2c-85c5-42a4ff152d54
    option fips-mode-rchecksum on
    option shared-brick-count 1
end-volume

volume rep_vol-trash
    type features/trash
    option trash-dir .trashcan
    option brick-path /data/rep_vol/brick
    option trash-internal-op off
    subvolumes rep_vol-posix
end-volume

volume rep_vol-changelog
    type features/changelog
    option changelog-brick /data/rep_vol/brick
    option changelog-dir /data/rep_vol/brick/.glusterfs/changelogs
    option changelog-notification off
    option changelog-barrier-timeout 120
    subvolumes rep_vol-trash
end-volume

volume rep_vol-bitrot-stub
    type features/bitrot-stub
    option export /data/rep_vol/brick
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
    option volume-uuid 55d9aec2-df92-4d2c-85c5-42a4ff152d54
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
    option index-base /data/rep_vol/brick/.glusterfs/indices
    option xattrop-dirty-watchlist trusted.afr.dirty
    option xattrop-pending-watchlist trusted.afr.rep_vol-
    subvolumes rep_vol-barrier
end-volume

volume rep_vol-quota
    type features/quota
    option volume-uuid rep_vol
    subvolumes rep_vol-index
end-volume
    type debug/io-stats
    option auth.addr./data/rep_vol/brick.allow *
    option auth-path /data/rep_vol/brick
    option unique-id /data/rep_vol/brick
    option volume-id 55d9aec2-df92-4d2c-85c5-42a4ff152d54
    option latency-measurement off

volume rep_vol-server
    type protocol/server
    option transport.socket.listen-port 49152
    option rpc-auth.auth-glusterfs on
    option rpc-auth.auth-unix on
    option rpc-auth.auth-null on
    option rpc-auth-allow-insecure on
    option transport-type tcp
    option transport.address-family inet
    option auth.login./data/rep_vol/brick.allow f3f51d99-7752-4de2-b7b0-8bdff0969cb5
    option auth.login.f3f51d99-7752-4de2-b7b0-8bdff0969cb5.password fafbe12d-1065-49f9-9009-e8a69d267d7a
    option auth-path /data/rep_vol/brick
    option auth.addr./data/rep_vol/brick.allow *
    option transport.socket.keepalive 1
    option transport.socket.ssl-enabled off
    option transport.socket.keepalive-time 20
    option transport.socket.keepalive-interval 2
    option transport.socket.keepalive-count 9
    option transport.listen-backlog 1024
    subvolumes /data/rep_vol/brick
end-volume
```

### info break

```
br create_fuse_mount 
br default_mkdir 
br posix_acl_mkdir 
br io_stats_mkdir 
br mdc_mkdir
br default_mkdir 
br dht_mkdir
br afr_mkdir
br dht_pt_mkdir
br afr_mkdir_wind 
br client_mkdir
```
