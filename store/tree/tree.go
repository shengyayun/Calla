package tree

import (
	"Calla/store/vo"
	"errors"
)

//Branch 枝
type Branch struct {
	father *Branch
	child  *Branch
	data   *Leaf
}

//Leaf 叶
type Leaf struct {
	data vo.Entry
	prev *Leaf
	next *Leaf
}

//Tree 树
type Tree struct {
	root []Branch
}

//New 返回一棵空树
func New() *Tree {
	return &Tree{root: make([]Branch, 0)}
}

//Put 添加数据
func (tree *Tree) Put(entry vo.Entry) error {
	return errors.New("ToDo")
}

//Get 添加数据
func (tree *Tree) Get(key string) (string, error) {
	return "", errors.New("ToDo")
}

//Del 删除数据
func (tree *Tree) Del(key string) error {
	return errors.New("ToDo")
}
