## glusterfs io-thread工作方式

| 作者 | 时间 |QQ技术交流群 |
| ------ | ------ |------ |
| perrynzhou@gmail.com |2020/12/01 |中国开源存储技术交流群(672152841) |

### io-thread介绍
- io-thread 是以translator的方式运行在glusterfsd的进程中，用于动态调整glusterfsd的IO线线程，线程名称是以glfs_iotwr名称开始
- 调整glusterfs的glfs_iotwr线程的参数是performance.io-thread-count。

### io-thread translator核心数据结构
- io-thread中定义了4中类型的操作优先级，每个操作都会对应
```
typedef enum {
    GF_FOP_PRI_UNSPEC = -1, /* Priority not specified */
    GF_FOP_PRI_HI = 0,      /* low latency */
    GF_FOP_PRI_NORMAL,      /* normal */
    GF_FOP_PRI_LO,          /* bulk */
    GF_FOP_PRI_LEAST,       /* least */
    GF_FOP_PRI_MAX,         /* Highest */
} gf_fop_pri_t;
```
- iot_conf 是io-thread translator的核心数据结构，线程的动态的创建和销毁是按照这个配置文件动态来操作的
```
struct iot_conf {
    pthread_mutex_t mutex;
    pthread_cond_t cond;

    int32_t max_count;  /* configured maximum */
    int32_t curr_count; /* actual number of threads running */
    int32_t sleep_count;

    int32_t idle_time; /* in seconds */

    struct list_head clients[GF_FOP_PRI_MAX];
    /*
     * It turns out that there are several ways a frame can get to us
     * without having an associated client (server_first_lookup was the
     * first one I hit).  Instead of trying to update all such callers,
     * we use this to queue them.
     */
    iot_client_ctx_t no_client[GF_FOP_PRI_MAX];

    int32_t ac_iot_limit[GF_FOP_PRI_MAX];
    int32_t ac_iot_count[GF_FOP_PRI_MAX];
    int queue_sizes[GF_FOP_PRI_MAX];
    int32_t queue_size;
    gf_atomic_t stub_cnt;
    pthread_attr_t w_attr;
    gf_boolean_t least_priority; /*Enable/Disable least-priority */

    xlator_t *this;
    size_t stack_size;
    gf_boolean_t down; /*PARENT_DOWN event is notified*/
    gf_boolean_t mutex_inited;
    gf_boolean_t cond_inited;

    int32_t watchdog_secs;
    gf_boolean_t watchdog_running;
    pthread_t watchdog_thread;
    gf_boolean_t queue_marked[GF_FOP_PRI_MAX];
    gf_boolean_t cleanup_disconnected_reqs;
};
```

### 线程的扩缩容

