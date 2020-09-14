## glusterfs客户端进程的statedump

```
//针对glusterfd/glusterfsd/glusterfs进程启动一个statedump
kill -SIGUSR1 {glusterd/glusterfsd/glusterfs-process-pod}

//statedump保存路径
/var/run/gluster/

```