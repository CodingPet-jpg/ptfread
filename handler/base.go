package handler

import (
	"container/list"
	"flag"
	"io/fs"
	"os"
)

type LinkedNode struct {
	*list.List
}

func NewLn() LinkedNode {
	return LinkedNode{List: list.New()}
}

// 在文件链表中所有节点的手顺列表都彼此互斥，所以遍历最后一个文件节点的手顺节点并拿条目和之前
// 所有文件节点的手顺节点比较时只要有一次匹配此函数就可以退出，无匹配时会遍历完全部的节点后退出
// 如果出现匹配的情况需要额外的操作，比如删除被匹配文件节点的相应手顺节点条目
func (ln LinkedNode) ComparedAppend(c Case) {
	current := ln.PushBack(c)
	content := current.Value.(Case)
	// 遍历所有代表文件的节点，其中包含了手顺列表
	for j := 0; j < ln.Len(); j++ {
		// 遍历最后一个节点的手顺列表
		// 和前面的文件节点中的手顺列表比较并移除相等的节点
		c2 := current.Value.(Case)
		for s, i := content.Back(), 0; content.Len() > i; i++ {
			if ele, ok := c2.Contain(s.Value.(string)); ok {
				c2.Remove(ele)
				return
			}

			s = s.Next()
		}
		// 如果当前节点的手顺列表为空则移除此文件节点
		if c2.Len() == 0 {
			ln.Remove(current)
			j++
		}
		current = current.Prev()
	}
}

var wd = flag.String("w", ".", "Set WorkDirectory")

var Fs fs.FS

type Case struct {
	Name string
	*list.List
}

func (c Case) Contain(s string) (*list.Element, bool) {
	crtElem := c.Back()
	var target *list.Element
	var ok bool
	for i := 0; i < c.Len(); i++ {
		if crtElem.Value.(string) == s {
			ok = true
			target = crtElem
			break
		}
		crtElem = crtElem.Next()
	}
	return target, ok
}

var BaseSheet string

func init() {
	if !flag.Parsed() {
		flag.Parse()
	}
	BaseSheet = "Sheet1"
	Fs = os.DirFS(*wd)
}

func GetWd() string {
	return *wd
}
