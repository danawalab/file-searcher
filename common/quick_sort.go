package common

import (
	"math/rand"

	"gitlab.danawa.com/devops/file-searcher/model"
)

// 퀵 소트 알고리즘
func QuickSortList(a []model.Index) []model.Index {

	if len(a) < 2 {
		return a
	}

	left, right := 0, len(a)-1

	pivot := rand.Int() % len(a)

	a[pivot], a[right] = a[right], a[pivot]

	for i, _ := range a {
		if a[i].ProductId < a[right].ProductId {
			a[left], a[i] = a[i], a[left]
			left++
		}
	}

	a[left], a[right] = a[right], a[left]

	QuickSortList(a[:left])
	QuickSortList(a[left+1:])

	return a

}
