package cache

import (
	"context"
	"time"

	"github.com/jsasg/gopkg/cache/memcached"
	"github.com/jsasg/gopkg/cache/redis"
)

type Store interface {
	Exists(ctx context.Context, key string) bool

	Scan(ctx context.Context, key string, val interface{}) error

	Get(ctx context.Context, key string) (string, error)

	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error

	Del(ctx context.Context, key string) error

	Client() interface{}

	Close() error
}

type Config struct {
	Default string
	Stores  struct {
		Redis     *redis.Config
		Memcached []string
	}
}

type Option func(*Config)

func WithRedisConfig(rconf *redis.Config) Option {
	return func(conf *Config) {
		conf.Default = "redis"
		conf.Stores.Redis = rconf
	}
}

func WithMemcachedConfig(addrs []string) Option {
	return func(conf *Config) {
		conf.Default = "memcached"
		conf.Stores.Memcached = addrs
	}
}

type Cache struct {
	defaultStore string
	stores       map[string]Store
}

func New(options ...Option) *Cache {
	var (
		config = &Config{}
		stores = make(map[string]Store)
	)

	for _, option := range options {
		option(config)
	}

	switch config.Default {
	case "redis":
		stores["redis"] = redis.New(config.Stores.Redis)
	case "memcached":
		stores["memcached"] = memcached.New(config.Stores.Memcached)
	default:
		stores[config.Default] = redis.New(config.Stores.Redis)
	}

	return &Cache{
		defaultStore: config.Default,
		stores:       stores,
	}
}

func (c *Cache) Default() Store {
	return c.stores[c.defaultStore]
}

func (c *Cache) Store(store string) Store {
	return c.stores[store]
}
