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
			c.JSON(http.StatusBadRequest, Resp[*string]{
				Code: http.StatusBadRequest,
				Msg:  "require parameter tid",
			})
			return
		}
		value, err := storage.Get(c.Request.Context(), tid)
		if err != nil {
			logger.Error("get failed", zap.Error(err))
			c.JSON(http.StatusNotFound, Resp[*string]{
				Code: http.StatusNotFound,
				Msg:  "not found",
			})
			return
		}
		c.JSON(http.StatusOK, Resp[string]{
			Code: 0,
			Msg:  "success",
			Data: value,
		})
	}
}

func ShareAPI(logger *zap.Logger, storage Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		content := c.PostForm("content")
		if len(content) == 0 {
			c.JSON(http.StatusBadRequest, Resp[*string]{
				Code: http.StatusBadRequest,
				Msg:  "require parameter content",
			})
			return
		}
		key, err := storage.Put(c.Request.Context(), content, time.Minute)
		if err != nil {
			logger.Error("put failed", zap.Error(err))
			c.JSON(http.StatusInternalServerError, Resp[*string]{
				Code: http.StatusInternalServerError,
				Msg:  "fail",
			})
			return
		}
		c.JSON(http.StatusOK, Resp[string]{
			Code: 0,
			Msg:  "success",
			Data: key,
		})
	}
}
