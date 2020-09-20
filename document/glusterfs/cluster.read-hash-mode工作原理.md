



### cluster.read-hash-mode工作原理

- read-hash-mode参数说明
```
[root@CentOS1 ~]$ gluster volume set help |grep cluster.read-hash-mode -A7
Option: cluster.read-hash-mode
Default Value: 1
Description: inode-read fops happen only on one of the bricks in replicate. AFR will prefer the one computed using the method specified using this option.
0 = first readable child of AFR, starting from 1st child.
1 = hash by GFID of file (all clients use same subvolume).
2 = hash by GFID of file and client PID.
3 = brick having the least outstanding read requests.
4 = brick having the least network ping latency.
```


- read-hash-mode 定义的类型

```
//在afr.h中定义了read-hash-mode的几个变量值
typedef enum {
		//对用0值
    AFR_READ_POLICY_FIRST_UP,
    //对应1
    AFR_READ_POLICY_GFID_HASH,
    //对应2
    AFR_READ_POLICY_GFID_PID_HASH,
    //对应3
    AFR_READ_POLICY_LESS_LOAD,
    //对应4
    AFR_READ_POLICY_LEAST_LATENCY,
    //对应5
    AFR_READ_POLICY_LOAD_LATENCY_HYBRID,
} afr_read_hash_mode_t;
```


- read-hash-mode 核心的实现函数
```
int afr_hash_child(afr_read_subvol_args_t *args, afr_private_t *priv,
               unsigned char *readable)
{
    uuid_t gfid_copy = {
        0,
    };
    pid_t pid;
    int child = -1;

    switch (priv->hash_mode) {
        case AFR_READ_POLICY_FIRST_UP:
            break;
        case AFR_READ_POLICY_GFID_HASH:
            gf_uuid_copy(gfid_copy, args->gfid);
            child = SuperFastHash((char *)gfid_copy, sizeof(gfid_copy)) %
                    priv->child_count;
            break;
        case AFR_READ_POLICY_GFID_PID_HASH:
            if (args->ia_type != IA_IFDIR) {
                /*
                 * Why getpid?  Because it's one of the cheapest calls
                 * available - faster than gethostname etc. - and
                 * returns a constant-length value that's sure to be
                 * shorter than a UUID. It's still very unlikely to be
                 * the same across clients, so it still provides good
                 * mixing.  We're not trying for perfection here. All we
                 * need is a low probability that multiple clients
                 * won't converge on the same subvolume.
                 */
                pid = getpid();
                memcpy(gfid_copy, &pid, sizeof(pid));
            }
            child = SuperFastHash((char *)gfid_copy, sizeof(gfid_copy)) %
                    priv->child_count;
            break;
        case AFR_READ_POLICY_LESS_LOAD:
            child = afr_least_pending_reads_child(priv, readable);
            break;
        case AFR_READ_POLICY_LEAST_LATENCY:
            child = afr_least_latency_child(priv, readable);
            break;
        case AFR_READ_POLICY_LOAD_LATENCY_HYBRID:
            child = afr_least_latency_times_pending_reads_child(priv, readable);
            break;
    }

    return child;
}

```


