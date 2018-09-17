package main

import (
	"servemap/zabbixUtil"
)

const (
	zabbixAddress = "http://monitor.demolx.com:8080/api_jsonrpc.php"
	zabbixUser = "admin"
	zabbixPassword = "LinkedSee@2017"
)

func main() {
	zabbix_client := zabbixUtil.ZabbixClient{zabbixAddress, ""}
	zabbix_client.Connecting(zabbixUser, zabbixPassword)
}