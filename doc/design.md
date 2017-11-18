### influxdb 组件依赖和功能

```
Server
    |---- *meta.Client                                          // 元数据管理
    |---- *tsdb.Store                                           // 提供存储
    |---- *query.QueryExecutor                                  // 执行查询
        |---- StatementExecutor                                 // 接口, 执行语句
    |---- *coordinator.PointsWriter                             // 数据写入入口
        |---- MetaClient
            |---- Database
            |---- RetentionPolicy
            |---- CreateShardGroup
        |---- TSDBStore
            |---- CreateShard
            |---- WriteToShard
    |---- *subscriber.Service                                   // 提供对 coordinator.PointsWriter.WritePointsPrivileged 的监听
        |---- MetaClient
            |---- Databases
            |---- WaitForDataChanged
    |---- *snaptshotter.Service                                 // 提供数据快照
        |---- MetaClient
            |---- encoding.BinaryMarshaler
                |---- MarshalBinary
            |---- Database
    |---- *monitor.Monitor                                      // 监控及记录
        |---- MetaClient
            |---- CreateDatabaseWithRetentionPolicy
            |---- Database
            |---- PointsWriter
                |---- WritePoints
    |---- Services
        |---- snapshotter.NewService                            // 提供数据快照传输
        |---- Server.Monitor                                    // 监控及记录
        |---- retention.NewService                              // 检查并实施 retention
        |---- httpd.NewService                                  // 提供 http 服务
        |---- storage.NewService                                // 提供对底层 tsdb.Store 的读服务
        |---- collectd.NewService
        |---- opentsdb.NewService
        |---- graphite.NewService
        |---- precreator.NewService                             // 定时检查并执行 shard precreation
        |---- udp.NewService                                    // 提供 udp 的数据输入服务
        |---- continuous_querier.NewService                     // 提供

```
