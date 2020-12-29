### libfuse 3.10源码编译

| 作者 | 时间 |QQ技术交流群 |
| ------ | ------ |------ |
| perrynzhou@gmail.com |2020/12/01 |中国开源存储技术交流群(672152841) |

- 下载编译需要安装的包

```

yum install python36 -y 


wget https://github.com/mesonbuild/meson/archive/0.56.0.tar.gz && tar 0.56.0.tar.gz && zxvf cd meson-0.56.0


// 安装ninja之前需要安装re2c-2.0.3 
wget  https://github.com/skvadrik/re2c/archive/2.0.3.tar.gz && tar zxvf 2.0.3.tar.gz && cd re2c-2.0.3 && autoreconf -i -W all && ./configure && make && make install 


//开始安装ninja
wget https://github.com/ninja-build/ninja/archive/v1.10.2.tar.gz && tar zxvf v1.10.2.tar.gz && cd ninja-1.10.2
./configure.py --bootstrap
cp ninja /usr/bin/

```

- 下载libfuse源码并编译

```
wget https://github.com/libfuse/libfuse/archive/fuse-3.10.0.tar.gz && tar zxvf fuse-3.10.0.tar.gz && cd libfuse-fuse-3.10.0

[root@CentOS7 /tmp/temp/meson-0.56.0]$ python3 ./meson.py  ../../libfuse-fuse-3.10.0/
[root@CentOS7 /tmp/temp/meson-0.56.0]$ ninja install

//如果需要重新编译则需要重新解压https://github.com/mesonbuild/meson/archive/0.56.0.tar.gz



 // 为了防止编译错误需要在libfuse-fuse-3.10.0/examples/meson.build中如下内容去掉
if not platform.endswith('bsd') and platform != 'dragonfly' and add_languages('cpp', required : false)
    executable('passthrough_hp', 'passthrough_hp.cc',
               dependencies: [ thread_dep, libfuse_dep ],
               install: false)
endif

```

- 编译例子

```
[root@CentOS7 /tmp/libfuse-fuse-3.10.0/example]$ gcc -Wall passthrough_ll.c  -o passthrough_ll -lfuse3 
```