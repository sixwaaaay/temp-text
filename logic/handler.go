package logic

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

func QueryAPI(storage Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.Query("tid")
		value, err := storage.Get(c.Request.Context(), rid)
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
		log.Println(content)
		key, err := storage.Put(c.Request.Context(), content, time.Minute)
		if err != nil {
			log.Println(err)
			c.String(http.StatusInternalServerError, "fail")
		}
		c.String(http.StatusOK, key)
	}
}
