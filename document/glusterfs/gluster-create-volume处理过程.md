## 涉及到基本命令处理的源文件
- rpc/rpc-lib/src/protocol-common.h:定义了cli命令的一些枚举变量，包括了gluster_cli_procnum、glusterd_mgmt_procnum、glusterd_brick_procnum等
- xlators/mgmt/glusterd/glusterd-handler.c:定义gd_svc_cli_actors、gd_svc_cli_trusted_actors，来针对客户端cli请求的处理rpcsvc_program数组，glusterd这些protocol-common中定义的枚举作为这些d_svc_cli_actors、gd_svc_cli_trusted_actors下标，数组中每一个rpcsvc_actor_t对应cli处理的一个函数

## gluster如何创建volume?

- gluster使用gluster volume create 命令进行volume的创建，使用gd_svc_cli_actors[GLUSTER_CLI_CREATE_VOLUME]这个rpcsvc_actor_t实现，就具体定义如下：
```
static rpcsvc_actor_t gd_svc_cli_actors[GLUSTER_CLI_MAXVALUE] = {
    [GLUSTER_CLI_CREATE_VOLUME] = {"CLI_CREATE_VOLUME",
                                   glusterd_handle_create_volume, NULL,
                                   GLUSTER_CLI_CREATE_VOLUME, DRC_NA, 0},
    [GLUSTER_CLI_START_VOLUME] = {"START_VOLUME",
                                  glusterd_handle_cli_start_volume, NULL,
                                  GLUSTER_CLI_START_VOLUME, DRC_NA, 0}
}
```

- 创建volume的函数核心的函数都做了什么?
  - glusterd_handle_create_volume:
  - __glusterd_handle_create_volume:
  - glusterd_op_begin_synctask：
  - gd_sync_task_begin:
  - gd_stage_op_phase:
  - gd_brick_op_phase:
  - gd_commit_op_phase:

- 创建过程描述：
  - 第一:
  
- 函数注释
  - [1.glusterd_handle_create_volume](./document/glusterfs调试.md)
  - [2.__glusterd_handle_create_volume](./document/glusterfs调试.md)
  - [3.glusterd_op_begin_synctask](./document/glusterfs调试.md)
  - [4.gd_sync_task_begin](./document/glusterfs调试.md)
  - [5.gd_stage_op_phase](./document/glusterfs调试.md)
  - [6.gd_brick_op_phase](./document/glusterfs调试.md)
  - [7.gd_commit_op_phase](./document/glusterfs调试.md)
