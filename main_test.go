package main

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"log"
	"net/http"
	"testing"
	"time"
)

func Test_ip(t *testing.T) {
	yellow := color.New(color.FgHiYellow).Add(color.Bold).SprintFunc()
	green := color.New(color.FgHiGreen).SprintFunc()
	cyan := color.New(color.FgHiCyan).SprintFunc()

	rsp, _ := http.Get("http://ipinfo.io")
	info := IpInfo{}
	json.NewDecoder(rsp.Body).Decode(&info)

	fmt.Println(green("国家: ") + cyan(info.Country) + green(" 城市: ") + cyan(info.City) + green(" 服务商: ") + cyan(info.Org))
	fmt.Println(green("项目地址:"), yellow("https://github.com/zhanghanyun/backtrace"))

}

func Test_trace(test *testing.T) {
	GlobalTestMode = true
	GlobalDebugMode = true

	Test_ip(test)

	var (
		s [16]string
		c = make(chan Result)
		t = time.After(time.Second * 10)
	)

	log.Println("正在测试三网回程路由...")

	for i := range ips {
		go trace(c, i)
	}

loop:
	for range s {
		select {
		case result := <-c:
			DebugLogPrintf("go to loop result %v %s ***** ", result.i, result.s)
			s[result.i] = result.s
		case <-t:
			DebugLogPrintf("*********go to loop result case <-t")
			break loop
		}
	}

	for _, r := range s {
		log.Println(r)
	}
	log.Println("测试完成!")
}

func tt_TestDuration(t *testing.T) {
}
