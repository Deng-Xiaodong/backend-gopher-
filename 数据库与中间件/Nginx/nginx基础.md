**Nginx常用功能**

- 正向代理
- 反向代理
- 负载均衡
- web缓存

# 一、功能

## 1、正向代理

`正向代理`是客户端知道要访问的最终目标服务器对象。但因为无法直接访问，所以将正向代理服务器作为跳板（中间人）去访问目标服务器。



**按代理是否解密htps报文分类**

- **隧道代理**：也就是透传代理。代理服务器只是在TCP协议上透传HTTPS流量，对于其代理的流量的具体内容不解密不感知。客户端和其访问的目的服务器做直接TLS/SSL交互。

- **中间人代理**：代理服务器解密HTTPS流量，对客户端利用自签名证书完成TLS/SSL握手，对目的服务器端完成正常TLS交互。在客户端-代理-服务器的链路中建立两段TLS/SSL会话。

  

## 2、反向代理

反向代理是客户端不知道目标服务器的情况下，通过向反向代理服务器发送不同的请求从而获取不同的服务。客户端只知道自己可以想反代理服务器拿东西，但是反代理服务器从哪里拿对客户端是透明的。一般下，不同的请求路径、不同的域名、不同的端口代表请求不同的服务。



![](https://www.runoob.com/wp-content/uploads/2018/08/1535725078-5993-20160202133724350-1807373891.jpg)



## 3、负载均衡

Nginx提供的负载均衡策略有2种：内置策略和扩展策略。内置策略为轮询，加权轮询，Ip hash。扩展策略，就天马行空，只有你想不到的没有他做不到的啦，你可以参照所有的负载均衡算法，给他一一找出来做下实现。

![](https://www.runoob.com/wp-content/uploads/2018/08/1535725078-1224-20160201162405944-676557632.jpg)

Ip hash算法，对客户端请求的ip进行hash操作，然后根据hash结果将同一个客户端ip的请求分发给同一台服务器进行处理，可以解决session不共享的问题。



## 4、web缓存

Nginx可以对不同的文件做不同的缓存处理，配置灵活，并且支持FastCGI_Cache，主要用于对FastCGI的动态程序进行缓存。配合着第三方的ngx_cache_purge，对制定的URL缓存内容可以的进行增删管理。



# 二、安装

通过docker安装

```sh
docker run --name some-nginx -p 443:443 
-v ./conf/nginx.conf:/etc/nginx/nginx.conf 
-v ./conf/ssl:/etc/nginx/ssl 
-v ./data/nginx:/usr/nginx/html 
-v ./logs/nginx:/var/log/nginx
-d nginx
```

- 使用最新的nginx镜像；
- 开放容器443端口与宿主机443端口映射；

- 做了四个用途的卷挂载（主配置、https证书配置、静态文件、日志）



# 三、配置

通过docker安装的nginx的默认主配置文件路径：

`/etc/nginx/nginx.conf`

**主配置中的六个主要模块：**

- 全局
- events
- http
- server
- location
- upstream

## 1、全局块

配置头到events之间，配置影响全局。

主要有以下可配置项：

- user ：nginx服务的用户（组）（不写默认为nobody）
- worker_process ：工作进程数；一般置为核心数
- 进程PID存放路径
- 错误日志存放路径及格式自定义
- 配置文件导入



```sh
# 用户
user nobody
# 工作进程数
worker_process 2
# 日志位置；应与日志挂载卷一致
error_log  /var/log/nginx/error.log notice/info/debug
#PID文件位置
pid  /var/log/nginx/nginx.pid
```



## 2、events块

events块配置主要影响用户与nginx服务器的网络连接。

常用配置有：

- 是否开启对多worker_process下的网络连接进行序列化
- 是否允许同时接收多个网络请求
- 选取哪种事件驱动模型
- 每个worker_process可以同时支持的最大连接数

```shell
 events  {
         use epoll;    #使用epoll模型。2.6及以上版本的系统内核，建议使用epo11模型以提高性能，实现I/O多路复用
         worker_connections  4096;  #每个工作进程处理4096个连接。默认值为1024。一般设置为2的次方
 }
 
 #epoll是Linux内核为处理大批句柄而作改进的poll，是Linux 下多路复用IO接口select/poll的增强版本，它能显著的减少程序在大量并发连接中只有少量活跃的情况下的系统CPU利用率。
 
 #如提高每个进程的连接数还需执行“ulimit -n 65535”命令临时修改本地每个进程可以同时打开的最大文件数。
 #在Linux平台上，在进行高并发TCP连接处理时，最高的并发数量都要受到系统对用户单一进程同时可打开文件数量的限制(这是因为系统为每个TCP连接都要创建一个socket句柄，每个socket句柄同时也是个文件句柄)。
 #可使用"ulimit -a"命令查看系统允许当前用户进程打开的文件数限制。
```



## 3、http块

http块的配置关系如下：

- http
  - http全局
  - server
    - server全局
    - location



**http全局**

文件引入、MIME-TYPE定义、日志自定义、连接超时时间、文件传输模式

**server**

这块和虚拟主机有密切关系，虚拟主机从用户角度看，和一台独立的硬件主机是完全一样的。每个 http 块可以包括多个 server 块，而每个 server 块就相当于一个虚拟主机。

**server全局**

$\color{green}监听端口$ ：如果有多个server块，则有可能会监听多个端口；一般为80或者443

$\color{green}代理服务器的域名或ip地址$ ：一般写服务器的域名或者IP地址；如果nginx部署在docker容器，也可以是容器名或者容器ip地址



**location**

一个 server 块可以配置多个 location 块。主要作用是根据请求地址路径的匹配，匹配成功进行特定的处理。对特定的请求进行处理。地址定向、数据缓存和应答控制等功能，

location常见配置指令：root、alias、proxy_pass

- root（根路径）
-  alias（虚拟路径）
- proxy_pass （反向代理配置）

```shell
 location / {
     root html;    #表示此时根目录为/usr/nginx/html 
     index index.html index.htm;
 }
 
 #此时访问路径和返回文件的关系为：
 http://192.168.72.10/index.html      --> /usr/nginx/html/index.html

 
 location /test {   
     root data;    #根目录，表示此时根目录为/data/
     index index.html index.htm；;
 }
 http://192.168.72.10/test/index.html --> /data/test/index.html
 
 location /test{
     alias /var/www/html;    #设置别名，即虚拟路径。别名是一个整体。  
     index index.html index.htm;
 }
 http://192.168.72.10/test/index.html --> /var/www/html/index.html

```



http配置示例：

```shell
 http {
   ##文件扩展名与文件类型映射表
   include mime.types;
   ##默认文件类型
   default_type  application/octet-stream;
   ##日志格式设定
   #log_format main '$remote_addr - $remote_user [$time_local] "$request" '
   #      '$status $body_bytes_sent "$http_referer"
   #     '"$http_user_agent" "S$http_x_forwarded_for"' ;
   ##访问日志位置
   #access_log logs/access.1og main;
   ##开启文件传输模式
   sendfile  on;
   ##减少网络报文段的数量
   #tcp_nopush  on;
   ##连接保持超时时间，单位是秒
   #keepalive_timeout 0;
   keepalive_timeout 65;
   ##gzip模块设置，设置是否开启gzip压缩输出
   #gzip  on;
   ##Web服务的监听配置
   server {
       ##监听地址及端口
       listen 80;
       ##站点域名，可以有多个，用空格隔开
       server_name www.yuji.com;
       ##网页的默认字符集
       charset utf-8;
       ##根目录配置
       location / {
           ##网站根目录的位置/usr/1ocal/nginx/html
           root html ;
           ##默认首页文件名
           index index.html index.php;
       }
       ##内部错误的反馈页面
       error_page 500 502 503 504 /50x.html;
       ##错误页面配置
       location = /50x.html {
            root html ;
            }
       }
   }
   
 ------------------------------------------------
 日志格式设定:
 $remote_addr与$http_x_forwarded_for 用以记录客户端的ip地址。
 $remote_addr：记录上一个请求消息发送端的IP（代理服务器的IP）。
 $http_x_forwarded_for ：会记录所有经过的服务器的IP地址。
 $remote_user：用来记录客户端用户名称。
 $time_local：用来记录访问时间与时区。
 $request：用来记录请求的url与http协议。
 $status：用来记录请求状态;成功是200。
 $body_bytes_sent：记录发送给客户端文件主体内容大小。
 $http_referer：用来记录从哪个页面链接访问过来的。
 $http_user_agent: 记录客户浏览器的相关信息。
 ​
 通常web服务器放在反向代理的后面，这样就不能获取到客户的IP地址了，通过$remote_add拿到的IP地址是反向代理服务器的iP地址。反向代理服务器在转发请求的http头信息中，可以增加x_forwarded_for信息，用以记录原有客户端的IP地址和原来客户端的请求的服务器地址。
```



# 三、进阶使用



## 1、 nginx用户（组）配置

**添加nginx用户,且不允许该用户登录**

```shell
useradd nginx -s /sbin/nologin -M
```

**添加nginx用户组，并把nginx用户加入nginx用户组**

```shell
groupadd nginx
usermod -G nginx nginx
```

**`nginx.conf`配置**

```shell
user nginx nginx
```



## 2、隐藏响应头的版本号

**配置文件中关闭版本号**

$\color{green}server\_tokens$

```shell
http {
     include       mime.types;
     default_type  application/octet-stream;
     server_tokens off;      #添加这一行，关闭版本号
     ......
 }
```



## 3、修改缓存时间

当Nginx将网页数据返回给客户端后，可设置缓存的时间，以方便在日后进行相同内容的请求时直接返回，避免重复请求，加快了访问速度。

一般针对静态网页设置，对动态网页不设置缓存时间。

$\color{green}expires$

```shell
 http {
 ......
     server {
     ...... 
         location / {
             root html;
             index index.html index.htm;
         }
         
         #加入新的 location，以图片作为缓存对象
         location ~* \.(gif|jpg|jepg|bmp|ico)$ {
             root html;
             expires 1d;             #指定缓存时间，1天
         }
 ......
     }
 }
 ## 响应头中包含 Cahce-Control:max-age=86400 表示缓存时间是 86400 秒。
```



## 4、日志切割

随着Nginx运行时间的增加，产生的日志也会逐渐增加，为了方便掌握Nginx的运行状态，需要时刻关注Nginx日志文件。太大的日志文件对监控是一个大灾难，不便于分析排查，需要定期的进行日志文件的切割。

**操作步骤：**

1. 创建旧日志存放目录
2. 通过mv命令将原有日志移动到日志目录中
3. kill -USR1 < PID>    重新生成日志文件
4. 删除30天之前的日志文件



```shell
 1. #写脚本
 [root@yuji ~]# vim /opt/fenge.sh
 #!/bin/bash
 # Filename: fenge.sh
 # nginx日志分割，按时间分割
 ​
 #显示前一天的时间
 day=$(date -d "-1 day" "+%Y%m%d")
 #旧日志文件目录
 logs_path="/var/log/nginx"
 #nginx进程的PID
 pid_path="/usr/local/nginx/logs/nginx.pid"
 ​
 #如果旧日志目录不存在，则创建日志文件目录
 [ -d $logs_path ] || mkdir -p $logs_path
 #将日志移动到旧日志目录，并重命名日志文件
 mv /usr/local/nginx/logs/access.log ${logs_path}/tt.com-access.log-$day
 #重建新日志文件
 kill -USR1 $(cat $pid_path) 
 #删除30天之前的日志文件
 find $logs_path -mtime +30 -exec rm -rf {} ;           
 ​
 2. #赋予执行权限，执行脚本。查看日志文件目录。
 [root@yuji ~]# chmod +x /usr/local/nginx/nginx_log.sh 
 [root@yuji ~]# /opt/fenge.sh
 [root@yuji ~]# ls /var/log/nginx/            //旧日志文件已被移动到设置好的目录
 tt.com-access.log-20220516
 [root@yuji ~]# ls /usr/local/nginx/logs/     //已重建新日志文件
 access.log  error.log  nginx.pid
 ​
 3. #编写计划任务，每天定点执行
 [root@localhost nginx]#crontab -e
 0 1 * * * /opt/fenge.sh
```



## 5、设置超时时间

HTTP有一个KeepAlive模式，它告诉web服务器在处理完一个请求后保持这个TCP连接的打开状态。若接收到来自同一客户端的其它请求，服务端会利用这个未被关闭的连接，而不需要再建立一个连接。

```sh
http {
 ...... 
     keepalive_timeout 65 180;       //设置连接超时时间    
 ...... 
 }
 
 keepalive_timeout
 指定KeepAlive的超时时间（timeout）。指定每个TCP连接最多可以保持多长时间，服务器将会在这个时间后关闭连接。
 Nginx的默认值是65秒，有些浏览器最多只保持 60 秒，所以可以设定为 60 秒。若将它设置为0，就禁止了keepalive 连接。
 第二个参数（可选的）指定了在响应头Keep-Alive:timeout=time中的time值。这个头能够让一些浏览器主动关闭连接，这样服务器就不必去关闭连接了。没有这个参数，Nginx 不会发送 Keep-Alive 响应头。


```





## 6、更新工作进程数

在高并发场景，需要启动更多的Nginx进程以保证快速响应，以处理用户的请求，避免造成阻塞。

```shell
#1 查看cup核数
cat /proc/cpuinfo |grep -c processor

#2、查看nginx主进程中包含几个工作进程
 ps aux | grep nginx
 
 #3、编辑配置文件，修改工作进程数
 worker_processes  2;        #修改为与CPU核数相同
```



## 7、配置页面压缩

Nginx的ngx_http_gzip_module压缩模块提供对文件内容压缩的功能。

允许Nginx服务器将输出内容在发送客户端之前进行压缩，以节约网站带宽，提升用户的访问体验，默认已经安装。

可在配置文件中加入相应的压缩功能参数对压缩性能进行优化。

```
 #1、修改配置文件
 vim /usr/local/nginx/conf/nginx.conf
 http {
 ...... 
    gzip on;                 #取消注释，开启gzip压缩功能
    gzip_min_length 1k;      #最小压缩文件大小
    gzip_buffers 4 64k;      #压缩缓冲区，大小为4个64k缓冲区
    gzip_http_version 1.1;   #压缩版本（默认1.1，前端如果是squid2.5请使用1.0）
    gzip_comp_level 6;       #压缩比率
    gzip_vary on;            #支持前端缓存服务器存储压缩页面
    
    #压缩类型，表示哪些网页文档启用压缩功能
    gzip_types text/plain text/javascript application/x-javascript text/css text/xml application/xml application/xml+rss image/jpg image/jpeg image/png image/gif application/x-httpd-php application/javascript application/json;
 ...... 
 }

```

## 8、配置防盗链

对一些静态资源检查源主机，只允许合法访问者访问

```shell
location ~* \.(jpg|gif|swf)$ {
            root  html;
            expires 1d;
            valid_referers none blocked *.tt.com tt.com;
                if ( $invalid_referer ) {
                  rewrite ^/ http://www.tt.com/error.png;
                }
         }
```



## 9、IP白名单、黑名单

设置只有指定IP/IP段才可以访问该网页，或指定IP/IP段不能访问。想要对哪个路径进行限制，就在location块下添加路径和控制规则。

**访问控制规则如下：**

- deny IP/IP段：拒绝某个IP或IP段的客户端访问。（黑名单）
- allow IP/IP段：允许某个IP或IP段的客户端访问。（白名单）
- 规则从上往下执行，匹配到则停止，不再往下匹配。

```shell
 [root@yuji ~]# vim /usr/local/nginx/conf/nginx.conf
 ...........
       server {
            location / {
                root   html;
                index  index.html index.htm;
                ##添加控制规则
                allow  192.168.72.0/24;    #允许192.168.72.0网段的客户端ip访问
                deny   all;            #拒绝其他ip客户端访问
            }
       }
 ​
```



## 10、配置错误页

语法

```shell
error_page code ...  response_url
code ... ：错误码列表
response_url: 页面跳转；可以是本地错误页面或者任意外部页面
```

## 11、location正则匹配

```shell
location [ = | ~ | ~* | ^~] url{

}
```

- `=` 表示精确匹配。只有请求的url路径与后面的字符串完全相等时，才会命中。使用 `=` 精确匹配可以加快查找的顺序。
- `^~` 表示如果该符号后面的字符是最佳匹配（前缀匹配），采用该规则，不再进行后续的查找。
- `~`表示执行一个正则匹配，区分大小写
- `~*` 表示执行一个正则匹配，不区分大小写
- 没有修饰符表示默认的前缀匹配；多个匹配则选最佳匹配



## 12、nginx环境变量



## 13、负载均衡

[负载均衡](https://blog.csdn.net/LL845876425/article/details/97621365)

当有多个相同服务的地址，可以使用`upstream`分流

```shell
upstream test_stream {
	server 127.0.0.1;
	server 127.0.0.2;
	server 127.0.0.3;
}
server{
	listen 80;
	location / {
		proxy_pass http://test_stream
	}
}
```

server指令拥有丰富的参数，其参数说明见下表：

| 参数         | 作用                                                         |
| ------------ | ------------------------------------------------------------ |
| weight       | 请求权重                                                     |
| max_fails    | 指定时间内，请求最大失败次数。默认1；0代表无限制             |
| fail_timeout | 指定时间；跟max_fails配合使用                                |
| down         | 标记服务不可用                                               |
| backup       | 备用服务器；当所有在线服务器都请求不成功（超时）时则会使用备用服务器 |

**均衡策略**

默认使用轮询策略

| 策略    | 用途                                                         |
| ------- | ------------------------------------------------------------ |
| hash    | 根据指定key哈希选择服务器；当key设为$request_url时，可以达到提高缓存命中率的作用。因为相同的请求路径和请求参数总会分到上一次服务过的服务器 |
| ip_hash | 根据请求的p的哈希选择服务器。同一个ip请求会代理到同一后端服务器。当需要移除其中一台服务器时建议使用down对该服务器停止分流，这样可以保留当前ip_hash策略的值 |
| sticky  | 根据Cookie将请求分流到后端服务器上，同一个cookie的请求会进入同一台服务器。如果当前服务器无法服务，则会重写“绑定”一台后端服务器 |



## 14、Https配置

- 配置私秘和证书
- http重定向https

```shell
server{
 	listen 443 ssl;
    ssl_certificate /etc/nginx/cert/cert.crt;
    ssl_certificate_key /etc/nginx/cert/private.key;
    ssl_session_timeout 5m;
    ssl_protocols TLSv1 TLSv1.1 TLSv1.2;
    ssl_ciphers ECDHE-RSA-AES128-GCM-SHA256:HIGH:!aNULL:!MD5:!RC4:!DHE;
    ssl_prefer_server_ciphers on;
}
server{
 	listen 80;
    server_name peadx.live *.peadx.live;
    rewrite ^/(.*)$ https://miaosha.peadx.live:443/$1 permanent;
}
```



## 15、rewrite

[重写](https://blog.csdn.net/qq_36095679/article/details/101277202)

# 四、其他

## 1、虚拟主机

[虚拟主机](https://juejin.cn/post/7096443628326748174)
