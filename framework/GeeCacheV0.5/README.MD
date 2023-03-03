版本：v0.5

# 功能
- 实现一致性哈希

# 介绍


# api介绍
## net.url
- func QueryEscape(s string) string
  - QueryEscape函数对s进行转码使之可以安全的用在URL查询里。
- func QueryUnescape(s string) (string, error)
  - QueryUnescape函数用于将QueryEscape转码的字符串还原。它会把%AB改为字节0xAB，将'+'改为' '。如果有某个%后面未跟两个十六进制数字，本函数会返回错误。

# 捋一捋这个框架思路
## 实现LRU算法(lru包)
- key是string类型，value是实现Value接口的类型，该接口只有一个方法，返回Len()表示该Value占的字节
- 如果配置了onEvicted函数，会在移除缓存数据时执行
## 整合lru算法，实现并发安全(cache包)
- 整合到cache包下的cache里面
  - 传进行的是key:string、value: ByteView
## 整合cache包的IMutexCache(group包)
- 提供获取的方法，获取不到就调用回调函数获取数据，并缓存起来，用户可以自定义回调函数
- 还有一个重要的变量，就是map类型的groups,可以说它是所有缓存的入口，没有它，本机的缓存就没了
## 上面三个的小结
- 通过上面三个包，最后整合处理的group包就是用户可以使用的了。
  - 用户根据NewGroup创建缓存分组，根据GetGroup获取缓存分组
  - 获取不到缓存也可以调用用户的回调函数Getter的Get函数获取，并把数据加入缓存中
## consistentHash包
- 这里面实现了插入哈希环和在哈希环寻找的简单算法，并没有实现加入节点需要转移缓存等操作，同时也没有删除节点操作
## group包(核心包
- 我们将http服务端加入这个包中，HTTPPool既可以提供web服务(ServeHTTP)，也可根据缓存的key获取别的主机的客户端服务(PickPeer)
- Group结构体整合HTTPPool，实现远程服务器获取数据服务