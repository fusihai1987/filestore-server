Centos7下Docker部署Ceph集群(nautilus最新版)
****************************
1 gateway默认端口7480

****************************
准备三个host节点
172.18.104.97
172.18.104.98
172.18.104.99
初始化ceph挂载目录(每个节点都执行)
rm -rf /etc/ceph /var/lib/ceph /var/log/ceph
mkdir -p /etc/ceph /var/lib/ceph /var/log/ceph
**启动monitor节点**
启动第一个Monitor节点(在172.18.104.97上执行)
docker run -d --net=host 
--restart always 
-v /etc/ceph:/etc/ceph 
-v /var/lib/ceph/:/var/lib/ceph/ 
-v /var/log/ceph/:/var/log/ceph/ 
-e MON_IP=172.18.104.97 
-e CEPH_PUBLIC_NETWORK=172.18.104.0/24 
--name="ceph-mon" 
ceph/daemon mon

# 然后将/etc/ceph目录拷到其他两个节点, 如用scp命令远程拷贝:
scp -r /etc/ceph root@172.18.104.98:/etc/
scp -r /etc/ceph root@172.18.104.99:/etc/
启动第二个monitor节点(在172.18.104.98上执行)
docker run -d --net=host 
--restart always 
-v /etc/ceph:/etc/ceph 
-v /var/lib/ceph/:/var/lib/ceph/ 
-v /var/log/ceph/:/var/log/ceph/ 
-e MON_IP=172.18.104.98 
-e CEPH_PUBLIC_NETWORK=172.18.104.0/24 
--name="ceph-mon" 
ceph/daemon mon
启动第三个monitor节点(在172.18.104.99上执行)
docker run -d --net=host 
--restart always 
-v /etc/ceph:/etc/ceph 
-v /var/lib/ceph/:/var/lib/ceph/ 
-v /var/log/ceph/:/var/log/ceph/ 
-e MON_IP=172.18.104.99 
-e CEPH_PUBLIC_NETWORK=172.18.104.0/24 
--name="ceph-mon" 
ceph/daemon mon


**启动mgr节点(每个节点执行同样的run命令)**
# mgr模块用于分担monitor部分扩展功能,减轻monitor负担
docker run -d --net=host 
--privileged=true 
--pid=host 
--name="ceph-mgr" 
--restart=always 
-v /etc/ceph:/etc/ceph 
-v /var/lib/ceph/:/var/lib/ceph/ 
ceph/daemon mgr
查看集群状态


**启动OSD节点(每个节点执行)**
# 导出osd用于连ceph集群的keyring
docker exec ceph-mon ceph auth get client.bootstrap-osd -o /var/lib/ceph/bootstrap-osd/ceph.keyring
# 创建osd的存储目录
mkdir -p /data/ceph/osd/vdb
# 启动osd
docker run -d 
--privileged=true 
--name=ceph-osdvdb 
--net=host 
-v /etc/ceph:/etc/ceph 
-v /var/lib/ceph/:/var/lib/ceph/ 
-v /data/ceph/osd/vdb:/var/lib/ceph/osd 
-e OSD_TYPE=directory 
-v /etc/localtime:/etc/localtime:ro 
ceph/daemon osd


**启动gateway节点**
# 导出rgw用于连接集群的keyring
docker exec ceph-mon ceph auth get client.bootstrap-rgw -o /var/lib/ceph/bootstrap-rgw/ceph.keyring
# 运行rgw节点, 可以启动一个或多个
docker run -d --net=host 
--privileged=true 
--name=ceph-rgw 
-v /var/lib/ceph/:/var/lib/ceph/ 
-v /etc/ceph:/etc/ceph 
-v /etc/localtime:/etc/localtime:ro 
-e RGW_NAME=rgw0 
ceph/daemon rgw



**启动dashboard可视化管理页面**
# 开启dashboard模块并禁用ssl(也可以用ssl, 需额外配置ssl证书)
docker exec ceph-rgw ceph mgr module enable dashboard
docker exec ceph-rgw ceph config set mgr mgr/dashboard/ssl false
docker exec ceph-rgw ceph mgr module disable dashboard
docker exec ceph-rgw ceph mgr module enable dashboard
# 设置UI管理的host:port, 登录名及密码
docker exec ceph-rgw ceph config set mgr mgr/dashboard/server_addr 172.18.104.97
docker exec ceph-rgw ceph config set mgr mgr/dashboard/server_port 18080
docker exec ceph-rgw ceph dashboard ac-user-create <自定义user> <自定义pwd> administrator

**创建一个user, 用于管理存储对象**
docker exec ceph-rgw radosgw-admin user create --uid=test_user --display-name=test_user --system