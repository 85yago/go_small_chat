package pkg_dbinit

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DBに格納するデータの定義
// https://gorm.io/ja_JP/docs/models.html#gorm-Model
type Message struct {
	gorm.Model
	Name      string `gorm:"not null"`
	Message   string `gorm:"not null"`
	IpAddress string `gorm:"not null"`
}

// gorm DBの生成
func DbInitialization() *gorm.DB {
	// https://gorm.io/ja_JP/docs/connecting_to_the_database.html
	// https://qiita.com/chan-p/items/cf3e007b82cc7fce2d81
	USER := os.Getenv("POSTGRES_USER")
	PASS := os.Getenv("POSTGRES_PASSWORD")
	DBNAME := "postgres"

	dsn := "host=db user=" + USER + " password=" + PASS + " dbname=" + DBNAME + " port=5432 sslmode=disable TimeZone=Asia/Tokyo"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&Message{})

	// TODO: デバッグ用
	// IDの連番の開始を1にリセットしている
	// db.Exec("SELECT setval ('messages_id_seq', 1, false)")s

	return db
}
