# start project
```
 mkdir ecommerce-yt
 cd .\ecommerce-yt
 go mod init github.com/neutron2025/ecommerce-yt
```
 # when finished
 ```
 docker-compose up -d

 go run main.go
 ```

 ## install docker on ubuntu
 ```
$ curl -fsSL get.docker.com -o get-docker.sh
$ sudo sh get-docker.sh --mirror Aliyun

$ sudo systemctl enable docker
$ sudo systemctl start docker

$ docker run --rm hello-world
```
## 配置国内镜像
```
如果/etc/docker/daemon.json文件不存在，创建该文件并添加以下内容：
{
  "registry-mirrors": [
    "https://docker.m.daocloud.io",
    "https://docker.mirrors.ustc.edu.cn",
    "https://docker.nju.edu.cn"
  ]
}
```
然后重启Docker服务：
```
sudo systemctl daemon-reload
sudo systemctl restart docker

```
### link
https://docker-practice.github.io/zh-cn/install/ubuntu.html