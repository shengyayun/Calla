package store

import (
	"Calla/store/dict"
	"Calla/store/vo"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

//Store 仓库
type Store struct {
	engine Engine
	path   string
	term   int32
	cursor int32
}

//NewStore 创建仓库实例
func NewStore(path string) *Store {
	eg := dict.New()
	return &Store{&eg, path, 1, 0} //dict中实现接口的方法为指针方法，所以dict的指针对象才满足接口Engine要求
}

//Load 通过wal日志加载数据
func (store *Store) Load() error {
	//查找wal目录下的全部文件
	files, err := ioutil.ReadDir(store.path)
	if err != nil {
		return err
	}
	//通过wal日志加载entry到内存
	for _, path := range files {
		if !strings.HasSuffix(path.Name(), ".wal") { //检查文件后缀
			continue
		}
		f, err := os.Open(store.path + "/" + path.Name())
		if err != nil {
			return err
		}
		buf := bufio.NewReader(f)
		for {
			line, err := buf.ReadBytes('\n')
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}
			wal := vo.Wal{}
			if err = json.Unmarshal(line, &wal); err != nil {
				return err
			}
			store.cursor = wal.ID
			switch wal.Method {
			case vo.WalMethodPut:
				if err = store.engine.Put(wal.Entry); err != nil {
					return err
				}
			case vo.WalMethodDel:
				if err = store.engine.Del(wal.Entry.Key); err != nil {
					return err
				}
			default:
				continue
			}
		}
	}
	return nil
}

//Put 操作
func (store *Store) Put(entry vo.Entry) error {
	if err := store.append(vo.WalMethodPut, entry); err != nil {
		return err
	}
	return store.engine.Put(entry)
}

//Get 操作
func (store *Store) Get(key string) (string, error) {
	return store.engine.Get(key)
}

//Del 操作
func (store *Store) Del(key string) error {
	if err := store.append(vo.WalMethodDel, vo.Entry{key, "", 0}); err != nil {
		return err
	}
	if err := store.engine.Del(key); err != nil {
		return err
	}
	return nil
}

//wal日志添加
func (store *Store) append(method vo.MethodType, entry vo.Entry) error {
	store.cursor++
	wal := vo.Wal{store.term, store.cursor, method, entry}
	path := store.path + "/" + time.Now().Format("2006010215") + ".wal"
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil && os.IsNotExist(err) { //文件不存在
		file, err = os.Create(path)
	}
	if err != nil {
		return err
	}
	defer file.Close()
	if str, err := wal.ToString(); err != nil {
		return err
	} else if _, err = file.Write(append([]byte(str), '\n')); err != nil {
		return err
	}
	return nil
}

//测试用
func (store *Store) Test(quit chan int) {
	if err := store.Put(vo.Entry{Key: "Greet", Value: "Hello", Expire: 0}); err != nil {
		fmt.Println(err)
	}
	fmt.Println(store.Get("Greet"))
	if err := store.Put(vo.Entry{Key: "Greet", Value: "World", Expire: 0}); err != nil {
		fmt.Println(err)
	}
	fmt.Println(store.Get("Greet"))
	if err := store.Del("Greet"); err != nil {
		fmt.Println(err)
	}
	fmt.Println(store.Get("Greet"))
	quit <- 1
}
