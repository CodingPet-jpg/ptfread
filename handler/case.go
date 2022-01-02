package handler

import (
	"container/list"
)

type (
	Case struct {
		Name string
		*list.List
	}

	CaseChain struct {
		*list.List
	}
)

func NewCaseChain() CaseChain {
	return CaseChain{List: list.New()}
}

func NewCase(path string) Case {
	return Case{Name: path, List: list.New()}
}

func (cc CaseChain) EliAppend(c Case) {
	// iterate c(Type Case)'s entry
	for entry := c.Front(); entry != nil; entry = entry.Next() {
		// iterate Case
		for e := cc.Front(); e != nil; e = e.Next() {
			// case waiting process
			current := e.Value.(Case)
			// if entry can be found in current case,then remove it
			if ent, ok := current.Contain(entry.Value.(string)); ok {
				current.Remove(ent)
				// determine if the case has zero entry
				if current.Len() == 0 {
					cc.Remove(e)
				}
				break
			}
		}
	}
	cc.PushBack(c)
}

func (c Case) Contain(s string) (target *list.Element, ok bool) {
	for entry := c.Front(); entry != nil; entry = entry.Next() {
		if entry.Value.(string) == s {
			ok = true
			target = entry
			return
		}
	}
	return
}
