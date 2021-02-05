## redix
redix 是解决单点redis切腾讯云redis集群的迁移工具，使用前须知:

单点源: src
目标集群: dest

1. 如果dest内没有数据，并且dest尚未正式给应用使用，那么可以直接采用腾讯云的[DTS方案](https://cloud.tencent.com/document/product/571/13748?from=information.detail.dts%20%E8%85%BE%E8%AE%AF%E4%BA%91)进行迁移

2. 如果src存在冷数据(就算没有人访问，也不会失效的键值。失效期远远大于同步期的键值。)，那么这些冷数据不会正常同步，需要用户自主解决冷数据迁移。

3. 迁移周期受缓存周期影响，如果不需要保障强一致，那么可以自行决定迁移周期。

4. redix只解决golang应用的redis迁移。


## 作用

1. 单点源src和目标集群dest不需要在迁移过程中停服,数据迁移十分平滑。

## 步骤概述
- 第一步，所有应用双向写入源和目标redis，读写返回值取源返回值。
- 第二步，保持15天+(假定你的热数据最大15天失效)。确保未同步的缓存正常失效，新缓存全量同步。
- 第三步，所有应用继续双向写入，但是读写返回值取目标redis
- 第四步，至此，写的是目标redis，读的也是目标redis，源redis正式退休，移除。

## 执行详情
### 1. 第一步，双写返源
```go
    src := redix.NewPool("x.x.x.x:6379", "password", 0)
    dest := redix.NewPool("y.y.y.y:6379", "password", 0)
    RedisPool = redix.NewRedisPoolx(
    	src, dest, func() string {
    		return "src" // src,dest
    })
```

### 2. 第二步，挂机等待旧key失效

### 3. 第三步,双写返dest
```go
    src := redix.NewPool("x.x.x.x:6379", "password", 0)
    dest := redix.NewPool("y.y.y.y:6379", "password", 0)
    RedisPool = redix.NewRedisPoolx(
    	src, dest, func() string {
    		return "dest" // src,dest
    })
```

### 4. 第四步，移除双写和src
```go
RedisPool = redix.NewPool("y.y.y.y:6379", "password", 0)
```
