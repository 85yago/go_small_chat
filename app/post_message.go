// クライアントが叩くpostMessageのサーバ側実装を記述する

package main

import (
	"app/pkg_dbinit"
	"regexp"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

// メッセージの最大文字数
const MSG_MAX_LEN = 100

// ユーザーネームの最大文字数
const NAME_MAX_LEN = 10

// 禁止メッセージを設定する正規表現文字列
const NG_MSG = "^[\\s　]*$"

// 禁止メッセージ判定用の正規表現をコンパイルする関数
func ngMsgCompile() *regexp.Regexp {
	return regexp.MustCompile(NG_MSG)
}

// 禁止メッセージか判定する関数
func isNgMsg(reg *regexp.Regexp, str string) bool {
	return reg.Match([]byte(str))
}

// APIのバックエンド側の実装
// クライアントがメッセージ投稿に使い、メッセージをサーバーに送信する
// サーバーは受け取ったメッセージをDBに登録
// DBへの登録に成功した場合はブロードキャストチャンネルにデータを流し込む
// DBへの登録に失敗した場合はエラーメッセージを返す
func postMessage(c *gin.Context, db *gorm.DB, ws *websocket.Conn, broadcastChan chan<- RetMessage, clientMsg ClientMessage) PostRetMessage {
	// 空欄のメッセージは拒否する
	// TODO:禁止メッセージコンパイルをmainで行うように変更する
	// 禁止メッセージの正規表現のコンパイル
	ngReg := ngMsgCompile()
	if isNgMsg(ngReg, clientMsg.Message) {
		return PostRetMessage{Status: "message error"}
	}
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
		ret.Status = "OK"
		if utf8.RuneCountInString(writedata.Message) > MSG_MAX_LEN {
			// Messageが最大文字数を超えている場合の処理
			ret.Status = "Your message or name is too large."
			writedata.Message = string([]rune(writedata.Message)[:MSG_MAX_LEN-1]) + "…"
		}
		if utf8.RuneCountInString(writedata.Name) > NAME_MAX_LEN {
			// Nameが最大文字数を超えている場合の処理
			ret.Status = "Your message or name is too large."
			writedata.Name = string([]rune(writedata.Name)[:NAME_MAX_LEN-1]) + "…"
		}
		// 書き込みに成功したらチャンネルにメッセージを流し込む
		// writedata にIDやcreatedAtが自動で付与される
		broadcastChan <- RetMessage{Name: writedata.Name, Message: writedata.Message, CreatedAt: writedata.CreatedAt, wsconn: ws}
	}

	return ret
}
