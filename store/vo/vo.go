package vo

import (
	"encoding/json"
	"time"
)

//MethodType 方法类型
type MethodType uint8

const (
	//WalMethodGet 获取数据
	WalMethodGet MethodType = 0
	//WalMethodPut 存放数据
	WalMethodPut MethodType = 1
	//WalMethodDel 删除数据
	WalMethodDel MethodType = 2
)

//Entry 存储对象
type Entry struct {
	Key    string `json:"key"`    //键
	Value  string `json:"value"`  //值
	Expire int64  `json:"expire"` //到期时间(0为永不到期)
}

//Wal 预写日志
type Wal struct {
	Term   int32      `json:"term"`   //任期
	ID     int32      `json:"id"`     //日志
	Method MethodType `json:"method"` //操作
	Entry  Entry      `json:"entry"`
}

//IsExpired 是否过期
func (entry *Entry) IsExpired() bool {
	return entry.Expire > 0 && entry.Expire < time.Now().Unix()
}

//ToString 转为字符串用来存储
func (wal *Wal) ToString() (string, error) {
	bytes, err := json.Marshal(wal)
	return string(bytes), err
}
