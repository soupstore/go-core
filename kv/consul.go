package kv

import (
	"encoding/json"
	consul "github.com/hashicorp/consul/api"
)

type ConsulStore struct {
	kv *consul.KV
}

func NewConsulStore() (*ConsulStore, error) {
	consulClient, err := consul.NewClient(consul.DefaultConfig())
	if err != nil {
		return nil, err
	}

	kv := consulClient.KV()

	return &ConsulStore{kv}, nil
}

func (c *ConsulStore) Get(key string, valuePtr interface{}) error {
	pair, _, err := c.kv.Get(key, nil)
	if err != nil {
		return err
	}

	switch v := valuePtr.(type) {
	case []byte:
		v = pair.Value
	case *string:
		*v = string(pair.Value)
	default:
		if err := json.Unmarshal(pair.Value, valuePtr); err != nil {
			return err
		}
	}

	return nil
}

func (c *ConsulStore) Set(key string, value []byte) error {
	_, err := c.kv.Put(&consul.KVPair{
		Key:   key,
		Value: value,
	}, nil)

	return err
}
