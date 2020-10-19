# 多副本情况下mount挂载目录如何选择可用的副本目录
| author | update |
| ------ | ------ | 
| perrynzhou@gmail.com | 2020/09/24 | 

## 场景
- 一个多副本（>=3)副本集群，如果一组副本对应的brick全部宕机或者磁盘损坏，这时候glusterfs如果恰好选择这组副本进行读写，这时候glusterfs就会有问题，比如数据找不到了。在进行挂载时候，glusterfs客户端是如何选择可用副本的，如果挂载还是选择已经宕机的副本组，那永远就挂载失败
## 梳理的基本函数

```
SuperFastHash
afr_hash_child
afr_read_subvol_select_by_policy
afr_read_subvol_decide
afr_lookup_done
afr_lookup_metadata_heal_check
afr_discover_cbk
afr_discover_do
afr_discover
afr_lookup
```