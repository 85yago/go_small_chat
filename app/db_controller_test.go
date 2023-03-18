// dbの読み書きのAPIを記述する

package main

import (
	"app/pkg_dbinit"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetTestTx(t *testing.T) *gorm.DB {
	USER := os.Getenv("POSTGRES_USER")
	PASS := os.Getenv("POSTGRES_PASSWORD")
	DBNAME := "postgres"

	dsn := "host=db user=" + USER + " password=" + PASS + " dbname=" + DBNAME + " port=5432 sslmode=disable TimeZone=Asia/Tokyo"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// トランザクションの開始
	tx := db.Begin()
	// 呼び出したテスト終了後のロールバック
	t.Cleanup(func() { db.Rollback() })

	return tx
}

func msg2str(m pkg_dbinit.Message) string {
	return "[@" + m.Name + "]" + "[" + m.CreatedAt.Format(time.RFC3339) + "]" + m.Message
}

// 受け取った1件のメッセージをDBに書きこむ
// MSG_MAX件を超えていたら論理削除する
// 引数message にID, createdAtの情報が付加される
func TestWrite_db(t *testing.T) {
	db := GetTestTx(t)

	var message pkg_dbinit.Message
	db.First(&message)
	println(msg2str(message))

	assert.Equal(t, 200, 200)
}

// 受け取ったスライスにdb内のメッセージをすべて突っ込む
func TestRead_db(t *testing.T) {
	// db := GetTestTx(t)
}
