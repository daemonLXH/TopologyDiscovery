package zabbixUtil

import "net"

// 基础结构
type baseStruct struct {
	JsonRPC 	string 		`json:"jsonrpc"`
	Id 			int 		`json:"id"`
}

// zabbix请求基础结构
type baseRequestStruct struct {
	baseStruct
	Method 		string 		`json:"method"`
}

// 登录参数
type userParams struct {
	User 	string `json:"user"`
	Password string `json:"password"`
}

// 登录请求
type LoginParams struct {
	baseRequestStruct
	Params  userParams `json:"params"`
}

// zabbix 客户端只需要在每次请求时携带token即可
type ZabbixClient struct {
	Url 		string
	Token 		string
}


// 其它查询请求基础结构
type QueryRequestParams struct {
	baseRequestStruct
	Auth string 		`json:"auth"`
	Params interface{}	`json:"params"`
}

type HostInfo struct {
	HostID string 	`json:"hostid"`
	Name string `json:"name"`
}

// zabbix 中host.get的查询结果
type HostResponse struct {
	baseStruct
	Result 		[]HostInfo 			`json:"result"`
}

type InterfaceResult struct {
	InterfaceId		string		`json:"interfaceid"`
	Ip 				net.IP 		`json:"ip"`
}

// Zabbix中interface.get的查询结果
type IpResponse struct {
	baseStruct
	Result 		[]InterfaceResult 	`json:"result"`
}

// 登录结果
type TokenResponse struct {
	baseStruct
	Result	string `json:"result"`
}
