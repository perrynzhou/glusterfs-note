## glusterfs 安装及卷创建使用

1、安装指定版本的glusterfs,采用以下脚本：

```shell
#!/bin/bash

basearch='$basearch'
releasever='$releasever'
str=$(cat <<EOF
[centos-gluster7]
gpgcheck = 0
mirrorlist = http://mirrorlist.centos.org?arch=$basearch&release=$releasever&repo=storage-gluster-7
name = CentOS-$releasever - Gluster 7

[centos-gluster6]
gpgcheck = 0
mirrorlist = http://mirrorlist.centos.org?arch=$basearch&release=$releasever&repo=storage-gluster-6
name = CentOS-$releasever - Gluster 6
EOF
)

echo "$str" > /etc/yum.repos.d/gluster.repo
version=7.2

rpm_packages=(
glusterfs-server
glusterfs-events
glusterfs-extra-xlators
glusterfs-geo-replication
glusterfs-libs
glusterfs-rdma
glusterfs
glusterfs-api
glusterfs-api-devel
glusterfs-cli
glusterfs-client-xlators
glusterfs-fuse
python2-gluster
)

for rpm in ${rpm_packages[@]}; do
    echo "yum install -y $rpm-$version"
    yum install -y $rpm-$version
done

systemctl enable glusterd
systemctl start glusterd
```

2、添加节点

```shell
#查看所有节点信息，显示时不包括本节点
gluster peer status 
#添加节点
gluster peer probe   NODE-NAME
#移除节点，需要提前将该节点上的brick移除
gluster peer detach  NODE-NAME  
```

3、创建卷

```shell
#hash卷创建,假如我们在172.25.71.117机器上有3块磁盘，然后用这个3块磁盘搭建一个卷
#创建brick目录
mkdir /data01/gfs/brick
mkdir /data02/gfs/brick
mkdir /data03/gfs/brick
#创建hash卷
gluster volume create test_volume 172.25.71.117:/data01/gfs/brick0 172.25.71.117:/data02/gfs/brick 172.25.71.117:/data03/gfs/brick
#启动卷
gluster volume start test_volume
#查看卷信息
gluster volume info test_volume

#创建复制卷，复制卷brick副本最好分布在不同的机器上，这样才能够保证当有机器宕机时候不会影响整个卷使用.
#创建3副本的复制卷
gluster volume create sz_cv_aimark_rep_vol replica 3 \
10.193.227.32:/sz_cv_aimark/data1/brick \
10.193.226.34:/sz_cv_aimark/data1/brick \
10.193.226.4:/sz_cv_aimark/data1/brick \
10.193.227.32:/sz_cv_aimark/data2/brick \
10.193.226.34:/sz_cv_aimark/data2/brick \
10.193.226.4:/sz_cv_aimark/data2/brick \
10.193.227.32:/sz_cv_aimark/data3/brick \
10.193.226.34:/sz_cv_aimark/data3/brick \
10.193.226.4:/sz_cv_aimark/data3/brick \
10.193.227.32:/sz_cv_aimark/data4/brick \
10.193.226.34:/sz_cv_aimark/data4/brick \
10.193.226.4:/sz_cv_aimark/data4/brick

#ec卷创建，系统推荐使用4+2的方式使用，所以一组brick最好是分布在6台机器上，不这样分布也是可以创建的
#下面的backup_szcv_ec2_vol卷分布在3台机器上
gluster volume create backup_szcv_ec2_vol disperse-data 4 redundancy 2 \
10.193.27.1:/backup_szcv_ec2/data1/brick \
10.193.27.2:/backup_szcv_ec2/data1/brick \
10.193.27.15:/backup_szcv_ec2/data1/brick \
10.193.27.1:/backup_szcv_ec2/data2/brick \
10.193.27.2:/backup_szcv_ec2/data2/brick \
10.193.27.15:/backup_szcv_ec2/data2/brick \
10.193.27.1:/backup_szcv_ec2/data3/brick \
10.193.27.2:/backup_szcv_ec2/data3/brick \
10.193.27.15:/backup_szcv_ec2/data3/brick \
10.193.27.1:/backup_szcv_ec2/data4/brick \
10.193.27.2:/backup_szcv_ec2/data4/brick \
10.193.27.15:/backup_szcv_ec2/data4/brick \
10.193.27.1:/backup_szcv_ec2/data5/brick \
10.193.27.2:/backup_szcv_ec2/data5/brick \
10.193.27.15:/backup_szcv_ec2/data5/brick \
10.193.27.1:/backup_szcv_ec2/data6/brick \
10.193.27.2:/backup_szcv_ec2/data6/brick \
10.193.27.15:/backup_szcv_ec2/data6/brick
```

4、挂载卷

```shell
#需要物理机安装客户端
#采用最初的脚步，安装以下组件
glusterfs-libs
glusterfs
glusterfs-fuse
glusterfs-client-xlators
#创建挂载目录
mkdir /data/glusterfs_sz_cv_aimark
#挂载卷
mount -t glusterfs -o acl,backup-volfile-servers=10.193.226.34:10.193.226.4 10.193.227.32:sz_cv_aimark_rep_vol /data/glusterfs_sz_cv_aimark
```

