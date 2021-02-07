### fixes storoge.reserve检测磁盘预留空闲时间间隔问题

| 作者 | 时间 |QQ技术交流群 |
| ------ | ------ |------ |
| perrynzhou@gmail.com |2020/12/01 |中国开源存储技术交流群(672152841) |

- 问题
  - 目前glusterfs有一个参数运行在storage/posix的xlator中的storage.reserve参数，这个参数用于预留每个brick的空闲空间，当空闲空间少于预留空间时候，在数据写入时候会提示“No Space Left”。这种情况可以很好规避当存储集群快要满的时候进行扩容，而不是一直写直到写满导致集群不可用。
  - 目前这个预留空间设置是通过storage.reserve来设置百分比或者大小，但是在storage/posix下检车每个操作都会去检查一个变量(posix_private->disk_reserve),这个变量在每次数据写入操作都会使用，而这个变量是在一个线程内定时更新的，实现是一个hard code方式，默认是5s.如果针对预留空间比较高的情况，这个时间间隔是不太合适，最好可以配置。可配置根据业务场景和需求来，所以才向社区提了这个pr.
- fork 最新glusterfs分支
    ```
    //这里的{username}是你从glusterfs社区fork过来的你自己的分支，比如我自己的perrynzhou,这里的username就是perrynzhou
    git clone git@github.com:${username}/glusterfs.git
    cd glusterfs/
    git remote add upstream git@github.com:gluster/glusterfs.git
    git fetch upstream
    git rebase upstream/devel
    git checkout -b perryn/dev
    //进行修改，修改完毕后执行
    git commit -m 'xxx'
    //修改提交信息，Change-Id可以从git log查看commit id，设置Change-Id 为commit id
    git commit --amend
    /---------------修改git message----------------/
    change to an option for time checking storage.reserve
    storage/posix: change to an option for time checking storage.reserve (#2120)
    
    
    Change-Id: 66e34de17b1e6e308a5b648fc2232f543047b40d
    Fixed: #2120
    Signed-off-by: perrynzhou <perrynzhou@gmail.com>
    /---------------修改git message----------------/    
    ./rfc.sh
    git push origin perryn/dev
    //通过社区的code view如果没有问题可以把perryn/dev分支代码merge到glusterfs社区分支，需要提一个PR
    //基于自己的分支来提一个PR，让社区的小伙伴进行code review
    ```
- 添加Gluster Cli执行的参数选项

    ```
   // 添加gluster cli支持的参数
   //glusterfs-8.3/xlators/mgmt/glusterd/src/glusterd-volume-set.c:2410
   struct volopt_map_entry glusterd_volopt_map[] = {
   {
          .key = "storage.reserve",
          .voltype = "storage/posix",
          .op_version = GD_OP_VERSION_3_13_0,
      },
      {
          .key = "storage.reserve-check-interval",
          .voltype = "storage/posix",
          .op_version = GD_OP_VERSION_4_0_0,
      }
    }
    ```

  

- 在storage/posix xlator中posix_private结构体添加存储客户端参数的字段

  ```
  // 在storage/posix xlator中的posix_private添加一个字段，用户存储来自客户端的传进来的storage.reserve-check-interval值
  //glusterfs-8.3/xlators/storage/posix/src/posix.h
  struct posix_private {
  	    double disk_reserve;
      /* seconds for check disk reversion  */
      uint32_t disk_reserve_check_interval;
  }
  ```

- 初始化reserve-check-interval的值

  ```
  // 需要在posix_init中初始化reserve-check-interval的值
  //glusterfs-8.3/xlator/storage/posix/src/posix-common.c
  int posix_init(xlator_t *this)
  {
  
      GF_OPTION_INIT("reserve", _private->disk_reserve, percent_or_size, out);
      //在这里添加在storage/posix xlator中的初始化值，这个值依赖于posix_options配置的默认值
      // add disk reserve internal seconds for reserve check
      GF_OPTION_INIT("reserve-check-interval", _private->disk_reserve_check_interval, uint32, out);
  }
  struct volume_options posix_options[] = {
       {.key = {"reserve-check-interval"},
       .type = GF_OPTION_TYPE_INT,
       .min = 1,
       .default_value = "5",
       .validate = GF_OPT_VALIDATE_MIN,
       .description =  "Interval in second to check disk reserve",
       .op_version = {GD_OP_VERSION_4_0_0},
       .flags = OPT_FLAG_SETTABLE | OPT_FLAG_DOC},
     }
  ```


