// メッセージの送受信を行う関数，構造体を記述する

package main

import (
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

// ブロードキャスト用のチャネル
type BroadChan struct {
	sync.RWMutex
	c chan RetMessage
}

// チャネルに送られたメッセージをブロードキャストする関数
func broadcastMsg(wsMap *WsMap, c *BroadChan) {
	for {
		// チャネルにメッセージが放り込まれるの待ち
		cmsg := <-c.c
		var msg SendData
		msg.DataType = BROADCAST
		msg.Data = &cmsg

		// map用のロック
		wsMap.RLock()
		for ws := range wsMap.m {
			// 送信先が送信元と同じならば
			if ws == msg.Data.(*RetMessage).wsconn {
				msg.Data.(*RetMessage).IsMe = true

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
				msg.Data.(*RetMessage).IsMe = false

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
