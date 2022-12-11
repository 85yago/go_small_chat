package main

import (
	"app/pkg_dbinit"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
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
	m map[*websocket.Conn]struct{}
}

// ブロードキャスト用のチャネル
type BroadChan struct {
	sync.RWMutex
	c chan RetMessage
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

// upgraderはHTTPをWSにするときに呼ばれる
// ここで許可するoriginや接続時間を設定する
var upgrader = websocket.Upgrader{
	// TODO: ここ絶対直すこと
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// クライアントとのwsの処理
func procClient(c *gin.Context, db *gorm.DB, ws *websocket.Conn, broadcastChan *BroadChan) {
	for {
		// メッセージを読む
		var clientMsg ClientMessage
		err := ws.ReadJSON(&clientMsg)
		if err != nil {
			fmt.Println(err)
			break
		}

		var retMsg any

		// 送られたjson読んでクライアントがどっちを呼び出してるか判定
		switch clientMsg.Method {
		case "getMessage":
			retMsg = getMessage(db)
		case "postMessage":
			retMsg = postMessage(c, db, ws, broadcastChan.c, clientMsg)
		default:
			retMsg = PostRetMessage{Status: "method error"}
		}

		broadcastChan.Lock()
		err = ws.WriteJSON(retMsg)
		broadcastChan.Unlock()
		// クライアントに返す
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

// cに送られたメッセージをブロードキャストする関数
func broadcastMsg(wsMap *WsMap, c *BroadChan) {
	for {
		// チャネルにメッセージが放り込まれるの待ち
		msg := <-c.c

		// map用のロック
		wsMap.RLock()
		for ws := range wsMap.m {
			// 送信先が送信元と同じならば
			if ws == msg.wsconn {
				msg.IsMe = true

				// ws用のロック
				c.Lock()

				// 各wsにメッセージを送る
				err := ws.WriteJSON(msg)
				if err != nil {
					fmt.Println(err)
					continue
				}

				c.Unlock()
			} else {
				msg.IsMe = false

				// 各wsにメッセージを送る
				err := ws.WriteJSON(msg)
				if err != nil {
					fmt.Println(err)
					continue
				}
			}
		}
		wsMap.RUnlock()
	}
}

// ws用のエンドポイントのハンドラ
func wshandler(db *gorm.DB, wsMap *WsMap, broadcastChan *BroadChan) func(*gin.Context) {
	return func(c *gin.Context) {
		// websocketで接続
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			fmt.Println(err)
			return
		}

		// r.GETが終わった時にクローズ（クローズハンドラがないから特に何もしない）
		defer ws.Close()

		// r.GETが終わった時にwsをwsMapから削除
		defer func(wsMap *WsMap) {
			wsMap.Lock()
			delete(wsMap.m, ws)
			wsMap.Unlock()
		}(wsMap)

		// ブロードキャスト用にソケットを保存
		wsMap.Lock()
		wsMap.m[ws] = struct{}{}
		wsMap.Unlock()

		// クライアントからのwebsocketを処理
		procClient(c, db, ws, broadcastChan)
	}
}

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
	internalIpList := readIpList()

	// ページを返す
	r.StaticFile("/chat", "../public/chat.html")
	r.StaticFile("/chat.js", "../public/chat.js")
	r.StaticFile("/", "../public/index.html")

	// /wsでハンドリング
	// ip制限をかけるミドルウェアも挟む
	r.GET("/ws", ipBan(internalIpList), wshandler(db, &wsMap, &broadcastChan))

	// 8080でリッスン
	r.Run(":8080")
}
