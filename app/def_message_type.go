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

// Message構造体の必要な情報のみの構造体
type RetMessage struct {
	Name      string
	Message   string
	CreatedAt time.Time
	IsMe      bool
	wsconn    *websocket.Conn
}

// getMessage関数が返す構造体
type GetRetMessage struct {
	Status  string // OKかerrorかのみ書き込む
	Count   int
	Message []RetMessage
}

// postMessage関数が返す構造体
type PostRetMessage struct {
	Status string // OKかerrorかのみ書き込む
}
