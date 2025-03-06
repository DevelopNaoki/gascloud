package cache

import (
	"fmt"
	"strconv"
	"time"

	"github.com/DevelopNaoki/gascloud/auth/internal/model"
)

type Cache interface {
	Get(key string) (string, error)
	Set(key string, value string) error
	SetWithTTL(key string, value string, ttl time.Duration) error
	Delete(key string) error
}

func NewCache(conf model.CacheConfig) (client Cache, err error) {
	switch conf.Driver {
	case "memcache":
		client = NewMemcached(conf.Host + ":" + strconv.Itoa(conf.Port))
	default:
		return client, fmt.Errorf("unsupported cache driver")
	}

	return client, nil
}
