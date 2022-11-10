package common

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"strings"

	"gitlab.danawa.com/devops/file-searcher/logging"
	"gitlab.danawa.com/devops/file-searcher/model"
	"golang.org/x/net/html/charset"
)

var (
	M runtime.MemStats
)

// 파일 오픈
func OpenFile(filePath string, flag int) *os.File {
	file, err := os.OpenFile(filePath, flag, 0777)
	if err != nil {
		logging.Warn(fmt.Sprint(err))
	}
	return file
}

func CheckMemory() {
	runtime.GC()
	runtime.ReadMemStats(&M)
}

func WriteCsv(csvWriter *csv.Writer, arg1 string, arg2 string) {
	x := []string{arg1, arg2}
	strWrite := [][]string{x}
	csvWriter.WriteAll(strWrite)
	csvWriter.Flush()
}

func WriteCsvAdded(csvWriter *csv.Writer, arg1 string, arg2 string, arg3 string) {
	x := []string{arg1, arg2, arg3}
	strWrite := [][]string{x}
	csvWriter.WriteAll(strWrite)
	csvWriter.Flush()
}

func GetLineBreakerCount(reader *Reader) ([]byte, int, error) {
	var lineBreakerNumber int

	// 파일을 줄 단위로 끊어낸다.
	// ex. line : <<<begin>>>\r\n
	line, err := reader.ReadSlice('\n')
	if err == nil {
		if line[len(line)-1] == '\n' {
			lineBreakerNumber = 1
			if len(line) > 1 && line[len(line)-2] == '\r' {
				lineBreakerNumber = 2
			}
		}
	}

	line = line[:len(line)-lineBreakerNumber]

	return line, lineBreakerNumber, err
}

// 인덱스 청크 파일을 만든다
func WriteIndexChunkFile(indexList []model.Index, w0 *csv.Writer, w1 *csv.Writer, gubun int) {
	var chunkFileWriter *csv.Writer

	if gubun%2 == 0 {
		chunkFileWriter = w0
	} else {
		chunkFileWriter = w1
	}

	// 청크파일 생성
	for _, line := range indexList {
		WriteCsv(chunkFileWriter, line.ProductId, strconv.Itoa(line.Position))
	}
}

// 이원화된 파일과 리더를 가져온다
func GetChunkFileReader(filepath string, indexChunkDivision int) (*os.File, *Reader, *os.File, *Reader) {
	// 두 개의 파일을 각각 한 줄씩 읽으면서 병합
	chunkFileA := OpenFile(filepath+".index.tmp.A."+strconv.Itoa(indexChunkDivision), os.O_RDWR)
	chunkFileAReader := NewReader(chunkFileA)
	chunkFileB := OpenFile(filepath+".index.tmp.B."+strconv.Itoa(indexChunkDivision), os.O_RDWR)
	chunkFileBReader := NewReader(chunkFileB)

	return chunkFileA, chunkFileAReader, chunkFileB, chunkFileBReader
}

// 이원화된 파일과 라이터를 가져온다
func GetChunkFileWriter(filepath string, indexChunkDivision int) (*os.File, *csv.Writer, *os.File, *csv.Writer) {
	chunkFileA := OpenFile(filepath+".index.tmp.A."+strconv.Itoa(indexChunkDivision), os.O_CREATE|os.O_WRONLY)
	chunkFileB := OpenFile(filepath+".index.tmp.B."+strconv.Itoa(indexChunkDivision), os.O_CREATE|os.O_WRONLY)
	chunkFileWriterA := csv.NewWriter(chunkFileA)
	chunkFileWriterB := csv.NewWriter(chunkFileB)
	return chunkFileA, chunkFileWriterA, chunkFileB, chunkFileWriterB
}

