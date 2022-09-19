package logic

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

func QueryAPI(storage Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		tid := c.Query("tid")
		if len(tid) == 0 {
			c.String(http.StatusBadRequest, "require parameter tid")
			return
		}
		value, err := storage.Get(c.Request.Context(), tid)
		if err != nil {
			log.Println(err.Error())
			c.String(http.StatusNotFound, "not found")
			return
		}
		c.String(http.StatusOK, value)
	}
}

func ShareAPI(storage Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		content := c.PostForm("content")
		if len(content) == 0 {
			c.String(http.StatusBadRequest, "require parameter content")
			return
		}
		log.Println(content)
		key, err := storage.Put(c.Request.Context(), content, time.Minute)
		if err != nil {
			log.Println(err)
			c.String(http.StatusInternalServerError, "fail")
		}
		c.String(http.StatusOK, key)
	}
}
