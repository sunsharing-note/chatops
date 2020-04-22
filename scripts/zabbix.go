package scripts

import (
	"code.rookieops.com/coolops/chatops/config"
	"code.rookieops.com/coolops/chatops/message"
	"code.rookieops.com/coolops/chatops/scripts/zabbix"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// GetHost 获取主机
func GetHost(api *zabbix.API, host string) (zabbix.ZabbixHost, error) {
	params := make(map[string]interface{}, 0)
	filter := make(map[string]string, 0)
	filter["host"] = host
	params["filter"] = filter
	params["output"] = []string{"hostid", "host"}
	params["select_groups"] = "extend"
	params["templated_hosts"] = 1
	ret, err := api.Host("get", params)

	// This happens if there was an RPC error
	if err != nil {
		return nil, err
	}

	// If our call was successful
	if len(ret) > 0 {
		return ret[0], err
	}

	// This will be the case if the RPC call was successful, but
	// Zabbix had an issue with the data we passed.
	return nil, &zabbix.ZabbixError{0, "", "Host not found"}
}

// GetUser 获取用户
func GetUser(api *zabbix.API, u string) ([]zabbix.ZabbixUser, error) {
	params := make(map[string]interface{}, 0)
	filter := make(map[string]string, 0)
	if u != "all" {
		filter["surname"] = u
		params["filter"] = filter
	}
	params["output"] = []string{"userid", "alias", "name", "surname"}
	user, err := api.User("get", params)
	if err != nil {
		return nil, err
	}
	return user, err
}

// GetGraph 获取图形
func GetGraph(api *zabbix.API, hostid int) ([]zabbix.ZabbixGraph, error) {
	params := make(map[string]interface{}, 0)
	params["output"] = "extend"
	params["hostids"] = hostid
	params["sortfield"] = "name"
	graph, err := api.Graph("get", params)
	if err != nil {
		return nil, err
	}
	return graph, nil
}

// GetAlert 获取告警
func GetAlert(api *zabbix.API, actionid string) ([]zabbix.ZabbixAlert,error){
	params := make(map[string]interface{}, 0)
	params["output"] = "extend"
	params["actionids"] = actionid
	alerts, err := api.Alert(params)
	if err != nil {
		return nil,err
	}
	return alerts,nil
}

// GetEvent 获取事件
func GetEvent(api *zabbix.API, hostids string) ([]zabbix.ZabbixEvent,error){
	params := make(map[string]interface{}, 0)
	params["output"] = "extend"
	params["hostids"] = hostids
	alerts, err := api.Event("get",params)
	if err != nil {
		return nil,err
	}
	return alerts,nil
}

// GetHistory 获取事件
func GetHistory(api *zabbix.API, hostids string) ([]zabbix.ZabbixHistoryItem,error){
	params := make(map[string]interface{}, 0)
	params["output"] = "extend"
	params["hostids"] = hostids
	params["sortfield"] = "clock"
	params["sortorder"] = "DESC"
	params["limit"] = 10
	alerts, err := api.History("get",params)
	if err != nil {
		return nil,err
	}
	return alerts,nil
}


func doZabbix(msg *message.Message) {
	// 连接zabbix
	api, err := zabbix.NewAPI(config.Setting.Zabbix.Url, config.Setting.Zabbix.UserName, config.Setting.Zabbix.PassWord)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 登录zabbix
	_, err = api.Login()
	if err != nil {
		fmt.Println(err)
		return
	}
	var resData string
	var tmp string
	content := msg.ReadMessageToString()
	if strings.Contains(content, "版本信息") {
		version, err := api.Version()

		if err != nil {
			fmt.Println("get zabbix version failed. err:", err)
			resData = "获取zabbix版本失败"
		}else{
			resData = version
		}
		tmp = "#### 顺风耳机器人\n" +
			//"##### 主机：" + ip + "\n" +
			"version:" + resData
	}

	if strings.Contains(content, "所有用户") {
		user, err := GetUser(api, "all")
		if err != nil {
			fmt.Println("get all user failed. err:", err)
			resData = "获取ZABBIX所有用户失败"
		}else{
			data, err := json.Marshal(&user)
			if err != nil {
				fmt.Println("marshal user data failed. err:",err)
				return
			}
			resData = string(data)
		}
		tmp = "#### 顺风耳机器人\n" +
			//"##### 主机：" + ip + "\n" +
			"所有用户:" + resData
	}
	reg := regexp.MustCompile(`\d+.\d+.\d+.\d+`)
	res := reg.FindAllString(content, -1)
	if strings.Contains(content, "主机信息") {
		host, err := GetHost(api, res[0])
		if err != nil {
			fmt.Printf("获取主机%s的监控信息失败, err:%s", res[0], err.Error())
			resData = fmt.Sprintf("获取主机%s的监控信息失败", res[0])
		}else{
			data, err := json.Marshal(&host)
			if err != nil {
				fmt.Println("marshal host data failed. err:",err)
				return
			}
			resData = string(data)
		}

		tmp = "#### 顺风耳机器人\n" +
			"##### 主机：" + res[0] + "\n" +
			"主机信息:" + resData
	}
	if strings.Contains(content,"报警信息"){
		alert, err := GetAlert(api, "3")
		if err != nil {
			fmt.Println("get alert info failed. err:",err)
			resData = "获取报警信息失败"
		}else{
			data, err := json.Marshal(&alert)
			if err != nil {
				fmt.Println("marshal alert data failed. err:",err)
				return
			}
			resData = string(data)
		}

		tmp = "#### 顺风耳机器人\n" +
			//"##### 主机：" + res[0] + "\n" +
			"报警信息:" + resData
	}
	if strings.Contains(content,"事件信息"){
		// 1、获取到hostids
		fmt.Println(res[0])
		host, err := GetHost(api, res[0])
		if err != nil {
			fmt.Println(err)
			return
		}
		if host["hostid"] != nil{
			events, err := GetEvent(api, host["hostid"].(string))
			if err != nil {
				fmt.Println(err)
				return
			}else{
				data, err := json.Marshal(&events)
				if err != nil {
					fmt.Println("marshal events data failed. err:",err)
					return
				}
				resData = string(data)
			}
			tmp = "#### 顺风耳机器人\n" +
				"##### 主机：" + res[0] + "\n" +
				"事件信息:" + resData
		}
	}
	if strings.Contains(content,"历史记录"){
		// 1、获取到hostids
		fmt.Println(res[0])
		host, err := GetHost(api, res[0])
		if err != nil {
			fmt.Println(err)
			return
		}
		if host["hostid"] != nil{
			history, err := GetHistory(api, host["hostid"].(string))
			if err != nil {
				fmt.Println(err)
				return
			}else{
				data, err := json.Marshal(&history)
				if err != nil {
					fmt.Println("marshal events data failed. err:",err)
					return
				}
				resData = string(data)
			}
			tmp = "#### 顺风耳机器人\n" +
				"##### 主机：" + res[0] + "\n" +
				"历史记录:" + resData
		}
	}
	msg.Header.Set("msgtype","markdown")
	msg.Body = strings.NewReader(tmp)
	message.OutChan <- msg
}
