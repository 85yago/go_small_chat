package main

import (
	"time"

	"github.com/gorilla/websocket"
)

// クライアントから送られてくるjsonパース用
type ClientMessage struct {
	Method  string
	Name    string
	Message string
}

// 通信用メッセージのデータ型を決める列挙型もどきの実装
type SendType string

const (
	BROADCAST SendType = "broadcast"
	GETMSG    SendType = "getReturn"
	POSTMSG   SendType = "postReturn"
)

type SendData struct {
	DataType SendType    `json:"type"`
	Data     interface{} `json:"data"`
}

// Message構造体の必要な情報のみの構造体
type RetMessage struct {
	Name      string    `json:"name"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"createtime"`
	IsMe      bool      `json:"isme"`
	wsconn    *websocket.Conn
}

// getMessage関数が返す構造体
type GetRetMessage struct {
	Status  string       `json:"status"` // OKかerrorかのみ書き込む
	Count   int          `json:"count"`
	Message []RetMessage `json:"messages"`
}

// postMessage関数が返す構造体
type PostRetMessage struct {
	Status string `json:"status"` // OKかerrorかのみ書き込む
}
