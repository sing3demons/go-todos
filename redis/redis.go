package redis

import (
	"context"
	"fmt"
	"sync"
	"time"

	redis "github.com/go-redis/redis/v8"
)

// ICacher is the interface for cache service
type ICacher interface {
	Set(key string, value interface{}, expire time.Duration) error
	MSet(kv map[string]interface{}) error
	Get(key string) (string, error)
	MGet(keys []string) ([]interface{}, error)
}

// Cacher is the struct for cache service
type Cacher struct {
	config      ICacherConfig
	clientMutex sync.Mutex
	client      *redis.Client
	oldClients  []*redis.Client
	subsribers  *sync.Map
}

// NewCacher return new Cacher
func NewCacher(config ICacherConfig) *Cacher {
	return &Cacher{
		config:     config,
		oldClients: nil,
		subsribers: &sync.Map{},
	}
}

func (cache *Cacher) newClient() *redis.Client {
	cfg := cache.config
	settings := cfg.ConnectionSettings()
	return redis.NewClient(&redis.Options{
		Addr:               cfg.Endpoint(),
		Password:           cfg.Password(),
		DB:                 cfg.DB(),
		PoolSize:           settings.PoolSize(),
		MinIdleConns:       settings.MinIdleConns(),
		MaxRetries:         settings.MaxRetries(),
		MinRetryBackoff:    settings.MinRetryBackoff(),
		MaxRetryBackoff:    settings.MaxRetryBackoff(),
		IdleTimeout:        settings.IdleTimeout(),
		IdleCheckFrequency: settings.IdleCheckFrequency(),
		PoolTimeout:        settings.PoolTimeout(),
		ReadTimeout:        settings.ReadTimeout(),
		WriteTimeout:       settings.WriteTimeout(),
	})
}

func (cache *Cacher) getClient() (*redis.Client, error) {
	cache.clientMutex.Lock()
	defer cache.clientMutex.Unlock()

	retriesDelayMs := cache.getRetriesDelayInMs()
	retries := -1
	for {
		retries++
		if retries > len(retriesDelayMs)-1 {
			return nil, fmt.Errorf("cacher: retry exceed limits")
		}

		client := cache.client
		if client == nil {
			client = cache.newClient()
			cache.client = client
		}

		_, err := client.Ping(context.Background()).Result()
		if err != nil {
			// Wait by retry delay then reset client and try connect again
			time.Sleep(time.Millisecond * time.Duration(retriesDelayMs[retries]))
			cache.client = nil
			continue
		}

		// If we can PING without error, just return
		return client, nil
	}
}

// Close close the redis client
func (cache *Cacher) Close() error {
	cache.clientMutex.Lock()
	defer cache.clientMutex.Unlock()

	// Close current client
	client := cache.client
	if client != nil {
		cache.client = nil

		err := client.Close()
		if err != nil {
			return err
		}

		// Close old clients
		for _, client := range cache.oldClients {
			err := client.Close()
			if err != nil {
				return err
			}
		}
		if len(cache.oldClients) > 0 {
			cache.oldClients = nil
		}
	}

	return nil
}

// getRetriesDelayInMs sum only 1 second
func (cache *Cacher) getRetriesDelayInMs() []int {
	return []int{200, 200, 200, 200, 200}
}
