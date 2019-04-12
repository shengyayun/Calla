package core

import (
	"Calla/store"
	"fmt"
)

func Run() error {
	//加载配置
	cfg, err := NewConfig();
	//服务退出管道
	quit := make(chan int, 1)
	if err != nil {
		return err
	}
	//初始化存储
	st := store.NewStore(cfg.wal)
	if err := st.Load(); err != nil {
		return err
	}
	//提供http服务
	go func() {
		ha := &HttpAccess{st}
		if err := ha.Listen(cfg.listen); err != nil {
			fmt.Println(err)
		}
		quit <- 1
	}()
	//服务退出
	<-quit
	return nil
}
