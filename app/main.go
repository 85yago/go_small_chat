package main

import (
	"app/pkg_dbinit"
	"log"

	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/acme/autocert"
)

func main() {
	// DBの初期化をする
	db := pkg_dbinit.DbInitialization()
	// ginの初期化
	r := gin.Default()

	// https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies
	r.SetTrustedProxies(nil)

	// 禁止メッセージの正規表現のコンパイル
	// ngRegはpost_message.goで定義されるグローバル変数
	ngReg = ngMsgCompile()

	// broadcast用のmap生やしてGETの中でwsを保存しておく
	var wsMap = WsMap{m: make(map[*websocket.Conn]WsMapData)}

	// ブロードキャスト用のチャネル
	var broadcastChan BroadChan
	broadcastChan.c = make(chan RetMessage)
	// ブロードキャスト用の関数
	go broadcastMsg(&wsMap, &broadcastChan)

	// 国内ipのリストを読み込む
	ipWhiteList := readIpList()

	// ページを返す
	r.StaticFile("/chat", "../public/chat.html")
	r.StaticFile("/chat.js", "../public/chat.js")
	r.StaticFile("/chat.css", "../public/chat.css")
	r.StaticFile("/", "../public/index.html")

	// /wsでハンドリング
	// ip制限をかけるミドルウェアも挟む
	r.GET("/ws", ipBan(ipWhiteList), wshandler(db, &wsMap, &broadcastChan))

	// TLS用の設定
	m := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("azi.f5.si"),
		Cache:      autocert.DirCache("/var/www/.cache"),
	}

	// 8080でリッスン
	// r.Run(":8080")
	// TODO: リッスンの仕方はあってる？
	log.Fatal(autotls.RunWithManager(r, &m))
}
