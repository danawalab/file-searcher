package test

import (
	"fmt"
	"os"
	"testing"

	"gitlab.danawa.com/devops/file-searcher/service"
	"gitlab.danawa.com/devops/file-searcher/setting"
)

// 디렉토리 리스닝, 색인 테스트
func TestDirectoryListening(t *testing.T) {
	configFilePath := "C:\\Users\\admin\\Desktop\\workspace\\file-searcher\\config\\config.yml"

	// 설정파일을 로드합니다.
	fmt.Println("SETTING_FILE_PATH : " + configFilePath)
	loadSettingError := setting.LoadSetting(configFilePath)
	if loadSettingError != nil {
		fmt.Println("[ERROR] ", loadSettingError)
		os.Exit(1)
	}

	service.Initalize(setting.GetServerEpFilePathConfig(), setting.GetServerEpFilePathConfig(), setting.GetServerWorkerConfig())
}

// 일괄 동적색인 테스트
func TestAllDirIndex(t *testing.T) {
	configFilePath := "C:\\Users\\admin\\Desktop\\workspace\\file-searcher\\config.yml"

	// 설정파일을 로드합니다.
	fmt.Println("SETTING_FILE_PATH : " + configFilePath)
	loadSettingError := setting.LoadSetting(configFilePath)
	if loadSettingError != nil {
		fmt.Println("[ERROR] ", loadSettingError)
		os.Exit(1)
	}

	service.IndexingAllFiles(setting.GetServerEpFilePathConfig(), setting.GetServerEpFilePathConfig(), setting.GetServerWorkerConfig())
}
