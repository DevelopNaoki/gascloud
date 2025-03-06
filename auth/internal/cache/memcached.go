package cache

import (
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"time"
)

type Memcached struct {
	client *memcache.Client
}

func NewMemcached(server string) *Memcached {
	return &Memcached{
		client: memcache.New(server),
	}
}

func (m *Memcached) Get(key string) (string, error) {
	item, err := m.client.Get(key)
	if err == memcache.ErrCacheMiss {
		return "", fmt.Errorf("cacheMiss")
	} else if err != nil {
		return "", err
	}
	return string(item.Value), nil
}

func (m *Memcached) Set(key string, value string) error {
	return m.client.Set(&memcache.Item{Key: key, Value: []byte(value)})
}

func (m *Memcached) SetWithTTL(key string, value string, ttl time.Duration) error {
	return m.client.Set(&memcache.Item{
		Key:        key,
		Value:      []byte(value),
		Expiration: int32(ttl.Seconds()),
	})
}

func (m *Memcached) Delete(key string) error {
	return m.client.Delete(key)
}
