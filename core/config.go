package core

import "os"

type Config struct {
	root string
	wal  string
	http string
}

func NewConfig() (*Config, error) {
	c := Config{}
	//文件目录
	c.root = "./data"
	c.wal = c.root + "/wal"
	//创建工作目录
	for _, path := range []string{c.root, c.wal} {
		_, err := os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				err := os.Mkdir(path, os.FileMode(0777))
				if err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		}
	}
	//http服务
	c.http = "127.0.0.1:8000"
	return &c, nil
}
