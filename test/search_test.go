package test

import (
	"fmt"
	"io"
	"log"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	"gitlab.danawa.com/devops/file-searcher/common"
	"gitlab.danawa.com/devops/file-searcher/service"
	"gitlab.danawa.com/devops/file-searcher/setting"
)

func TestSearch(t *testing.T) {
	shopName := "TR305"
	configFilePath := "C:\\Users\\admin\\Desktop\\workspace\\file-searcher\\config.yml"
	productId := "58220447"
	filepath := "C:\\Users\\admin\\Desktop\\workspace\\file-searcher\\data\\tmp\\2022\\TR305-1665101700-1.txt"
	start := ""
	end := ""
	renew := ""

	// 설정파일 로딩
	fmt.Println("SETTING_FILE_PATH : " + configFilePath)
	loadSettingError := setting.LoadSetting(configFilePath)
	if loadSettingError != nil {
		fmt.Println("[ERROR] ", loadSettingError)
		os.Exit(1)
	}

	// 리스트 로딩
	service.Load(filepath)

	// 해당 키워드 검색
	service.SearchID(shopName, productId, start, end, renew)
}

// 전체 일치 여부 테스트..
func TestSearchTotal(t *testing.T) {
	shopName := "TP304"
	configFilePath := "C:\\Users\\admin\\Desktop\\workspace\\file-searcher\\config.yml"
	filepath := "C:\\Users\\admin\\Desktop\\workspace\\file-searcher\\data\\tmp\\2022\\TP304-1665637500-2.txt"
	date := ""
	renew := ""

	// 설정파일 로딩
	fmt.Println("SETTING_FILE_PATH : " + configFilePath)
	loadSettingError := setting.LoadSetting(configFilePath)
	if loadSettingError != nil {
		fmt.Println("[ERROR] ", loadSettingError)
		os.Exit(1)
	}

	// 리스트 로딩
	service.Load(filepath)

	// 색인 대상 파일
	file := common.OpenFile(filepath, os.O_RDONLY)
	sourceFileReader := common.NewReader(file)

	var header int

	var count int

	start := time.Now()

	r := new(big.Int)
	fmt.Println(r.Binomial(1000, 10))

	for {
		header++

		line, _, err, _ := sourceFileReader.ReadLine()
		if err != nil {
			if err == io.EOF {
				elapsed := time.Since(start)
				log.Printf("총 소요 시간 %s", elapsed)
				log.Printf("결과와 불일치 갯수 %d", count)
				break
			} else {
				fmt.Println(err.Error())
			}
		}

		if header == 1 {
			continue
		}

		pid := strings.Split(string(line), "^")

		// 해당 키워드 검색
		data, _ := service.SearchID(shopName, pid[0], date, date, renew)

		if len(string(line)) != len(data[0].Source) {
			fmt.Println("불일치..")
			fmt.Println(string(line))
			fmt.Println(data[0].Source)
			break
		}
	}
}

func TestSearchMultiLine(t *testing.T) {
	shopName := "TH201"
	configFilePath := "C:\\Users\\admin\\Desktop\\workspace\\file-searcher\\config.yml"
	filepath := "C:\\Users\\admin\\Desktop\\workspace\\file-searcher\\data\\tmp\\2022\\TH201-1665562089-2.txt"
	date := ""
	renew := ""

	// 설정파일 로딩
	fmt.Println("SETTING_FILE_PATH : " + configFilePath)
	loadSettingError := setting.LoadSetting(configFilePath)
	if loadSettingError != nil {
		fmt.Println("[ERROR] ", loadSettingError)
		os.Exit(1)
	}

	// 리스트 로딩
	service.Load(filepath)

	// 색인 대상 파일
	file := common.OpenFile(filepath, os.O_RDONLY)
	sourceFileReader := common.NewReader(file)

	var docuOpen bool
	var multiline, id string
	var count int

	start := time.Now()

	r := new(big.Int)
	fmt.Println(r.Binomial(1000, 10))

	for {
		line, _, err, _ := sourceFileReader.ReadLine()
		if err != nil {
			if err == io.EOF {
				elapsed := time.Since(start)
				log.Printf("총 소요 시간 %s", elapsed)
				log.Printf("결과와 불일치 갯수 %d", count)
				break
			} else {
				fmt.Println(err.Error())
			}
		}

		// 단위 문서 시작 부분을 만났을 때
		if strings.Contains(string(line), "<<<begin>>>") {
			docuOpen = true
		}

		if docuOpen {
			if strings.Contains(string(line), "<<<mapid>>>") {
				lines := strings.Split(string(line), "<<<mapid>>>")
				if len(lines) >= 2 {
					id = lines[1]
				} else {
					fmt.Println(len(lines))
					break
				}
			}

			multiline += string(line)
		}

		// 단위 문서 시작 끝부분을 만났을 때
		if strings.Contains(string(line), "<<<ftend>>>") {
			// 해당 키워드 검색
			data, _ := service.SearchID(shopName, id, date, date, renew)

			if data != nil {
				if len(multiline) != len(data[0].Source) {
					count++
				}
			} else {
				break
			}

			docuOpen = false
			multiline = ""
		}
	}
}
