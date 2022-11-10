package common

import "gitlab.danawa.com/devops/file-searcher/model"

func BinarySearch(array []model.MiniIndex, target string) (int, bool) {
	var isExist bool

	r := -1 // not found
	start := 0
	end := len(array) - 1
	for start <= end {
		mid := (start + end) / 2
		if array[mid].ProductId == target {
			r = mid // found
			break
		} else if array[mid].ProductId < target {
			start = mid + 1
		} else if array[mid].ProductId > target {
			end = mid - 1
		}
	}

	if r == -1 {
		if end == -1 {
			r = 0
		} else {
			// 없을 경우 최종적으로 읽은 부분을 리턴한다
			r = end
		}
	} else {
		isExist = true
	}

	// 찾은 인덱스를 리턴한다
	return r, isExist
}

func BinarySearchString(array []string, target string) int {
	r := -1 // not found
	start := 0
	end := len(array) - 1
	for start <= end {
		mid := (start + end) / 2
		if array[mid] == target {
			r = mid // found
			break
		} else if array[mid] < target {
			start = mid + 1
		} else if array[mid] > target {
			end = mid - 1
		}
	}

	if r == -1 {
		if end == -1 {
			r = 0
		} else {
			// 없을 경우 최종적으로 읽은 부분을 리턴한다
			r = end
		}
	}

	// 찾은 인덱스를 리턴한다
	return r
}
