## event-threads设定后都做了什么
| 作者 | 时间 |QQ技术交流群 |
| ------ | ------ |------ |
| perrynzhou@gmail.com |2021/03/27 |672152841 |
### event-threads 参数说明

- client.event-threads:指定客户端多个event线程并行处理，这个线程数调大可以让请求处理更快一些，设定的最大值是32.

  ```
  Option: client.event-threads
  Default Value: 2
  Description: Specifies the number of event threads to execute in parallel. Larger values would help process responses faster, depending on available processing power. Range 1-32 threads.
  ```

  

- server.event-threads:指定客户端多个event线程并行处理，这个线程数调大可以让请求处理更快

  ```
  Option: server.event-threads
  Default Value: 2
  Description: Specifies the number of event threads to execute in parallel. Larger values would help process responses faster, depending on available processing power.
  ```

### client.event-threads分析

- client.event-threads作用在xlator的类型是protocol/client，也就是在客户端协议层

  ```
  // options 在glusterfs/xlator/protocol/client/src/client.c中
  struct volume_options options[] ={
  {.key = {"event-threads"},
       .type = GF_OPTION_TYPE_INT,
       .min = 1,
       .max = 32,
       .default_value = "2",
       .description = "Specifies the number of event threads to execute "
                      "in parallel. Larger values would help process"
                      " responses faster, depending on available processing"
                      " power. Range 1-32 threads.",
       .op_version = {GD_OP_VERSION_3_7_0},
  }
  ```

  



### server.event-threads分析

- server.event-threads作用在xlator的类型是protocol/server，也就是在服务端协议层

  ```
  // server_options 定义在glusterfs/xlator/protocol/server/src/server.c
  struct volume_options server_options[] = {
  {.key = {"event-threads"},
       .type = GF_OPTION_TYPE_INT,
       .min = 1,
       .max = 1024,
       .default_value = "2",
       .description = "Specifies the number of event threads to execute "
                      "in parallel. Larger values would help process"
                      " responses faster, depending on available processing"
                      " power.",
       .op_version = {GD_OP_VERSION_3_7_0},
       .flags = OPT_FLAG_SETTABLE | OPT_FLAG_DOC}
   }
  ```

## event-threads的工作线程更新的实现

```
// 针对glusterfs中的客户端和服务端网络事件监听和处理的工作线程的添加或者减少的处理函数
int event_reconfigure_threads_epoll(struct event_pool *event_pool, int value)
{
    // 函数的部分代码已经省略
    int i;
    int ret = 0;
    pthread_t t_id;
    int oldthreadcount;
    struct event_thread_data *ev_data = NULL;

    // 获取当前已经配置的event-thread值
    oldthreadcount = event_pool->eventthreadcount;

    if (event_pool_dispatched_unlocked(event_pool) &&
        (oldthreadcount < value))
    {
       // 如果配置新的event-threads大于现在的event-threads,则添加监听和处理网络事件的线程
        for (i = oldthreadcount; i < value; i++)
        {
            /* Start a thread if the index at this location
                 * is a 0, so that the older thread is confirmed
                 * as dead */
            if (event_pool->pollers[i] == 0)
            {
                ev_data = GF_CALLOC(1, sizeof(*ev_data),
                                    gf_common_mt_event_pool);
                if (!ev_data)
                {
                    continue;
                }
				
                ev_data->event_pool = event_pool;
                // 设置当前工作线程中的event_index
                ev_data->event_index = i + 1;
				// 创建工作线程
                ret = gf_thread_create(&t_id, NULL, event_dispatch_epoll_worker, ev_data, "epoll%03hx", i & 0x3ff);
				//分离工作线程线程
                pthread_detach(t_id);
                event_pool->pollers[i] = t_id;
            }
        }
    }

    // 更新当前event_pool中的工作线程数
    event_pool->eventthreadcount = value;

    return 0;
}

// 事件的工作线程的处理函数，用于处理当前接受到的网络事件。如果针对glusterfs中的客户端和服务端的处理网络事件的工作线程进行减小，根据传进去的 ev_data->event_index和event_pool中的可用工作线程数进行比较，如果ev_data->event_index <event_pool->eventthreadcount则工作线程退出
static void *event_dispatch_epoll_worker(void *data)
{
   // 函数的部分代码已经省略
    struct epoll_event event;
    int ret = -1;
    struct event_thread_data *ev_data = data;
    struct event_pool *event_pool;
    int myindex = -1;
    int timetodie = 0, gen = 0;
    struct list_head poller_death_notify;
    struct event_slot_epoll *slot = NULL, *tmp = NULL;

    GF_VALIDATE_OR_GOTO("event", ev_data, out);

    event_pool = ev_data->event_pool;
    myindex = ev_data->event_index;
	// 新增一个工作线程，需要针对全局的event_pool中的activethreadcount更新
    pthread_mutex_lock(&event_pool->mutex);
    {
        event_pool->activethreadcount++;
    }
    pthread_mutex_unlock(&event_pool->mutex);

    for (;;)
    {
        // 循环判断当前线程的index和总的线程数，根据结果来选择是执行处理请求退出
        if (event_pool->eventthreadcount < myindex)
        {
            /* ...time to die, thread count was decreased below
             * this threads index */
            /* Start with extra safety at this point, reducing
             * lock conention in normal case when threads are not
             * reconfigured always */
            pthread_mutex_lock(&event_pool->mutex);
            {
                if (event_pool->eventthreadcount < myindex)
                {
                    while (event_pool->poller_death_sliced)
                    {
                        pthread_cond_wait(&event_pool->cond,
                                          &event_pool->mutex);
                    }

                    INIT_LIST_HEAD(&poller_death_notify);
                    /* if found true in critical section,
                     * die */
                    event_pool->pollers[myindex - 1] = 0;
                    event_pool->activethreadcount--;
                    timetodie = 1;
                    gen = ++event_pool->poller_gen;
                    list_for_each_entry(slot, &event_pool->poller_death,
                                        poller_death)
                    {
                        event_slot_ref(slot);
                    }

                    list_splice_init(&event_pool->poller_death,
                                     &poller_death_notify);
                    event_pool->poller_death_sliced = 1;
                    pthread_cond_broadcast(&event_pool->cond);
                }
            }
            pthread_mutex_unlock(&event_pool->mutex);
            if (timetodie)
            {

                list_for_each_entry_safe(slot, tmp, &poller_death_notify,
                                         poller_death)
                {
                    __event_slot_unref(event_pool, slot, slot->idx);
                }

                list_splice(&poller_death_notify,
                            &event_pool->poller_death);
                event_pool->poller_death_sliced = 0;
                pthread_cond_broadcast(&event_pool->cond);

                goto out;
            }
        }

        ret = epoll_wait(event_pool->fd, &event, 1, -1);

        if (ret == 0)
            /* timeout */
            continue;

        if (ret == -1 && errno == EINTR)
            /* sys call */
            continue;

        ret = event_dispatch_epoll_handler(event_pool, &event);
        if (ret)
        {
            gf_smsg("epoll", GF_LOG_ERROR, 0, LG_MSG_DISPATCH_HANDLER_FAILED,
                    NULL);
        }
    }
out:
    if (ev_data)
        GF_FREE(ev_data);
    return NULL;
}
```
