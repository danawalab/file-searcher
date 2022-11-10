package parser

import (
	"io"
	"strconv"
	"strings"

	"gitlab.danawa.com/devops/file-searcher/common"
	"gitlab.danawa.com/devops/file-searcher/model"
	"gitlab.danawa.com/devops/file-searcher/setting"
)

type CustomParser struct {
}

// 멀티라인 전용 커스텀 파서..
func (customParser CustomParser) MultiLineParse(format map[string]string, reader *common.Reader, epFilePath, tmpFilePath, fileName string) {
	var startPosition, endPosition, lineBreakerCount, indexLineLength, chunkGubun, lineCount int
	var docuOpen bool
	var productId string
	var err error
	var byteLine []byte
	var items []model.Index

	const ValueIdx = 1

	// 파일을 쪼갤 청크 단위
	indexChunkDivision := setting.GetFileSortDivision()

	// 경량화 인덱스 색인 간격
	indexFileInterval := setting.GetIndexFileInterval()

	// 이원화 처리를 위해 두 개의 인덱스 청크 파일을 준비한다
	chunkFileA, chunkFileWriterA, chunkFileB, chunkFileWriterB := common.GetChunkFileWriter(tmpFilePath+"/"+fileName, indexChunkDivision)

	shopName := strings.Split(fileName, "-")[0]

	for {
		// 읽기 시작..
		byteLine, _, err, lineBreakerCount = reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				// 나머지를 모두 청크 파일에 작성
				if len(items) > 0 {
					sortIndexList := common.QuickSortList(items)
					common.WriteIndexChunkFile(sortIndexList, chunkFileWriterA, chunkFileWriterB, chunkGubun)
				}

				chunkFileA.Close()
				chunkFileB.Close()

				// 병합 정렬 시작
				common.WriteIndexFileMerge(tmpFilePath+"/"+fileName, indexChunkDivision, lineCount)

				// 인덱싱 완료되면 청크파일은 모두 지워준다
				common.ConvertToIndexFile(fileName, epFilePath, tmpFilePath+"/"+fileName)

				// 경량화 색인 파일을 생성한다
				common.CreateMiniIndexFile(shopName, epFilePath, fileName, indexFileInterval)
			}
			break
		}

		// 단위 문서 시작 부분을 만났을 때
		if strings.Contains(string(byteLine), format["startWord"]) {
			// 해당 위치 기록
			startPosition = endPosition
			docuOpen = true
		}

		// 읽은 라인 기록
		endPosition += len(byteLine) + lineBreakerCount

		if docuOpen {
			// 라인을 읽을 때마다 ID를 찾는다
			if strings.Contains(string(byteLine), format["idWord"]) {
				productId = strings.Split(string(byteLine), format["idWord"])[ValueIdx]
			}

			// 단위 문서 종료 부분을 만났을 때
			if strings.Contains(string(byteLine), format["endWord"]) {
				lineCount++

				// 리스트로 메모리에 할당한다
				items = append(items, model.Index{ProductId: productId, Position: startPosition})

				// 특정 갯수마다 작성한다
				if len(items)%indexChunkDivision == 0 {
					// 정렬 (퀵 소트)..
					sortIndexList := common.QuickSortList(items)
					common.WriteIndexChunkFile(sortIndexList, chunkFileWriterA, chunkFileWriterB, chunkGubun)
					chunkGubun++
					// 임시 변수 초기화
					items = []model.Index{}
				}

				indexLineLength += len([]byte(productId+","+strconv.Itoa(startPosition))) + 1

				docuOpen = false
			}
		}
	}
}
