package main

import "fmt"

type Comparable interface {
	LessThan() Comparable
}

type MyDouble float32

// 如果Mi真的实现了Comparable接口,应该可以返回所有的Comparable接口实现类,比如MyBool或者MyDouble

func (Md MyDouble) LessThan() Comparable {
	return MyDouble(2)
}

type Playable interface {
	String() string
}

type Eatable interface {
	Eat()
}

type ComponentUnion struct {
	Playable
	Eatable
}

func (cu *ComponentUnion) Template() {
	cu.String()
	cu.Eat()
}

var anything interface{}

func main() {
	var t = 1
	var i interface{} = t
	t2 := i.(int)
	t2 = 3
	fmt.Println(t2)
	fmt.Println(i)
	fmt.Print(t)

}
