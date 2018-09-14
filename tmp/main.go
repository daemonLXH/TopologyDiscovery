package main

import (
	"os"
	"log"
	"fmt"
	"runtime"
)

func main() {
	fmt.Println(" **  开始检测系统信息 **")
	hostName, err := os.Hostname()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("主机名: ", hostName)
	fmt.Println("CPU: ", runtime.NumCPU())
}

