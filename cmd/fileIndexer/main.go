package main

import (
	"fmt"
	"os"

	"gitlab.danawa.com/devops/file-searcher/service"
	"gitlab.danawa.com/devops/file-searcher/setting"
)

func main() {
	if len(os.Args) < 3 {
		panic("설정파일 위치, 인덱스를 생성할 파일명을 입력하세요")
	}

	configFilePath := os.Args[1]
	filename := os.Args[2]

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
}
