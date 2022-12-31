//go:build release

package main

import (
	"log"

	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/acme/autocert"
)

// リリース用関数
func RunServer(r *gin.Engine) {
	r.StaticFile("/chat.js", "/var/public/chat.js")

	// ginのリリースモード
	gin.SetMode(gin.ReleaseMode)

	// TLS用の設定
	m := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("azi.f5.si"),
		Cache:      autocert.DirCache("/var/www/.cache"),
	}

	// 443でリッスン
	log.Fatal(autotls.RunWithManager(r, &m))
}
