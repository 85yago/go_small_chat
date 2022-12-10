package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// クライアントから送られてくるjsonパース用
type ClientMessage struct {
	Method  string
	Name    string
	Message string
}

// wsの保存用
type WsMap struct {
	sync.RWMutex
	m map[int]*websocket.Conn
}

// Message構造体の必要な情報のみの構造体
type RetMessage struct {
	Name        string
	Message     string
	CreatedTime time.Time
}

// getMessage関数が返す構造体
type GetRetMessage struct {
	Count   int
	Message []RetMessage
}

// postMessage関数が返す構造体
type PostRetMessage struct {
	Status string // OKかerrorかのみ書き込む
}

var upgrader = websocket.Upgrader{
	// TODO: ここ絶対直すこと
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// クライアントとのwsの処理
func procClient(c *gin.Context, ws *websocket.Conn) {
	for {
		// メッセージを読む
		mt, message, err := ws.ReadMessage()
		if err != nil {
			fmt.Println(err)
			break
		}

		// jsonのパース
		var clientMsg ClientMessage
		err = json.Unmarshal(message, &clientMsg)
		if err != nil {
			fmt.Println(err)
			break
		}

		// 送られたjson読んでクライアントがどっちを呼び出してるか判定
		switch clientMsg.Method {
		case "getMessage":
			// TODO: getMessage関数に変える
			message = []byte("get")
		case "postMessage":
			// TODO: postMessage関数に変える、ClientMessage型で渡す
			message = []byte("post")
		default:
			// TODO: methodエラーを入れる
			message = []byte(c.ClientIP())
		}

		// クライアントに返す
		err = ws.WriteMessage(mt, message)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

// cに送られたメッセージをブロードキャストする関数
func broadcastMsg(wsMap *WsMap, c <-chan RetMessage) {
	// このmapをgoroutineで回してbroadcast、これは更新があったら回すのを生やすって感じでよさそう？要検討

	// チャネルにメッセージが放り込まれるの待ち
	// interface定義してちゃんとそっちでやるとWriteJSONが使えると思う
	mess := <-c

	wsMap.RLock()
	for _, ws := range wsMap.m {
		err := ws.WriteMessage(websocket.TextMessage, []byte(mess))
		if err != nil {
			fmt.Println(err)
			break
		}
	}
	wsMap.RUnlock()

	// closeをどうするか->とりあえず後回し、wsが生きてるか確認できるのでは？
	// ping pongで確認できるじゃん、送る前に確認する or 別ルーチンで確認と削除
}

func main() {
	// DBの初期化をする
	// db := pkg_dbinit.DbInitialization()

	// ginの初期化
	r := gin.Default()

	// broadcast用のmap生やしてGETの中でwsを保存しておく
	var wsIndex int = 0
	var wsMap = WsMap{m: make(map[int]*websocket.Conn)}

	// ブロードキャスト用のチャネル
	broadcastChan := make(chan RetMessage)
	// ブロードキャスト用の関数
	go broadcastMsg(&wsMap, broadcastChan)

	r.GET("/", func(c *gin.Context) {
		// websocketで接続
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer ws.Close()
		// ブロードキャスト用にソケットを保存
		wsMap.Lock()
		wsMap.m[wsIndex] = ws
		wsMap.Unlock()

		// クライアントからのwebsocketを処理
		procClient(c, ws)
	})

	r.Run(":8080") // 8080でリッスン
}
