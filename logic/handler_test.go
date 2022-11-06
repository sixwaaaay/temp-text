package logic

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"
)

type MockStorage struct {
	mapping    map[string]string
	raiseError bool
	key        uint64
}

func (m *MockStorage) Put(_ context.Context, value string, _ time.Duration) (key string, err error) {
	if m.raiseError {
		return "", errors.New("error")
	}
	id := strconv.FormatUint(m.key, 10)
	m.key++
	m.mapping[id] = value
	return id, nil
}

func (m *MockStorage) Get(_ context.Context, key string) (value string, err error) {
	if m.raiseError {
		return "", errors.New("error")
	}
	v, ok := m.mapping[key]
	if !ok {
		return "", errors.New("not exist")
	}
	return v, nil
}

func TestShareAPIOk(t *testing.T) {
	w, c, r := httpTestHelper()
	c.Request, _ = http.NewRequest("POST", "/share", nil)
	c.Request.PostForm = url.Values{
		"content": []string{"test"},
	}
	storage := &MockStorage{
		mapping:    map[string]string{},
		raiseError: false,
		key:        0,
	}
	r.POST("/share", ShareAPI(zap.L(), storage))
	r.ServeHTTP(w, c.Request)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "0", w.Body.String())
}
func TestShareAPIParamError(t *testing.T) {
	w, c, r := httpTestHelper()
	// following request missing param content
	c.Request, _ = http.NewRequest("POST", "/share", nil)
	c.Request.PostForm = url.Values{
		//"content": []string{""},
	}
	storage := &MockStorage{
		mapping:    map[string]string{},
		raiseError: false,
		key:        0,
	}
	r.POST("/share", ShareAPI(zap.L(), storage))
	r.ServeHTTP(w, c.Request)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "require parameter content", w.Body.String())
}
func TestShareAPIFail(t *testing.T) {
	w, c, r := httpTestHelper()
	c.Request, _ = http.NewRequest("POST", "/share", nil)
	c.Request.PostForm = url.Values{
		"content": []string{"test"},
	}
	storage := &MockStorage{
		mapping:    map[string]string{},
		raiseError: false,
		key:        0,
	}
	r.POST("/share", ShareAPI(zap.L(), storage))
	storage.raiseError = true
	r.ServeHTTP(w, c.Request)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "fail", w.Body.String())
}

func TestQueryAPIOk(t *testing.T) {
	storage := &MockStorage{
		mapping:    map[string]string{},
		raiseError: false,
		key:        0,
	}
	testVal := "a quick fox jumps over a lazy dog"
	key, _ := storage.Put(context.Background(), testVal, time.Second)
	w, c, r := httpTestHelper()
	c.Request, _ = http.NewRequest("GET", fmt.Sprintf("/query?tid=%s", key), nil)
	r.GET("/query", QueryAPI(zap.L(), storage))
	r.ServeHTTP(w, c.Request)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, testVal, w.Body.String())
}

func TestQueryAPIFail(t *testing.T) {
	storage := &MockStorage{
		mapping:    map[string]string{},
		raiseError: false,
		key:        0,
	}
	testVal := "a quick fox jumps over a lazy dog"
	key, _ := storage.Put(context.Background(), testVal, time.Second)
	w, c, r := httpTestHelper()
	c.Request, _ = http.NewRequest("GET", fmt.Sprintf("/query?tid=%s", key), nil)
	r.GET("/query", QueryAPI(zap.L(), storage))
	storage.raiseError = true
	r.ServeHTTP(w, c.Request)
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "not found", w.Body.String())
}

func TestQueryAPIParamError(t *testing.T) {
	storage := &MockStorage{
		mapping:    map[string]string{},
		raiseError: false,
		key:        0,
	}
	testVal := "a quick fox jumps over a lazy dog"
	key, _ := storage.Put(context.Background(), testVal, time.Second)
	w, c, r := httpTestHelper()
	// following request do not contain the 'tid' param
	c.Request, _ = http.NewRequest("GET", fmt.Sprintf("/query?x=%s", key), nil)
	r.GET("/query", QueryAPI(zap.L(), storage))
	r.ServeHTTP(w, c.Request)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "require parameter tid", w.Body.String())
}

// httpTestHelper 返回用于测试的三个http相关对象
func httpTestHelper() (*httptest.ResponseRecorder, *gin.Context, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	r := gin.New()
	return w, c, r
}
