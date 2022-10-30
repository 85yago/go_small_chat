package main

import (
	"app/pkg_dbinit"
	"fmt"
	"strconv"
)

func main() {
	db := pkg_dbinit.DbInitialization()

	delete_err := db.Where("1 = 1").Delete(&pkg_dbinit.Message{})
	fmt.Println(delete_err)
	delete_err = db.Unscoped().Where("1 = 1").Delete(&pkg_dbinit.Message{})
	fmt.Println(delete_err)
	var message pkg_dbinit.Message
	// https://gorm.io/ja_JP/docs/delete.html#%E8%AB%96%E7%90%86%E5%89%8A%E9%99%A4%E3%81%95%E3%82%8C%E3%81%9F%E3%83%AC%E3%82%B3%E3%83%BC%E3%83%89%E3%82%92%E5%8F%96%E5%BE%97%E3%81%99%E3%82%8B
	db.Unscoped().Last(&message)
	fmt.Println("db.Unscoped().Last(&message)")
	fmt.Println(message.ID)
	fmt.Println(message.Name)
	fmt.Println(message.Message)
	fmt.Println(message.IpAddress)
	fmt.Println(message.CreatedAt)

	db.Create(&pkg_dbinit.Message{Name: "hoge", Message: "†💩† THE GRAVE OF UNCHI " + strconv.Itoa(int(message.ID)+1), IpAddress: "192.0.2.0"})

	var messages []pkg_dbinit.Message

	db.Find(&messages)
	for _, m := range messages {
		fmt.Println(m.ID)
		fmt.Println(m.Name)
		fmt.Println(m.Message)
		fmt.Println(m.IpAddress)
		fmt.Println(m.CreatedAt)
	}

}
