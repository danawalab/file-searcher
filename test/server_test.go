package test

import (
	"fmt"
	"os"
	"testing"

	"gitlab.danawa.com/devops/file-searcher/service"
	"gitlab.danawa.com/devops/file-searcher/setting"
)

func TestDeleteSchedule(t *testing.T) {
	configFilePath := "C:\\Users\\admin\\Desktop\\workspace\\file-searcher\\config.yml"

	// 설정파일 로딩
	fmt.Println("SETTING_FILE_PATH : " + configFilePath)
	loadSettingError := setting.LoadSetting(configFilePath)
	if loadSettingError != nil {
		fmt.Println("[ERROR] ", loadSettingError)
		os.Exit(1)
	}

	service.ScheduleFileDeleteInit(setting.GetServerEpFilePathConfig(), setting.GetFileDeletePeriodDay())
}
