package main

import (
	"app/pkg_dbinit"
	"fmt"

	"gorm.io/gorm"
)

// クライアントに返す用の構造体にメッセージを入れて返す
func getMessage(db *gorm.DB) GetRetMessage {
	// DBを読み込む
	var rawMessages []pkg_dbinit.Message
	err := read_db(db, &rawMessages)
	if err != nil {
		fmt.Println(err)
		// エラーならerrorだけ書いて返す
		return GetRetMessage{Status: err.Error()}
	}

	// retにメッセージを入れる
	var ret GetRetMessage
	for _, mess := range rawMessages {
		ret.Message = append(ret.Message, RetMessage{Name: mess.Name, Message: mess.Message, CreatedTime: mess.CreatedAt})
	}
	ret.Count = len(ret.Message)

	return ret
}
