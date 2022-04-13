package db

import (
    "time"
)

type Config struct {
    Driver                  string        // 选择数据库种类，默认mysql,postgres
    DSN                     string        // DSN地址: mysql://username:password@tcp(127.0.0.1:3306)/mysql?charset=utf8mb4&collation=utf8mb4_general_ci&parseTime=True&loc=Local&timeout=1s&readTimeout=3s&writeTimeout=3s
    Debug                   bool          // 是否开启调试，默认不开启，开启后并加上export EGO_DEBUG=true，可以看到每次请求，配置名、地址、耗时、请求数据、响应数据
    MaxIdleConnections      int           // 最大空闲连接数，默认10
    MaxOpenConnections      int           // 最大活动连接数，默认100
    ConnectionKeepAlive     time.Duration // 连接的最大存活时间，默认300s
    SlowLogThreshold        time.Duration // 慢日志阈值，默认500ms
    EnableMetric            bool          // 是否开启监控，默认开启
    EnableTrace             bool          // 是否开启链路追踪，默认开启
    EnableDetailSQL         bool          // 记录错误sql时,是否打印包含参数的完整sql语句，select * from aid = ?;
    EnableLogAccess         bool          // 是否开启，记录请求数据
    EnableLogAccessRequest  bool          // 是否开启记录请求参数
    EnableLogAccessResponse bool          // 是否开启记录响应参数
}
