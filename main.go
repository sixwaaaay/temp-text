package main

import (
	"context"
	ginzap "github.com/gin-contrib/zap"
	"github.com/sixwaaaay/temp-text/logic"
	ginprom "github.com/zsais/go-gin-prometheus"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type Conf struct {
	Redis struct {
		Addr     []string `json:"addr"`
		Password string   `json:"password"`
	}
	ApiAddr string `json:"ApiAddr"`
}

func main() {
	fx.New(
		fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: logger}
		}),
		fx.Provide(
			NewLogger,
			NewConfig,
			NewStorage,
			NewHandlers,
			NewRouter,
			NewServer,
		),
		fx.Invoke(NewServer),
	).Run()
}

func NewServer(lc fx.Lifecycle, logger *zap.Logger, router *gin.Engine, conf *Conf) *http.Server {
	server := &http.Server{Addr: conf.ApiAddr, Handler: router}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logger.Fatal("listen: %s\n", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			return server.Shutdown(ctx)
		},
	})
	return server
}

func NewLogger() *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	return logger
}

type Handler struct {
	Method  string
	Path    string
	Handler gin.HandlerFunc
}

func NewHandlers(logger *zap.Logger, storage logic.Storage) []Handler {
	return []Handler{
		{
			Method:  "GET",
			Path:    "/query",
			Handler: logic.QueryAPI(logger, storage),
		},
		{
			Method:  "POST",
			Path:    "/share",
			Handler: logic.ShareAPI(logger, storage),
		},
		{
			Method: "GET",
			Path:   "/ping",
			Handler: func(c *gin.Context) {
				c.String(http.StatusOK, "pong")
			},
		},
	}
}

func NewRouter(logger *zap.Logger, handlers []Handler) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	router.Use(ginzap.RecoveryWithZap(logger, true))
	p := ginprom.NewPrometheus("gin")
	p.Use(router)
	for _, handler := range handlers {
		router.Handle(handler.Method, handler.Path, handler.Handler)
	}
	return router
}

func NewConfig(logger *zap.Logger) *Conf {
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		logger.Panic("read config failed", zap.Error(err))
	}
	conf := Conf{}
	err = viper.Unmarshal(&conf)
	if err != nil {
		logger.Panic("unmarshal config failed", zap.Error(err))
	}
	return &conf
}

func NewStorage(conf *Conf, logger *zap.Logger) logic.Storage {
	return logic.NewDefaultStorage(conf.Redis.Addr, conf.Redis.Password, logger)
}
