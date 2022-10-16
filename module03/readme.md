1. 编写 Dockerfile 将练习 2.2 编写的 httpserver 容器化,构建本地镜像

   [Dockerfile](Dockerfile)

   构建本地镜像

   ```
   docker build -t hecatefu/httpserver:1.0-alpine3.16 .
   ```
   
   执行结果

   ```
   root@fu-ubuntu:/home/hecate/module03# docker build -t hecatefu/httpserver:1.0-alpine3.16 .
   Sending build context to Docker daemon  6.656kB
   Step 1/8 : FROM golang:1.19.2-alpine3.16 AS build
    ---> f9a40cb7e8ec
   Step 2/8 : ENV GO111MODULE=on     CGO_ENABLED=0     GOOS=linux     GOARCH=amd64
    ---> Using cache
    ---> 6702db3bd0d8
   Step 3/8 : WORKDIR /go/src/httpserver
    ---> Using cache
    ---> 70ce2a1f013e
   Step 4/8 : COPY ./httpserver/* ./
    ---> Using cache
    ---> 3baf7986d293
   Step 5/8 : RUN go build -o /bin/httpserver
    ---> Using cache
    ---> 17cbefadee75
   Step 6/8 : FROM alpine:3.16
   3.16: Pulling from library/alpine
   213ec9aee27d: Already exists 
   Digest: sha256:bc41182d7ef5ffc53a40b044e725193bc10142a1243f395ee852a8d9730fc2ad
   Status: Downloaded newer image for alpine:3.16
    ---> 9c6f07244728
   Step 7/8 : COPY --from=build /bin/httpserver /bin/httpserver
    ---> abc7873ccbbf
   Step 8/8 : ENTRYPOINT ["/bin/httpserver"]
    ---> Running in 695b2d95f2e1
   Removing intermediate container 695b2d95f2e1
    ---> b18bd784fdec
   Successfully built b18bd784fdec
   Successfully tagged hecatefu/httpserver:1.0-alpine3.16
   root@fu-ubuntu:/home/hecate/module03# docker images
   REPOSITORY            TAG                 IMAGE ID       CREATED         SIZE
   hecatefu/httpserver   1.0-alpine3.16      b18bd784fdec   6 seconds ago   12MB
   hecatefu/httpserver   latest              102c57e305e9   4 hours ago     6.46MB
   hecatefu/httpserver   test                102c57e305e9   4 hours ago     6.46MB
   <none>                <none>              17cbefadee75   4 hours ago     375MB
   golang                1.19.2-alpine3.16   f9a40cb7e8ec   9 days ago      352MB
   alpine                3.16                9c6f07244728   2 months ago    5.54MB
   ```

2. 将镜像推送至 docker 官方镜像仓库

   先在 hub.docker.com 注册账号，创建仓库 httpserver, 然后在本地登录远程仓库

   ```
   root@fu-ubuntu:/home/hecate/module03# docker login
   Login with your Docker ID to push and pull images from Docker Hub. If you don't have a Docker ID, head over to https://hub.docker.com to create one.
   Username: hecatefu
   Password: 
   WARNING! Your password will be stored unencrypted in /root/.docker/config.json.
   Configure a credential helper to remove this warning. See
   https://docs.docker.com/engine/reference/commandline/login/#credentials-store
   
   Login Succeeded
   ```

   然后执行推送仓库命令

   ```
   root@fu-ubuntu:/home/hecate/module03# docker push hecatefu/httpserver:1.0-alpine3.16
   The push refers to repository [docker.io/hecatefu/httpserver]
   2cbb849a33cf: Pushed 
   994393dc58e7: Mounted from library/alpine 
   1.0-alpine3.16: digest: sha256:88414a391cf751c5adeb1fa2facbd747356cda060ed626660736f8d4382a6ebf size: 739
   ```

3. 通过 docker 命令本地启动 httpserver
   
   启动命令

   ```
   docker run --name httpserver -p 80:8080 -itd hecatefu/httpserver:1.0-alpine3.16
   ```

   执行结果

   ```
   root@fu-ubuntu:/home/hecate/module03# docker run --name httpserver -p 80:8080 -itd hecatefu/httpserver:1.0-alpine3.16
   a6d6baccb8b5191811ed3d137906ce4b1bac309517669239cd78dead60d1b993
   root@fu-ubuntu:/home/hecate/module03# docker ps
   CONTAINER ID   IMAGE                                COMMAND             CREATED         STATUS         PORTS                                   NAMES
   a6d6baccb8b5   hecatefu/httpserver:1.0-alpine3.16   "/bin/httpserver"   2 minutes ago   Up 2 minutes   0.0.0.0:80->8080/tcp, :::80->8080/tcp   httpserver
   ```
   
   验证容器

   ```
   root@fu-ubuntu:~# curl -vv http://localhost/
   *   Trying 127.0.0.1:80...
   * TCP_NODELAY set
   * Connected to localhost (127.0.0.1) port 80 (#0)
   > GET / HTTP/1.1
   > Host: localhost
   > User-Agent: curl/7.68.0
   > Accept: */*
   > 
   * Mark bundle as not supporting multiuse
   < HTTP/1.1 403 Forbidden
   < Accept: */*
   < User-Agent: curl/7.68.0
   < Version: 
   < Date: Sun, 16 Oct 2022 12:09:08 GMT
   < Content-Length: 11
   < Content-Type: text/plain; charset=utf-8
   < 
   * Connection #0 to host localhost left intact
   hello world
   ```

4. 通过 nsenter 进入容器查看 IP 配置

   查看容器进程pid

   ```
   root@fu-ubuntu:~# docker inspect -f {{.State.Pid}} httpserver
   5739
   ```

   nsenter查看ip配置

   ```
   root@fu-ubuntu:~# nsenter -t 5739 -n ip addr
   1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
       link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
       inet 127.0.0.1/8 scope host lo
          valid_lft forever preferred_lft forever
   26: eth0@if27: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue state UP group default 
       link/ether 02:42:ac:11:00:02 brd ff:ff:ff:ff:ff:ff link-netnsid 0
       inet 172.17.0.2/16 brd 172.17.255.255 scope global eth0
          valid_lft forever preferred_lft forever
   ```