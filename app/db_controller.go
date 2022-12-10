package main

import (
	"app/pkg_dbinit"

	"gorm.io/gorm"
)

const MSG_MAX = 30

// 受け取った1件のメッセージをDBに書きこむ
// MSG_MAX件を超えていたら論理削除する
// 引数message にID, createdAtの情報が付加される
func write_db(db *gorm.DB, message *pkg_dbinit.Message) error {
	// メッセージを書きこみ
	err := db.Create(message).Error
	if err != nil {
		return err
	}
	// 最新MSG_MAX件から溢れたものを論理削除
	// DELETE FROM Messages WHERE id NOT in ( SELECT id FROM Messages ORDER BY CreatedAt ASC LIMIT 30 )
	type Id struct {
		ID uint
	}
	var search_ids []Id
	err = db.Select("id").Model(&pkg_dbinit.Message{}).Order("created_at ASC").Limit(MSG_MAX).Find(&search_ids).Error
	if err != nil {
		return err
	}
	var messages []pkg_dbinit.Message
	err = db.Not(search_ids).Delete(&messages).Error
	if err != nil {
		return err
	}

	return err
}

// 受け取ったスライスにdb内のメッセージをすべて突っ込む
func read_db(db *gorm.DB, messages *[]pkg_dbinit.Message) error {

	err := db.Find(messages).Error

	return err
}
