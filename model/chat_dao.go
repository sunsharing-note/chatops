package model

import "github.com/gomodule/redigo/redis"

var MyChatDao *ChatDao

type ChatDao struct {
	pool *redis.Pool
}

// NewChatDao 初始化操作
func NewChatDao(pool *redis.Pool)*ChatDao{
	return &ChatDao{pool:pool}
}

// Set 设置值
func (c *ChatDao)Set(key,value string)(err error){
	conn := c.pool.Get()
	defer conn.Close()
	if _, err = conn.Do("Set", key, value);err != nil{
		return
	}
	return
}

// Get 获取值
func (c *ChatDao)Get(key string)(value string,err error){
	conn := c.pool.Get()
	defer conn.Close()
	if get,err:=conn.Do("Get",key);err != nil{
		return "",err
	}else{
		value,err = redis.String(get,err)
		if err != nil{
			return "",err
		}
		return value,err
	}
}

// Delete 删除值
func (c *ChatDao)Delete(key string)(err error){
	conn := c.pool.Get()
	defer conn.Close()
	if _, err = conn.Do("DEL", key);err != nil{
		return
	}
	return
}