//go:build !release

package main

import (
	"github.com/gin-gonic/gin"
)

// デバッグ用関数
func RunServer(r *gin.Engine) {
	r.StaticFile("/chat.js", "/var/public/chat_dev.js")

	// 127.0.0.1:8080でリッスン
	r.Run(":8080")
}
