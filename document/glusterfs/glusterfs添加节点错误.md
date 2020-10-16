## glusterfs添加节点成功后状态不正常

1、现象:

添加节点成功，但是状态为 Accepted peer request (Connected)

172.25.78.11 集群添加172.25.78.12 成功,在此机器上显示正常

```
Number of Peers: 1

Hostname: 172.25.78.12
Uuid: ee4be557-73c6-4bc2-aabf-04616bcb2b10
State: Peer in Cluster (Connected)
```

在172.25.78.12上状态不正常

```
Number of Peers: 1

Hostname: bogon
Uuid: 2adc400b-15f4-4135-a2f8-101774c53ee5
State: Accepted peer request (Connected)
```

2、原因:

网络配置问题,没有把hostname暴露出去，导致采集到默认的bogon 

3、解决：

a.修改172.25.78.12 /var/lib/glusterd/peers/的配置文件

修改2adc400b-15f4-4135-a2f8-101774c53ee5文件的hostname改为对应的IP

```shell
uuid=2adc400b-15f4-4135-a2f8-101774c53ee5
state=3
####在此把bogon改为IP
hostname1=bogon
```

b. 重启glusterfs后解决

