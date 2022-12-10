package main

import (
	"app/pkg_dbinit"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// APIのバックエンド側の実装
// クライアントがメッセージ投稿に使い、メッセージをサーバーに送信する
// サーバーは受け取ったメッセージをDBに登録
// DBへの登録に成功した場合はブロードキャストチャンネルにデータを流し込む
// DBへの登録に失敗した場合はエラーメッセージを返す
func postMessage(c *gin.Context, db *gorm.DB, broadcastChan chan<- RetMessage, clientMsg ClientMessage) PostRetMessage {
	// 終了メッセージを格納する変数
	// ("OK" or エラーメッセージ)
	var ret PostRetMessage

	// DBに書き込むデータをclientMessageから転写
	var writedata pkg_dbinit.Message
	writedata.Name = clientMsg.Name
	writedata.Message = clientMsg.Message
	writedata.IpAddress = c.ClientIP()

	// DBに書き込み
	err := write_db(db, &writedata)
	if err != nil {
		ret.Status = err.Error()
	} else {
		// 書き込みに成功したらチャンネルにメッセージを流し込む
		ret.Status = "OK"
		// !writedata にIDやcreatedAtが自動で付与されるのか不明
		broadcastChan <- RetMessage{Name: writedata.Name, Message: writedata.Message, CreatedAt: writedata.CreatedAt}
	}

	return ret
}
