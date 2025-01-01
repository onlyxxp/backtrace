package main

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"net/http"
	"os"
	"testing"
)

func Test_ip(t *testing.T) {
	green := color.New(color.FgHiGreen).SprintFunc()
	cyan := color.New(color.FgHiCyan).SprintFunc()

	rsp, _ := http.Get("http://ipinfo.io")
	info := IpInfo{}
	json.NewDecoder(rsp.Body).Decode(&info)

	fmt.Println(green("国家: ") + cyan(info.Country) + green(" 城市: ") + cyan(info.City) + green(" 服务商: ") + cyan(info.Org))

}

func Test_main(test *testing.T) {
	GlobalTestMode = true
	GlobalDebugMode = true

	main()
}

func TestParams(t *testing.T) {
	// 获取所有命令行参数
	args := os.Args

	// 第一个参数是程序名称，从第二个参数开始是传入的命令行参数
	// 可以根据需要进行进一步的处理
	for i, arg := range args {
		fmt.Printf("Argument %d: %s\n", i, arg)
	}
}
