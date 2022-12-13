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

	// broadcast用のmap生やしてGETの中でwsを保存しておく
	var wsMap = WsMap{m: make(map[*websocket.Conn]struct{})}

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
	r.StaticFile("/", "../public/index.html")

	// /wsでハンドリング
	// ip制限をかけるミドルウェアも挟む
	r.GET("/ws", ipBan(ipWhiteList), wshandler(db, &wsMap, &broadcastChan))

	// 8080でリッスン
	r.Run(":8080")
}
