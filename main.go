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

	var (
		s [16]string
		c = make(chan Result)
		t = time.After(time.Second * 10)
	)

	green := color.New(color.FgHiGreen).SprintFunc()
	cyan := color.New(color.FgHiCyan).SprintFunc()
	log.Println("正在测试三网回程路由...")

	rsp, _ := http.Get("http://ipinfo.io")
	info := IpInfo{}
	json.NewDecoder(rsp.Body).Decode(&info)

	fmt.Println(green("国家: ") + cyan(info.Country) + green(" 城市: ") + cyan(info.City) + green(" 服务商: ") + cyan(info.Org))

	for i := range ips {
		go trace(c, i)
	}

loop:
	for range s {
		select {
		case o := <-c:
			s[o.i] = o.s
		case <-t:
			break loop
		}
	}

	for _, r := range s {
		fmt.Println(r)
	}
	log.Println(green("测试完成!"))
}
