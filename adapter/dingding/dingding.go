package dingding

import (
	"code.rookieops.com/coolops/chatops/config"
	"code.rookieops.com/coolops/chatops/message"
	"code.rookieops.com/coolops/chatops/middleware"
	"code.rookieops.com/coolops/chatops/model"
	"code.rookieops.com/coolops/chatops/myredis"
	"code.rookieops.com/coolops/chatops/scripts"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"io/ioutil"
	"strconv"
)

// 定义一个钉钉结构体
type Dingtalk struct {
}

var (
	Dingding *Dingtalk
	redisPool *redis.Pool
	content string
	)

// 加签
func signature(ts int64, secret string) string {
	strToSign := fmt.Sprintf("%d\n%s", ts, secret)
	hmac256 := hmac.New(sha256.New, []byte(secret))
	hmac256.Write([]byte(strToSign))
	data := hmac256.Sum(nil)
	return base64.StdEncoding.EncodeToString(data)
}

// 初始化dingtalk
func NewDingtalk() *Dingtalk {
	return &Dingtalk{}
}

//var MySession sessions.Session

func (d *Dingtalk) DingDing(c *gin.Context) {
	// 获取body里的请求参数
	//fmt.Println(c.Request.Header)
	HttpSign := c.Request.Header.Get("Sign")
	HttpTimestamp := c.Request.Header.Get("Timestamp")
	// timestamp 与系统当前时间戳如果相差1小时以上，则认为是非法的请求。
	tsi, err := strconv.ParseInt(HttpTimestamp, 10, 64)
	if err != nil {
		middleware.Logger().Error("请求头可能未附加时间戳信息!!")
	}

	data, _ := ioutil.ReadAll(c.Request.Body)
	fmt.Println("---body/--- \r\n " + string(data))

	sign := signature(tsi, config.Setting.DingDing.AppSecret)

	// 校验成功
	if HttpSign == sign {
		// 开启协程，监听消息发送
		go d.listenOutChanMsg()

		var body incoming
		err := json.Unmarshal(data, &body)
		if err != nil {
			fmt.Println(err)
			return
		}

		// 初始化Dingtalk
		Dingding = NewDingtalk()
		senderNick := c.Request.Header.Get("senderNick")
		// 初始化redis和MyChatDao
		redisPool = myredis.RedisPool(config.Setting.Redis.IpAddr)
		model.MyChatDao = model.NewChatDao(redisPool)
		// 从redis中取值
		//myredis.MyPool = myredis.RedisPool()
		//defer myredis.MyPool.Close()
		//rdsConn := myredis.MyPool.Get()
		getName, _ := model.MyChatDao.Get("name")
		getData, _ := model.MyChatDao.Get("data")
		//getName, _ := redis.String(rdsConn.Do("get", "name"))
		//getData, _ := redis.String(rdsConn.Do("get", "data"))
		fmt.Println(getName)
		if getName == senderNick && getData != "" {
			// 将起拼接到现有的前面
			content = getData + body.Text.Content
		} else {
			//_, _ = rdsConn.Do("DEL", "data")
			//_, _ = rdsConn.Do("SET", "name", senderNick)
			if err := model.MyChatDao.Delete("data");err!=nil{
				middleware.Logger().Error("delete data from redis failed,",err)
				return
			}
			if err := model.MyChatDao.Set("name", senderNick);err != nil{
				middleware.Logger().Error("set name to redis failed,",err)
				return
			}
			content = body.Text.Content
		}
		fmt.Println(content)
		msg := message.NewMessage(content)
		msg.Sender = body.SenderId
		msg.Header.Set("sender", body.SenderNick)

		// 可以剔除
		//resMsg := message.NewMessage("收到，马上处理.....")
		//resMsg.Header.Set("msgtype","text")
		//message.OutChan <- resMsg

		scripts.RunCommand(msg)
	}
}
