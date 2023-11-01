
# Cache 
### 基于一致性算法的分布式缓存

cache是根据极客兔兔大佬编写的7days-golang中的GeeCache开发的。
其特点：  
* 单机缓存和基于 HTTP 的分布式缓存  
* 最近最少访问(Least Recently Used, LRU) 缓存策略  
* 使用 Go 锁机制防止缓存击穿  
* 使用一致性哈希选择节点，实现负载均衡  
* 使用 protobuf 优化节点间二进制通信  
* 使用 log 封装的日志  
* 中心节点具备服务发现  
### Distributed Caching Based on Consistency Algorithm

Cache is developed based on GeeCache in 7days golang, written by Geek Rabbit.
Its characteristics:  
* Single machine caching and HTTP based distributed caching  
* Least Recently Used (LRU) cache policy  
* Using the Go lock mechanism to prevent cache breakdown  
* Using consistent hashing to select nodes for load balancing  
* Optimizing binary communication between nodes using protobuf  
* Using log encapsulated logs  
* The central node has service discovery  

## API Reference

### Center Node
```http
  POST /registry
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `port` | `string` | **Required**. port of node to register|

```http
  GET /get/${key}
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `key`      | `string` | **必需**. key of value to fetch |

### 中心节点
```http
  POST /registry
```

| 参数   | 类型     | 说明                |
| :-------- | :------- | :------------------------- |
| `port` | `string` | **必需**. 端口注册|

```http
  GET /get/${key}
```

| 参数 | 类型     | 说明                       |
| :-------- | :------- | :-------------------------------- |
| `key`      | `string` | **必需**. 通过key查找值 |

### Node
Waiting for updates  
待跟新



## Related Links
1. [7 day to golang]( https://github.com/geektutu/7days-golang/tree/master)                                             分布式缓存 GeeCache
## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=jarvis0919/Cache&type=Date)](https://star-history.com/#jarvis0919/Cache&Date)