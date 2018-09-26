/*
zabbix工具，提供zabbix连接，查询全部主机的功能
 */
package zabbixUtil

import (
	"net/http"
	"encoding/json"
	"log"
	"bytes"
	"net"
)

// 登录并获取token
func (z *ZabbixClient) Connecting(userName string, password string) error {
	// 这里有没有更好的写法？有点繁琐
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

	data := TokenResponse{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return err
	}
	z.Token = data.Result
	log.Printf("zabbix 登录成功，token: %s", z.Token)
	return nil
}

// 初始化一些不变的参数
func (z *ZabbixClient) CreateBaseParams() QueryRequestParams {
	var params = QueryRequestParams{}
	params.JsonRPC = "2.0"
	params.Id = 1
	params.Auth = z.Token

	return params
}

// 以主机ID查询IP
func (z *ZabbixClient) QueryHostIP(hostID string) net.IP {
	body := z.CreateBaseParams()
	body.Method = "hostinterface.get"
	body.Auth = z.Token
	body.Params = map[string]interface{}{
		"output": []string{"ip"},
		"hostids": hostID,
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

	data := IpResponse{}
	json.NewDecoder(resp.Body).Decode(&data)
	return data.Result[0].Ip
}

// 查询全部主机
func (z *ZabbixClient) QueryHosts() ([]net.IP, error) {
	body := z.CreateBaseParams()
	body.Method = "host.get"
	body.Params = map[string]interface{}{
		"output": []string{"hostid", "name"},
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

	data := HostResponse{}
	json.NewDecoder(resp.Body).Decode(&data)

	result :=[]net.IP{}
	for _, host := range data.Result {
		ip := z.QueryHostIP(host.HostID)
		result = append(result, ip)
	}

	return result, nil
}
