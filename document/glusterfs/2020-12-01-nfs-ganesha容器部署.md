## nfs-ganesha容器方案



## 解决问题

- nfs-ganesha是一个nas网关，服务接入的是glusterfs服务，底层通过gfapi和glusterfs进行交互
- 基于物理机方式进行部署，需要解决nfs-ganesha多实例部署的端口隔离，同时对外提供服务需要提供负载均衡，说到这里就可以使用k8s的service来做负载均衡，每个机器部署gluster卷的一个或者多个nfs-ganesha实例（nfs-ganesha如果是ipv4是走111端口，如果是ipv6是走2049端口)


## nfs-ganehsa镜像

- 首先需要构建nfs-ganesha进行，具体镜像构建和物理机安装的方式一致，镜像完成后需要push到镜像仓库
- nfs-ganesha的pod资源定义
```
//为了让nfs-ganesha实例可以进行mount需要super权限和d-bus权限
# cat pod-template.yaml 
apiVersion: v1
kind: Pod
metadata:
  name: test2-pod
spec:
  containers:
    - name: test-container
      image: xxx/nfs-ganesha
      command: 
      - /usr/sbin/init
      - /usr/bin/ganesha.nfsd -L /var/log/ganesha/ganesha.log -f /etc/ganesha/ganesha.conf -N NIV_EVENT
      volumeMounts:
      - name: config-volume
        mountPath: /etc/ganesha
      securityContext:
        privileged: True
  volumes:
    - name: config-volume
      configMap:
        name: speech-v5-rep-vol-config
  restartPolicy: Never
```

## nfs-ganesha配置文件存储到configmap

```
# cat speech_v5_rep_vol.conf 
NFS_CORE_PARAM {
        mount_path_pseudo = true;
        Protocols = 3,4;
}

EXPORT_DEFAULTS {
        Access_Type = RW;
}

EXPORT{
    Export_Id = 101 ;   
    Path = "/mnt/speech_v5_rep_vol";  

    FSAL {
        name = GLUSTER;
        hostname = "10.193.226.12"; 
        volume = "speech_v5_rep_vol";  
    }

    Access_type = RW;    
    Squash = No_root_squash; 
    Disable_ACL = TRUE;  
    Pseudo = "/speech_v5_rep_vol";  
    Protocols = "3","4" ;  
    Transports = "UDP","TCP" ; 
    SecType = "sys";    
}
```
```
// 这里会建立一个configmap资源，然后--from-file=ganesha.conf=speech_v5_rep_vol.conf中第一个ganesha.conf定义为这个文件内容的文件名称，第二个speech_v5_rep_vol是nfs-ganesha针对一个卷的配置
# kubectl create configmap speech-v5-rep-vol-config --from-file=ganesha.conf=speech_v5_rep_vol.conf 

//查看kubectl configmap中内容
# kubectl describe configmap
Name:         speech-v5-rep-vol-config
Namespace:    default
Labels:       <none>
Annotations:  <none>

Data
====
ganesha.conf:
----
NFS_CORE_PARAM {
        mount_path_pseudo = true;
        Protocols = 3,4;
}

EXPORT_DEFAULTS {
        Access_Type = RW;
}

EXPORT{
    Export_Id = 101 ;   
    Path = "/mnt/speech_v5_rep_vol";  

    FSAL {
        name = GLUSTER;
        hostname = "10.191.13.12"; 
        volume = "speech_v5_rep_vol";  
    }

    Access_type = RW;    
    Squash = No_root_squash; 
    Disable_ACL = TRUE;  
    Pseudo = "/speech_v5_rep_vol";  
    Protocols = "3","4" ;  
    Transports = "UDP","TCP" ; 
    SecType = "sys";    
}

Events:  <none>
```
## nfs-ganesha的service定义

```
# cat service-template.yaml 
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-service
spec:
  selector:
    matchLabels:
      run: load-balancer-example
  replicas: 2
  template:
    metadata:
      labels:
        run: load-balancer-example
    spec:
      containers:
        - name: test-container
          image: registry/romai_dev/nfs-ganesha
          imagePullPolicy: Always
          ports:
          - containerPort: 2049
          command: ["/usr/sbin/init"] 
          lifecycle:
           postStart:
            exec:
              command: ["sh","-c","/usr/bin/ganesha.nfsd -L /var/log/ganesha/ganesha.log -f /etc/ganesha/ganesha.conf -N NIV_EVENT"]
          volumeMounts:
          - name: config-volume
            mountPath: /etc/ganesha
          securityContext:
            privileged: True
      volumes:
        - name: config-volume
          configMap:
            name: speech-v5-rep-vol-config
```

```
kubectl create -f service-template.yaml
```
## nfs-ganesha导出服务

```
# kubectl expose deployment/test-service
# kubectl get svc
NAME                    TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)    AGE
test-service            ClusterIP   10.43.186.234   <none>        2049/TCP   22m
```

## 启动一个客户端pod进行挂载nfs-ganesha

```
# cat client-pod.yaml 
apiVersion: v1
kind: Pod
metadata:
  name: client-pod
spec:
  containers:
    - name: test-container
      image: xxx/nfs-ganesha
      command: 
      - /usr/sbin/init
      securityContext:
        privileged: True
  restartPolicy: Never
  
# kubectl create -f client-pod.yaml
# kubectl exec -ti client-pod   bash
[root@client-pod /]#  mount -t nfs4  -o  port=2049,vers=4,proto=tcp   test-service:/speech_v5_rep_vol /mnt/speech_v5_rep_vol
```

