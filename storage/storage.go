package storage

import "github.com/garyburd/redigo/redis"

// Storage 存储器接口
type Storage interface {
	Add(k, v string) error
	Check(k, v string) bool
	Count(k string) int
}

// Redis redis存储器
type Redis struct {
	redis.Conn
}

// Add 新增k中的内容
func (conn Redis) Add(k, v string) error {
	return nil
}

// Check 检查k,v是否存在
func (conn Redis) Check(k, v string) bool {
	return true
}

// Count 统计k中v的数量
func (conn Redis) Count(k string) int {
	return 0
}

// NewRedis 创建一个redis存储器
func NewRedis(CONN redis.Conn) Storage {
	return Redis{
		Conn: CONN,
	}
}
