# 1. 安装

## etcd安装

**windows**

下载解压，配置bin环境

[Releases · etcd-io/etcd (github.com)](https://github.com/etcd-io/etcd/releases/)

**测试**

```shel
etcd --version
```

# 2. 客户端使用(golang)

## package

```go
go get github.com/coreos/etcd/clientv3
```





## 获取客户端

```go
c,_ := clientv3.New(clientv3.Config{
		Endpoints: []string{"localhost:2379"},
	})
```



## 增删改查

```go
c.Put(ctx,key,value,option)
c.Get(ctx,key,option)
c.Delete(ctx,key,option)
```

## 租约

**获取租约**

设置租约时间为10s

```go
lea, _ := c.Grant(context.TODO(), 10)
```

**添加键值对时绑定一个租约**

当租约过期该键值对自动删除

```go
c.Put(context.TODO(), serviceKey, serviceValue, clientv3.WithLease(lea.ID))
```

**续约**

保持租约一直不过期

```go
c.KeepAlive(context.TODO(), lea.ID)
```



## 监控

**监控满足前缀的关键字，并根据事件类型完成对应业务**

```go
watch := c.Watch(context.TODO(), prefix, clientv3.WithPrefix())

for watchChan := range watch {
    for _, event := range watchChan.Events {
        switch event.Type {
        case mvccpb.PUT:
            println("插入" + string(event.Kv.Key))
            serverList[string(event.Kv.Key)] = string(event.Kv.Value)
        case mvccpb.DELETE:
            println("删除" + string(event.Kv.Key))
            delete(serverList, string(event.Kv.Key))
        }
    }
}
```

