package redis

import (
	"context"
	"encoding/json"
	"time"

	redis "github.com/go-redis/redis/v8"
)

func (cache *Cacher) Set(key string, value interface{}, expire time.Duration) error {
	c, err := cache.getClient()
	if err != nil {
		return err
	}

	str, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = c.Set(context.Background(), key, str, expire).Err()
	if err != nil {
		if err == redis.Nil {
			// Key does not exists
			return nil
		} else {
			return err
		}
	}

	return nil
}

func (cache *Cacher) MSet(kv map[string]interface{}) error {
	c, err := cache.getClient()
	if err != nil {
		return err
	}

	pairs := []interface{}{}
	for k, v := range kv {

		str, ok := v.(string)
		// Check empty string if value string
		if ok && len(str) == 0 {
			pairs = append(pairs, k, "")
			continue
		}
		// If value is string, not pass it to json.Marshal
		if len(str) > 0 {
			pairs = append(pairs, k, str)
			continue
		}

		strbyte, err := json.Marshal(v)
		if err != nil {
			return err
		}
		pairs = append(pairs, k, strbyte)
	}

	err = c.MSet(context.Background(), pairs...).Err()
	if err != nil {
		return err
	}
	return nil
}

func (cache *Cacher) Get(key string) (string, error) {
	c, err := cache.getClient()
	if err != nil {
		return "", err
	}

	val, err := c.Get(context.Background(), key).Result()
	if err == redis.Nil {
		// Key does not exists
		return "", nil
	} else if err != nil {
		return "", err
	}

	return val, nil
}

func (cache *Cacher) MGet(keys []string) ([]interface{}, error) {
	c, err := cache.getClient()
	if err != nil {
		return nil, err
	}

	vals, err := c.MGet(context.Background(), keys...).Result()
	if err == redis.Nil {
		// Key does not exists
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return vals, nil
}

// Expires set expiration for objects in cache
// if there is error happen, just return last error
func (cache *Cacher) Expires(keys []string, expire time.Duration) error {
	return cache.expires(keys, expire)
}

// Expire set expiration for object in cache
func (cache *Cacher) Expire(key string, expire time.Duration) error {
	return cache.expires([]string{key}, expire)
}

// Expires set expiration for objects in cache
// if there is error happen, just return last error
func (cache *Cacher) expires(keys []string, expire time.Duration) error {
	c, err := cache.getClient()
	if err != nil {
		return err
	}

	var lastErr error
	for _, key := range keys {
		err = c.Expire(context.Background(), key, expire).Err()
		if err != nil {
			if err == redis.Nil {
				// Key does not exists
				return nil
			} else {
				lastErr = err
			}
		}
	}
	return lastErr
}

// Del the cache by keys
func (cache *Cacher) Del(keys ...string) error {
	if len(keys) == 0 {
		return nil
	}

	c, err := cache.getClient()
	if err != nil {
		return err
	}

	// Delete 10000 items per page
	pageLimit := 10000
	from := 0
	to := pageLimit

	for {
		// Lower bound
		if from >= len(keys) {
			break
		}
		// Upper bound
		if to > len(keys) {
			to = len(keys)
		}

		delKeys := keys[from:to]
		if len(delKeys) == 0 {
			break
		}

		_, err = c.Del(context.Background(), delKeys...).Result()
		if err != nil {
			if err == redis.Nil {
				continue
			} else {
				return err
			}
		}
		from += pageLimit
		to += pageLimit
	}

	return nil
}
