# AMQP(高级消息队列协议)

可以简单理解为是一个带路由的生产者消费者模型

## 交换机

- Name
- Durability （消息代理重启后，交换机是否还存在）
- Auto-delete （当所有与之绑定的消息队列都完成了对此交换机的使用后，删掉它）
- Arguments（依赖代理本身）

### 默认交换机

### 直连交换机

- 单播路由

- 消息和队列通过交换机的一个路由键绑定

### 扇形交换机

- 广播路由
- 消息和队列通过交换机整机绑定

### 主题交换机

- 多播路由

- 消息和队列通过交换机的一个主题绑定

  相当于直连是精致匹配，主题是模糊匹配

### 头交换机

- 多关键字路由
  1. 不使用路由键，多关键字信息存放与消息头部
  2. 支持两种匹配模式：all、any
  3. 不局限于字符串，可以整数或者字典



## 队列

<font color=red>队列相当于生产线，生产者在队列上游，消费者在队列下游，生产者和消费者共同操作队列对象</font>

- Name
- Durable
- Exclusive 独占队列只能由声明它们的连接访问，并且在连接关闭时将被删除。 当尝试声明、绑定、使用或删除同名队列时，其他连接上的通道将收到错误
- Auto-delete （当最后一个消费者退订后队列被删除）
- Arguments

### 队列名

- 自定义或代理默认
- 通道能记住当前最后一次使用的队列名，后序可以使用空字符代表
- 使用`amq`开头的队列名会报错

### 队列持久化

代理重启是否能恢复队列上也同样开启了持久化的消息

### 队列绑定

将队列通过什么路由键绑定到什么类型的交换机上

<font color=green>一个队列可以绑定多个路由键</font>，类似于一个IP能绑定多个域名

如果某交换机没有被任何队列绑定或者消息到达交换机是找不到满足的队列，则消息可能会被销毁或者<font color=green>返还发布者</font>

## 消息

### 内容

**属性**

- 类型与编码
- 路由键 
- 是否持久化
- 优先权
- 时间戳与有效期
- 发送者id

**负载**

消息内容本身，代理不能查看其内容，只负责转发

### 回执

#### 消息确认

- 不开启确认机制或者说自动确认模式，消息从队列发送出后立即删除

- 开启确认机制或者说主动确认模式

  ACK时机

  - 收到即ACK
  - 本地备份后ACK
  - 处理完才ACK

- 如果消费者没来得及ACK就挂掉了，则该消费者已经拿到的消息，队列会将其发送给别的消费者；如果此时已经没有消费者了则会得到有新的消费者加入再投送



#### 拒绝消息

当一个消息到达消费者手里的时候，消费者认为此时还是消费该消息的最好时机又或者消息本身有问题则可以选择拒绝该消息。被拒绝的消息，消费者可以指定该消息是继续留在队列中还是立即销毁

#### negative acknowledgements

？？？

### 预取

<font color=red>在一个队列有多个消费者的情境下</font>

可以控制流入每一个消费者的速度从而达到负载均衡和提高吞吐量的作用

`Qos controls how many messages or how many bytes the server will try to keep on the network for consumers before receiving delivery acks. The intent of Qos is to make sure the network buffers stay full between the server and client.`

## 消费者

- 队列里的消息主动推送到消费者手里
- 消费者自己跑过来问队列拿消息
- 一个队列可以注册多个平等消费者或者注册一个独享型消费者





## 连接

### 连接

- TCP长连接
- SSL保护
- 优雅关闭

### 通道

- 一个通道是一个虚拟连接
- 多个通道基于同一个·TCP长连接
- 通道需要通道号区分彼此，因为都在同一个TCP连接中发送与接收
- 通道与通道之间是隔离的

### 虚拟主机

- 一个主机包含用户、用户组、交换机、队列等资源

- 内核资源（装了AMQP环境的一台机器）只有一份，通过虚拟主机（addr/vhost_i）来做到环境隔离

- 类似于对于同一个ip，可以通过不同的域名映射来做到虚拟化

  又或者对于同一个ip下的同一个域名，可以通过不同的端口做到虚拟化





# 客户端实战（golang版）

## 基础

**安装**

```shell
go get -u github.com/streadway/amqp
```

**建立连接**

<font color=red>`amqp.Dial`</font>

TCP长连接

```go
url:="amqp://guest:guest@localhost:5672/vhostName"
//协议://主机地址：通信端口|/虚拟主机名
conn,err:=amqp.Dial(url)
failError(err,"failed to connect to rabbitmq")
defer conn.Close()
```

**建立通道**

<font color=red>`conn.Channel`</font>

基于通道通信，多个通道共用一个TCP连接

```go
ch,err:=conn.Channel()
failError(err,"failed to open a channel")
defer ch.Close()
```

**声明交换机**

<font color=red>`ch.ExchangeDeclare`</font>

函数签名

- name：交换机名字
- kind：交换机类型。支持四种
  - "direct", 
  - fanout",
  - "topic" 
  -  "headers"
