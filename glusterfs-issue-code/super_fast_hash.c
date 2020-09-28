/*************************************************************************
    > File Name: super_fast_hash.c
  > Author:perrynzhou 
  > Mail:perrynzhou@gmail.com 
  > Created Time: Monday, September 28, 2020 AM09:15:56
 ************************************************************************/

#include <stdio.h>
#include <stdint.h>
#include <stdlib.h>
#include <assert.h>
#include <fcntl.h>
#include <uuid/uuid.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <pthread.h>
#define get16bits(d) (*((const uint16_t *)(d)))

typedef struct thread_ctx_t
{
  uint32_t id;
  pthread_t tid;
  uint64_t *counter;
  uint64_t uid_count;
  uint32_t thread_count;
  char *uid_save_path;
} thread_ctx;
uint32_t super_fast_hash(const char *data, int32_t len)
{
  uint32_t hash = len, tmp;
  int32_t rem;

  if (len <= 1 || data == NULL)
    return 1;

  rem = len & 3;
  len >>= 2;

  /* Main loop */
  for (; len > 0; len--)
  {
    hash += get16bits(data);
    tmp = (get16bits(data + 2) << 11) ^ hash;
    hash = (hash << 16) ^ tmp;
    data += 2 * sizeof(uint16_t);
    hash += hash >> 11;
  }

  /* Handle end cases */
  switch (rem)
  {
  case 3:
    hash += get16bits(data);
    hash ^= hash << 16;
    hash ^= data[sizeof(uint16_t)] << 18;
    hash += hash >> 11;
    break;
  case 2:
    hash += get16bits(data);
    hash ^= hash << 11;
    hash += hash >> 17;
    break;
  case 1:
    hash += *data;
    hash ^= hash << 10;
    hash += hash >> 1;
  }

  /* Force "avalanching" of final 127 bits */
  hash ^= hash << 3;
  hash += hash >> 5;
  hash ^= hash << 4;
  hash += hash >> 17;
  hash ^= hash << 25;
  hash += hash >> 6;

  return hash;
}

void *thread_func(void *ctx)
{
  thread_ctx *thd = (thread_ctx *)ctx;
  //char *uuid_str=(char *)calloc(37,sizeof(char));
  for (uint64_t i = 0; i < thd->uid_count; i++)
  {
    // typedef unsigned char uuid_t[16];
    uuid_t uuid;

    // generate
    uuid_generate(uuid);
    // char uuid_str[38] = {'\0'};      // ex. "1b4e28ba-2fa1-11d2-883f-0016d3cca427" + "\0"
    //uuid_unparse_lower(uuid, uuid_str);
    //printf("address :%p,generate uuid=%s\n",&uuid, uuid_str);
    uint32_t index = super_fast_hash((char *)&uuid, sizeof(uuid)) % thd->thread_count;
    //  printf("address :%p,super_fast_hash uuid=%s\n",&uuid, uuid_str);
    __sync_fetch_and_add(&thd->counter[index], 1);
  }
  /*
  if(uuid_str !=NULL)
  {
    free(uuid_str);
  }
  */
}
//uuid_generate_time_safe
//argv[1], worker thread
//argv[2],number of uuid
int main(int argc, const char *argv[])
{
  uint32_t thread_count = atoi(argv[1]);
  uint64_t *counter = calloc(atoi(argv[1]), sizeof(uint64_t));
  assert(counter != NULL);
  thread_ctx *threads = (thread_ctx *)calloc(thread_count, sizeof(thread_ctx));
  assert(threads != NULL);
  uint64_t uid_count = atoi(argv[2]);
  for (uint32_t i = 0; i < thread_count; i++)
  {
    threads[i].id = i;
    threads[i].counter = counter;
    threads[i].uid_count = uid_count;
    threads[i].thread_count = thread_count;
    pthread_create(&threads[i].tid, NULL, (void *)&thread_func, &threads[i]);
  }
  for (uint32_t i = 0; i < thread_count; i++)
  {
    pthread_join(threads[i].tid, NULL);
  }
  for (uint32_t i = 0; i < thread_count; i++)
  {
    fprintf(stdout, "index[%d]:%ld\n", i, counter[i]);
  }
  return 0;
}