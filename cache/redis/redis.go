package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Select          int
	Username        string
	Password        string
	MinIdleConns    int
	MaxIdleConns    int
	ConnMaxLifetime int
	Address         string
	Port            int
}

type Redis struct {
	Cli *redis.Client
}

func New(config *Config) *Redis {
	options := &redis.Options{
		Addr: "127.0.0.1:2379",
		DB:   config.Select,
	}

	if config.Address != "" {
		options.Addr = fmt.Sprintf("%s:2379", config.Address)
	}
	if config.Port != 0 {
		ridx := strings.Index(options.Addr, ":")
		options.Addr = fmt.Sprintf("%s:%d", config.Address[:ridx], config.Port)
	}
	options.Username = config.Username
	options.Password = config.Password
	options.MinIdleConns = config.MinIdleConns
	options.MaxIdleConns = config.MaxIdleConns
	options.ConnMaxLifetime = time.Duration(config.ConnMaxLifetime)

	return &Redis{
		Cli: redis.NewClient(options),
	}
}

func (r *Redis) Exists(ctx context.Context, key string) bool {
	exists, _ := r.Cli.Exists(ctx, key).Result()
	return exists > 0
}

func (r *Redis) Scan(ctx context.Context, key string, val interface{}) error {
	b, err := r.Cli.Get(ctx, key).Bytes()

	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, val); err != nil {
		return err
	}

	return nil
}

func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	return r.Cli.Get(ctx, key).Result()
}

func (r *Redis) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return r.Cli.Set(ctx, key, value, ttl).Err()
}

func (r *Redis) Del(ctx context.Context, key string) error {
	return r.Cli.Del(ctx, key).Err()
}

func (r *Redis) Client() interface{} {
	return r.Cli
}

func (r *Redis) Close() error {
	return r.Cli.Close()
}
