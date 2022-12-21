// websocket通信関係の関数，構造体を記述する

package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

// 各wsに持たせるデータ
// iphash : ipのMD5ハッシュ値（を文字列化したもの）
type WsMapData struct {
	// iphash string
}

// wsの保存用
type WsMap struct {
	sync.RWMutex
	m map[*websocket.Conn]WsMapData
}

// upgraderはHTTPをWSにするときに呼ばれる
// ここで許可するoriginや接続時間を設定する
var upgrader = websocket.Upgrader{
	// TODO: ここ絶対直すこと
	CheckOrigin: func(r *http.Request) bool {
		// r.URL.
		return true
	},
}

// 未実装
// // ipのハッシュを計算して返す関数
// func ipToMd5(ip string) string {
// 	bip := []byte(ip+"salt")
// 	md5ip := md5.Sum(bip)
// 	return hex.EncodeToString(md5ip[:])
// }

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
		wsMap.m[ws] = WsMapData{}
		wsMap.Unlock()

		// クライアントからのwebsocketを処理
		procClient(c, db, ws, broadcastChan)
	}
}
