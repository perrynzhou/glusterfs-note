
### Gluster 如何限制brick的预留空间

| 作者 | 时间 |QQ技术交流群 |
| ------ | ------ |------ |
| perrynzhou@gmail.com |2020/12/01 |中国开源存储技术交流群(672152841) |

#### 为什么要要限制glusterfs brick？

- 由于glusterfs是通过客户端计算来决定去操作远程哪一个brick的数据，这个哈希计算可能会导致多个相同规格的磁盘的使用情况不一致，比如有10块盘，其中有1块盘使用量达到90%。剩余其他的磁盘使用才80%，如果这样导致使用达到99%，最后直到glusterfsd的进程crash(glusterfsd定期会写一个日期的字符串来验证glusterfsd进程对应磁盘是否健康，一旦写入发现磁盘剩余空间无法写入的时候，glusterfsd就自杀了)

##### 是否有一个比较好的规避的办法？

- 是有的，在glusterfsd运行时候可以设定storage.reserve和storage.reserve-size，前者是设定百分比，后者是设定大小。但是2个参数设定后，glusterfsd会运行一个thread每5s来操作检查一下，free(剩余空间)和storage.reserve或者storage.reserve-size设定大小比较，如果free<=storage.reserve或者free<=storage.reserve-size都会写入失败，报出“No space left on device”错误


##### storage.reserve和storage.reserve-size是否有一定的风险?

- 这里谈不上是风险，站在自己角度应该是一个bug,磁盘剩余空间检查每5s一次，上一次和这一次检测时间间隔，用户来一个非常大的文件写入，有非常大的概率会把birck写爆，然后glusterfs heal进程来检查磁盘健康，发现无法写入，glusterfsd自杀了。所以站在自己角度应该磁盘剩余空间函数posix_disk_space_check_thread_proc最好是1s一次，这样减少了brick被写满的概率


##### glusterfs 设定磁盘空间保留
```

// 按照百分比对brick进行设定
gluster volume set test storage.reserve percentage
Option: storage.reserve
Description: Percentage of disk space to be reserved. Set to 0 to disable
Option: storage.reserve-size
Description: If set, priority will be given to storage.reserve-size over storage.reserve

```

##### glusterfs磁盘剩余空间检查实现
```
int posix_spawn_disk_space_check_thread(xlator_t *xl)
{
	   ret = gf_thread_create(&priv->disk_space_check, NULL,
                               posix_disk_space_check_thread_proc, xl,
                               "posix_reserve");
}
static void *posix_disk_space_check_thread_proc(void *data) {
	  posix_disk_space_check(this);
}

void posix_disk_space_check(xlator_t *this)
{
    struct posix_private *priv = NULL;
    char *subvol_path = NULL;
    int op_ret = 0;
    double size = 0;
    double percent = 0;
    struct statvfs buf = {0};
    double totsz = 0;
    double freesz = 0;

    GF_VALIDATE_OR_GOTO("posix-helpers", this, out);
    priv = this->private;
    GF_VALIDATE_OR_GOTO(this->name, priv, out);

    subvol_path = priv->base_path;

	// 获取磁盘的总大小
    op_ret = sys_statvfs(subvol_path, &buf);

    if (op_ret == -1) {
        gf_msg(this->name, GF_LOG_ERROR, errno, P_MSG_STATVFS_FAILED,
               "statvfs failed on %s", subvol_path);
        goto out;
    }
	// 如果使用的是storage.reserve，则priv->disk_unit值就是'p'
    if (priv->disk_unit == 'p') {
        percent = priv->disk_reserve;
        totsz = (buf.f_blocks * buf.f_bsize);
        size = ((totsz * percent) / 100);
    } else {
   	// 如果使用的是storage.reserve-size，则直接是大小，单位是字节
        size = priv->disk_reserve;
    }
	// 计算磁盘剩余的空间大小
    freesz = (buf.f_bfree * buf.f_bsize);
    // 如果剩余空间小于size,则表示磁盘满了，但是不影响heal线程来检查磁盘的状态
    if (freesz <= size) {
        priv->disk_space_full = 1;
    } else {
        priv->disk_space_full = 0;
    }
out:
    return;
}

  GF_OPTION_RECONF("reserve", priv->disk_reserve, options, percent_or_size,out);

// 每个posix操作都会执行   DISK_SPACE_CHECK_AND_GOTO(frame, priv, xdata, ret, ret, unlock) 来检查 disk的空间空间，设置brick对外提供的容量
```
