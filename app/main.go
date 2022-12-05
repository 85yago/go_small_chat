package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	// TODO: ここ絶対直すこと
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	// TODO: DBの初期化をする
	r := gin.Default()

	// broadcast用のmap生やしてGETの中でwsを保存しておく
	var wsMap map[int]*websocket.Conn

	// TODO: broadcast用の関数を書いてチャネルを渡しておく、でそれに更新来たメッセージを投げる

	// このmapをgoroutineで回してbroadcast、これは更新があったら回すのを生やすって感じでよさそう？要検討
	go func(wsMap *map[int]*websocket.Conn, c chan string) {
		// チャネルにメッセージが放り込まれるの待ち
		// interface定義してちゃんとそっちでやるとWriteJSONが使えると思う
		mess := <-c
		for _, ws := range *wsMap {
			err := ws.WriteMessage(websocket.TextMessage, []byte(mess))
			if err != nil {
				fmt.Println(err)
				break
			}
		}
	}(&wsMap) // 引数不足、チャネルを作って渡すこと
	// closeをどうするか->とりあえず後回し、wsが生きてるか確認できるのでは？
	// ping pongで確認できるじゃん、送る前に確認する or 別ルーチンで確認と削除

	r.GET("/", func(c *gin.Context) {
		// websocketで接続
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer ws.Close()

		for {
			// メッセージを読む
			mt, message, err := ws.ReadMessage()
			if err != nil {
				fmt.Println(err)
				break
			}

			// jsonのパース
			var rawJson map[string]any
			err = json.Unmarshal(message, &rawJson)
			if err != nil {
				fmt.Println(err)
				break
			}

			// クライアントがどっちを呼び出してるか判定
			switch rawJson["method"] {
			case "getMessage":
				// TODO: getMessage関数に変える
				message = []byte("get")
			case "postMessage":
				// TODO: postMessage関数に変える
				message = []byte("post")
			default:
				message = []byte(c.ClientIP())
			}

			// クライアントに返す
			err = ws.WriteMessage(mt, message)
			if err != nil {
				fmt.Println(err)
				break
			}
		}
	})

	r.Run(":8080") // 8080でリッスン
}
