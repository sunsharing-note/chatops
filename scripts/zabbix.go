package scripts

import (
	"code.rookieops.com/coolops/chatops/adapter/dingding"
	"code.rookieops.com/coolops/chatops/config"
	"code.rookieops.com/coolops/chatops/scripts/zabbix"
	"fmt"
	"strings"
)

func GetHost(api *zabbix.API, host string) (zabbix.ZabbixHost, error) {
	params := make(map[string]interface{}, 0)
	filter := make(map[string]string, 0)
	filter["host"] = host
	params["filter"] = filter
	params["output"] = "extend"
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
	return nil, &zabbix.ZabbixError{0,"","Host not found"}
}

func GetUser(api *zabbix.API,u string) (zabbix.ZabbixUser,error){
	params := make(map[string]interface{}, 0)
	filter := make(map[string]string, 0)
	filter["surname"] = u
	params["filter"] = filter
	params["output"] = "extend"
	user, err := api.User("get", params)
	if err != nil {
		return nil,err
	}
	return user,err
}

func GetGraph(api *zabbix.API,hostid int)([]zabbix.ZabbixGraph,error){
	params := make(map[string]interface{}, 0)
	params["output"] = "extend"
	params["hostids"]=hostid
	params["sortfield"] = "name"
	graph, err := api.Graph("get", params)
	if err != nil {
		return nil,err
	}
	return graph,nil
}

func doZabbix(content string){
	// 连接zabbix
	api, err := zabbix.NewAPI(config.Setting.Zabbix.Url, config.Setting.Zabbix.UserName, config.Setting.Zabbix.PassWord)
	if err != nil {
		fmt.Println(err)
		return
	}
	var resData string
	if strings.Contains(content,"zabbix版本"){
		version, err := api.Version()

		if err != nil {
			fmt.Println("get zabbix version failed. err:",err)
			resData = "获取zabbix版本失败"
		}
		resData = version
	}

	msg := "#### 顺风耳机器人\n" +
		//"##### 主机：" + ip + "\n" +
		"##### 内容：\n\n" +
		resData
	dingding.SendMsgToDingTalk("markdown", msg)
}