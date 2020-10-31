### Debug Glusterfs

| author | update |
| ------ | ------ |
| perrynzhou@gmail.com | 2020/09/24 |

#### Download source

```shell
git clone https://github.com/gluster/glusterfs.git
```



#### Install Deps

- centos7
```
yum install autoconf automake bison cmockery2-devel dos2unix flex fuse-devel glib2-devel libacl-devel libaio-devel libattr-devel libcurl-devel libibverbs-devel librdmacm-devel libtirpc-devel libtool libxml2-devel lvm2-devel make openssl-devel pkgconfig pyliblzma python-devel python-eventlet python-netifaces libxml2-devel python-paste-deploy python-simplejson python-sphinx python-webob pyxattr readline-devel rpm-build sqlite-devel systemtap-sdt-devel tar userspace-rcu-devel openssl-devel  -y
```
- centos8

```
yum -y install gcc gcc-c++ make expat-devel autoconf automake libtool flex expat-devel bison openssl-devel libuuid-devel libacl-devel libxml2-devel libtirpc-devel gcc gcc-c++ make expat-devel autoconf automake libtool rdma-core-devel readline-devel libaio-devel python3 rpcbind

wget https://github.com/thkukuk/rpcsvc-proto/releases/download/v1.4/rpcsvc-proto-1.4.tar.gz && tar -xf rpcsvc-proto-1.4.tar.gz && cd rpcsvc-proto-1.4 &&./configure && make -j4 && make install

wget https://github.com/urcu/userspace-rcu/archive/v0.7.16.tar.gz -O userspace-rcu-0.7.16.tar.gz && tar -xf userspace-rcu-0.7.16.tar.gz  && cd userspace-rcu-0.7.16 && ./bootstrap && ./configure && make &&make install
find /usr/ -name \*liburcu-bp.so\* && echo '/usr/local/lib' > /etc/ld.so.conf.d/userspace-rcu.conf
```
- debian

```
# sudo apt-get install make automake autoconf libtool flex bison pkg-config libssl-dev libxml2-dev python-dev libaio-dev libibverbs-dev librdmacm-dev libreadline-dev liblvm2-dev libglib2.0-dev liburcu-dev libcmocka-dev libsqlite3-dev libacl1-dev uuid uuid-dev
```


#### Build With Debug

```shell
# cd glusterfs-6.5
# ./autogen.sh
#CFLAGS="-ggdb3 -O0"  ./configure  --enable-debug --enable-gnfs 
# make -j4
# make install
```



#### Create Mount Point

```shell
# vgcreate --physicalextentsize 128K gfs_test_vg /dev/sdb
# lvcreate -L 2G --name gfs_test_lv gfs_test_vg
# lvdisplay
  --- Logical volume ---
  LV Path                /dev/gfs_test_vg/gfs_test_lv
  LV Name                gfs_test_lv
  VG Name                gfs_test_vg
  LV UUID                1uw4WG-cdFi-vtHL-vyCZ-qluF-xknL-1Z8mY3
  LV Write Access        read/write
  LV Creation host, time localhost.localdomain, 2019-10-14 16:40:58 +0800
  LV Status              available
  LV Size                2.00 GiB
  Current LE             16384
  Segments               1
  Allocation             inherit
  Read ahead sectors     auto
  - currently set to     8192
# echo "/dev/gfs_test_vg/gfs_test_lv   /data/glusterfs/test_vol/brick1   xfs     defaults        0 0" | tee --append /etc/fstab

```



#### Begin Gdb

