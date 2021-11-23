#### ①构建本地镜像
  参照Dockerfile和Makefile
  
#### ②将镜像推送到Docker Hub：mengfu521/httpserver
```
make realse
```

#### ③通过Docker命令本地启动httpserver
```
fumeng@fumengdeMacBook-Pro httpserver % docker run -it -d -P mengfu521/httpserver:v1.0
7ede587cf5ce4bd80dcd2715b2a52b8dc3ae5f3a4ec76dcbf4ae4170cc5cfa94
```
```
fumeng@fumengdeMacBook-Pro httpserver % docker ps
CONTAINER ID   IMAGE                       COMMAND                  CREATED         STATUS         PORTS                   NAMES
7ede587cf5ce   mengfu521/httpserver:v1.0   "/bin/sh -c /httpser…"   3 seconds ago   Up 2 seconds   0.0.0.0:55000->80/tcp   boring_taussig
```
```
fumeng@fumengdeMacBook-Pro httpserver % curl 0.0.0.0:55000
===================Details of the http request header:============
User-Agent=[curl/7.64.1]
Accept=[*/*]
VERSION=1.15.6
```

#### ④查看容器IP，因Mac OS无法使用nsenter命令，用其他方法取得ip
方法1:  执行docker inspect命令
```
fumeng@fumengdeMacBook-Pro httpserver % docker inspect --format='{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' 7ede587cf5ce
172.17.0.2
```

方法2:  进入容器，查看/etc/hosts配置
```
fumeng@fumengdeMacBook-Pro httpserver % docker exec -it 7ede587cf5ce bin/bash
root@7ede587cf5ce:/# cat  /etc/hosts
127.0.0.1       localhost
::1     localhost ip6-localhost ip6-loopback
fe00::0 ip6-localnet
ff00::0 ip6-mcastprefix
ff02::1 ip6-allnodes
ff02::2 ip6-allrouters
172.17.0.2      7ede587cf5ce
```
