package store

import (
	"Calla/store/vo"
)

//Engine 存储引擎的接口
type Engine interface {
	Put(entry vo.Entry) error
	Get(key string) (string, error)
	Del(key string) error
}
