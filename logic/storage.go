package logic

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/sony/sonyflake"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type Storage interface {
	Put(ctx context.Context, value string, duration time.Duration) (key string, err error)
	Get(ctx context.Context, key string) (value string, err error)
}

type defaultStorage struct {
	redisCli *redis.ClusterClient // redis cli
	sf       *sonyflake.Sonyflake // unique id generator
	logger   *zap.Logger
}

func NewDefaultStorage(addr []string, password string, logger *zap.Logger) Storage {
	cli := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    addr,
		Password: password,
	})
	sf := sonyflake.NewSonyflake(
		sonyflake.Settings{
			StartTime:      time.Time{},
			MachineID:      nil,
			CheckMachineID: nil,
		})
	return &defaultStorage{
		cli,
		sf,
		logger,
	}
}

// Put 保存值
func (d *defaultStorage) Put(ctx context.Context, value string, duration time.Duration) (key string, err error) {
	id, err := d.sf.NextID()
	if err != nil {
		d.logger.Error("error to generate id", zap.Error(err))
		return key, errors.New("server error")
	}
	key = strconv.FormatUint(id, 10)
	err = d.redisCli.Set(ctx, key, value, duration).Err()
	if err != nil {
		d.logger.Error("error to set key", zap.Error(err))
		return "", errors.New("server error")
	}
	return
}

// Get 获取相关值
func (d *defaultStorage) Get(ctx context.Context, key string) (value string, err error) {
	value, err = d.redisCli.Get(ctx, key).Result()
	return
}
