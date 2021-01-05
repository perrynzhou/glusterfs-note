## event-threads设定后都做了什么
| 作者 | 时间 |QQ技术交流群 |
| ------ | ------ |------ |
| perrynzhou@gmail.com |2020/12/01 |中国开源存储技术交流群(672152841) | 
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

  