- 配置reserve-check-interval重新加载

  ```
  // glusterfs-8.3/xlator/storage/posix/src/posix-common.c
  int posix_reconfigure(xlator_t *this, dict_t *options)
  {
      GF_OPTION_RECONF("reserve", priv->disk_reserve, options, percent_or_size,
                       out);
      // add config reconfig
      GF_OPTION_RECONF("reserve-check-interval", priv->disk_reserve_check_interval, options, uint32,
                       out);
  }
  ```


- 修改posix_disk_space_check_thread_proc函数中的休眠时间

    ```
      /********glusterfs 8.x版本修改方式********/

    //glusterfs-8.3/xlator/storage/posix/src/posix-helpers.c
    static void *posix_disk_space_check_thread_proc(void *data)
    {
      xlator_t *this = NULL;
      struct posix_private *priv = NULL;
      int ret = -1;
  
      this = data;
      priv = this->private;
  
      gf_msg_debug(this->name, 0,
                   "disk-space thread started, "
                   "interval = %d seconds",
                   interval);
      while (1) {
          ret = sleep(priv->reserve-check-interval);
          if (ret > 0)
              break;
          pthread_setcancelstate(PTHREAD_CANCEL_DISABLE, NULL);
          // 磁盘检查函数
          posix_disk_space_check(this);
          if (!priv->disk_space_check_active)
              goto out;
          pthread_setcancelstate(PTHREAD_CANCEL_ENABLE, NULL);
      }
  
      return NULL;
    }
  
 
    /********glusterfs 9.x版本以及以上版本********/
    int posix_spawn_disk_space_check_thread(xlator_t *this)
    {
	    //修改gf_thread_create函数由原来传递的glusterfs_ctx_t修改为xlator,每次storage.reserve重新开启都会在Debug日志中打印reserve-check-interval时间
	    ret = gf_thread_create(&ctx->disk_space_check, NULL,posix_ctx_disk_thread_proc, this,"posixctxres");
    }

    // 实际做磁盘检查的线程工作，posix_ctx_disk_thread_proc原来传递的是glusterfs_ctx_t指针修改为xlator指针，在函数内部从xlator指针获取glusterfs_ctx_t指针
    static void *posix_ctx_disk_thread_proc(void *data)
    {
      struct posix_private *priv = NULL;
      glusterfs_ctx_t *ctx = NULL;
      priv = this->private;
      ctx = this->ctx;
      pthread_mutex_lock(&ctx->xl_lock);
      {
        while (ctx->diskxl_count > 0) {
            list_for_each_entry(pthis, &ctx->diskth_xl, list)
            {
                timespec_now_realtime(&sleep_till);
                // 每次累计priv->disk_reserve_check_interval;直到超时priv->disk_reserve_check_interval
                sleep_till.tv_sec += priv->disk_reserve_check_interval;
                (void)pthread_cond_timedwait(&ctx->xl_cond, &ctx->xl_lock,
                                         &sleep_till);
            }
        }
        pthread_mutex_unlock(&ctx->xl_lock);
    
        return NULL;
    }
    ```
  
  

- 修改的文件

    ```
    [root@CentOS glusterfs]# git status
    # On branch perryn/storage-reserve-option-dev
    # Changes not staged for commit:
    #   (use "git add <file>..." to update what will be committed)
    #   (use "git checkout -- <file>..." to discard changes in working directory)
    #
    #       modified:   xlators/mgmt/glusterd/src/glusterd-volume-set.c
    #       modified:   xlators/storage/posix/src/posix-common.c
    #       modified:   xlators/storage/posix/src/posix-helpers.c
    #       modified:   xlators/storage/posix/src/posix.h
    ```