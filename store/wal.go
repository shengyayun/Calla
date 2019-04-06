package store

import (
	"encoding/json"
	"time"
)

type MethodType uint8

const (
	WAL_METHOD_GET MethodType = 0 //get
	WAL_METHOD_PUT MethodType = 1 //put
	WAL_METHOD_DEL MethodType = 2 //del
)

type Entry struct {
	Key    string `json:"key"`    //键
	Value  string `json:"value"`  //值
	Expire int64  `json:"expire"` //到期时间(0为永不到期)
}

type Wal struct {
	Term   int32      `json:"term"`   //任期
	Id     int32      `json:"id"`     //日志id
	Method MethodType `json:"method"` //操作
	Entry  *Entry     `json:"entry"`
}

//key不存在
type NullError struct {
	err error
}

//是否过期
func (entry *Entry) IsExpired() bool {
	return entry.Expire > 0 && entry.Expire < time.Now().Unix()
}

//转为字符串用来存储
func (wal *Wal) ToString() (string, error) {
	bytes, err := json.Marshal(wal)
	return string(bytes), err
}
