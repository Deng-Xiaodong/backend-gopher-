# 0、 Docker安装

[官方](https://docs.docker.com/engine/install/centos/)

`centos7`

**安装**

```shell
#0 更新centos yum
sudo yum update -y && sudo yum upgrade -y
#1 centos自带的软件包太旧，需要安装扩展包epel
sudo yum install -y epel-release
	
#2 卸载旧版本
sudo yum remove docker \
                  docker-client \
                  docker-client-latest \
                  docker-common \
                  docker-latest \
                  docker-latest-logrotate \
                  docker-logrotate \
                  docker-engine
                    
#4 安装docker
sudo yum install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

#5 查看docker版本
docker version
```



`option`

换源

```shell
# 如需换yum源
wget -O /etc/yum.repos.d/CentOS-Base.repo http://mirrors.aliyun.com/repo/Centos-7.repo

#如需设置docker源
sudo yum install -y yum-utils
sudo yum-config-manager \
    --add-repo \
    https://download.docker.com/linux/centos/docker-ce.repo
    
#安装成功后，如需换镜像源
	#1 创建daemon.json文件
	vim /etc/docker/daemon.json
	#2 写入
    {
        "registry-mirrors":[
                            "https://hub-mirror.c.163.com/",
                            "https://docker.mirrors.ustc.edu.cn/"
                            ]
    }

	#3重启daemon和docker服务
    systemctl daemon-reload
    systemctl restart docker

```



**docker服务启动与关闭**

```shell
#启动
systemctl start docker
#开机自启动
systemctl enable docker
#关闭
systemctl stop docker
#重启
systemctl restart docker
```



**更新docker版本**

```shell
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
```



**添加当前用户到docker用户组**

这样当前用户操作docker无需sudo权限

```shell
#如当前无docker用户组则需要先添加docker用户组
sudo groupadd  docker
#添加当前用户到docker用户组
sudo usermod -aG docker $USER
```



# 1、Docker基本操作

## 1.1 镜像操作

```shell
#通过关键字从docker hub 中检索镜像
docker search keyword
#拉取镜像
docker pull 镜像名:版本 #默认为latest
# 查看当前已有镜像
docker images
#删除镜像，可删除多个
docker rmi 镜像名/镜像ID ...

# 从dockerfile制作镜像
docker build -f based_dockerfile -d 镜像名 版本号 # `.`为latest
```



## 1.2 容器操作

**启动容器**

有以下重点：

- `--name`：指定容器名
- `-p`：指定宿主机到容器的端口映射( $\color{red}宿主机端口:容器端口$)
- `-v`：挂载卷( $\color{red}宿主机目录:容器目录$)
- `-e`：设置容器的环境变量
  - `-d`：后台启动
- 不同容器的特定配置

```shell
#启动redis容器
docker run 
--name some-redis
-p 6379:6379 
-v ./conf/redis.conf:/etc/conf.redis.conf 
-d redis 
redis-server /etc/redis/redis.conf
#  后台启动了一个名为some-redis的容器。其中做了端口映射、目录/文件挂载。最后还规定了redis以配置文件形式启动

#redis.conf
appendonly yes
requirepass 123456


#启动mysql
ocker run 
‐p 3306:3306 
‐‐name MySQL 
‐e MYSQL_ROOT_PASSWORD=123456 
‐d mysql:8.0.15
# ‐e MYSQL_ROOT_PASSWORD=123456: 指定 root 的密码
```



**常规容器操作**

```shell
#交互式运行容器，进入到容器内部
docker exec -it 容器名/容器ID /bin/bash
# 查看运行中的容器
docker ps

# 查看所有的容器 
docker ps ‐a 

# 停止运行中的容器 
docker stop 容器名/容器ID

# 启动容器 
docker start 容器名/容器ID

# 删除一个容器 
docker rm 容器名/容器ID

# 查看容器日志
docker logs 容器名/容器ID
```

**容器外执行容器内命令**

一般用于写shell脚本

```shell
#执行容器内的shell脚本
docker exec 容器标识 /bin/bash -c "chmod +x workdir/shell/start.sh"
docker exec 容器标识 workdir/shell/start.sh
或者
docker exec 容器标识 /bin/bash -c "workdir/shell/start.sh"


#执行多条命令;使用分号隔开
docker exec 容器标识 /bin/bash -c "cmd1;cmd2;cmd3;"

```

## 1.3 卷操作

**挂载卷的好处**

- 容器与主机之间共享数据
- 容器与容器之间共享数据

**卷的类型**

- 命名数据卷：提前指定与之挂载的主机目录
- 匿名数据卷：没有指定主机目录；dockers会自动挂载到` /var/lib/docker/volumes/dbd07daa4e40148b11`

**常见需要挂载的目录类型**

- 配置文件 ：一般docker有默认位置；挂载时需要写对目录
- 数据文件 ：一般自定义；容器内使用时需要注意用对正确的挂载位置
- 日志文件： 一般自定义；容器内使用时需要注意用对正确的挂载位置

**挂载注意事项**

写绝对路径是最后的选择



**卷的操作**

```shell
#容器启动时挂载
-v #进行挂载操作

#查看所有卷
docker volume ls #
#删除卷
docker volume rm xxx #
#清理不用的卷 
docker volume prune -f #-f 强制清除

```

**docker cp 操作**

`docker cp` 能实现容器与主机之间双向的文件复制

```shell
#将宿主机的文件复制到容器内
docker cp 主机目录/文件 容器目录/文件

#将容器内的文件复制到宿主机
docker cp 容器目录/文件 主机目录/文件 
```



# 2、Dockerfile

## 2.1 指令

**FROM** ：定义基础镜像

```dockerfile
FROM nginx
```



**RUN**：构建镜像过程中执行的命令

**CMD/ENTRYPOINT**：该镜像构建的容器在运行后执行的初始命令

​	命令有以下两种格式：

- `Shell`  **RUN/CMD/ENTRYPOINT** $\color{red} <命令行命令>$
- `Exec`   **RUN/CMD/ENTRYPOINT  ["可执行文件","参数1","参数2"]**

```dockerfile
RUN mkdir somedir
ENTRYPOINT ["/bin/bash","./start.sh"]
CMD ["./miaosha","config.json"]
```



**ENV**：设置环境变量；作用周期可以是镜像构建过程或容器内部

```dockerfile
ENV <key> <value>
ENV <key1>=<value1> <key2>=<value2> ...
```



**COPY/ADD**：从上下文目录中复制文件或者目录到容器的指定路径

```dockerfile
COPY [--chown=<user>:<group>] <源路径> <容器路径>
ADD [--chown=<user>:<group>] <源路径> <容器路径>
```



**VOLUME**：定义匿名数据卷。在启动容器时忘记挂载数据卷，会自动挂载到匿名卷

```dockerfile
VOLUME ["<路径1>", "<路径2>"...]
VOLUME <路径>
```



**WORKDIR** 指定工作目录。指定的工作目录，必须是提前创建好的。

```dockerfile
WORKDIR <工作目录路径>
```



## 2.2 镜像制作

`docker build`

```shell
docker build -f based_dockerfile -d miaosha .
```



# 3、Docker Compose

## 3.1 安装

```shell
#方法一(省事)
sudo yum update
sudo yum install docker-compose-plugin


#方法二(推荐)
#https://docs.docker.com/compose/release-notes/ compose 版本（事先查好最新版本，在url替换就行）
# 版本信息在 https://docs.docker.com/compose/release-notes/
sudo curl -L "https://github.com/docker/compose/releases/download/2.16.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

#验证
docker-compose --version
```

## 3.2 docker-compose

一般使用步骤：

- 准备好dockerfile（不需要则直接使用docker hub 上的镜像）
- 新建并配置好`docker-compose.yml`
- `docker-compose up -d`构建并运行所有容器
- 出现错误则`docker-compose down`停止并删除所有相关容器

## 3.3 文件配置

`version` ：指定compose版本

`services`：

`image`：指定镜像

`build`：指定构建镜像的dockerfile

- context：上下文路径。
- dockerfile：指定构建镜像的 Dockerfile 文件名。
- args：添加构建参数，这是只能在构建过程中访问的环境变量。
- labels：设置构建镜像的标签。

```yaml
version: "3.7"
services:
  webapp:
    build:
      context: ./dir
      dockerfile: Dockerfile-alternate
      args:
        buildno: 1
      labels:
        - "com.example.description=Accounting webapp"
        - "com.example.department=Finance"
        - "com.example.label-with-empty-value"
```

`container_name`：指定容器名

`enviroment`：添加环境变量，可以使用数组或字典、任何布尔值，布尔值需要用引号引起来

```yaml
environment:
  RACK_ENV: development
  SHOW: 'true'
```

`command`：覆盖容器启动的默认命令

- shell
- exec

[执行多条命令](https://blog.csdn.net/whatday/article/details/108863389)

`entrypoint`：覆盖容器默认的 entrypoint

`depends_on`：设置依赖关系。

- docker-compose up ：以依赖性顺序启动服务。在以下示例中，先启动 db 和 redis ，才会启动 web。
- docker-compose up SERVICE ：自动包含 SERVICE 的依赖项。在以下示例中，docker-compose up web 还将创建并启动 db 和 redis。
- docker-compose stop ：按依赖关系顺序停止服务。在以下示例中，web 在 db 和 redis 之前停止。

```yaml
version: "3.7"
services:
  web:
    build: .
    depends_on:
      - db
      - redis
  redis:
    image: redis
  db:
    image: postgres
```

**注意：web 服务不会等待 redis db 完全启动 之后才启动。**

`restart`：重启策略，默认为no

- no：是默认的重启策略，在任何情况下都不会重启容器。
- always：容器总是重新启动。
- on-failure：在容器非正常退出时（退出状态非0），才会重启容器。
- unless-stopped：在容器退出时总是重启容器，但是不考虑在Docker守护进程启动时就已经停止了的容器

`volumes`：将主机的数据卷**或着文件**挂载到容器里。



`networks`：配置容器连接的网络

```yaml
services:
  some-service:
    networks:
      some-network:
        aliases:
         - alias1
      other-network:
        aliases:
         - alias2
networks:
  some-network:
    # Use a custom driver
    driver: custom-driver-1
  other-network:
    # Use a custom driver which takes special options
    driver: custom-driver-2
```

**aliases** ：同一网络上的其他容器可以使用服务名称或此别名来连接到对应容器的服务

## 3.4 服务启动

```shell
docker-compose up #启动所有服务
docker-compose up service_name ...
# 如果服务有很强的依赖关系，为了避免因为被依赖项还没启动而出现错误，可以使用这种方式启动服务
#先启动被依赖性再启动其他
```

# 4. 使用踩过的坑

【1】 容器里，pid=1的进程被kill后容器会自动退出。pid=1是容器进入容器允许的第一个进程，有可能是

- docker run 指定的启动进程
- dockerfile里的 RUN 、ENTRYPOINT等
- dockercompse 里的 command

【2】启动容器时没有指定任何启动进程，容器启动后会自动退出，这种情况下可以配置 <font color=red>`tty=true` </font>让容器不会自动退出



# 5. 常见容器的参数配置和使用

