package main

import (
	"fmt"
	"github.com/fatih/color"
	"net"
	"slices"
	"strings"
)

type Result struct {
	i int
	s string
}

var (
	ips = []string{"219.141.140.10", "202.106.195.68", "221.179.155.161",
		"101.227.191.14", "139.226.226.2", "120.204.34.85",
		"59.36.213.86", "210.21.196.6", "120.196.165.24",
		"61.139.2.69", "119.6.6.6", "211.137.96.205",
		"113.246.68.203", "111.183.255.222", "113.207.25.138", "117.28.254.129"}
	names = []string{"北京电信", "北京联通", "北京移动",
		"上海电信", "上海联通", "上海移动",
		"广州电信", "广州联通", "广州移动",
		"成都电信", "成都联通", "成都移动",
		"长沙电信", "咸宁电信", "重庆电信", "福建电信"}

	m = map[string]string{
		"AS4134":    "电信163  [普通线路]",
		"AS4809":    "电信CN2  [优质线路]",
		"AS4837":    "联通4837 [普通线路]",
		"AS9929":    "联通9929 [优质线路]",
		"AS58807":   "移动CMIN2[优质线路]",
		"AS9808":    "移动CMI  [普通线路]",
		"AS58453":   "移动CMI  [普通线路]",
		"AS-CTG-CN": "电信CTGCN[优化线路]",
		"跳墙":        "路由bug  [直连线路]"}
)

func trace(ch chan Result, i int) {
	//log.Printf("trace begin for %v", i)
	hops, err := Trace(net.ParseIP(ips[i]))
	if err != nil {
		s := fmt.Sprintf("%v %-15s %v", names[i], ips[i], err)
		ch <- Result{i, s}
		return
	}

	lastIpUnknow := ""

	for _, h := range hops {
		//DebugLogPrintf("ip[%v] hops[%v] -- %v ", ips[i], i, h)

		for _, n := range h.Nodes {
			lastIpUnknow = n.IP.String()
			asn := ipAsn(n.IP.String())
			as := m[asn]
			var c func(a ...interface{}) string
			switch asn {
			case "":
				continue //没找到asn
			case "AS9929":
				c = color.New(color.FgHiYellow).Add(color.Bold).SprintFunc()
			case "AS4809":
				c = color.New(color.FgHiMagenta).Add(color.Bold).SprintFunc()
			case "AS58807":
				c = color.New(color.FgHiBlue).Add(color.Bold).SprintFunc()
			case "AS-CTG-CN":
				c = color.New(color.FgHiCyan).Add(color.Bold).SprintFunc()
			case "跳墙":
				c = color.New(color.FgHiGreen).Add(color.Bold).SprintFunc()
			default:
				c = color.New(color.FgWhite).Add(color.Bold).SprintFunc()
			}

			var duration int64 = 0
			for _, rtt := range n.RTT {
				if rtt.Milliseconds() > duration {
					duration = rtt.Milliseconds()
				}
			}

			//找到asn
			s := fmt.Sprintf("%v %-15s %-23s,  rtt::%4vms", names[i], ips[i], c(as), duration)
			ch <- Result{i, s}
			return
		}
	}

	// 没找到asn
	c := color.New(color.FgRed).Add(color.Bold).SprintFunc()
	s := fmt.Sprintf("%v %-15s %v :%-15s", names[i], ips[i], c("未知线路.."), lastIpUnknow)

	//for hop, h := range hops {
	//	for node, n := range h.Nodes {
	//		lastIpUnknow = n.IP.String()
	//		//找到asn
	//		log.Printf(".  .  .  . %v %v %v", hop, node, lastIpUnknow)
	//	}
	//}

	ch <- Result{i, s}
}

func ipAsn(ip string) string {

	switch {
	case strings.HasPrefix(ip, "59.43"):
		return "AS4809"
	case strings.HasPrefix(ip, "202.97"):
		return "AS4134"
	case strings.HasPrefix(ip, "218.105") || strings.HasPrefix(ip, "210.51"):
		return "AS9929"
	case strings.HasPrefix(ip, "219.158") || strings.HasPrefix(ip, "221.194"):
		return "AS4837"
	case strings.HasPrefix(ip, "223.120.19") || strings.HasPrefix(ip, "223.120.17") || strings.HasPrefix(ip, "223.120.16"):
		return "AS58807"
	case strings.HasPrefix(ip, "223.118") || strings.HasPrefix(ip, "223.119") || strings.HasPrefix(ip, "223.120") || strings.HasPrefix(ip, "223.121"):
		return "AS58453"
	//case strings.HasPrefix(ip, "203.22"):
	//	return "AS-CTG-CN"
	case slices.Contains(ips, ip) || strings.HasPrefix(ip, "219.141.140") || strings.HasPrefix(ip, "202.96.209") || strings.HasPrefix(ip, "58.60.188") || strings.HasPrefix(ip, "61.139.2"):
		return "跳墙"

	// case strings.HasPrefix(ip, "129.250"):  NTT
	// 	return "AS2914"
	default:
		return ""
	}
}
