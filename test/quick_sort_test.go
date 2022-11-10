package test

import (
	"fmt"
	"testing"

	"gitlab.danawa.com/devops/file-searcher/common"
	"gitlab.danawa.com/devops/file-searcher/model"
)

var str, str2 []model.Index

func TestQuickSort(t *testing.T) {
	str = append(str, model.Index{ProductId: "0001", Position: 0})
	str = append(str, model.Index{ProductId: "923456", Position: 0})
	str = append(str, model.Index{ProductId: "12345678", Position: 0})
	str = append(str, model.Index{ProductId: "1111", Position: 0})
	str = append(str, model.Index{ProductId: "22345678", Position: 0})
	str = common.QuickSortList(str)

	str2 = append(str2, model.Index{ProductId: "4949", Position: 0})
	str2 = append(str2, model.Index{ProductId: "1000", Position: 0})
	str2 = append(str2, model.Index{ProductId: "10000", Position: 0})
	str2 = append(str2, model.Index{ProductId: "100000", Position: 0})
	str2 = append(str2, model.Index{ProductId: "1000000", Position: 0})
	str2 = common.QuickSortList(str2)

	// 성공..
	for _, ele := range str {
		fmt.Println(ele)
	}

	// // 성공..
	for _, ele := range str2 {
		fmt.Println(ele)
	}
}
