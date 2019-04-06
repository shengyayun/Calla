package core

import (
	"Calla/store"
)

func Run() error {
	//加载配置
	c, err := NewConfig()
	if err != nil {
		return err
	}
	//初始化存储
	s := store.NewStore(c.wal)
	err = s.Load()
	if err != nil {
		return err
	}
	//return s.Put(&store.Entry{"name", "shengyayun", 0})
	return nil
}
