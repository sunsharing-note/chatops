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

// Set 将key写入redis
func (c *ChatDao)Set(key,value string)(err error){
	conn := c.pool.Get()
	defer conn.Close()
	if _, err = conn.Do("Set", key, value);err != nil{
		return
	}
	return
}

// Get 从redis中根据key值获取值
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

// Delete 从redis中删除key
func (c *ChatDao)Delete(key string)(err error){
	conn := c.pool.Get()
	defer conn.Close()
	if _, err = conn.Do("DEL", key);err != nil{
		return
	}
	return
}

// 清除数据
func (c *ChatDao)ClearData(data []string)(err error){
	for _, v := range data{
		if err = c.Delete(v);err !=nil{
			return
		}
	}
	return
}