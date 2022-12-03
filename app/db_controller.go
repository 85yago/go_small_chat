package main

import (
	"app/pkg_dbinit"
	"strconv"

	"gorm.io/gorm"
)

const MSG_MAX = 30

func write_db(db gorm.DB, message pkg_dbinit.Message) string {
	db.Create(message)

	// 最新MSG_MAX件から溢れたものを論理削除
	db.Not(db.Order("CreatedAt desc limit " + strconv.Itoa(MSG_MAX))).Delete(&pkg_dbinit.Message{})

	return "ここにエラーメッセージ"
}

func read_db(db gorm.DB)