// 인덱스 청크 파일을 작성하는 메소드
func WriteIndexFileMerge(filepath string, indexChunkDivision int, lineCount int) {
	var Aerr, Berr error
	var Aid, Bid string
	var fileSwap bool

	for {
		// 간격보다 작은 경우라면 스킵한다
		if lineCount == indexChunkDivision {
			break
		}

		var totalReadCount, AreadCount, BreadCount, Apos, Bpos int

		// 이전의 이원화된 파일을 바라본다
		fileA, readerA, fileB, readerB := GetChunkFileReader(filepath, indexChunkDivision)
		defer fileA.Close()
		defer fileB.Close()

		// 이전 병합 갯수를 저장
		chunkCountLimit := indexChunkDivision

		// 병합할 청크 갯수
		indexChunkDivision *= 2

		// 작성할 청크 파일을 가져온다
		chunkFileA, chunkFileWriterA, chunkFileB, chunkFileWriterB := GetChunkFileWriter(filepath, indexChunkDivision)
		defer chunkFileA.Close()
		defer chunkFileB.Close()

		// A 첫번째 라인 선언
		Aid, Apos, Aerr = readLineToIndex(readerA)
		if Aerr != io.EOF {
			AreadCount++
		}

		// B 첫번째 라인 선언
		Bid, Bpos, Berr = readLineToIndex(readerB)
		if Berr != io.EOF {
			BreadCount++
		}

		for {
			if (Aerr != io.EOF || Berr != io.EOF) && totalReadCount+AreadCount+BreadCount < lineCount {
				if Aid >= Bid {
					Bid, Bpos, Berr, BreadCount, Aid, Apos, Aerr, AreadCount = mergeChunkInto(Bid, Aid, Bpos, Apos, BreadCount, Berr, Aerr, chunkFileWriterA, chunkFileWriterB, fileSwap, AreadCount, chunkCountLimit, readerB, readerA)
				} else if Aid < Bid {
					Aid, Apos, Aerr, AreadCount, Bid, Bpos, Berr, BreadCount = mergeChunkInto(Aid, Bid, Apos, Bpos, AreadCount, Aerr, Berr, chunkFileWriterA, chunkFileWriterB, fileSwap, BreadCount, chunkCountLimit, readerA, readerB)
				}
			} else {
				fileA.Close()
				fileB.Close()
				chunkFileA.Close()
				chunkFileB.Close()
				totalReadCount = 0
				break
			}

			// 한 단위가 마무리되었을 때
			if ((AreadCount >= chunkCountLimit) && (BreadCount >= chunkCountLimit)) || AreadCount+BreadCount >= lineCount {
				totalReadCount += AreadCount + BreadCount
				AreadCount, BreadCount = 0, 0
				fileSwap = !fileSwap

				if Aerr != io.EOF {
					Aid, Apos, Aerr = readLineToIndex(readerA)
					AreadCount++
				}

				if Berr != io.EOF {
					Bid, Bpos, Berr = readLineToIndex(readerB)
					BreadCount++
				}
			}
		}

		// 전체 라인 수보다 청크 단위크기가 커지면 종료한다
		if lineCount < indexChunkDivision {
			break
		}
	}
}

// 청크 파일을 만든다
func writeChunkFile(line model.Index, w0 *csv.Writer, w1 *csv.Writer, switchFlag bool) {
	var chunkFileWriter *csv.Writer

	if switchFlag {
		chunkFileWriter = w0
	} else {
		chunkFileWriter = w1
	}

	// 청크파일 생성
	WriteCsv(chunkFileWriter, line.ProductId, strconv.Itoa(line.Position))
}

// 가장 큰 청크 파일이 인덱스 파일이 된다..
func ConvertToIndexFile(fileName, epfilepath, tmpfilepath string) {
	var maxSize int64
	var maxFileName string
	var files []fs.FileInfo
	var err error

	targetDir := strings.Split(tmpfilepath, fileName)[0]

	files, err = ioutil.ReadDir(targetDir)
	if err != nil {
		logging.Warn(err.Error())
	}

	// 제일 큰 파일 구하기
	for _, file := range files {
		// 디렉토리 파일명 <-> tmon-2020201 파일명 비교

		if strings.Contains(file.Name(), fileName) && strings.Contains(file.Name(), "tmp") {

			res, err := os.Stat(targetDir + file.Name())
			if err != nil {
				logging.Error(fmt.Sprint("[Stat Check ERROR] ", err))
			}

			if res != nil {
				if maxSize < res.Size() {
					maxSize = res.Size()
					maxFileName = targetDir + file.Name()
				}
			}
		}
	}

	e := os.Rename(maxFileName, epfilepath+".index")
	if e != nil {
		logging.Error(fmt.Sprint("[Rename ERROR] ", e))
	}

	files, _ = ioutil.ReadDir(targetDir)

	// 청크파일을 모두 제거한다
	for _, file := range files {
		// 파일의 절대경로
		if strings.Contains(file.Name(), fileName) && strings.Contains(file.Name(), "tmp") {
			path := fmt.Sprintf(targetDir + file.Name())
			os.Remove(path)
		}
	}
}

func readLineToIndex(reader *Reader) (string, int, error) {
	var id string
	var pos int
	var err error

	line, _, err, _ := reader.ReadLine()
	if err == io.EOF {
		return "", 0, err
	}

	parseLine := strings.Split(string(line), ",")

	if len(parseLine) > 1 {
		id = parseLine[0]
		pos, err = strconv.Atoi(parseLine[1])
	}

	return id, pos, err
}

