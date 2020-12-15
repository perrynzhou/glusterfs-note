

### 基于glusterfs的nfs-ganesha方案

| 作者 | 时间 |QQ技术交流群 |
| ------ | ------ |------ |
| perrynzhou@gmail.com |2020/12/01 |中国开源存储技术交流群(672152841) |

#### 背景

- glusterfs fuse是走fuse这一层，fuse一旦hang住，就需要重新挂载，在我们实践的方案中稳定性比较差。如果该方案在容器内，glusterfs进程容易编程D状态，D状态最终只能重启机器解决。
- 选择nfs-ganesha是因为它的操作结合glusterfs api,绕开了fuse层，数据操作都是走网络，在fuse层面提供了可用性和稳定性


####  Glusterfs版本

```
# gluster --version
glusterfs 7.2
Repository revision: git://git.gluster.org/glusterfs.git
```

####  卷信息

```
Volume Name: speech_v5_rep_vol
Type: Distributed-Replicate
Volume ID: 1ce29325-b71f-45b9-a5e4-2a506420de13
Status: Started
Snapshot Count: 0
Number of Bricks: 12 x 3 = 36
Transport-type: tcp
Bricks:
Brick1: 10.193.226.10:/test_v5_vol/data1/brick
Brick2: 10.193.226.11:/test_v5_vol/data1/brick
Brick3: 10.193.226.12:/test_v5_vol/data1/brick
Brick4: 10.193.226.10:/test_v5_vol/data2/brick
Brick5: 10.193.226.11:/test_v5_vol/data2/brick
Brick6: 10.193.226.12:/test_v5_vol/data2/brick
Brick7: 10.193.226.10:/test_v5_vol/data3/brick
Brick8: 10.193.226.11:/test_v5_vol/data3/brick
Brick9: 10.193.226.12:/test_v5_vol/data3/brick
Brick10: 10.193.226.10:/test_v5_vol/data4/brick
Brick11: 10.193.226.11:/test_v5_vol/data4/brick
Brick12: 10.193.226.12:/test_v5_vol/data4/brick
Brick13: 10.193.226.10:/test_v5_vol/data5/brick
Brick14: 10.193.226.11:/test_v5_vol/data5/brick
Brick15: 10.193.226.12:/test_v5_vol/data5/brick
Brick16: 10.193.226.10:/test_v5_vol/data6/brick
Brick17: 10.193.226.11:/test_v5_vol/data6/brick
Brick18: 10.193.226.12:/test_v5_vol/data6/brick
Brick19: 10.193.226.10:/test_v5_vol/data7/brick
Brick20: 10.193.226.11:/test_v5_vol/data7/brick
Brick21: 10.193.226.12:/test_v5_vol/data7/brick
Brick22: 10.193.226.10:/test_v5_vol/data8/brick
Brick23: 10.193.226.11:/test_v5_vol/data8/brick
Brick24: 10.193.226.12:/test_v5_vol/data8/brick
Brick25: 10.193.226.10:/test_v5_vol/data9/brick
Brick26: 10.193.226.11:/test_v5_vol/data9/brick
Brick27: 10.193.226.12:/test_v5_vol/data9/brick
Brick28: 10.193.226.10:/test_v5_vol/data10/brick
Brick29: 10.193.226.11:/test_v5_vol/data10/brick
Brick30: 10.193.226.12:/test_v5_vol/data10/brick
Brick31: 10.193.226.10:/test_v5_vol/data11/brick
Brick32: 10.193.226.11:/test_v5_vol/data11/brick
Brick33: 10.193.226.12:/test_v5_vol/data11/brick
Brick34: 10.193.226.10:/test_v5_vol/data12/brick
Brick35: 10.193.226.11:/test_v5_vol/data12/brick
Brick36: 10.193.226.12:/test_v5_vol/data12/brick
Options Reconfigured:
performance.readdir-ahead: on
performance.cache-size: 32GB
performance.io-thread-count: 32
server.event-threads: 32
cluster.readdir-optimize: on
performance.rda-cache-limit: 1024MB
features.shard-block-size: 256MB
features.shard: on
transport.address-family: inet
storage.fips-mode-rchecksum: on
nfs.disable: on
performance.client-io-threads: off
cluster.enable-shared-storage: disable
```
#### 安装nfs-ganesha

```
// 在客户端安装如下这些包
// 找三台后端大内存机器，每台机器必须安装如下包
yum install epel-release centos-release-nfs-ganesha -y
yum install nfs-ganesha nfs-ganesha-gluster -y
rpm -qa|grep nfs-ganesha
```

#### 配置启动nfs-ganesha

- 每个节点配置nfs-ganesha
```
//10.193.18.141,10.193.18.142,10.193.18.143都节点需要配置
vi /etc/ganesha/ganesha.conf 

NFS_CORE_PARAM {
        mount_path_pseudo = true;
        Protocols = 3,4;
}

EXPORT_DEFAULTS {
        Access_Type = RW;
}

EXPORT{
    Export_Id = 101 ;   
    Path = "/mnt/nfs";  

    FSAL {
        name = GLUSTER;
        //host等于当前glusterfs后端节点的IP 
        hostname = "10.193.226.12"; 
        //卷信息
        volume = "speech_v5_rep_vol";  
    }

    Access_type = RW;    
    Squash = No_root_squash; 
    Disable_ACL = TRUE;  
    //导出的的目录，比如10.1.1.1:/xxx,这个xxx就是speech_v5_rep_vol
    Pseudo = "/speech_v5_rep_vol";  
    Protocols = "3","4" ;  
    Transports = "UDP","TCP" ; 
    SecType = "sys";    
}

```

- 每个节点启动nfs-ganesha

```
systemctl start nfs-ganesha
systemctl status nfs-ganesha
showmount -e localhost 

tail -f /var/log/ganesha/ganesha.log
```

#### 客户端节点挂载

```
yum install –y nfs-utils

//支持挂载子目录
mount -t nfs4  10.193.18.141:/speech_v5_rep_vol  /mnt/nfs
```
