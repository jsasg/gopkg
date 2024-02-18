package memcached

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

type Memcached struct {
	Cli *memcache.Client
}

func New(address []string) *Memcached {
	return &Memcached{
		Cli: memcache.New(address...),
	}
}

func (m *Memcached) Exists(ctx context.Context, key string) bool {
	if item, err := m.Cli.Get(key); err == nil && item != nil {
		return true
	}

	return false
}

func (m *Memcached) Scan(ctx context.Context, key string, val interface{}) error {
	item, err := m.Cli.Get(key)
	if err != nil {
		return err
	}

	return json.Unmarshal(item.Value, val)
}

func (m *Memcached) Get(ctx context.Context, key string) (string, error) {
	item, err := m.Cli.Get(key)
	if err != nil {
		return "", err
	}

	return string(item.Value), nil
}

func (m *Memcached) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	var val []byte

	if reflect.TypeOf(value).Kind() == reflect.Struct ||
		reflect.TypeOf(value).Kind() == reflect.Ptr ||
		reflect.TypeOf(value).Kind() == reflect.Slice ||
		reflect.TypeOf(value).Kind() == reflect.Array ||
		reflect.TypeOf(value).Kind() == reflect.Map {
		val, _ = json.Marshal(value)
	} else {
		switch v := value.(type) {
		case string:
			val = []byte(v)
		default:
			val = []byte(fmt.Sprint(v))
		}
	}

	return m.Cli.Set(&memcache.Item{Key: key, Value: val, Expiration: int32(ttl / time.Second)})
}

func (m *Memcached) Del(ctx context.Context, key string) error {
	return m.Cli.Delete(key)
}

func (m *Memcached) Client() interface{} {
	return m.Cli
}

func (m *Memcached) Close() error {
	return m.Cli.Close()
}
