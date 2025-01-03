package main

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"log"
	"net/http"
	"time"
)

var GlobalTestMode bool = false
var GlobalDebugMode bool = false

func DebugLogPrintf(format string, v ...any) {
	if GlobalDebugMode {
		log.Printf(format, v...)
	}
}

type IpInfo struct {
	Ip      string `json:"ip"`
	City    string `json:"city"`
	Region  string `json:"region"`
	Country string `json:"country"`
	Org     string `json:"org"`
}

func main() {
	var timeoutSecond = 10
	if GlobalTestMode {
		//timeoutSecond = 7
		ips = ips[:0]
		ips = append(ips, "219.141.140.10")
	} else {
		network_info()
	}

	var (
		s       [16]string
		c       = make(chan Result)
		timeout = time.After(time.Second * time.Duration(timeoutSecond))
	)

	for i := range ips {
		go trace(c, i)
	}

loop:
	for index, _ := range s {
		select {
		case o := <-c:
			s[o.i] = o.s
		case <-timeout:
			DebugLogPrintf("")
			DebugLogPrintf("")
			DebugLogPrintf("~~~~~loop loop case <-t: time out %vsecond", timeoutSecond)
			DebugLogPrintf("~~~~~loop loop case <-t: time out %vsecond", timeoutSecond)
			break loop
		}

		if (index + 1) == len(ips) {
			DebugLogPrintf("index reach, break loop")
			break loop
		}
	}

	fmt.Println("")
	fmt.Println("")
	for index, r := range s {
		fmt.Println(r)

		if (index + 1) == len(ips) {
			DebugLogPrintf("index reach, break Println")
			break
		}
	}
	log.Println(color.New(color.FgHiGreen).SprintFunc()("测试完成!"))
}

func network_info() {
	green := color.New(color.FgHiGreen).SprintFunc()
	cyan := color.New(color.FgHiCyan).SprintFunc()
	log.Println("正在测试三网回程路由...")

	rsp, _ := http.Get("http://ipinfo.io")
	info := IpInfo{}
	json.NewDecoder(rsp.Body).Decode(&info)

	fmt.Println(green("国家: ") + cyan(info.Country) + green(" 城市: ") + cyan(info.City) + green(" 服务商: ") + cyan(info.Org))
}
