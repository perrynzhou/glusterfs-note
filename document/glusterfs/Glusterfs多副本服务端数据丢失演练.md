
### Glusterfs多副本服务端数据丢失演练

| 作者 | 时间 |QQ技术交流群 |
| ------ | ------ |------ |
| perrynzhou@gmail.com |2020/12/01 |中国开源存储技术交流群(672152841) |

- 三副本模式，采用read-hash-mode=2,以文件gfid和mount的pid进行哈希计算然一个副本进行文件操作

```
root@132.21.60.10 ~ $ gluster volume info
 
Volume Name: ssd_rep_vol
Type: Replicate
Volume ID: 483e6440-6ab7-4a5b-bd6d-1af22fa221c7
Status: Started
Snapshot Count: 0
Number of Bricks: 1 x 3 = 3
Transport-type: tcp
Bricks:
Brick1: 132.21.60.10:/glusterfs/rep_vol/data1/brick
Brick2: 172.25.78.12:/glusterfs/rep_vol/data1/brick
Brick3: 172.25.78.13:/glusterfs/rep_vol/data1/brick
Options Reconfigured:
cluster.read-hash-mode: 2
cluster.granular-entry-heal: enable
cluster.heal-timeout: 10
performance.client-io-threads: off
nfs.disable: on
transport.address-family: inet
storage.fips-mode-rchecksum: on
network.ping-timeout: 1010
cluster.self-heal-daemon: enable
```

- 客户端挂载信息
```
root@172.25.78.25 /mnt/ssd_rep_vol $ df -h
Filesystem                        Size  Used Avail Use% Mounted on
devtmpfs                           63G     0   63G   0% /dev
tmpfs                              63G     0   63G   0% /dev/shm
tmpfs                              63G   19M   63G   1% /run
tmpfs                              63G     0   63G   0% /sys/fs/cgroup
/dev/sda2                        1014M  154M  861M  16% /boot
tmpfs                              13G     0   13G   0% /run/user/5989
tmpfs                              13G     0   13G   0% /run/user/5001
172.25.78.19:train_vol            139T   73T   66T  53% /mnt/train_vol
172.25.78.13:ssd_rep_vol          100G   89G   12G  89% /mnt/ssd_rep_vol

//当前有几个文件和目录，这次我们删除ddd这个文件
root@172.25.78.25 /mnt/ssd_rep_vol $ ls -l
total 832
-rwxr-xr-x 1 root root 212632 Oct 20 14:23 data1
-rwxr-xr-x 1 root root 212632 Oct 20 14:26 data2
-rwxr-xr-x 1 root root 212632 Oct 20 15:39 ddd
drwxr-xr-x 2 root root     22 Oct 20 14:09 fuck
drwxr-xr-x 4 root root     87 Oct 19 18:24 public
-rwxr-xr-x 1 root root 212632 Oct 19 17:48 test_gfs
```

- 在Brick1上的删除ddd文件
```
//查看客户端挂载情况
root@132.21.60.10 /glusterfs/rep_vol/data1/brick $ df -h
Filesystem                Size  Used Avail Use% Mounted on
devtmpfs                   63G     0   63G   0% /dev
tmpfs                      63G     0   63G   0% /dev/shm
tmpfs                      63G   11M   63G   1% /run
tmpfs                      63G     0   63G   0% /sys/fs/cgroup
/dev/mapper/centos-root    50G   11G   40G  21% /
/dev/sdb1                1014M  150M  865M  15% /boot
/dev/mapper/centos-home   690G   33M  690G   1% /home
/dev/mapper/ssd_vg-lvol0  100G   88G   13G  88% /glusterfs/rep_vol/data1
/dev/mapper/ssd_vg-lvol1   50G   33M   50G   1% /glusterfs/dht_vol/data1
tmpfs                      13G     0   13G   0% /run/user/5989
tmpfs                      13G     0   13G   0% /run/user/5001
132.21.60.10:ssd_rep_vol  100G   89G   12G  89% /mnt/ssd_rep_vol
//未删除文件之前的情况
root@132.21.60.10 /glusterfs/rep_vol/data1/brick $ ls -l
total 832
-rwxr-xr-x 2 root root 212632 Oct 20 14:23 data1
-rwxr-xr-x 2 root root 212632 Oct 20 14:26 data2
-rwxr-xr-x 2 root root 212632 Oct 20 15:39 ddd
drwxr-xr-x 2 root root     22 Oct 20 14:09 fuck
drwxr-xr-x 4 root root     87 Oct 19 18:24 public
-rwxr-xr-x 2 root root 212632 Oct 19 17:48 test_gfs

//手动删除ddd和data2两个文件
root@132.21.60.10 /glusterfs/rep_vol/data1/brick $ rm -rf ddd data2
root@132.21.60.10 /glusterfs/rep_vol/data1/brick $ ls -l
total 416
-rwxr-xr-x 2 root root 212632 Oct 20 14:23 data1
drwxr-xr-x 2 root root     22 Oct 20 14:09 fuck
drwxr-xr-x 4 root root     87 Oct 19 18:24 public
-rwxr-xr-x 2 root root 212632 Oct 19 17:48 test_gfs
```

- 再次查看客户端的挂载
```
//虽然已经在brick1上进行了ddd和data2的文件删除，但是挂载客户端的选择的subvolume不是brick1.是brick2或者brick3其中一个(这个选择与read-hash-mode有关系)，因此我们可以在客户端看到ddd和data2
root@172.25.78.25 /mnt/ssd_rep_vol $ ls -l
total 832
-rwxr-xr-x 1 root root 212632 Oct 20 14:23 data1
-rwxr-xr-x 1 root root 212632 Oct 20 14:26 data2
-rwxr-xr-x 1 root root 212632 Oct 20 15:39 ddd
drwxr-xr-x 2 root root     22 Oct 20 14:09 fuck
drwxr-xr-x 4 root root     87 Oct 19 18:24 public
-rwxr-xr-x 1 root root 212632 Oct 19 17:48 test_gfs
```
- 修复brick1已经删除的文件
```
//使用这个命令触发自动修复，并且修复必须在客户端(mount)进行这个操作
root@172.25.78.25 /mnt/ssd_rep_vol $ find /mnt/ssd_rep_vol/ddd -type f -print0 |xargs -0 head -c1  > /dev/null
root@172.25.78.25 /mnt/ssd_rep_vol $ find /mnt/ssd_rep_vol/data2 -type f -print0 |xargs -0 head -c1  >/dev/null
```
- 查看brick1修复情况
```
//在brick1上查看ddd和data1，已经恢复了
root@132.21.60.10 /glusterfs/rep_vol/data1/brick $ ls -l
total 832
-rwxr-xr-x 2 root root 212632 Oct 20 14:23 data1
-rwxr-xr-x 2 root root 212632 Oct 20 14:26 data2
-rwxr-xr-x 2 root root 212632 Oct 20 15:39 ddd
drwxr-xr-x 2 root root     22 Oct 20 14:09 fuck
drwxr-xr-x 4 root root     87 Oct 19 18:24 public
-rwxr-xr-x 2 root root 212632 Oct 19 17:48 test_gfs
root@132.21.60.10 /glusterfs/rep_vol/data1/brick $ 
```