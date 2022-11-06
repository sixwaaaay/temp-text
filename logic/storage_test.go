package logic

import (
	"context"
	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"log"
	"testing"
	"time"
)

func TestPut(t *testing.T) {
	s := miniredis.RunT(t)
	storage := NewDefaultStorage([]string{s.Addr()}, "", zap.L())
	ctx := context.Background()
	key, err := storage.Put(ctx, "test", time.Second)
	assert.NoError(t, err)
	assert.NotEmpty(t, key)
	log.Println(key)
}

func TestGetInTime(t *testing.T) {
	s := miniredis.RunT(t)
	storage := NewDefaultStorage([]string{s.Addr()}, "", zap.L())
	ctx := context.Background()
	v := "test"
	key, err := storage.Put(ctx, v, time.Second)
	assert.NoError(t, err)
	assert.NotEmpty(t, key)
	value, err := storage.Get(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, v, value)
}

func TestGetNotInTime(t *testing.T) {
	s := miniredis.RunT(t)
	storage := NewDefaultStorage([]string{s.Addr()}, "", zap.L())
	ctx := context.Background()
	v := "test"
	key, err := storage.Put(ctx, v, time.Second)
	assert.NoError(t, err)
	assert.NotEmpty(t, key)
	s.FastForward(time.Second) // to expire
	value, err := storage.Get(ctx, key)
	assert.Error(t, err)
	assert.Empty(t, value)
}
