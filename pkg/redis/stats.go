package redis

type Stats struct {
	TotalCount int64 // 总连接数量
	IdleCount  int64 // 空闲连接数量
	StaleCount int64 // 失效连接数量
	Hits       int64 // 命中连接数量
	Misses     int64 // 未命中连接数量
	Timeouts   int64 // 超时连接数量
}
