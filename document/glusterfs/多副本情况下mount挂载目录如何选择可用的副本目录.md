# 多副本情况下mount挂载目录如何选择可用的副本目录
| author | update |
| ------ | ------ | 
| perrynzhou@gmail.com | 2020/09/24 | 

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