//go:build !release

package main

import (
	"github.com/gin-gonic/gin"
)

// デバッグ用関数
func runServer(r *gin.Engine) {
	// 127.0.0.1:8080でリッスン
	r.Run(":8080")
}
