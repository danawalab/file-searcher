package test

import (
	"fmt"
	"testing"

	"gitlab.danawa.com/devops/file-searcher/common"
)

func TestBinarySearch(t *testing.T) {
	items := []string{"AAAA", "BBBB", "CCCC", "DDDD", "EEEE", "FFFF", "GGGG", "HHHH", "III"}

	fmt.Println("있는 경우 검색 테스트 : ", items[common.BinarySearchString(items, "AAAA")])
	fmt.Println("있는 경우 검색 테스트 : ", items[common.BinarySearchString(items, "BBBB")])
	fmt.Println("있는 경우 검색 테스트 : ", items[common.BinarySearchString(items, "CCCC")])
	fmt.Println("있는 경우 검색 테스트 : ", items[common.BinarySearchString(items, "DDDD")])
	fmt.Println("있는 경우 검색 테스트 : ", items[common.BinarySearchString(items, "III")])

	// 결과
	// 있는 경우 검색 테스트 :  AAAA
	// 있는 경우 검색 테스트 :  BBBB
	// 있는 경우 검색 테스트 :  CCCC
	// 있는 경우 검색 테스트 :  DDDD
	// 있는 경우 검색 테스트 :  III

	fmt.Println("없는 경우 검색 테스트 : ", items[common.BinarySearchString(items, "Z")])
	fmt.Println("없는 경우 검색 테스트 : ", items[common.BinarySearchString(items, "K")])
	fmt.Println("없는 경우 검색 테스트 : ", items[common.BinarySearchString(items, "GGG")])
	fmt.Println("없는 경우 검색 테스트 : ", items[common.BinarySearchString(items, "I")])
	fmt.Println("없는 경우 검색 테스트 : ", items[common.BinarySearchString(items, "AAA")])

	// 결과
	// 없는 경우 검색 테스트 :  III
	// 없는 경우 검색 테스트 :  III
	// 없는 경우 검색 테스트 :  FFFF
	// 없는 경우 검색 테스트 :  HHHH
	// 없는 경우 검색 테스트 :  AAAA
}
