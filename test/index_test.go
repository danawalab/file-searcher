package test

import (
	"fmt"
	"os"
	"testing"

	"gitlab.danawa.com/devops/file-searcher/service"
	"gitlab.danawa.com/devops/file-searcher/setting"
)

// 리스트에 적재 테스트
func TestLoad(t *testing.T) {
	configFilePath := "C:\\Users\\admin\\Desktop\\workspace\\file-searcher\\config.yml"
	fileName := "C:\\Users\\admin\\Desktop\\workspace\\file-searcher\\data\\tmp\\2022\\lenovo\\lenovo-1661825057-0.txt.index.index"

	// 설정파일 로딩
	fmt.Println("SETTING_FILE_PATH : " + configFilePath)
	loadSettingError := setting.LoadSetting(configFilePath)
	if loadSettingError != nil {
		fmt.Println("[ERROR] ", loadSettingError)
		os.Exit(1)
	}

	// 리스트 로딩
	service.Load(fileName)
}

func TestIndexer(t *testing.T) {
	configFilePath := "C:\\Users\\admin\\Desktop\\workspace\\file-searcher\\config.yml"
	filepath := "C:\\Users\\admin\\Desktop\\workspace\\file-searcher\\data\\ep\\2022\\TP90D-1665432000-1.txt"
	tmpPath := "C:\\Users\\admin\\Desktop\\workspace\\file-searcher\\data\\tmp"

	// 설정파일 로딩
	fmt.Println("SETTING_FILE_PATH : " + configFilePath)
	loadSettingError := setting.LoadSetting(configFilePath)
	if loadSettingError != nil {
		fmt.Println("[ERROR] ", loadSettingError)
		os.Exit(1)
	}

	service.Indexing(filepath, tmpPath)
}