```
void *iot_worker(void *data)
{
    iot_conf_t *conf = NULL;
    xlator_t *this = NULL;
    call_stub_t *stub = NULL;
    struct timespec sleep_till = {
        0,
    };
    int ret = 0;
    int pri = -1;
    gf_boolean_t bye = _gf_false;

    conf = data;
    this = conf->this;
    THIS = this;

    for (;;) {
        pthread_mutex_lock(&conf->mutex);
        {
            if (pri != -1) {
                conf->ac_iot_count[pri]--;
                pri = -1;
            }
            while (conf->queue_size == 0) {
                if (conf->down) {
                    bye = _gf_true; /*Avoid sleep*/
                    break;
                }

                clock_gettime(CLOCK_REALTIME_COARSE, &sleep_till);
                sleep_till.tv_sec += conf->idle_time;

                conf->sleep_count++;
                ret = pthread_cond_timedwait(&conf->cond, &conf->mutex,
                                             &sleep_till);
                conf->sleep_count--;

                if (conf->down || ret == ETIMEDOUT) {
                    bye = _gf_true;
                    break;
                }
            }

            if (bye) {
                if (conf->down || conf->curr_count > IOT_MIN_THREADS) {
                    conf->curr_count--;
                    if (conf->curr_count == 0)
                        pthread_cond_broadcast(&conf->cond);
                    gf_msg_debug(conf->this->name, 0,
                                 "terminated. "
                                 "conf->curr_count=%d",
                                 conf->curr_count);
                } else {
                    bye = _gf_false;
                }
            }

            if (!bye)
                stub = __iot_dequeue(conf, &pri);
        }
        pthread_mutex_unlock(&conf->mutex);

        if (stub) { /* guard against spurious wakeups */
            if (stub->poison) {
                gf_log(this->name, GF_LOG_INFO, "Dropping poisoned request %p.",
                       stub);
                call_stub_destroy(stub);
            } else {
                call_resume(stub);
            }
            GF_ATOMIC_DEC(conf->stub_cnt);
        }
        stub = NULL;

        if (bye)
            break;
    }

    return NULL;
}

int do_iot_schedule(iot_conf_t *conf, call_stub_t *stub, int pri)
{
    int ret = 0;

    pthread_mutex_lock(&conf->mutex);
    {
        __iot_enqueue(conf, stub, pri);

        pthread_cond_signal(&conf->cond);

        ret = __iot_workers_scale(conf);
    }
    pthread_mutex_unlock(&conf->mutex);

    return ret;
}
int
iot_schedule(call_frame_t *frame, xlator_t *this, call_stub_t *stub)
{
    int ret = -1;
    gf_fop_pri_t pri = GF_FOP_PRI_MAX - 1;
    iot_conf_t *conf = this->private;

    if ((frame->root->pid < GF_CLIENT_PID_MAX) &&
        (frame->root->pid != GF_CLIENT_PID_NO_ROOT_SQUASH) &&
        conf->least_priority) {
        pri = GF_FOP_PRI_LEAST;
        goto out;
    }

    switch (stub->fop) {
        case GF_FOP_OPEN:
        case GF_FOP_STAT:
        case GF_FOP_FSTAT:
        case GF_FOP_LOOKUP:
        case GF_FOP_ACCESS:
        case GF_FOP_READLINK:
        case GF_FOP_OPENDIR:
        case GF_FOP_STATFS:
        case GF_FOP_READDIR:
        case GF_FOP_READDIRP:
        case GF_FOP_GETACTIVELK:
        case GF_FOP_SETACTIVELK:
        case GF_FOP_ICREATE:
        case GF_FOP_NAMELINK:
            pri = GF_FOP_PRI_HI;
            break;

        case GF_FOP_CREATE:
        case GF_FOP_FLUSH:
        case GF_FOP_LK:
        case GF_FOP_INODELK:
        case GF_FOP_FINODELK:
        case GF_FOP_ENTRYLK:
        case GF_FOP_FENTRYLK:
        case GF_FOP_LEASE:
        case GF_FOP_UNLINK:
        case GF_FOP_SETATTR:
        case GF_FOP_FSETATTR:
        case GF_FOP_MKNOD:
        case GF_FOP_MKDIR:
        case GF_FOP_RMDIR:
        case GF_FOP_SYMLINK:
        case GF_FOP_RENAME:
        case GF_FOP_LINK:
        case GF_FOP_SETXATTR:
        case GF_FOP_GETXATTR:
        case GF_FOP_FGETXATTR:
        case GF_FOP_FSETXATTR:
        case GF_FOP_REMOVEXATTR:
        case GF_FOP_FREMOVEXATTR:
        case GF_FOP_PUT:
            pri = GF_FOP_PRI_NORMAL;
            break;

        case GF_FOP_READ:
        case GF_FOP_WRITE:
        case GF_FOP_FSYNC:
        case GF_FOP_TRUNCATE:
        case GF_FOP_FTRUNCATE:
        case GF_FOP_FSYNCDIR:
        case GF_FOP_XATTROP:
        case GF_FOP_FXATTROP:
        case GF_FOP_RCHECKSUM:
        case GF_FOP_FALLOCATE:
        case GF_FOP_DISCARD:
        case GF_FOP_ZEROFILL:
        case GF_FOP_SEEK:
            pri = GF_FOP_PRI_LO;
            break;

        case GF_FOP_FORGET:
        case GF_FOP_RELEASE:
        case GF_FOP_RELEASEDIR:
        case GF_FOP_GETSPEC:
            break;
        case GF_FOP_IPC:
        default:
            return -EINVAL;
    }
out:
    gf_msg_debug(this->name, 0, "%s scheduled as %s priority fop",
                 gf_fop_list[stub->fop], iot_get_pri_meaning(pri));
    if (this->private)
        ret = do_iot_schedule(this->private, stub, pri);
    return ret;
}

int iot_create(call_frame_t *frame, xlator_t *this, loc_t *loc, int32_t flags,
           mode_t mode, mode_t umask, fd_t *fd, dict_t *xdata)
{
    IOT_FOP(create, frame, this, loc, flags, mode, umask, fd, xdata);
    return 0;
}
```

