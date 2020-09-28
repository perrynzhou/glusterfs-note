## glusterfs fuse实现

| author | update |
| ------ | ------ |
| perrynzhou@gmail.com | 2020/09/24 |

### glusterfs-fuse模块请求流转

- fuse是一个简单的C/S协议，由两部分组成，一部分是linux kernel，另外一部分是文件系统的daemon;在glusterfs中linux kernel是作为客户端(mount端)，glusterfs daemon作为服务端从/dev/fuse写入和读取数据，发送到远端的glusterfsd提供文件操作的相关服务
- fuse读取数据和写入数据的标准头部结构
```
//glusterfs damon读取消息是从fuse_in_header开始
struct fuse_in_header {
	//数据的总大小，包括了fuse_in_header大小
	uint32_t	len;
	//操作类型
	uint32_t	opcode;
	//请求的唯一标识
	uint64_t	unique;
	//文件系统对象的唯一标识
	uint64_t	nodeid;
	//发起请求的进程uid
	uint32_t	uid;
	//发起请求的进程gid
	uint32_t	gid;
	//发起请求的进程pid
	uint32_t	pid;
	uint32_t	padding;
    };

//glusterfs处理请求结束后，处理请求回写，是从fuse_out_header开始的
struct fuse_out_header {
	//写入文件描述符的总大小
	uint32_t	len;
	//是否遇到错误，0值表示没有遇到错误
	int32_t		error;
	//对应请求的唯一标识
	uint64_t	unique;
};
```
- notify_kernel_loop
- fuse_thread_proc:该函数循环往复从/dev/fuse读取数据，然后转发给glusterfs中的mount/fuse中的函数处理
```
static void *fuse_thread_proc(void *data)
{
    for (;;) {
        sys_readv(priv->fd, iov_in, 2);
        //msg中保存的写入的数据
        void *msg = iov_in[1].iov_base;
        //此次写操作的header
        finh = (fuse_in_header_t *)iov_in[0].iov_base;
        fuse_async_t *fasync= iov_in[0].iov_base + iov_in[0].iov_len;
        fasync->finh = finh;
        fasync->msg = msg;
        fasync->iobuf = iobuf;
        //this是当前glusterfs mount/fuse的xlator
        gf_async(&fasync->async, this, fuse_dispatch);
    }
}
```
- fuse_async_t：从/dev/fuse读到的数据，传递给fuse_async_t，然后转发给对应glusterfs mount/fuse对应的OP函数
```
typedef struct _fuse_async {
    //io buffer结构
    struct iobuf *iobuf;
    //此次操作OP的header
    fuse_in_header_t *finh;
    //写入的数据
    void *msg;
    gf_async_t async;
} fuse_async_t;
```