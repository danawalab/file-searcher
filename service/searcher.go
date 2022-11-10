package service

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jasonlvhit/gocron"
	"github.com/saintfish/chardet"
	"gitlab.danawa.com/devops/file-searcher/common"
	"gitlab.danawa.com/devops/file-searcher/logging"
	"gitlab.danawa.com/devops/file-searcher/model"
	"gitlab.danawa.com/devops/file-searcher/setting"
)

var (
	fileMap map[string][]model.MiniIndex
)

const (
	SOURCE_FILE_INDEX = 1
	INDEX_FILE_INDEX  = 2
	// Jan 2 15:04:05 2006 MST
	// 1   2  3  4  5    6  -7
	LAYOUT = "20060102"
)

// 색인이 완료된 문서를 메모리에 적재한다
func Load(file string) {
	defer func() {
		if r := recover(); r != nil {
			logging.Error(fmt.Sprint("Recovered in Load", r))
		}
	}()
	var filename string

	if !strings.Contains(file, ".index.index") {
		filename = file + ".index.index"
	} else {
		filename = file
	}

	indexfile := common.OpenFile(filename, os.O_RDWR)
	defer indexfile.Close()

	// 인덱스 파일 맵에 적재
	indexToList(indexfile, file)
}

// 경량화 인덱스 파일을 읽어서 맵리스트로 변환하는 메소드
func indexToList(file *os.File, fileName string) {
	defer func() {
		if r := recover(); r != nil {
			logging.Error(fmt.Sprint("Recovered in indexToList", r))
		}
	}()

	reader := common.NewReader(file)

	for {
		line, _, err, _ := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				common.CheckMemory()
				logging.Info(fmt.Sprintf("%s 적재 완료, 현재 메모리 사용량 = %d byte \n", fileName, common.M.Alloc))
				break
			}
		}
		splitLines := strings.Split(string(line), ",")

		if len(splitLines) > 2 {
			productId := splitLines[0]

			// 상품에 대한 포지션
			originPosition, err := strconv.Atoi(splitLines[1])
			if err != nil {
				logging.Warn(fmt.Sprint("origin index ", fileName, err))
			}

			// 인덱스에 대한 포지션
			indexPosition, err := strconv.Atoi(splitLines[2])
			if err != nil {
				logging.Warn(fmt.Sprint("index position ", fileName, err))
			}

			indexMutex.Lock()

			if fileMap == nil {
				fileMap = make(map[string][]model.MiniIndex)
			}

			fileName = strings.Replace(fileName, ".index.index", "", -1)

			fileMap[fileName] = append(fileMap[fileName], model.MiniIndex{ProductId: productId, OriginPosition: originPosition, IndexPosition: indexPosition})

			indexMutex.Unlock()
		} else {
			// 생성된 인덱스가 없을 경우, 파싱이 잘못 되었을 경우
			break
		}
	}
}

// Map에서 해당 상품코드를 검색하기 위한 메소드
func SearchID(shop, productId, startDate, endDate, renew string) ([]model.Search, error) {
	var result []model.Search

	// 타겟 디렉토리
	filePath := setting.GetServerEpFilePathConfig()

	// 해당 업체 파일명 가져오기
	fileList, err := getMiniIndexFileList(shop, filePath, startDate, endDate, renew)

	for _, fileName := range fileList {
		containsLine := findIndexLine(fileName, fileMap[fileName], productId)

		// 결국 필요한건 인덱스 라인 한줄
		if containsLine != "" {
			if setting.IsCustomParse(shop) {
				result = append(result, model.Search{Source: findSourceFileMulti(setting.GetShopParseConfigMap(shop), fileName, containsLine), Timestamp: getTimestamp(fileName), Renew: getRenewValue(fileName)})
			} else {
				result = append(result, model.Search{Source: findSourceFileSingle(setting.GetShopParseConfigMap(shop), containsLine, fileName), Timestamp: getTimestamp(fileName), Renew: getRenewValue(fileName)})
			}
		}
	}

	return result, err
}

func findIndexLine(fileName string, indexList []model.MiniIndex, productId string) string {
	var targetIndex int
	var indexLine string
	var isExist bool

	// 이진 탐색 진행
	targetIndex, isExist = common.BinarySearch(indexList, productId)

	if len(indexList) > targetIndex {
		// 존재여부 체크하여 분기
		if isExist {
			indexLine = convertToCSVString(indexList[targetIndex])
		} else {
			// 없으면 실패 지점부터 찾아들어간다
			indexLine = searchAtIndex(indexList[targetIndex], productId, fileName)
		}
	}

	return indexLine
}

