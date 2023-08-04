# 1  环境安装

## 1.1 protoc

从 [Protobuf Releases](https://github.com/protocolbuffers/protobuf/releases)下载最新版发布包安装。

```shell
# 下载安装包
$ wget https://github.com/protocolbuffers/protobuf/releases/download/v21.12/protoc-21.12-linux-x86_64.zip

# 解压到/usr/local目录下（其他目录也可以）
$ sudo 7z x protoc-21.12-linux-x86_64.zip -o/usr/local
```

解压后在`/usr/local`目录下有两个子目录 `bin`和`include`，并将该`bin`加入到环境变量中。

```shell
# 打开当前用户配置文件
$ vim ~/.profile

# 添加环境变量到.profile
export PATH=$PATH:/usr/local/bin
# 执行生效
$ source ~/.profile
```

如果能正常显示版本，则表示安装成功。

```go
$ protoc --version
libprotoc 3.21.12
```



## 1.2 两个工具 protoc-gen-go & protoc-gen-go-grpc

`protoc-gen-go`和`protoc-gen-go-grpc`将`.proto`文件生成golang消息以及grpc客户端和服务端代码。

```sh
go get -u google.golang.org/protobuf/cmd/protoc-gen-go
go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
#or
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

两个命令分别将`proto-gen-go`和`proto-gen-go-grpc`安装到`$GOPATH/bin`目录下，所以也务必将该`bin`目录加入到环境变量中。



# 2 定义消息类型

## 2.1  .proto

[proto3]([Language Guide (proto3)  | Protocol Buffers  | Google Developers](https://developers.google.com/protocol-buffers/docs/proto3#simple))

Protobuf 在 `.proto` 定义需要处理的结构化数据，可以通过 `protoc` 工具，将 `.proto` 文件转换为 C、C++、Golang、Java、Python 等多种语言的代码，兼容性好，易于使用。

### 2.1.1 常用类型



- 协议版本、包目录、包名

  ```protobuf
  syntax="proto3";
  option go_package="grpc/demo/service";
  package service;
  ```

  

- string, int, float double, bool, bytes等

- 数组

  关键字：`repeated`

  ```protobuf
  repeated string cards =1
  ```

  

- 枚举类型

  ```protobuf
  enum Corpus {
    CORPUS_UNSPECIFIED = 0;
    CORPUS_UNIVERSAL = 1;
    CORPUS_WEB = 2;
    CORPUS_IMAGES = 3;
    CORPUS_LOCAL = 4;
    CORPUS_NEWS = 5;
    CORPUS_PRODUCTS = 6;
    CORPUS_VIDEO = 7;
  }
  message SearchRequest {
    string query = 1;
    int32 page_number = 2;
    int32 result_per_page = 3;
    Corpus corpus = 4;
  }
  ```

  

- 保留类型：被保留的序号或者值不能使用

  ```protobuf
  enum Foo {
    reserved 2, 15, 9 to 11, 40 to max;
    reserved "FOO", "BAR";
  }
  ```

- 引用其他消息

  同一package直接使用；引用其他package的消息可使用`import`。

  要使用protobuf官方的proto文件，暂时的方法只能将`/usr/local/include/google`复制到当前工作目录（并且需要授予读取权限1）

- 嵌套类型

  ```protobuf
  message SearchResponse {
    message Result {
      string url = 1;
      string title = 2;
      repeated string snippets = 3;
    }
    repeated Result results = 1;
  }
  ```

  

- Any

  Any 可以表示不在 .proto 中定义任意的内置类型。

  ```protobuf
  import "google/protobuf/any.proto";
  
  message ErrorStatus {
    string message = 1;
    repeated google.protobuf.Any details = 2;
  }
  ```

- Oneof

- Map

  ```protobuf
  // map<key_type, value_type> map_field = N;
  map<string, Project> projects = 3;
  ```

  

### 2.1.2 protoc代码生成



```shell
protoc -I=PROTO_DIR --go_out=OUT_DIR --go-grpc_out=OUT_DIR *.proto
```

- -I相当于`proto_path`，指定proto的目录，不指定则使用当前目录
- `go_out`产出Message的代码
- `go-grpc_out`产出Service以及GRPC客户端服务端代码



# 3 GRPC通信

四种通信方式：

| 服务类型     | 特点                                                         | 应用举例 |
| ------------ | ------------------------------------------------------------ | -------- |
| 单次模式     | 一次rpc调用                                                  | 一般场景 |
| 服务端流模式 | 客户端发送一次请求，服务端以流的形式回复多条消息，直到主动关闭 | 下载文件 |
| 客户端流模式 | 客户端以流的形式发送多次请求，直到该流关闭后服务端回复一条消息 | 上传文件 |
| 双向流模式   | 双端都持有流，同时接收和发送                                 | 实时聊天 |

`proto`代码

```protobuf
service HelloService {
  rpc SayHello1 (HelloRequest) returns (HelloResponse);
  // 客户端流式 RPC
  rpc SayHello2 (stream  HelloRequest) returns (HelloResponse);
  // 服务端流式 RPC
  rpc SayHello3 (HelloRequest) returns (stream HelloResponse);
  // 双向流式 RPC
  rpc SayHello4 (stream  HelloRequest) returns (stream HelloResponse);

```

## 3.1 应用示例

`demo.proto`

```protobuf
syntax="proto3";
option go_package="grpc/demo/service";
package service;

enum gender{
  male=0;
  female=1;
}
message Profile{
  string username =1;
  string password =2;
}
//基本类型、枚举类型、自定义类型、数组、字典、任意类型、其他类型（时间）
message DemoRequest{
  string name =1;
  int32 age=2;
  gender gd=3;
  Profile pf =4;
  repeated int32 cards=5;
}

message DemoReply{
  string msg=1;
}
service Echo{
  rpc EchoHello1 (DemoRequest) returns (DemoReply);
  rpc EchoHello2 (DemoRequest) returns (stream DemoReply);
  rpc EchoHello3 (stream DemoRequest) returns (DemoReply);
  rpc EchoHello4 (stream DemoRequest) returns (stream DemoReply);
}
```



`server.go`

服务端的main函数

```go
package main

func main() {
	//1 监听
	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatal("listen failed\n", err.Error())
	}
	//2 ssl认证
	cred, credErr := credentials.NewServerTLSFromFile("/home/ecs-user/root/cert/cert.crt", "/home/ecs-user/root/cert/private.key")
	if credErr != nil {
		log.Fatal("access ssl failed\n", credErr.Error())
	}
	//3 初始化安全的grpc服务端并注册消息服务
	s := grpc.NewServer(grpc.Creds(cred))
	service.RegisterEchoServer(s, &Server{})

	//4 开启grpc服务
	if grpcErr := s.Serve(lis); grpcErr != nil {
		log.Fatal(grpcErr)
	}
}

//服务端还需要定义一个结构体并实现XXXServer接口
type Server struct {
	service.UnimplementedEchoServer
}
```

`XXXServer`接口

```go
type EchoServer interface {
	EchoHello1(context.Context, *DemoRequest) (*DemoReply, error)
	EchoHello2(*DemoRequest, Echo_EchoHello2Server) error
	EchoHello3(Echo_EchoHello3Server) error
	EchoHello4(Echo_EchoHello4Server) error
	mustEmbedUnimplementedEchoServer()
}
```





`client.go`

客户端的main函数

```go
package main

func main() {

	//1 获取ssl证书
	cred, err := credentials.NewClientTLSFromFile("/home/ecs-user/root/cert/cert.crt", "*.peadx.live")
	if err != nil {
		log.Fatal("access ssl failed\n", err)
	}
	//2 呼叫服务端
	conn, dialErr := grpc.Dial(":9090", grpc.WithTransportCredentials(cred))
	if dialErr != nil {
		log.Fatal("dial failed\n", dialErr)
	}

	//3 初始化grpc客户端
	echoClient := service.NewEchoClient(conn)

	//4 选择要远程调用的服务
	//echo2(echoClient)
	echo3(echoClient)

}
```



### 3.1.1 Simple GRPC

客户端发送一个远程调用请求，而后接收一个回复

服务端接收一个请求，并回复一个结果

```go
//client.go
func echo1(echoClient service.EchoClient)  {
	reply, err := echoClient.EchoHello1(context.Background(), makeDemoRequest(25))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("recv %s", reply.Msg)
}

//server.go
func (s *Server) EchoHello1(ctx context.Context, req *service.DemoRequest) (*service.DemoReply, error) {
	fmt.Printf("accpet  %s", req)
	return &service.DemoReply{Msg: "hello"}, nil
}
```





### 3.1.2  Server Stream

客户端发送一个远程调用请求，而后从流中读取所有服务端的回复

服务端接收到一个请求，并将结果以流的形式多次发送

```go
//client.go
func echo2(echoClient service.EchoClient) {
	stream, callErr := echoClient.EchoHello2(context.Background(), makeDemoRequest(24))
	if callErr != nil {
		log.Fatal("call failed\n", callErr)
	}
	count := 0
	for {
		reply, recvErr := stream.Recv()
		if recvErr != nil {
			if recvErr == io.EOF {
				log.Printf("recv all %d reply", count)
				//_ = stream.CloseSend()
				break
			}
			log.Fatal(recvErr)
		}
		fmt.Printf("recv %s", reply.Msg)
		count++
	}
}

//server.go
func (s *Server) EchoHello2(req *service.DemoRequest, stream service.Echo_EchoHello2Server) error {
	fmt.Printf("accpet  %s", req)
	res := []service.DemoReply{
		{Msg: "reply created"},
		{Msg: "reply created"},
		{Msg: "reply created"},
		{Msg: "reply created"},
	}

	for _, data := range res {
		stream.Send(&data)
	}
	return nil
}
```



### 3.1.3 Client Stream

客户端以流的形式发送多个请求，而后接收服务端一个结果回复

服务端从流中读取所有的请求，并回复一个结果

```go
//client.go
func echo3(echoClient service.EchoClient) {
	stream, err := echoClient.EchoHello3(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	age := int32(24)
	for age <= 27 {
		err := stream.Send(makeDemoRequest(age))
		if err != nil {
			continue
		}
		age++
	}
	reply, recvErr := stream.CloseAndRecv()
	if recvErr != nil {
		log.Fatal(err)
	}
	fmt.Printf("end %s\n", reply.Msg)

}

//server.go
func (s *Server) EchoHello3(stream service.Echo_EchoHello3Server) error {

	count := 0
	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				stream.SendAndClose(&service.DemoReply{Msg: fmt.Sprintf("receive %d msg", count)})
				return nil
			}
			return err
		}
		log.Printf("reveive req:%s", req)
		count++
	}
}
```



### 3.1.4 Two-Sides Stream

（抢占式使用同一个流对象）

双向流的互动形式可以很多样

- 可以交替使用流对象，一发一回，类型聊天



# 4 GRPC原理

## 4.1 Protocol Buffers

[编码方式](https://blog.csdn.net/daaikuaichuan/article/details/105639884)



## 4.2 HTTP2



