package dict

import (
	"Calla/store/vo"
	"errors"
)

//Dict 基于map
type Dict struct {
	list map[string]vo.Entry
}

//New 返回一棵空树
func New() Dict {
	return Dict{list: make(map[string]vo.Entry)}
}

//Put 添加数据
func (dict Dict) Put(entry vo.Entry) error {
	dict.list[entry.Key] = entry
	return nil
}

//Get 添加数据
func (dict Dict) Get(key string) (string, error) {
	for k, v := range dict.list {
		if k == key {
			if !v.IsExpired() {
				return v.Value, nil
			}
			break
		}
	}
	return "", errors.New(key + "'s value is not found")
}

//Del 删除数据
func (dict Dict) Del(key string) error {
	delete(dict.list, key)
	return nil
}
