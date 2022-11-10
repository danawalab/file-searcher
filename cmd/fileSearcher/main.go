package main

import (
	"fmt"
	"os"

	"gitlab.danawa.com/devops/file-searcher/service"
	"gitlab.danawa.com/devops/file-searcher/setting"
)

func main() {
	if len(os.Args) < 3 {
		panic("검색할 파일명과 키워드를 입력하세요")
	}

	configFilePath := os.Args[1]
	filename := os.Args[2]
	shopName := os.Args[3]
	productId := os.Args[4]
	date := os.Args[5]
	renew := os.Args[6]

	// 설정파일 로딩
	fmt.Println("SETTING_FILE_PATH : " + configFilePath)
	loadSettingError := setting.LoadSetting(configFilePath)
	if loadSettingError != nil {
		fmt.Println("[ERROR] ", loadSettingError)
		os.Exit(1)
	}

	filepath := setting.GetServerEpFilePathConfig() + filename
	tmppath := setting.GetServerTempFilePathConfig()
	service.Indexing(filepath, tmppath)

	// 리스트 로딩
	service.Load(filename)

	// 해당 키워드 검색
	service.SearchID(shopName, productId, date, date, renew)
}
