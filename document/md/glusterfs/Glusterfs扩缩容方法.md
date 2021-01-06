## Glusterfs扩缩容方法


| 作者             | 时间       | QQ技术交流群                      |
| ---------------- | ---------- | --------------------------------- |
| 357884202@qq.com | 2020/12/01 | 中国开源存储技术交流群(672152841) |

- 扩容

```shell
#查看状态
gluster peer status
#扩容（添加机器 172.25.78.11）
gluster peer probe 172.25.78.11

#给卷添加brick
#先停止卷
gluster volume stop gold_vol
#复制卷扩容
gluster volume add-brick gold_vol replica 3 172.25.78.11:/data03/brick0 172.25.78.12:/data03/brick0 172.25.78.13:/data03/brick0
#hash卷扩容
gluster volume add-brick hash_vol 172.25.78.11:/data03/brick0 172.25.78.12:/data03/brick0 172.25.78.13:/data03/brick0
#启动
gluster volume start gold_vol

#平衡卷
gluster volume rebalance gold_vol start
gluster volume rebalance gold_vol status

##只做修复链接
gluster volume rebalance gold_vol fix-layout start

#数据容量大的卷，扩容缩容数据均衡会耗时很久，影响性能。建议一次性创建好卷，不在扩容.
#单个卷容量控制在100T以内最好
```



- 缩容

```shell
#查看各机器状态
gluster peer status
#hash卷缩容
gluster volume remove-brick hash_vol 172.25.78.11:/data03/brick0 start
gluster volume remove-brick hash_vol 172.25.78.11:/data03/brick0 status
gluster volume remove-brick hash_vol 172.25.78.11:/data03/brick0 commit

#复制卷缩绒, 其中11和12 互为备份
gluster volume remove-brick hash_vol replica 2 172.25.78.11:/data03/brick0 172.25.78.12:/data03/brick0 start
gluster volume remove-brick hash_vol replica 2 172.25.78.11:/data03/brick0 172.25.78.12:/data03/brick0 status
gluster volume remove-brick hash_vol replica 2 172.25.78.11:/data03/brick0 172.25.78.12:/data03/brick0 commit

#缩容不用进行平衡，复制卷和hash卷不用进行平衡,会自动进行数据平衡
```


