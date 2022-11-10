package service

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"regexp"
	"strings"
	"time"

	"gitlab.danawa.com/devops/file-searcher/common"
	"gitlab.danawa.com/devops/file-searcher/logging"
	"gitlab.danawa.com/devops/file-searcher/model"
	"gitlab.danawa.com/devops/file-searcher/parser"
	"gitlab.danawa.com/devops/file-searcher/setting"
)

// 메소드 리플렉션
func getMethod(shop string) reflect.Value {
	v := reflect.New(reflect.TypeOf(parser.CustomParser{}))
	// 커스텀 업체마다 고유한 메소드를 가짐
	// return v.MethodByName(strings.Title(shop) + "Parse")
	return v.MethodByName("MultiLineParse")
}

// 문서를 인덱스 파일로 변환한다
func Indexing(epFilePath, tmpFilePath string) {
	reg, _ := regexp.Compile("([A-Z0-9]+)-([0-9]+)-([0-9]+)")

	// ex. tmon-239293211
	fileName := reg.FindString(epFilePath)
	shop := strings.Split(fileName, "-")[0]

	if setting.IsCustomParse(shop) {
		writeCustomIndexingFile(fileName, epFilePath, tmpFilePath)
	} else {
		writeSingleIndexFile(fileName, epFilePath, tmpFilePath)
	}
}

// 싱글라인(csv, tsv, ...) 파일에 대해 인덱스 파일을 생성한다
func writeSingleIndexFile(fileName, epFilePath, tmpFilePath string) {
	// 시작 시간
	startTime := time.Now()

	sourceFile, reader := getTargetFile(epFilePath)
	defer sourceFile.Close()

	var err error
	var productId string
	var addLine []byte
	var startPosition, lineBreakerCount, chunkGubun, lineCount int
	var items []model.Index

	var result bool

	// 라인이 너무 길면 다 못읽고 hasPrefix + 추가 바이트 형태가 됨..
	var isPrefix bool

	shopName := strings.Split(fileName, "-")[0]

	singleLineParser := parser.NewSingleLineParser(shopName)
	shopConfig := setting.GetShopParseConfig(shopName)

	// 파일을 쪼갤 청크 단위
	indexChunkDivision := setting.GetFileSortDivision()

	// 경량화 인덱스 색인 간격
	indexFileInterval := setting.GetIndexFileInterval()

	// 이원화 처리를 위해 두 개의 인덱스 청크 파일을 준비한다
	chunkFileA, chunkFileWriterA, chunkFileB, chunkFileWriterB := common.GetChunkFileWriter(tmpFilePath+"/"+fileName, indexChunkDivision)

	if sourceFile != nil {
		for {
			var line []byte

			// 읽기 시작..
			line, isPrefix, err, lineBreakerCount = reader.ReadLine()
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
					result = common.CreateMiniIndexFile(shopName, epFilePath, fileName, indexFileInterval)

					// 경과 시간
					elapsedTime := time.Since(startTime)

					if result {
						Load(epFilePath)
					}

					logging.Info(fmt.Sprint("파일 색인 시간 : " + epFilePath + " " + elapsedTime.String()))
				}
				break
			} else {
				lineCount++
			}

			// 헤더가 있을 때 첫째 라인은 건너뛴다
			if shopConfig.Header && startPosition == 0 {
				startPosition += len(line) + lineBreakerCount
				lineCount--
				continue
			}

			productId = singleLineParser.SingleParse(string(line))

			if productId != "" {
				items = append(items, model.Index{ProductId: productId, Position: startPosition})
			}

			if isPrefix {
				// 길어서 다 못읽었던 한 줄을 모두 읽어 버린다.
				for isPrefix && err == nil {
					addLine, isPrefix, err, lineBreakerCount = reader.ReadLine()
					line = append(line, addLine...)
				}
			}

			// 특정 갯수마다 작성한다
			if len(items)%indexChunkDivision == 0 {
				// 정렬 (퀵 소트)..
				sortIndexList := common.QuickSortList(items)
				common.WriteIndexChunkFile(sortIndexList, chunkFileWriterA, chunkFileWriterB, chunkGubun)
				chunkGubun++
				// 임시 변수 초기화
				items = []model.Index{}
			}

			// 읽은 라인만큼 포지션 증가
			startPosition += len(line) + lineBreakerCount
		}
	} else {
		logging.Warn("색인 대상 파일이 존재하지 않습니다.")
	}
}

// 멀티라인(커스텀) 파일에 대해 인덱스 파일을 생성한다
func writeCustomIndexingFile(fileName, epFilePath, tmpFilePath string) {
	shopName := strings.Split(fileName, "-")[0]

	sourceFile, reader := getTargetFile(epFilePath)
	defer sourceFile.Close()

	parserMethod := getMethod(shopName)
	setting := setting.GetShopParseConfigMap(shopName)

	if sourceFile != nil {
		// 커스텀 파서 메소드 호출
		if parserMethod.IsValid() {
			parserMethod.Call([]reflect.Value{reflect.ValueOf(setting), reflect.ValueOf(reader), reflect.ValueOf(epFilePath), reflect.ValueOf(tmpFilePath), reflect.ValueOf(fileName)})
			Load(epFilePath)
		}
	} else {
		logging.Warn("색인 대상 파일이 존재하지 않습니다.")
	}
}

// 필요한 파일을 가져온다
func getTargetFile(filepath string) (*os.File, *common.Reader) {
	// 색인 대상 파일
	sourceFile := common.OpenFile(filepath, os.O_RDWR)
	sourceFileReader := common.NewReader(sourceFile)

	return sourceFile, sourceFileReader
}