- durable：是否开启持久化
- autoDelete：是否自动删除

```go
/*签名*/
func (ch *Channel) ExchangeDeclare(name string, kind string, durable bool, autoDelete bool, internal bool, noWait bool, args Table) err

err:=ch.ExchangeDeclare("fanout_logs", "fanout", false, false, false, false, nil) //internal noWait args 默认为false或nil
```



**声明队列**

<font color=red>`ch.QueueDeclare`</font>

函数签名

- name：队列名字
- durable
- autoDelete
- exclusive：是否为独占队列

```go
/*签名*/
func (ch *Channel) QueueDeclare(name string, durable bool, autoDelete bool, exclusive bool, noWait bool, args Table) (Queue, error)
```



**绑定队列到交换机**

<font color=red>`ch.ExchangeDeclare`</font>

函数签名

- name：即将被绑定的队列名
- key：路由键
- exchange：交换机名字。空字符串使用默认交换机

```go
/*签名*/
func (ch *Channel) QueueBind(name string, key string, exchange string, noWait bool, args Table) error
```



**推送信息**

<font color=red>`ch.Publish`</font>

函数签名

- exchange：推送到哪个交换机上
- key：消息绑定的路由键
- msg：消息本身，使用`Publishing`结构体表示，重要字段包含有：
  - Headers：头交换机需要用到
  - ContentType：一般设为` "text/plain"`
  - Body：序列化后的字节数组
  - UserId & AppId：发送方身份
  - ReplyTo：回调队列名，可以实现接收消费者返回的结果



```GO
/*签名*/
func (ch *Channel) Publish(exchange string, key string, mandatory bool, immediate bool, msg Publishing) error

type Publishing struct {
    Headers         Table
    ContentType     string
    ContentEncoding string
    DeliveryMode    uint8
    Priority        uint8
    CorrelationId   string
    ReplyTo         string
    Expiration      string
    MessageId       string
    Timestamp       time.Time
    Type            string
    UserId          string
    AppId           string
    Body            []byte
}
```



**消费信息**

<font color=red>`ch.Consumer`</font>

函数签名

- queue：消费哪个队列上的信息
- consumer：消费者名字
- **autoAck**：当设为false时，在接收到消息后需要找准时机向服务端回应
- <-chan Delivery：接收消息的通道。拿到该通道后需要for range 读取通道消息并消费

```go
func (ch *Channel) Consume(queue string, consumer string, autoAck bool, exclusive bool, noLocal bool, noWait bool, args Table) (<-chan Delivery, error)
```

**回应**

<font color=red>`msg.Ack`</font>

`multiple`：除了回应当前信息，是否连同同一通道上尚未回应的消息也一同回应

```go
func (d Delivery) Ack(multiple bool) error
```



**拒绝**

<font color=red>`msg.Reject`</font>

`requeue`：

- true：拒绝的消息可以发送给其他队列
- false：拒绝的消息马上销毁

```go
func (d Delivery) Reject(requeue bool) error
```



**吞吐量**

<font color=red>`msg.Qos`</font>



## 实战

### 工作队列

经典的生产者消费者模型

路由不重要，因为只需要一个作为中转

使用直连或者默认交换机



### 发布订阅

典型的广播模型

生产者生产的消息需要发送给所有的订阅者

使用扇形交换机

生产者的消息便会广播到所有绑定的队列



### 路由

路由是一个典型的点到点模型

所以顺其自然使用直连交换机

同时，如果一个队列希望能接收不同路由键的消息，可以绑定多个路由键，类似一个ip有多个域名，对这些域名的请求都会被这个ip接收到。

**但是，如果希望一个队列能接收一类消息，而这类消息的不同具体类型是不定的甚至是无限的，这时直连交换机就无法满足**



### 主题交换机

队列的绑定路由键是一个**正则表达式**

格式：<font color=red>word1.word2...wordn</font>

每个单词用点号隔开

- `*` (星号) 用来表示一个单词
- `#` (井号) 用来表示任意数量（零个或多个）单词



### RPC

生产者发送的消息可以指定回调队列名

但是每一个消息都新建一个回调队列成本太高

所以可以使用一个共用的回调队列，但引入新的问题

不能区分回调的信息属于哪个rpc请求的，所以可以为每个rpc请求设置唯一序列号





步骤：

- 当客户端启动的时候，它创建一个匿名独享的回调队列。
- 在RPC中，客户端发送带有两个属性的消息：一个是设置回调队列的 *reply_to* 属性，另一个是设置唯一值的 *correlation_id* 属性。
- 将请求发送到一个 *rpc_queue* 队列中。
- RPC工作者（又名：服务器）等待请求发送到这个队列中来。当请求出现的时候，它执行他的工作并且将带有执行结果的消息发送给*reply_to*字段指定的队列。
- 客户端等待回调队列里的数据。当有消息出现的时候，它会检查*correlation_id*属性。如果此属性的值与请求匹配，将它返回给应用。