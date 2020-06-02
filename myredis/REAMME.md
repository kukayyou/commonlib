# 1.简述
本包主要是对Redis常用操作的封装，包括连接池、发布订阅等。

# 2.使用介绍
## 2.1 连接池的使用
* 创建集群模式无密码的连接池：ucredis.NewRedisClusterConnPool(connsCount, clusterAddrs)
* 创建集群模式有密码的连接池：ucredis.NewRedisClusterConnPoolWithPassword(connsCount, clusterAddrs, password)
* 创建哨兵模式的连接池：ucredis.NewSentinelPool(name, clusterAddrs, connsCount)

## 2.2 执行命令
* 直接调用连接池的Do方法：cluster.Do(key, "set", key, "val")