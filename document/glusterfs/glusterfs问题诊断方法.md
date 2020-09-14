

# glusterfs问题诊断方法
### perf查看gluterfs相关进程函数

```
// 列举出当前可以采集的指标集合
perf list

//采集进程112547 CPU 时间消耗分析
perf record -e cpu-clock -g -p 112547

//分析采集到的数据
perf report -i perf.data
```
###  glusterfs客户端进程的statedump

- perf查看glusterfs函数消耗
  
- 生成statedump信息
```
//针对glusterfd/glusterfsd/glusterfs进程启动一个statedump
kill -SIGUSR1 {glusterd/glusterfsd/glusterfs-process-pod}

//statedump保存路径
/var/run/gluster/
```