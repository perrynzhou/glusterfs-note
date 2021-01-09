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
```


#### nfs-ganesha源码编译安装

```
$ cd nfs-ganesha/src && vi CMakeLists.txt 
 // 添加如下内容
 set(CMAKE_C_FLAGS "${CMAKE_C_FLAGS} -O0 -ggdb")
$ cd nfs-ganesha/src/build 
$ cmake -DUSE_FSAL_GLUSTER=ON ../
$ make -j32 && make install
$ whereis ganesha.nfsd
ganesha: /usr/bin/ganesha.nfsd /usr/lib64/ganesha /etc/ganesha /usr/local/etc/ganesha /usr/libexec/ganesha
```

#### nfs-ganehsa服务启动

```
$ systemctl daemon-reload
$ systemctl start  nfs-ganesha 
$ systemctl stop nfs-ganesha 
```