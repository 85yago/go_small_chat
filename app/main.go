package main

import (
	"app/pkg_dbinit"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func main() {
	// DBの初期化をする
	db := pkg_dbinit.DbInitialization()
	// ginの初期化
	r := gin.Default()

	// https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies
	r.SetTrustedProxies(nil)

	// 国内ipのリストを読み込む
	ipWhiteList := readIpList()
	// ip制限をかけるミドルウェアを挟む
	r.Use(ipBan(ipWhiteList))

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

	// ページを返す
	r.StaticFile("/chat", "/var/public/chat/chat.html")
	r.StaticFile("/chat.js", "/var/public/chat/chat.js")
	r.StaticFile("/chat.css", "/var/public/chat/chat.css")
	r.StaticFile("/azi.png", "/var/public/chat/azi.png")
	r.StaticFile("/favicon.ico", "/var/public/favicon.ico")
	r.StaticFile("/", "/var/public/index.html")

	// /wsでハンドリング
	r.GET("/ws", wshandler(db, &wsMap, &broadcastChan))

	runServer(r)
}
