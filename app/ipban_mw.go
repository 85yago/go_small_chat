package main

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"github.com/gin-gonic/gin"
)

// 国内IPのリストをテキストファイルから読み込む関数
// テキストファイルのルール:
//   - テキストファイルの名前は "internal_ip_list"
//   - 各行につきひとつ，CIDR形式のIPを記述する
//   - 空行は無視される
//   - #で始まる行は無視される
func readIpList() []*net.IPNet {
	var ipset []*net.IPNet
	fp, err := os.Open("internal_ip_list")
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)

	for scanner.Scan() {

		// 空行はスキップ
		if scanner.Text() == "" {
			continue
		}
		// #で始まる行はコメントとしてスキップ
		if scanner.Text()[0] == '#' {
			println(scanner.Text())
			continue
		}

		// CIDR形式のIP記述を順次格納していく
		_, ipnet, err := net.ParseCIDR(scanner.Text())
		if err != nil {
			panic(err)
		}
		ipset = append(ipset, ipnet)
	}

	// ローカルホストも許可する
	_, localip, err := net.ParseCIDR("127.0.0.0/8")
	if err != nil {
		panic(err)
	}
	ipset = append(ipset, localip)

	return ipset
}

// 国内のIPか判定する
func isInternalIp(ip net.IP, iplist []*net.IPNet) bool {
	for _, ipnet := range iplist {
		if ipnet.Contains(ip) {
			return true
		}
	}
	return false
}

// IPを弾くためのMiddleware
// ginのmiddlewareは/wsへの接続時に呼ばれる
func ipBan(iplist []*net.IPNet) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := ctx.ClientIP()

		// 国内IPリストに載っていなければwsを繋がない
		if !isInternalIp(net.ParseIP(ip), iplist) {
			fmt.Printf("I was accessed from invalid IP : ")
			fmt.Println(ip)
			ctx.JSON(418, "I'm a teapot")
			ctx.Abort()
		}
	}
}
