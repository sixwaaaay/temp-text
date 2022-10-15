package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sixwaaaay/temp-text/grace"
	"github.com/sixwaaaay/temp-text/logic"
	"github.com/spf13/viper"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

type Conf struct {
	Redis struct {
		Addr     []string `json:"addr"`
		Password string   `json:"password"`
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

	p := ginprometheus.NewPrometheus("gin")
	p.Use(r)

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	storage := logic.NewDefaultStorage(conf.Redis.Addr, conf.Redis.Password)
	r.GET("/query", logic.QueryAPI(storage))
	r.POST("/share", logic.ShareAPI(storage))

	srv := &http.Server{Addr: conf.ApiAddr, Handler: r}

	endless := grace.NewEndless(func() { // 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}, func() {
		log.Println("Shutdown Server ...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatal("Server Shutdown:", err)
		}
		log.Println("Server exiting")
	}, func() chan os.Signal {
		quit := make(chan os.Signal)
		signal.Notify(quit, os.Interrupt)
		return quit
	})
	endless.Run()
}
