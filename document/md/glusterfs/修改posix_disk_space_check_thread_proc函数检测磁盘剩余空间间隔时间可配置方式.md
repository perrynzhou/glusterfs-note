### 修改posix_disk_space_check_thread_proc函数检测磁盘剩余空间间隔时间可配置方式

| 作者 | 时间 |QQ技术交流群 |
| ------ | ------ |------ |
| perrynzhou@gmail.com |2020/12/01 |中国开源存储技术交流群(672152841) |

- glusterfs-8.3/libglusterfs/src/glusterfs/common-utils.h

  ```
  #define STORAGE_CHECK_TIMEOUT "5"
  ```

  

- glusterfs-8.3/xlators/mgmt/glusterd/src/glusterd-volume-set.c:2410

  ```
      {
          .key = "storage.reserve",
          .voltype = "storage/posix",
          .op_version = GD_OP_VERSION_3_13_0,
      },
        {
          .key = "storage.reserve-check-timeout",
          .voltype = "storage/posix",
          .value = STORAGE_CHECK_TIMEOUT,
          .op_version = GD_OP_VERSION_3_13_0,
      },
  ```

  

- glusterfs-8.3/xlators/storage/posix/src/posix.h

  ```
  struct posix_private {
  	 int32_t disk_reserve_check_timeout;
  }
  ```

- glusterfs-8.3/xlator/storage/posix/src/posix-common.c

  ```
  int
  posix_reconfigure(xlator_t *this, dict_t *options)
  {
   GF_OPTION_RECONF("reserve-check-timeout", priv->disk_reserve, options, int32,
                       out);
  }
  ```

- glusterfs-8.3/xlator/storage/posix/src/posix.c

  ```
  int
  posix_init(xlator_t *this)
  {
  ret = dict_get_int32(this->options, "reserve-check-timeout",
                           &_private->disk_reserve_check_timeout);
      if (ret == -1) {
          gf_msg(this->name, GF_LOG_ERROR, 0, P_MSG_INVALID_OPTION_VAL,
                 "'disk_reserve_check_timeout' takes only integer "
                 "values");
          goto out;
      }
  }
  ```

- glusterfs-8.3/xlator/storage/posix/src/posix-helpers.c

  ```
  static void *
  posix_disk_space_check_thread_proc(void *data)
  {
      xlator_t *this = NULL;
      struct posix_private *priv = NULL;
      uint32_t interval = 0;
      int ret = -1;
  
      this = data;
      priv = this->private;
  
      interval = priv->disk_reserve_check_timeout;
      gf_msg_debug(this->name, 0,
                   "disk-space thread started, "
                   "interval = %d seconds",
                   interval);
      while (1) {
          /* aborting sleep() is a request to exit this thread, sleep()
           * will normally not return when cancelled */
          ret = sleep(interval);
          if (ret > 0)
              break;
          /* prevent thread errors while doing the health-check(s) */
          pthread_setcancelstate(PTHREAD_CANCEL_DISABLE, NULL);
  
          /* Do the disk-check.*/
          posix_disk_space_check(this);
          if (!priv->disk_space_check_active)
              goto out;
          pthread_setcancelstate(PTHREAD_CANCEL_ENABLE, NULL);
      }
  
  out:
      gf_msg_debug(this->name, 0, "disk space check thread exiting");
      LOCK(&priv->lock);
      {
          priv->disk_space_check_active = _gf_false;
      }
      UNLOCK(&priv->lock);
  
      return NULL;
  }
  ```

  