######  Debug Volume Info
- Set BreakPoint
```
#gdb gluster 
(gdb)set args volume info
(gdb)br main
```
- Core Execute Path
```
main   gluster/cli/src/cli.c:797
	cli_cmds_register gluster/cli/src/cli-cmd.c:208
		cli_cmd_volume_register  gluster/cli/src/cli_cmd_volume_register:3585
			cli_cmd_register  gluster/cli/src/registry.c:356
				cli_cmd_ingest    gluster/cli/src/registry.c:313

	cli_input_init gluster/cli/src/cli.c:862
		cli_batch    gluster/cli/input.c:22
			cli_cmd_process gluster/cli/src/cli-cmd.c:87
				cli_cmd_volume_info_cbk(state->tree.root.cbkfn) gluster/cli/src/cli-cmd-volume.c:3355
					cli_rpc_prog->proctable[GLUSTER_CLI_GET_VOLUME] {
					    //gluster/cli/src/cli-rpc-ops.c:12152
							struct rpc_clnt_procedure gluster_cli_actors[GLUSTER_CLI_MAXVALUE] = {
	    					[GLUSTER_CLI_NULL] = {"NULL", NULL},
	    					[GLUSTER_CLI_PROBE] = {"PROBE_QUERY", gf_cli_probe},
	    					[GLUSTER_CLI_DEPROBE] = {"DEPROBE_QUERY", gf_cli_deprobe},
	    					[GLUSTER_CLI_LIST_FRIENDS] = {"LIST_FRIENDS", gf_cli_list_friends},
	    					[GLUSTER_CLI_UUID_RESET] = {"UUID_RESET", gf_cli3_1_uuid_reset},
	    					[GLUSTER_CLI_UUID_GET] = {"UUID_GET", gf_cli3_1_uuid_get},
	    					[GLUSTER_CLI_CREATE_VOLUME] = {"CREATE_VOLUME", gf_cli_create_volume},
	    					[GLUSTER_CLI_DELETE_VOLUME] = {"DELETE_VOLUME", gf_cli_delete_volume},
	    					[GLUSTER_CLI_START_VOLUME] = {"START_VOLUME", gf_cli_start_volume},
	    					[GLUSTER_CLI_STOP_VOLUME] = {"STOP_VOLUME", gf_cli_stop_volume},
	    					[GLUSTER_CLI_RENAME_VOLUME] = {"RENAME_VOLUME", gf_cli_rename_volume},
	    					[GLUSTER_CLI_DEFRAG_VOLUME] = {"DEFRAG_VOLUME", gf_cli_defrag_volume},
	    					[GLUSTER_CLI_GET_VOLUME] = {"GET_VOLUME", gf_cli_get_volume},
	    					[GLUSTER_CLI_GET_NEXT_VOLUME] = {"GET_NEXT_VOLUME", gf_cli_get_next_volume},
	    					[GLUSTER_CLI_SET_VOLUME] = {"SET_VOLUME", gf_cli_set_volume},
	    					[GLUSTER_CLI_ADD_BRICK] = {"ADD_BRICK", gf_cli_add_brick},
	    					[GLUSTER_CLI_REMOVE_BRICK] = {"REMOVE_BRICK", gf_cli_remove_brick},
	    					[GLUSTER_CLI_REPLACE_BRICK] = {"REPLACE_BRICK", gf_cli_replace_brick},
	    					[GLUSTER_CLI_LOG_ROTATE] = {"LOG ROTATE", gf_cli_log_rotate},
	    					[GLUSTER_CLI_GETSPEC] = {"GETSPEC", gf_cli_getspec},
	    					[GLUSTER_CLI_PMAP_PORTBYBRICK] = {"PMAP PORTBYBRICK", gf_cli_pmap_b2p},
	    					[GLUSTER_CLI_SYNC_VOLUME] = {"SYNC_VOLUME", gf_cli_sync_volume},
	    					[GLUSTER_CLI_RESET_VOLUME] = {"RESET_VOLUME", gf_cli_reset_volume},
	    					[GLUSTER_CLI_FSM_LOG] = {"FSM_LOG", gf_cli_fsm_log},
	    					[GLUSTER_CLI_GSYNC_SET] = {"GSYNC_SET", gf_cli_gsync_set},
	    					[GLUSTER_CLI_PROFILE_VOLUME] = {"PROFILE_VOLUME", gf_cli_profile_volume},
	    					[GLUSTER_CLI_QUOTA] = {"QUOTA", gf_cli_quota},
	    					[GLUSTER_CLI_TOP_VOLUME] = {"TOP_VOLUME", gf_cli_top_volume},
	    					[GLUSTER_CLI_GETWD] = {"GETWD", gf_cli_getwd},
	    					[GLUSTER_CLI_STATUS_VOLUME] = {"STATUS_VOLUME", gf_cli_status_volume},
	    					[GLUSTER_CLI_STATUS_ALL] = {"STATUS_ALL", gf_cli_status_volume_all},
	    					[GLUSTER_CLI_MOUNT] = {"MOUNT", gf_cli_mount},
	    					[GLUSTER_CLI_UMOUNT] = {"UMOUNT", gf_cli_umount},
	    					[GLUSTER_CLI_HEAL_VOLUME] = {"HEAL_VOLUME", gf_cli_heal_volume},
	    					[GLUSTER_CLI_STATEDUMP_VOLUME] = {"STATEDUMP_VOLUME",
	    					                                  gf_cli_statedump_volume},
	    					[GLUSTER_CLI_LIST_VOLUME] = {"LIST_VOLUME", gf_cli_list_volume},
	    					[GLUSTER_CLI_CLRLOCKS_VOLUME] = {"CLEARLOCKS_VOLUME",
	    					                                 gf_cli_clearlocks_volume},
	    					[GLUSTER_CLI_COPY_FILE] = {"COPY_FILE", gf_cli_copy_file},
	    					[GLUSTER_CLI_SYS_EXEC] = {"SYS_EXEC", gf_cli_sys_exec},
	    					[GLUSTER_CLI_SNAP] = {"SNAP", gf_cli_snapshot},
	    					[GLUSTER_CLI_BARRIER_VOLUME] = {"BARRIER VOLUME", gf_cli_barrier_volume},
	    					[GLUSTER_CLI_GET_VOL_OPT] = {"GET_VOL_OPT", gf_cli_get_vol_opt},
	    					[GLUSTER_CLI_BITROT] = {"BITROT", gf_cli_bitrot},
	    					[GLUSTER_CLI_ATTACH_TIER] = {"ATTACH_TIER", gf_cli_attach_tier},
	    					[GLUSTER_CLI_TIER] = {"TIER", gf_cli_tier},
	    					[GLUSTER_CLI_GET_STATE] = {"GET_STATE", gf_cli_get_state},
	    					[GLUSTER_CLI_RESET_BRICK] = {"RESET_BRICK", gf_cli_reset_brick},
	    					[GLUSTER_CLI_REMOVE_TIER_BRICK] = {"DETACH_TIER", gf_cli_remove_tier_brick},
	    					[GLUSTER_CLI_ADD_TIER_BRICK] = {"ADD_TIER_BRICK", gf_cli_add_tier_brick}};
	
						struct rpc_clnt_program cli_prog = {
						    .progname = "Gluster CLI",
						    .prognum = GLUSTER_CLI_PROGRAM,
						    .progver = GLUSTER_CLI_VERSION,
						    .numproc = GLUSTER_CLI_MAXVALUE,
						    .proctable = gluster_cli_actors,
						};
					}
						gf_cli_get_volume   gluster/cli/src/cli-rpc-ops.c:4579
							gf_cli_get_volume_cbk  gluster/cli/src/cli-rpc-ops.c:819
```

