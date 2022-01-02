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

func (chain CaseChain) EliAppend(c Case) {
	// iterate entry of case waiting append
	for entry := c.Front(); entry != nil; entry = entry.Next() {
		// iterate all case in current chain
		for tcase := chain.Front(); tcase != nil; tcase = tcase.Next() {
			// current case in chain which waiting to be processed
			current := tcase.Value.(Case)
			// if entry can be found in current case,then remove this entry
			if ent, ok := current.Contain(entry.Value.(string)); ok {
				current.Remove(ent)
				// each time entry being removed,we determine if the case has zero entry
				if current.Len() == 0 {
					chain.Remove(tcase)
				}
				break
			}
		}
	}
	chain.PushBack(c)
}

// determine if the entry which is identical to s can be found in case,return entry and true otherwise nil and false

func (tcase Case) Contain(s string) (target *list.Element, ok bool) {
	for entry := tcase.Front(); entry != nil; entry = entry.Next() {
		if entry.Value.(string) == s {
			ok = true
			target = entry
			return
		}
	}
	return
}
