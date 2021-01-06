

##  根据Glusterfs gfid定位到文件具体路径

| 作者             | 时间       | QQ技术交流群                      |
| ---------------- | ---------- | --------------------------------- |
| 357884202@qq.com | 2020/12/01 | 中国开源存储技术交流群(672152841) |

- glusterfs的具体数据是存在 brick/.glusterfs 目录下的

- 具体文件形式是通过gfid 来命名的

```
  例如：gfid为： 0bb2b1e2-bc73-4d88-885b-b6c3884666fc 文件。   他在./glusterfs 的目录为:  0b/b2 
  
  1~2 为一级目录，3~4位为2级目录
  如果： 0bb2b1e2-bc73-4d88-885b-b6c3884666fc 映射的是目录，则它是一个指向这个目录真实位置的符号链接。
  
  如果：0bb2b1e2-bc73-4d88-885b-b6c3884666fc 映射的是文件，则它是一个指向这个文件真实位置的硬链接。
```

- 找到gfid 对应的文件

 ```shell
#查找到inode
ls -i 0bb2b1e2-bc73-4d88-885b-b6c3884666fc
8594219085 0bb2b1e2-bc73-4d88-885b-b6c3884666fc

#先到brick目录，通过inode 找到对应的路径
cd /data/brick/ 
find ./ -inum "8594219085" ! -path \*.glusterfs/\*
#执行结果为：
./11103398/input/cat/aishell_vgglstm_input_1850w/data/all_ark/tr.ark
 ```

- 通过文件大小查找文件

```
  find . -type f -size +100G 
```

- 通过客户端文件名，找到对应的服务端存储位置方法

```shell
#xxx是客户端挂载后的文件
getfattr -n trusted.glusterfs.pathinfo -e text  xxx/xxx

getfattr -n trusted.glusterfs.pathinfo -e text  /data/glusterfs_cv/11070574/dev.yaml
getfattr: Removing leading '/' from absolute path names
# file: data/glusterfs_cv/110/dev.yaml
trusted.glusterfs.pathinfo="(<DISTRIBUTE:test_rep_vol-dht> (<REPLICATE:test_rep_vol-replicate-0> <POSIX(/hz_cv_vol/data1/brick):10-193-103-35:/hz_cv_vol/data1/brick/110/dev.yaml>
<POSIX(/hz_cv_vol/data1/brick):10-193-103-28:/hz_cv_vol/data1/brick/110/dev.yaml> 
<POSIX(/hz_cv_vol/data1/brick):10-193-103-32:/hz_cv_vol/data1/brick/110/dev.yaml>))"

#表明了位置在/hz_cv_vol/data1/brick 这个brick上
```


