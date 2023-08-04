# RPC

## rpc协议是什么？

RPC(Remote Procedure Call Protocol)远程调用

### RPC五要素

- **协议格式**：<font color=blue>消息的数据结构</font>。例如HTTP协议规定了请求和响应的头部、正文、请求方法、状态码等等内容。同样，RPC协议中，一般请求需要包含<font color=Chocolate>方法名和参数列表</font>，响应需要包含<font color=Chocolate>返回值和异常信息</font>等等内容。
- **数据编码**：RPC协议需要对传输的数据进行编码和解码，以便客户端和服务端能正确地理解和处理数据。 比如GRPC的`Protocol Buffers`。
- **传输协议**：常用有TCP、UDP、HTTP等。实现细节上还有考虑网络IO模型。
- **接口定义**：RPC协议需要定义一组接口，用户描述客户端和服务端之间的交互方式。接口通常包括服务名、方法名、参数列表、返回值类型等等。
- **序列化和反序列化**：有了编码方式还需要有相对应的编码解码器。客户端需要将请求序列化为指定的编码格式；服务端接收到请求数据也需要反序列化为具体类型的值才能正确理解和处理请求。

### 为什么要用RPC

业务发展路线：<font color=Blue>单体架构-->非关联业务拆分-->业务必须交互</font>

单体架构将所有的业务都部署在同一个服务器上。随着业务流量增长，开始将互不关联的业务分别抽取作为单独的服务。业务逻辑越来越复杂，必不可少需要做服务与服务之间的交互，于是RPC协议就能派上用场了。

## RPC原理

### RPC调用流程

![](https://ask.qcloudimg.com/http-save/yehe-8197675/c97ab362cba7c866a9ce9f84c663c0b0.png?imageView2/2/w/1620)

1. client以本地调用的方式调用远程服务
2. client stub 接收到调用后负责将方法、参数、调用ID等组装并编码成能够进行网络传输的消息体
3. client stub 找到服务地址，并将消息发送给服务器
4. server stub 接收到消息后进行解码
5. server stub 根据解码结果调用本地的服务
6. 本地服务执行后将结果返回给server stub
7. server stub 将返回结果进行编码并发送到消费方
8. client stub 接收到消息进行解码
9. client得到最终结果

### 确定消息的数据结构

**客户端**的请求消息结构一般需要包括以下内容：

- 接口名称
- 方法名
- 参数类型&参数值
- 超时时间
- requestID（请求的唯一标识）

**服务端**返回的消息结构一般包括以下内容：

- 返回值
- 错误消息
- requestID

### 服务发现

**为什么需要服务发现?**

- <font color=blue>帮助分布式系统定位寻</font>：客户端在上线时，需要知道有哪些服务器可以正常提供服务。
- 实时感知服务器提供方的状态，<font color=blue>监控服务提供方的上下线</font>；为了高可用性，在生产环境中服务提供方通常以集群的方式对外提供服务。集群里的服务IP可能会变化，客户端需要及时获取可用的服务结点信息。

<font color=red>服务发现的本质，就是完成接口跟服务提供商IP之间的映射</font>

**怎么实现服务发现**

答案是<font color=blue>注册中心</font>



**nginx的服务发现和注册中心**

<font color=red>nginx</font>是一个反向代理组件，**nginx需要知道应用服务器的地址是什么**，这样才能将请求透传到应用服务器上，这其实是一个**服务发现**的过程。

那么nginx是怎样实现的呢？<font color=blue>它把应用服务器的地址配置在了文件中</font>。这样做的缺点是：

- 需要紧急扩容时，需要停止服务、修改nginx配置再重启所有的客户端进程
- 一旦有服务器出现故障，也需要修改客户端的配置后重启。服务修复后也无法自动恢复，也需要修改配置重启。



注册中心能解决以上问题。目前可使用的注册中心有ZooKeeper、ETCD、Nacos等

**注册中心的基本功能有两点**

- 提供服务地址的存储
- <font color=red>当存储内容发生改变时，可以将变更的内容推送给客户端</font>

#### RPC中的服务发现

对于RPC框架，注册中心应该怎么组织服务信息，客户端才能按需拿到可以服务的服务器地址呢？

<font color=blue>接口名作为key，服务地址集合作为value</font>

下图描述了服务注册与发现的过程：

![](https://img-blog.csdnimg.cn/a2ef0f89848a45a69233cabd9d991845.png?x-oss-process=image/watermark,type_ZHJvaWRzYW5zZmFsbGJhY2s,shadow_50,text_Q1NETiBAT2NlYW4mJlN0YXI=,size_12,color_FFFFFF,t_70,g_se,x_16)

- 服务端启动后向注册中心注册服务，并定时发送心跳
- 客户端拿着接口名去注册中心索取可服务地址列表
- 拿到地址后就可以向服务器发送远程调用了

# GRPC

## Protocol Buffers

GRPC使用的编码方式。protocol buffers 通过<font color=blue>IDL(接口描述语言)</font>定义消息结构，使用<font color=blue>二进制格式</font>进行编解码

[编码方式](https://blog.csdn.net/daaikuaichuan/article/details/105639884)

- `protobuf`将消息里的每个字段进行**编码**后，再利用**T-L-V或者T-V的方式进行数据存储**。
- `protobuf`对于**不同类型的数据会使用不同的编码和存储方式**。
- `protobuf`的**编码和存储方式是其性能优越、数据体积小的原因**。

![](https://img-blog.csdnimg.cn/20200420225609764.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L2RhYWlrdWFpY2h1YW4=,size_16,color_FFFFFF,t_70)

### Varints编码

Varints编码**使用<font color=red>msb</font>表示当前是否是最后一个字节**，如果**msb为1表示后面还有字节，如果msb为0表示当前字节是最后一个字节**。因为加入了msb，所以**每一个字节低七位表示具体的数值**。

<font color=blue>【**Varints编码过程如下**】：</font>

1. 先将整数转换为不带前导0的二进制表示：<font color=red>`100101100`</font>。
2. 每一个字节的低七位存储数值：<font color=red>`0000010 0101100`</font>。
3. 转换为小端模式：<font color=red>`0101100 0000010`</font>。
4. 添加`msb`标志位：<font color=red>`10101100 00000010`</font>。

### ZigZag编码

以负数-11为例，其二进制在计算机中是用补码表示的，整数原码为：<font color=red>`00001011`</font>，反码为：<font color=red>`11110100`</font>，补码（反码加1）为：<font color=red>`11110101`</font>。

<font color=blue>【ZigZag编码过程如下】：</font>

原数左移一位（左移时丢弃最高位，低位补0）：<font color=red>`11101010`</font>。
原数右移31位（右移时符号位不变，高位补符号位）：<font color=red>`11111111`</font>。
将上述两个二进制数取异或：<font color=red>`00010101`</font>。

### Length-delimited编码

```protobuf
message Test2 {
  required string b = 2;
}
// 设置b的值为"testing"，则编码之后为：
12 07 74 65 73 74 69 6e 67
```



![](https://img-blog.csdnimg.cn/20200420222844913.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L2RhYWlrdWFpY2h1YW4=,size_16,color_FFFFFF,t_70)

###  存储方式

![](https://img-blog.csdnimg.cn/20200420223045164.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L2RhYWlrdWFpY2h1YW4=,size_16,color_FFFFFF,t_70)

## HTTP2