// 인덱스에서 타겟을 찾아내는 메소드
func searchAtIndex(midx model.MiniIndex, productId, fileName string) string {
	var result string
	var readCount int

	interval := setting.GetIndexFileInterval()
	indexfile := common.OpenFile(fileName+".index", os.O_RDWR)
	defer indexfile.Close()

	reader := common.NewReader(indexfile)

	// 인덱스 시작 위치 이동
	indexfile.Seek(int64(midx.IndexPosition), 0)

	for {
		line, _, err, _ := reader.ReadLine()
		readCount++
		if err != nil {
			if err == io.EOF {
				break
			}
		}

		indexline := strings.Split(string(line), ",")
		if len(indexline) < 2 {
			break
		}

		if indexline[0] == productId {
			result = string(line)
			break
		}

		// 다음 한 간격까지만 읽도록 한다.
		if readCount == interval {
			break
		}
	}

	return result
}

// 싱글 라인 파일 검색
func findSourceFileSingle(setting map[string]string, containsLine, fileName string) string {
	var convertLine string

	filePosition, err := strconv.ParseInt(strings.Split(containsLine, ",")[1], 10, 64)
	if err != nil {
		logging.Error(fmt.Sprint("[ERROR] ", err))
	}

	file := common.OpenFile(fileName, os.O_RDWR)
	defer file.Close()

	// 파일 시작을 기준점으로 해당 포지션부터 찾아서 한 라인만 읽는다
	file.Seek(filePosition, 0)

	reader := common.NewReader(file)
	line, _, err, _ := reader.ReadLine()
	if err != nil {
		logging.Error(fmt.Sprint("[ERROR] ", err))
	}

	detector := chardet.NewTextDetector()
	result, err := detector.DetectBest(line)
	if err != nil {
		logging.Error(fmt.Sprint("[GET FORMAT ERROR] ", err))
	}

	if result.Charset != "UTF-8" {
		convertLine = common.ConvrtToUTF8(string(line), "euc-kr")
	} else {
		convertLine = string(line)
	}

	return convertLine
}

// 멀티 라인 파일 검색
func findSourceFileMulti(setting map[string]string, fileName, containsLine string) string {
	var result, convertLine string

	idxPos := getPositionToInt64(containsLine, SOURCE_FILE_INDEX)

	// 인덱스 파일
	indexFile, sourceFile, _, sourceFileReader := getTargetSearchFile(fileName)
	defer indexFile.Close()
	defer sourceFile.Close()

	// 시작 위치 설정
	sourceFile.Seek(idxPos, 0)

	// 시작점부터 마지막 구분자까지 읽어낸다
	for {
		line, _, err, _ := sourceFileReader.ReadLine()
		if err != nil {
			logging.Error(fmt.Sprint("[ERROR] ", err))
		}

		detector := chardet.NewTextDetector()
		data, err := detector.DetectBest(line)
		if err != nil {
			logging.Error(fmt.Sprint("[GET FORMAT ERROR] ", err))
		}
		if data.Charset != "UTF-8" {
			convertLine = common.ConvrtToUTF8(string(line), "euc-kr")
		} else {
			convertLine = string(line)
		}

		if err != nil {
			if err == io.EOF {
				return result
			} else {
				logging.Error(fmt.Sprint("[ERROR] ", err))
			}
		}
		result += convertLine
		if strings.Contains(convertLine, setting["endWord"]) {
			break
		}
	}
	return result
}

// 인덱스 파일, 원본 파일을 가져온다
func getTargetSearchFile(fileName string) (*os.File, *os.File, *common.Reader, *common.Reader) {
	indexFile := common.OpenFile(fileName+".index", os.O_RDWR)
	indexFileReader := common.NewReader(indexFile)
	sourceFile := common.OpenFile(fileName, os.O_RDWR)
	sourceFileReader := common.NewReader(sourceFile)
	return indexFile, sourceFile, indexFileReader, sourceFileReader
}

func getPositionToInt64(line string, index int) int64 {
	position, err := strconv.ParseInt(strings.Split(string(line), ",")[index], 10, 64)
	if err != nil {
		logging.Error(fmt.Sprint("[ERROR] ", err))
	}
	return position
}

func convertToCSVString(str model.MiniIndex) string {
	return str.ProductId + "," + strconv.Itoa(str.OriginPosition) + "," + strconv.Itoa(str.IndexPosition)
}

