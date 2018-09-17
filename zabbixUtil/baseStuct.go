package zabbixUtil


type baseStruct struct {
	JsonRPC 	string 		`json:"jsonrpc"`
	Id 			int 		`json:"id"`
}

type baseRequestStruct struct {
	baseStruct
	Method 		string 		`json:"method"`
}

type userParams struct {
	User 	string `json:"user"`
	Password string `json:"password"`
}

type LoginParams struct {
	baseRequestStruct
	Params  userParams `json:"params"`
}

type ZabbixClient struct {
	Url 		string
	Token 		string
}


type HostInfo struct {
	HostID string 	`json:"hostid"`
	Name string `json:"name"`
}


type HostResult struct {
	baseStruct
	Result 		[]HostInfo 			`json:"result"`
}

type StandResult struct {
	baseStruct
	Result 		map[string]interface{} 	`json:"result"`
}
