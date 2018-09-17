/*
zabbix工具，提供zabbix连接，查询全部主机的功能
 */
package zabbixUtil

import (
	"net/http"
	"encoding/json"
	"log"
	"bytes"
	"github.com/pkg/errors"
	"fmt"
	"net"
)


func (z *ZabbixClient) Connecting(userName string, password string) error {
	var body LoginParams
	body.JsonRPC = "2.0"
	body.Method = "user.login"
	body.Params.User = userName
	body.Params.Password = password
	body.Id = 1

	content, err := json.MarshalIndent(body, "", "  ")
	if err != nil {
		return err
	}

	log.Printf("login params: %s ", content)
	resp, err := http.Post(z.Url, "application/json", bytes.NewReader(content))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data := map[string]interface{}{}
	json.NewDecoder(resp.Body).Decode(&data)
	token := data["result"]

	if value, ok := token.(string); ok{
		z.Token = value
	} else {
		return errors.New(fmt.Sprintf("错误信息：%v", data))
	}

	log.Println("zabbix 登录成功.")
	return nil
}

func (z *ZabbixClient) QueryHostIP(hostID string) net.IP {
	body :=map[string]interface{}{
		"jsonrpc": "2.0",
		"method": "hostinterface.get",
		"auth": z.Token,
		"id": 1,
		"params": map[string]interface{}{
			"output": []string{"ip"},
			"hostids": hostID,
		},
	}

	content, err := json.Marshal(body)
	if err != nil {
		return nil
	}

	resp, err := http.Post(z.Url, "application/json", bytes.NewReader(content))
	if err != nil {
		return nil
	}

	defer resp.Body.Close()

	data := StandResult{}
	json.NewDecoder(resp.Body).Decode(&data)

	log.Println(data)
	log.Println(data.Result)



	if value, ok :=data.Result["ip"].(net.IP); ok {
		return value
	}
	return nil

}

func (z *ZabbixClient) QueryHosts() ([]net.IP, error) {
	body := map[string]interface{}{
		"jsonrpc": "2.0",
		"method": "host.get",
		"auth": z.Token,
		"id": 1,
		"params": map[string]interface{}{
			"output": []string{"hostid", "name"},
		},
	}
	content, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(z.Url, "application/json", bytes.NewReader(content))
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


//func main() {
//	zabbixClient := ZabbixClient{url:zabbixAddress}
//	err := zabbixClient.connecting(zabbixUser, zabbixPassword)
//	if err != nil {
//		log.Fatalln(err)
//	}
//	zabbixClient.QueryHosts()
//}