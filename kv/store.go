package kv

type Store interface {
	Get(key string, valuePtr interface{}) error
	Set(key string, value []byte) error
}
