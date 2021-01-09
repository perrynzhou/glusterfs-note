## glusterfs gfapi如何工作的


| 作者 | 时间 |QQ技术交流群 |
| ------ | ------ |------ |
| perrynzhou@gmail.com |2020/12/01 |中国开源存储技术交流群(672152841) |

### glusterfs gfapi的代码路径

```
[root@linuxzhou /home/perrynzhou/data/Source/open/glusterfs/api]$ tree ./
./
├── examples
│   ├── autogen.sh
│   ├── configure.ac
│   ├── getvolfile.py
│   ├── glfsxmp.c
│   ├── Makefile.am
│   └── README
├── Makefile.am
└── src
    ├── gfapi.aliases
    ├── gfapi.map
    ├── gfapi-messages.h
    ├── glfs.c
    ├── glfs-fops.c
    ├── glfs.h
    ├── glfs-handleops.c
    ├── glfs-handles.h
    ├── glfs-internal.h
    ├── glfs-master.c
    ├── glfs-mem-types.h
    ├── glfs-mgmt.c
    ├── glfs-resolve.c
    ├── Makefile.am
    └── README.Symbol_Versions

```

### 客户端和服务端初始化流程

- 客户端核心函数说明

```
//glusterfs api glfs_init函数的核心逻辑
GFAPI_SYMVER_PUBLIC_DEFAULT(glfs_init, 3.4.0)
int pub_glfs_init(struct glfs *fs) {
   //初始化glusterfs gfapi xlator
   ret = glfs_init_common(fs);
   //等待返回结果
   ret = glfs_init_wait(fs);
}
//gfapi xlator的创建
int glfs_init_common(struct glfs *fs) {
	//创建一个gfapi xlator并且fs->ctx->master = gfapi xlator
 	ret = create_master(fs);
 	//初始化网络IO的工作线程
    ret = gf_thread_create(&fs->poller, NULL, glfs_poller, fs, "glfspoll");
    //通过指定的volfile_server或者volfile来获取服务端的volume的spec
    ret = glfs_volumes_init(fs);
    fs->dev_id = gf_dm_hashfn(fs->volname, strlen(fs->volname));
}

int glfs_volumes_init(struct glfs *fs)
{
	//判断volfile或者volfile_server是否为空
 	if (!vol_assigned(cmd_args))
        return -1;
    if (cmd_args->volfile_server) {
    	//volfile_server不为空通过volfile_server获取volume spec
        ret = glfs_mgmt_init(fs);
        goto out;
    }
	//获取volfile信息
    fp = get_volfp(fs);
	//从volfile加载gfapi xlator,启动gfapi xlator
    ret = glfs_process_volfp(fs, fp);
out:
    return ret;
}

//获取volume的spec
int glfs_mgmt_init(struct glfs *fs)
{
	//注册网络IO事件，包括连接和断开需要做对应的事情
	rpc_clnt_register_notify(rpc, mgmt_rpc_notify, THIS);
}
static int mgmt_rpc_notify(struct rpc_clnt *rpc, void *mydata, rpc_clnt_event_t event,void *data)
{
switch (event) {
        case RPC_CLNT_DISCONNECT:
        //do somthuing
        	break;
        //连接到服务端的24007以后,执行glfs_volfile_fetch函数
        case RPC_CLNT_CONNECT:
            ret = glfs_volfile_fetch(fs);
            break;
        default:
        	break;
}

//clnt_handshake_procs这里定义了客户端和服务端握手时候，对应服务端GF_HNDSK_GETSPEC需要做的事情，GF_HNDSK_GETSPEC在服务端也会有对应的处理函数
static char *clnt_handshake_procs[GF_HNDSK_MAXVALUE] = {
    [GF_HNDSK_GETSPEC] = "GETSPEC",
};

//从glusterd获取volume的spec
int glfs_volfile_fetch(struct glfs *fs)
{
  //函数执行完毕后需要执行glfs_mgmt_getspec_cbk出处理mgmt_submit_request返回的结果
   ret = mgmt_submit_request(&req, frame, ctx, &clnt_handshake_prog,
                              GF_HNDSK_GETSPEC, glfs_mgmt_getspec_cbk,
                              (xdrproc_t)xdr_gf_getspec_req);
}







//接下来通过获取的spec来加载服务
int glfs_mgmt_getspec_cbk(struct rpc_req *req, struct iovec *iov, int count,
                      void *myframe)
{
    //当前struct glfs实例中volfile的重新配置，判断原来的volfile不可用时候，加载新的tmpfp中的voilfile内容
	gf_volfile_reconfigure(fs->oldvollen, tmpfp, fs->ctx, fs->oldvolfile);
	//从volfile加载gfapi xlator,启动gfapi xlator
	ret = glfs_process_volfp(fs, tmpfp);
}
```

- 服务端
```
// 客户端在提交 服务端GF_HNDSK_GETSPEC 对应的处理函数server_getspec
static rpcsvc_actor_t gluster_handshake_actors[GF_HNDSK_MAXVALUE] = {
    [GF_HNDSK_GETSPEC] = {"GETSPEC", server_getspec,NULL,GF_HNDSK_GETSPEC, DRC_NA, 0}
    }
    
// 客户端调用mgmt_submit_request函数后，会执行到server_getspec。
int server_getspec(rpcsvc_request_t *req)
{
    return glusterd_big_locked_handler(req, __server_getspec);
}
```

