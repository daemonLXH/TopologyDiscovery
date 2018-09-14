/*
zabbix工具，提供zabbix连接，查询全部主机的功能
 */
package main

import (
	"net/http"
	"encoding/json"
	"log"
	"bytes"
	"github.com/pkg/errors"
	"fmt"
	"net"
)

const (
	zabbixAddress = "http://monitor.demolx.com:8080/api_jsonrpc.php"
	zabbixUser = "admin"
	zabbixPassword = "LinkedSee@2017"
)

type userParams struct {
	User 	string `json:"user"`
	Password string `json:"password"`
}


type zabbixRequestContent struct {
	Jsonrpc		string				`json:"jsonrpc"`
	Method 		string				`json:"method"`
	Params 		userParams			`json:"params"`
	Id 			int					`json:"id"`
}


type ZabbixClient struct {
	url 		string
	token 		string
}

type HostInfo struct {
	HostID string 	`json:"hostid"`
	Name string `json:"name"`
}

type HostResult struct {
	JsonRPC		string 				`json:"jsonrpc"`
	Id   		int 				`json:"id"`
	Result 		[]HostInfo 			`json:"result"`
}


func (z *ZabbixClient) connecting(userName string, password string) error {
	body := zabbixRequestContent{
		Jsonrpc: "2.0",
		Method: "user.login",
		Params: userParams{userName,
			 password,
		},
		Id: 1,
	}
	content, err := json.MarshalIndent(body, "", "  ")
	if err != nil {
		return err
	}
	resp, err := http.Post(z.url, "application/json", bytes.NewReader(content))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data := map[string]interface{}{}
	json.NewDecoder(resp.Body).Decode(&data)
	token := data["result"]

	if value, ok := token.(string); ok{
		z.token = value
	} else {
		return errors.New(fmt.Sprintf("错误信息：%v", data))
	}

	log.Println("zabbix 登录成功.")
	return nil
}

func (z *ZabbixClient) QueryHostIP(hostID string) net.IP {

}

func (z *ZabbixClient) QueryHosts() ([]net.IP, error) {
	body := map[string]interface{}{
		"jsonrpc": "2.0",
		"method": "host.get",
		"auth": z.token,
		"id": 1,
		"params": map[string]interface{}{
			"output": []string{"hostid", "name"},
		},
	}
	content, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(z.url, "application/json", bytes.NewReader(content))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	data := HostResult{}
	json.NewDecoder(resp.Body).Decode(&data)

	result :=[]net.IP{}
	for _, host := range data.Result {
		ip := z.QueryHostIP(host.HostID)
		result = append(result, ip)
	}

	return result, nil
}


func main() {
	zabbixClient := ZabbixClient{url:zabbixAddress}
	err := zabbixClient.connecting(zabbixUser, zabbixPassword)
	if err != nil {
		log.Fatalln(err)
	}
	zabbixClient.QueryHosts()
}