// 청크파일을 작성하는 메소드, 이원화된 파일이 있고 한 파일을 다 읽으면 반대쪽도 대칭으로 읽는다
func mergeChunkInto(targetId, counterId string, targetPos, counterPos, targetReadCount int, targetErr, counterErr error, targetChunkFileWriter, counterChunkFileWriter *csv.Writer, fileSwap bool, counterReadCount int, chunkCountLimit int, targetReader, counterReaderB *Reader) (string, int, error, int, string, int, error, int) {
	if len(targetId) > 0 {
		writeChunkFile(model.Index{ProductId: targetId, Position: targetPos}, targetChunkFileWriter, counterChunkFileWriter, fileSwap)
	}

	// 적었으므로 타겟을 하나 읽는다
	if targetReadCount < chunkCountLimit {
		targetId, targetPos, targetErr = readLineToIndex(targetReader)
		targetReadCount++

		if targetErr != io.EOF {
			if targetReadCount >= chunkCountLimit {
				inputedTarget := false

				if targetId < counterId {
					writeChunkFile(model.Index{ProductId: targetId, Position: targetPos}, targetChunkFileWriter, counterChunkFileWriter, fileSwap)
					inputedTarget = true
				}

				for {
					if !inputedTarget {
						if targetId < counterId {
							writeChunkFile(model.Index{ProductId: targetId, Position: targetPos}, targetChunkFileWriter, counterChunkFileWriter, fileSwap)
							writeChunkFile(model.Index{ProductId: counterId, Position: counterPos}, targetChunkFileWriter, counterChunkFileWriter, fileSwap)
							inputedTarget = true
						} else {
							writeChunkFile(model.Index{ProductId: counterId, Position: counterPos}, targetChunkFileWriter, counterChunkFileWriter, fileSwap)
						}
					} else {
						writeChunkFile(model.Index{ProductId: counterId, Position: counterPos}, targetChunkFileWriter, counterChunkFileWriter, fileSwap)
					}

					// 허용량을 초과하면 반대편을 청크파일 단위까지 읽는다
					if counterReadCount < chunkCountLimit {
						counterId, counterPos, counterErr = readLineToIndex(counterReaderB)
						if counterErr != io.EOF {
							counterReadCount++
						} else {
							if !inputedTarget {
								writeChunkFile(model.Index{ProductId: targetId, Position: targetPos}, targetChunkFileWriter, counterChunkFileWriter, fileSwap)
								inputedTarget = true
							}
							break
						}
					} else {
						if !inputedTarget {
							writeChunkFile(model.Index{ProductId: targetId, Position: targetPos}, targetChunkFileWriter, counterChunkFileWriter, fileSwap)
							inputedTarget = true
						}
						break
					}
				}
			}
		} else if targetErr == io.EOF {
			// 타겟 파일의 끝에 도달했을때
			for {
				writeChunkFile(model.Index{ProductId: counterId, Position: counterPos}, targetChunkFileWriter, counterChunkFileWriter, fileSwap)

				// 반대편 파일 끝까지 읽는다
				counterId, counterPos, counterErr = readLineToIndex(counterReaderB)
				if counterErr != io.EOF {
					counterReadCount++
				} else {
					break
				}
			}
		}
	}

	return targetId, targetPos, targetErr, targetReadCount, counterId, counterPos, counterErr, counterReadCount
}

// 경량화 색인 파일을 작성하는 메소드
func CreateMiniIndexFile(shopName, filepath, fileName string, indexFileInterval int) bool {
	var indexFilePath, productId, position string
	var lineCount, startPosition int

	targetDir := strings.Split(filepath, fileName)[0]

	files, err := ioutil.ReadDir(targetDir)
	if err != nil {
		logging.Error(fmt.Sprint("[Read File ERROR] ", err))
	}

	for _, file := range files {
		// 인덱싱할 경량화 파일 검색
		if strings.Contains(file.Name(), fileName) && strings.Contains(file.Name(), ".index") && strings.Contains(file.Name(), shopName) {
			indexFilePath = targetDir + file.Name()
			break
		}
	}

	indexFile := OpenFile(indexFilePath, os.O_RDWR)
	indexFileReader := NewReader(indexFile)

	miniIndexFile := OpenFile(filepath+".index.index", os.O_CREATE|os.O_WRONLY)
	miniIndexFileWriter := csv.NewWriter(miniIndexFile)

	defer indexFile.Close()
	defer miniIndexFile.Close()

	for {
		line, _, err, _ := indexFileReader.ReadLine()
		lineCount++

		// 지정한 횟수마다 경량화 인덱스 파일에 작성한다
		lineArr := strings.Split(string(line), ",")

		if len(lineArr) == 2 {
			productId = lineArr[0]
			position = lineArr[1]
		}

		// 단, 첫번째 라인은 무조건 써준다
		if lineCount == 1 || lineCount%indexFileInterval == 0 {
			WriteCsvAdded(miniIndexFileWriter, productId, position, strconv.Itoa(startPosition))
		}

		startPosition += len([]byte(productId+","+position)) + 1

		if err == io.EOF {
			return true
		} else if err != nil {
			logging.Warn(err.Error())
			return false
		}
	}
}

// 파일 이름 가져오기
func GetFileName() string {
	return ""
}

// UTF8로 변경하는 코드
func ConvrtToUTF8(str string, origEncoding string) string {
	strBytes := []byte(str)
	byteReader := bytes.NewReader(strBytes)
	reader, readErr := charset.NewReaderLabel(origEncoding, byteReader)
	if readErr != nil {
		logging.Error(fmt.Sprint("[convrtToUTF8 ERROR] ", readErr))
	}

	strBytes, readErr = ioutil.ReadAll(reader)
	if readErr != nil {
		logging.Error(fmt.Sprint("[strRead ERROR] ", readErr))
	}

	return string(strBytes)
}
