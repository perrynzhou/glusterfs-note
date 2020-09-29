## glusterfs 代码提交流程

- git配置

  ```
  git config --global user.name "perrynzhou"
  git config --global user.email "perrynzhou@gmail.com"
  ssh-keygen -t rsa -C "perrynzhou@gmail.com"
  ```

- 登陆[Gerrit](http://review.gluster.org/)和在github上认证

  ```
  https://review.gluster.org/#/dashboard/self
  ```
- 在红色框内配置自己github邮箱和sshkey,配置完毕后需要登陆邮箱地址进行确认
  ![avatar](../images/config_2020929_111226.jpg)

- 克隆代码
```
git clone ssh://perrynzhou@review.gluster.org/glusterfs.git
```
- checkout分支
```
cd glusterfs
git checkout -b perryn/{问题}-dev
```
- 修改代码,并提交
```
 git add . --all
 git commit -m "fixed xxx issue"
 ./rfc.sh
 //接着输入这个问题关联的issue的id,比如https://github.com/gluster/glusterfs/issues/1499 这个issue的id就是1499
```
