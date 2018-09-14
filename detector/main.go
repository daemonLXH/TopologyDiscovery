// 每分钟读取一次连接信息
// 每10分钟从服务器上同步一次IP列表
// 本地记录已上传的连接关系，如果有新的连接关系，每5分钟上传一次。

package main

import (
	"log"
	"io/ioutil"
	"strings"
)

const (
	PROC_TCP = "tcp"
	PROC_TCP6 = "proc/net/tcp6"
)

type validData struct {
	LocalAddress string
	RemoteAddress string
	Status string
}


var STATE = map[string]string{
	"01": "ESTABLISHED",
	"02": "SYN_SENT",
	"03": "SYN_RECV",
	"04": "FIN_WAIT1",
	"05": "FIN_WAIT2",
	"06": "TIME_WAIT",
	"07": "CLOSE",
	"08": "CLOSE_WAIT",
	"09": "LAST_ACK",
	"0A": "LISTEN",
	"0B": "CLOSING",
}



type Process struct {
	User	string
	Name 	string
	Pid 	string
	Exe 	string
	State 	string
	Ip 		string
	Port 	int64
	ForeignIp 	string
	ForeignPort	string
}


func getData(t string) []string {
	var proc_t string

	switch t {
	case "tcp":
		proc_t = PROC_TCP
	case "tcp6":
		proc_t = PROC_TCP6

	default:
		log.Fatalln("只能查询tcp连接。")
	}

	data, err := ioutil.ReadFile(proc_t)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(data), "\n")

	return lines[1:len(lines)]
}

func deleteEmpty(line []string) []string {
	result := []string{}
	for _, s := range line {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		result = append(result, s)
	}
	return result
}


func main() {
	data := getData("tcp")

	var result = []validData{}
	for _, line := range data {
		line_array := deleteEmpty(strings.Split(line, " "))
		if len(line_array) < 2 {
			continue
		}

		valid_data := validData{line_array[1], line_array[2], line_array[3]}
		log.Println(valid_data)
		result = append(result, valid_data)
	}
}