## Lustre架构简介

#### lustre核心组件有哪些?

- lustre是一个分布式集群文件系统，Lustre客户端(client)、后端对象存储Object Storage Targets (ost)、Meta-data Service服务(MDS).客户端的涉及到的数据数据读写IO都是通过 ost服务进行；文件的元数据操作通过的mds进行。

- Lustre的架构视图如下:

  ![lustre-arc](G:\lustre简介\lustre-arc.JPG)

- Lustre IO交互视图

  ![lustre-io-interact](G:\lustre简介\lustre-io-interact.JPG)


#### Object Storage Targets (ost)主要作用是什么？
- 在整个lustre集群文件中，数据对象的IO服务是通过OST来提供。ost是存储客户端输入的数据。lustre的集群的namespace是通过mds服务管理，而mds用来管理整个集群的inode信息。在linux中inode可以是文件、目录、特殊设备。lustre通过mds中的inode来呈现用户就可以看到的文件，inode可以表示ost中的数据相关的元数据信息。
- Lustre设计ost的主要目的是在block申请来存储数据对象。ost分为二层，第一层network，为数据对象存储提供基本的网络服务；第二层是Object-Base Disk server、Lock server、Object-Base Disk，这一层是用于在lustre ost的服务端用来存储数据对象。

![ost-internal](G:\lustre简介\ost-internal.JPG)



#### Meta-data Service 主要做了什么？

