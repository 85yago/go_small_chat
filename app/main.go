package main

import (
	"app/pkg_dbinit"
	"encoding/json"
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

// Message構造体の必要な情報のみの構造体
type RetMessage struct {
	Name        string
	Message     string
	CreatedTime time.Time
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
func procClient(c *gin.Context, db *gorm.DB, ws *websocket.Conn) {
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
		// TODO: WriteJSONを使うこと
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

	for {
		// チャネルにメッセージが放り込まれるの待ち
		// interface定義してちゃんとそっちでやるとWriteJSONが使えると思う
		mess := <-c

		// map用のロック
		wsMap.RLock()
		for ws := range wsMap.m {
			// 各wsにメッセージを送る
			err := ws.WriteJSON(mess)
			if err != nil {
				fmt.Println(err)
				continue
			}
		}
		wsMap.RUnlock()
	}
}

// ws用のエンドポイントのハンドラ
func wshandler(db *gorm.DB, wsMap *WsMap) func(*gin.Context) {
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
		procClient(c, db, ws)
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
	broadcastChan := make(chan RetMessage)
	// ブロードキャスト用の関数
	go broadcastMsg(&wsMap, broadcastChan)

	// /wsでハンドリング
	r.GET("/ws", wshandler(db, &wsMap))

	// 8080でリッスン
	r.Run(":8080")
}
