package logic

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func QueryAPI(logger *zap.Logger, storage Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		tid := c.Query("tid")
		if len(tid) == 0 {
			c.String(http.StatusBadRequest, "require parameter tid")
			return
		}
		value, err := storage.Get(c.Request.Context(), tid)
		if err != nil {
			logger.Error("get failed", zap.Error(err))
			c.String(http.StatusNotFound, "not found")
			return
		}
		c.String(http.StatusOK, value)
	}
}

func ShareAPI(logger *zap.Logger, storage Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		content := c.PostForm("content")
		if len(content) == 0 {
			c.String(http.StatusBadRequest, "require parameter content")
			return
		}
		key, err := storage.Put(c.Request.Context(), content, time.Minute)
		if err != nil {
			logger.Error("put failed", zap.Error(err))
			c.String(http.StatusInternalServerError, "fail")
		}
		c.String(http.StatusOK, key)
	}
}