func getMiniIndexFileList(shopName, filePath, startDate, endDate, renew string) ([]string, error) {
	var items []string

	startTime, terr := time.Parse(LAYOUT, startDate)
	if terr != nil {
		if startDate != "" {
			logging.Error(fmt.Sprint("[SET TIME ERROR] ", terr))
		}
	}

	endTime, terr := time.Parse(LAYOUT, endDate)
	if terr != nil {
		if startDate != "" {
			logging.Error(fmt.Sprint("[SET TIME ERROR] ", terr))
		}
	}

	// 타겟 디렉토리 전체 파일 조사
	err := filepath.WalkDir(filePath,
		func(path string, entry fs.DirEntry, err error) error {
			info, errEntry := entry.Info()
			if errEntry != nil {
				return errEntry
			}

			// 해당 업체의 원본 파일 탐색
			if !info.IsDir() && strings.Contains(info.Name(), shopName) && !strings.Contains(info.Name(), ".index") && !strings.Contains(info.Name(), ".tmp") {
				// 날짜 검색 조건
				if startDate != "" && endDate != "" {
					// 날짜조건이 있을 경우
					// EP파일에서 날짜데이터 추출
					itemTime := time.Unix(getTimestamp(path), 0)

					// 해당 조건의 파일만 수집한다
					if (itemTime.After(startTime) || (itemTime.Year() == startTime.Year() && itemTime.Month() == startTime.Month() && itemTime.Day() == startTime.Day())) && (itemTime.Before(endTime) || (itemTime.Year() == endTime.Year() && itemTime.Month() == endTime.Month() && itemTime.Day() == endTime.Day())) {
						if renew != "" {
							if getRenewValue(path) == renew {
								// 해당 조건의 파일만 수집한다
								items = append(items, path)
							}
						} else {
							items = append(items, path)
						}
					}
				} else {
					// 전체(1), 갱신(2) 조건이 있다면
					if renew != "" {
						if getRenewValue(path) == renew {
							// 해당 조건의 파일만 수집한다
							items = append(items, path)
						}
					} else {
						items = append(items, path)
					}
				}
			}

			return nil
		})

	if err != nil {
		logging.Error(fmt.Sprint("[GET MINIINDEX FILE] ", err))
	}

	return items, err
}

// 타임스탬프 조회
func getTimestamp(path string) int64 {
	var timestamp string

	reg, _ := regexp.Compile("([A-Z0-9]+)-([0-9]+)")
	// ex. tmon-239293211
	fileName := reg.FindString(path)
	timeArr := strings.Split(fileName, "-")
	if len(timeArr) > 1 {
		timestamp = timeArr[1]
	}

	i, _ := strconv.ParseInt(timestamp, 10, 64)
	return i
}

// 전체, 갱신여부 조회
func getRenewValue(path string) string {
	var renew string

	reg, _ := regexp.Compile("([A-Z0-9]+)-([0-9]+)-([0-9]+)")
	// ex. tmon-239293211-1
	fileName := reg.FindString(path)
	renewArr := strings.Split(fileName, "-")

	if len(renewArr) > 1 {
		renew = renewArr[2]
	}

	return renew
}

// 삭제 스케줄 시작
func ScheduleFileDeleteInit(targetDir string, limitDay int) {
	// 한 시간마다 수행하도록 한다.
	gocron.Every(1).Day().At("00:05").Do(deleteTask, targetDir, limitDay)
	// start
	<-gocron.Start()
}

// 오래된 파일 삭제 로직
func deleteTask(targetDir string, limitDay int) {
	// 로그 교체
	logging.Rotate()

	deleteTime := time.Now().AddDate(0, 0, -1*limitDay)

	// ex. 2022, 2023
	deleteYY := fmt.Sprintf("%d", deleteTime.Year())

	// ex. 1101, 0115, 0613..
	deleteMMDD := fmt.Sprintf("%02d%02d", deleteTime.Month(), deleteTime.Day())

	path := filepath.Join(targetDir, deleteYY, deleteMMDD)

	logging.Info("DELETE SCHEDULE START ... Target : " + path)

	filepath.WalkDir(path,
		func(path string, entry fs.DirEntry, err error) error {
			if strings.HasSuffix(path, ".txt") {
				// 해당 시간이 지난 파일은 메모리에서 내려준다
				delete(fileMap, path)
				if err != nil {
					logging.Error(fmt.Sprint("[OLD DATA REMOVE ERROR] ", err))
				} else {
					logging.Info(fmt.Sprint("[OLD DATA REMOVED] ", path))
				}
			}
			return nil
		})
}
