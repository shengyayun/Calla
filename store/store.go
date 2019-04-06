package store

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type Store struct {
	list   map[string]*Entry
	path   string
	term   int32
	cursor int32
}

func NewStore(path string) *Store {
	return &Store{make(map[string]*Entry), path, 1, 0}
}

//通过wal日志加载数据
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
			wal := Wal{}
			if err = json.Unmarshal(line, &wal); err != nil {
				return err
			}
			store.cursor = wal.Id
			switch wal.Method {
			case WAL_METHOD_PUT:
				store.list[wal.Entry.Key] = wal.Entry
			case WAL_METHOD_DEL:
				delete(store.list, wal.Entry.Key)
			default:
				continue
			}
		}
	}
	return nil
}

//put操作
func (store *Store) Put(entry *Entry) error {
	if err := store.append(WAL_METHOD_PUT, entry); err != nil {
		return err
	}
	store.list[entry.Key] = entry
	return nil
}

//get操作
func (store *Store) Get(key string) (string, error) {
	for k, v := range store.list {
		if k == key {
			if !v.IsExpired() {
				return v.Value, nil
			}
			break
		}
	}
	return "", errors.New(key + "'s value is not found")
}

//del操作
func (store *Store) Del(key string) error {
	if err := store.append(WAL_METHOD_DEL, &Entry{key, "", 0}); err != nil {
		return err
	}
	delete(store.list, key)
	return nil
}

func (store *Store) append(method MethodType, entry *Entry) error {
	store.cursor++
	wal := Wal{store.term, store.cursor, method, entry}
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
