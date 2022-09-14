package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"tempMsg/logic"
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
	err = r.Run(conf.ApiAddr)
	if err != nil {
		log.Panicln(err)
	}
}
