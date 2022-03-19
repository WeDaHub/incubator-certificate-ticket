package main

import (
	"context"
	"fmt"
	"net/http"
	"qcloud-tools/certificate"
	"qcloud-tools/core"
	"qcloud-tools/core/config"
	"qcloud-tools/services"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())

	go core.SignalHandler(cancel)

	// 开启一个定时器
	go certificate.TickerSchedule(ctx)

	http.HandleFunc("/login", services.CheckLogin)
	http.HandleFunc("/sync/add", services.AddSync)
	http.HandleFunc("/info/add", services.AddDomain)
	http.HandleFunc("/", services.GetList)

	addr := fmt.Sprintf(":%d", config.QcloudTool.Http.Port)

	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Println(err)
	}
}
