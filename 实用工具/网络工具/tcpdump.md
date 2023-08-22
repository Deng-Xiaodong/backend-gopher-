**简介**

网络层数据包捕获工具。支持对协议、主机、网络（段）、端口、方向等过滤，并支持<font color=red>and、or、not</font>等逻辑语句帮助去除无用的消息。

# 基本使用

**不指定任何参数**

监听第一块网卡经过的所有数据包。主机上可能不止一块网卡，所以通常需要指定网卡

```shell
tcpdump
```

**监听特定网卡**

`-i interface`

```shell
//查看网卡信息
netstat

tcpdump -i en0
```

**监听特定主机**

```shell
tcpdump host hostname
```

**监听特定端口**

```shell
tcpdump port 8080
```

**监听网络段**

```shell
tcpdump net 192.168.0.1/24
```

**监听端口段**

```shell
tcpdump port 80:90
```

**指定方向**

`src`：指定来源主机（谁过来）

`dst`：指定目的主机（去找谁）

```shell
tcpdump src host hostname
tcpdump dst host hostname
```

## 使用例子

**监听指定来源在本机指定端口上的tcp报文**

```shell
tcpdump tcp -i eh0 port 90 and src host 47.107.46.198
```

**监两特定主机之间的通信**

```shell
tcpdump ip host hostname1 and hostname2
tcpdump ip host hsotname1 and ! hostname2
```

**稍微详细的例子**

```shell
tcpdump tcp -i et0 -t -s 0 -c 100 dst port !22 and net 192.168.0.1/24 -w ./target.cap
```

- `-i`：指定网卡
- `-t`：不显示时间
- `-s`：指定数据包大小
- `-c`：指定抓取包数量
- `dst prot !22`：不抓取 目标地址为\*.\*.\*.\*:22的数据包
- `net 192.168.0.1/24`：指定网段
- `-w`：输出数据写到本地文件