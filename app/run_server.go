//go:build release

package main

import (
	"log"

	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/acme/autocert"
)

// リリース用関数
func runServer(r *gin.Engine) {
	// TLS用の設定
	m := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("azi.f5.si"),
		Cache:      autocert.DirCache("/var/www/.cache"),
	}

	log.Println("runnin release mode")

	// 443でリッスン
	log.Fatal(autotls.RunWithManager(r, &m))
}
