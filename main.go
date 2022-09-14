package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"os/signal"
	"tempMsg/logic"
	"time"
)

type Conf struct {
	Redis struct {
		Addr     string `json:"addr"`
		Password string `json:"password"`
		DB       int    `json:"db"`
	}
	ApiAddr string `json:"ApiAddr"`
}

func main() {
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Panicln(err)
	}
	conf := Conf{}
	err = viper.Unmarshal(&conf)
	if err != nil {
		log.Panicln(err)
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	storage := logic.NewDefaultStorage(conf.Redis.Addr, conf.Redis.Password, conf.Redis.DB)
	r.GET("/query", logic.QueryAPI(storage))
	r.POST("/share", logic.ShareAPI(storage))

	srv := &http.Server{
		Addr:    conf.ApiAddr,
		Handler: r,
	}
	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
