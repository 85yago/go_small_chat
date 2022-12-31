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
	r.StaticFile("/chat", "/var/public/chat.html")
	r.StaticFile("/chat.js", "/var/public/chat.js")
	r.StaticFile("/chat.css", "/var/public/chat.css")
	r.StaticFile("/", "/var/public/index.html")

	// /wsでハンドリング
	// ip制限をかけるミドルウェアも挟む
	r.GET("/ws", ipBan(ipWhiteList), wshandler(db, &wsMap, &broadcastChan))

	runServer(r)
}
