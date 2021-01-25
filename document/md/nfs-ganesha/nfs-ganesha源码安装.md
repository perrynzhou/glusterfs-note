### nfs-ganesha源码编译安装

| 作者 | 时间 |QQ技术交流群 |
| ------ | ------ |------ |
| perrynzhou@gmail.com |2020/12/20 |中国开源存储技术交流群(672152841) |

#### nfs-ganesha源码下载

```
$ git clone  https://hub.fastgit.org/nfs-ganesha/nfs-ganesha.git
$ cd nfs-ganesha && git submodule update --init
$ git branch -r
  origin/1.5.x
  origin/HEAD -> origin/next
  origin/V2.1-stable
  origin/V2.2-stable
  origin/V2.3-stable
  origin/V2.4-stable
  origin/V2.5-stable
  origin/V2.6-stable
  origin/V2.7-stable
  origin/V2.8-stable
  origin/V3-stable
  origin/fixv3
  origin/gh-pages
  origin/next
  origin/test
$ git checkout V3-stable
```

#### nfs-ganesha安装依赖

```
// 在centos上安装
 $ yum install  git cmake autoconf libtool bison flex doxygen openssl-devel  krb5-libs krb5-devel libuuid-devel nfs-utils -y
 
 or 
 
 // 在debian/ubuntu上安装
 $ apt install uuid-dev  nfs-kernel-server git cmake autoconf libtool bison flex doxygen openssl-dev  libkrb5-3 krb5-dev
 
 // 需要安装glusterfs中gfapi
 $ yum install centos-release-gluster7.noarch -y  
 $ yum install glusterfs-fuse.x86_64 glusterfs-api.x86_64 glusterfs-libs.x86_64  glusterfs-api-devel.x86_64   glusterfs-api.x86_64 -y

 //或者源码安装glusterfs,后续需要指定glusterfs的安装路径，nfs-ganesha需要检查安装的系统包
```

#### nfs-ganesha源码编译安装

- 安装nfs-ganesha依赖

```
yum install -y centos-release-nfs-ganesha30
yum install -y avahi-libs checkpolicy cups-libs gssproxy keyutils libicu   libverto-libevent libwbclient nfs-utils policycoreutils-python-utils psmisc   python3-audit python3-libsemanage python3-policycoreutils  python3-pyyaml python3-setools quota quota-nls  samba-client-libs  samba-common samba-common-libs                             

```
- glusterfs 安装检查

```
// nfs-ganesha指定了GLUSTER_PREFIX安装目录，会检查GLUSTER_PREFIX/下面的lib、local/lib、local/lib64、include文件夹下面的关于glusterfs安装信息
// 如果是源码安装glusteerfs，默认是安装在/usr/local下面
IF(USE_FSAL_GLUSTER)
  IF(GLUSTER_PREFIX)
    set(GLUSTER_PREFIX ${GLUSTER_PREFIX} CACHE PATH "Path to Gluster installation")
    LIST(APPEND CMAKE_PREFIX_PATH "${GLUSTER_PREFIX}")
    LIST(APPEND CMAKE_LIBRARY_PATH "${GLUSTER_PREFIX}/lib")
    LIST(APPEND CMAKE_LIBRARY_PATH "${GLUSTER_PREFIX}/local/lib")
    LIST(APPEND CMAKE_LIBRARY_PATH "${GLUSTER_PREFIX}/local/lib64")
    LIST(APPEND CMAKE_REQUIRED_INCLUDES "${GLUSTER_PREFIX}/include")
  ELSE()
    set(GLUSTER_PREFIX "/usr" CACHE PATH "Path to Gluster installation")
  ENDIF()
ENDIF()
```
- 编译
```
$ git clone  https://github.com/nfs-ganesha/nfs-ganesha.git
$ cd nfs-ganesha && git submodule update --init --recursive

$ cd nfs-ganesha/src && vi CMakeLists.txt 
 // 修改如下内容，添加 类似于"set(CMAKE_C_FLAGS "${CMAKE_C_FLAGS} -O0 -ggdb3")"

 if (LINUX)
    set(PLATFORM "LINUX")
    set(CMAKE_C_FLAGS "${CMAKE_C_FLAGS} -O0 -ggdb3 -D_LARGEFILE64_SOURCE -D_FILE_OFFSET_BITS=64")
    set(CMAKE_C_FLAGS "${CMAKE_C_FLAGS} -O0 -ggdb3 -fno-strict-aliasing")
    set(OS_INCLUDE_DIR "${PROJECT_SOURCE_DIR}/include/os/linux")
    find_library(LIBDL dl)  # module loader
endif(LINUX)


$ mkdir nfs-ganhesha/src/build && cd nfs-ganesha/src/build 
// GLUSTER_PREFIX=/usr/local 指定glusterfs安装目录
# cmake -DUSE_FSAL_GLUSTER=ON -DGLUSTER_PREFIX=/usr/local ../ 
$ make -j32 && make install
$ whereis ganesha.nfsd
ganesha: /usr/bin/ganesha.nfsd /usr/lib64/ganesha /etc/ganesha /usr/local/etc/ganesha /usr/libexec/ganesha
```

#### 添加service

```
vi /usr/lib/systemd/system/nfs-ganesha.service

[Unit]
Description=NFS-Ganesha file server
Documentation=http://github.com/nfs-ganesha/nfs-ganesha/wiki
After=rpcbind.service nfs-ganesha-lock.service
Wants=rpcbind.service nfs-ganesha-lock.service
Conflicts=nfs.target

After=nfs-ganesha-config.service
Wants=nfs-ganesha-config.service

[Service]
Type=forking
# Let systemd create /run/ganesha, /var/log/ganesha and /var/lib/nfs/ganesha
# directories
RuntimeDirectory=ganesha
LogsDirectory=ganesha
StateDirectory=nfs/ganesha
EnvironmentFile=-/run/sysconfig/ganesha
ExecStart=/bin/bash -c "${NUMACTL} ${NUMAOPTS} /usr/bin/ganesha.nfsd ${OPTIONS} ${EPOCH}"
ExecReload=/bin/kill -HUP $MAINPID
ExecStop=/bin/dbus-send --system   --dest=org.ganesha.nfsd --type=method_call /org/ganesha/nfsd/admin org.ganesha.nfsd.admin.shutdown

[Install]
WantedBy=multi-user.target
Also=nfs-ganesha-lock.service
```
#### nfs-ganehsa服务启动

```
$ systemctl daemon-reload
$ systemctl start  nfs-ganesha 
$ systemctl stop nfs-ganesha 
```