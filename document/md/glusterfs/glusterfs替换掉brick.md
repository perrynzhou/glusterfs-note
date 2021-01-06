## glusterfs替换掉brick


| 作者             | 时间       | QQ技术交流群                      |
| ---------------- | ---------- | --------------------------------- |
| 357884202@qq.com | 2020/12/01 | 中国开源存储技术交流群(672152841) |

### 分布式卷替换方法

```shell
#找一块和需要替换的磁盘空间大小一致的新磁盘
gluster volume info test_vol
Volume Name: test_vol
Type: Distribute
Volume ID: 638ed2b9-928e-4d7c-982d-be7123a823f1
Status: Started
Snapshot Count: 0
Number of Bricks: 2
Transport-type: tcp
Bricks:
Brick1: 192.168.78.12:/test_vol/data1/brick
Brick2: 192.168.78.12:/test_vol/data2/brick
Options Reconfigured:
transport.address-family: inet
nfs.disable: on

#新增一个brick
gluster volume add-brick test_vol 192.168.78.12:/test_vol/data3/brick

#删除就的brick
gluster volume remove-brick test_vol 192.168.78.12:/test_vol/data1/brick start
#查看状态
gluster volume remove-brick test_vol 192.168.78.12:/test_vol/data1/brick status
#提交
gluster volume remove-brick test_vol 192.168.78.12:/test_vol/data1/brick commit
```



### **多副本卷故障brick替换**

```shell
gluster volume info test_rep_vol

Volume Name: test_rep_vol
Type: Distributed-Replicate
Volume ID: 50b9747d-591e-4c85-9f22-c94aa851adb8
Status: Started
Snapshot Count: 0
Number of Bricks: 22 x 3 = 66
Transport-type: tcp
Bricks:
Brick1: 192.168.66.21:/test_rep_vol/data1/brick
Brick2: 192.168.66.22:/test_rep_vol/data1/brick
Brick3: 192.168.66.23:/test_rep_vol/data1/brick
Brick4: 192.168.66.21:/test_rep_vol/data2/brick
Brick5: 192.168.66.22:/test_rep_vol/data2/brick
Brick6: 192.168.66.23:/test_rep_vol/data2/brick
Brick7: 192.168.66.21:/test_rep_vol/data3/brick
Brick8: 192.168.66.22:/test_rep_vol/data3/brick
Brick9: 192.168.66.23:/test_rep_vol/data3/brick
Brick10: 192.168.66.21:/test_rep_vol/data4/brick
Brick11: 192.168.66.22:/test_rep_vol/data4/brick
Brick12: 192.168.66.23:/test_rep_vol/data4/brick
Brick13: 192.168.66.21:/test_rep_vol/data5/brick
Brick14: 192.168.66.22:/test_rep_vol/data5/brick
Brick15: 192.168.66.23:/test_rep_vol/data5/brick
Brick16: 192.168.66.21:/test_rep_vol/data6/brick
Brick17: 192.168.66.22:/test_rep_vol/data6/brick
Brick18: 192.168.66.23:/test_rep_vol/data6/brick
Brick19: 192.168.66.21:/test_rep_vol/data7/brick
Brick20: 192.168.66.22:/test_rep_vol/data7/brick
Brick21: 192.168.66.23:/test_rep_vol/data7/brick
Brick22: 192.168.66.21:/test_rep_vol/data8/brick
Brick23: 192.168.66.22:/test_rep_vol/data8/brick
Brick24: 192.168.66.23:/test_rep_vol/data8/brick
Brick25: 192.168.66.21:/test_rep_vol/data9/brick
Brick26: 192.168.66.22:/test_rep_vol/data9/brick
Brick27: 192.168.66.23:/test_rep_vol/data9/brick
Brick28: 192.168.66.21:/test_rep_vol/data11/brick
Brick29: 192.168.66.22:/test_rep_vol/data11/brick
Brick30: 192.168.66.23:/test_rep_vol/data11/brick
Brick31: 192.168.69.44:/test_rep_vol/data1/brick
Brick32: 192.168.72.1:/test_rep_vol/data1/brick
Brick33: 192.168.72.4:/test_rep_vol/data1/brick
Brick34: 192.168.69.44:/test_rep_vol/data2/brick
Brick35: 192.168.72.1:/test_rep_vol/data2/brick
Brick36: 192.168.72.4:/test_rep_vol/data2/brick
Brick37: 192.168.69.44:/test_rep_vol/data3/brick
Brick38: 192.168.72.1:/test_rep_vol/data3/brick
Brick39: 192.168.72.4:/test_rep_vol/data3/brick
Brick40: 192.168.69.44:/test_rep_vol/data4/brick
Brick41: 192.168.72.1:/test_rep_vol/data4/brick
Brick42: 192.168.72.4:/test_rep_vol/data4/brick
Brick43: 192.168.69.44:/test_rep_vol/data5/brick
Brick44: 192.168.72.1:/test_rep_vol/data5/brick
Brick45: 192.168.72.4:/test_rep_vol/data5/brick
Brick46: 192.168.69.44:/test_rep_vol/data6/brick
Brick47: 192.168.72.1:/test_rep_vol/data6/brick
Brick48: 192.168.72.4:/test_rep_vol/data6/brick
Brick49: 192.168.69.44:/test_rep_vol/data7/brick
Brick50: 192.168.72.1:/test_rep_vol/data7/brick
Brick51: 192.168.72.4:/test_rep_vol/data7/brick
Brick52: 192.168.69.44:/test_rep_vol/data8/brick
Brick53: 192.168.72.1:/test_rep_vol/data8/brick
Brick54: 192.168.72.4:/test_rep_vol/data8/brick
Brick55: 192.168.69.44:/test_rep_vol/data9/brick
Brick56: 192.168.72.1:/test_rep_vol/data9/brick
Brick57: 192.168.72.4:/test_rep_vol/data9/brick
Brick58: 192.168.69.44:/test_rep_vol/data11/brick
Brick59: 192.168.72.1:/test_rep_vol/data11/brick
Brick60: 192.168.72.4:/test_rep_vol/data11/brick
Brick61: 192.168.69.44:/test_rep_vol/data10/brick
Brick62: 192.168.72.1:/test_rep_vol/data10/brick
Brick63: 192.168.72.4:/test_rep_vol/data10/brick
Brick64: 192.168.66.21:/test_rep_vol/data10/brick
Brick65: 192.168.66.22:/test_rep_vol/data10/brick
Brick66: 192.168.66.23:/test_rep_vol/data10/brick
Options Reconfigured:
transport.address-family: inet
storage.fips-mode-rchecksum: on
nfs.disable: on
performance.client-io-threads: off

#整个卷的192.168.69.44机器故障，在此机器上的磁盘都要替换掉
#找一台和192.168.69.44规格一抹一样的机器，安装好glusterfs,准备好磁盘
#把192.168.51.147添加到集群
gluster peer probe 192.168.51.147
#替换故障brick
gluster volume replace-brick test_rep_vol 192.168.69.44:/test_rep_vol/data1/brick 192.168.51.147:/test_rep_vol/data1/brick commit force

#继续替换
gluster volume replace-brick test_rep_vol 192.168.69.44:/test_rep_vol/data2/brick 192.168.51.147:/test_rep_vol/data2/brick commit force

#通过在192.168.51.147 上执行df -h 可以看到磁盘数据正在同步
#查看执行完成,通过heal 没有不一致了就表示同步完成
gluster volume heal test_rep_vol info
#查看替换结果
gluster volume info test_rep_vol

```


