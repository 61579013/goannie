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
	_, err := conn.Do("SADD", k, v)
	if err != nil {
		return err
	}
	return nil
}

// Check 检查k,v是否存在
func (conn Redis) Check(k, v string) bool {
	isVid, _ := redis.Int(conn.Do("SISMEMBER", k, v))
	if isVid == 0 {
		return false
	}
	return true
}

// Count 统计k中v的数量
func (conn Redis) Count(k string) int {
	resInt, _ := redis.Int(conn.Do("SCARD", k))
	return resInt
}

// NewRedis 创建一个redis存储器
func NewRedis(CONN redis.Conn) Storage {
	return Redis{
		Conn: CONN,
	}
}
