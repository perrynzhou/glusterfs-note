## event-threads设定后都做了什么
| 作者                 | 时间       | QQ技术交流群                      |
| -------------------- | ---------- | --------------------------------- |
| perrynzhou@gmail.com | 2020/12/01 | 中国开源存储技术交流群(672152841) |
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

  ```shell
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

- 核心函数

  ```c
  // 在客户端设置这个参数，最终会调用protocol/client/src/client.c中的reconfigure函数
  // 通过设定 gluster volume set single-dht-vol client.event-threads 4 
  int reconfigure(xlator_t *this, dict_t *options)
  {
        // 实际调用的是这个方法 xlator_reconfigure_rec
  	   GF_OPTION_RECONF("event-threads", new_nthread, options, int32, out);
         //设定procotol/client 这个xlator网络发送和接受时候的event_pool工作线程数据，调整这个参数确实可以增大处理网络请求的效率
         
     	   ret = client_check_event_threads(this, conf, conf->event_threads,new_nthread)
     	   {
     		 	  conf->event_threads = new;
                gf_event_reconfigure_threads(this->ctx->event_pool,conf->event_threads)
                {
                  // 这个函数是实际设定针对event_pool的工作线程，线程数最大是EVENT_MAX_THREADS(1024)
                	event_reconfigure_threads_epoll(struct event_pool *event_pool, int value);
                }
     		}
  }
  ```

  

### server.event-threads分析

- server.event-threads作用在xlator的类型是protocol/server，也就是在服务端协议层

  ```shell
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

- 核心函数

  ```c
  // 设定服务端网络收发事件的work,通过 gluster volume set single-dht-vol server.event-threads 16 进行设定
  // gluster volume set single-dht-vol server.event-threads 16 命令异步执行，只要客户端发出后，服务端监听到这个网络事件后返回
  int server_reconfigure(xlator_t *this, dict_t *options)
  {
   	GF_OPTION_RECONF("event-threads", new_nthread, options, int32, out);
      ret = server_check_event_threads(this, conf, new_nthread)
      {
          // 设置服务端监听网络事件的work数
      	gf_event_reconfigure_threads(pool, target){
      		event_reconfigure_threads_epoll(struct event_pool *event_pool, int value);
      	}
      }
  }
  ```

  
