package inmemlib

import (
	"encoding/json"
	"errors"

	"github.com/coocood/freecache"
)

const MB = 1024 * 1024

type InMemLibInterface interface {
	Set(key string, value interface{}) error
	Get(key string, unmarshalFn func(val []byte) error) (bool, error)
}

type InMemLib struct {
	client *freecache.Cache
}

func New() InMemLib {
	client := freecache.NewCache(10 * MB)
	return InMemLib{
		client: client,
	}
}

func (m InMemLib) Set(key string, value interface{}) error {
	// This is a custom wrapper, allowing us to add custom logs or metrics here.
	bkey := []byte(key)
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return m.client.Set(bkey, data, 0)
}

func (m InMemLib) Get(key string, unmarshalFn func(val []byte) error) (bool, error) {
	// This is a custom wrapper, allowing us to add custom logs or metrics here.
	bkey := []byte(key)
	val, err := m.client.Get(bkey)
	if errors.Is(err, freecache.ErrNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, unmarshalFn(val)
}
