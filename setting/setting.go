package setting

import (
	"io/ioutil"

	"gitlab.danawa.com/devops/file-searcher/model"
	"gopkg.in/yaml.v2"
)

var (
	shopSetting          = make(map[string]map[string]model.Setting)
	shopSettingInterface = make(map[string]map[string]map[string]string)
)

func LoadSetting(filePath string) error {
	raw, err := ioutil.ReadFile(filePath)
	if err == nil {
		_ = yaml.Unmarshal(raw, &shopSetting)
		_ = yaml.Unmarshal(raw, &shopSettingInterface)
	}

	return err
}

func GetShopParseConfig(shopCode string) model.Setting {
	var result model.Setting

	// 포맷이 없으면 디폴트 처리..
	if (shopSetting["shop"][shopCode] == model.Setting{}) {
		result = shopSetting["shop"]["DEFAULT"]
	} else {
		result = shopSetting["shop"][shopCode]
	}

	return result
}

func GetShopParseConfigMap(shopCode string) map[string]string {
	var result map[string]string

	// 포맷이 없으면 디폴트 처리..
	if (shopSetting["shop"][shopCode] == model.Setting{}) {
		result = shopSettingInterface["shop"]["DEFAULT"]
	} else {
		result = shopSettingInterface["shop"][shopCode]
	}
	return result
}

func IsCustomParse(shopCode string) bool {
	return shopSetting["shop"][shopCode].Custom
}

func GetServerEpFilePathConfig() string {
	return shopSetting["server"]["config"].Epfilepath
}

func GetServerWorkerConfig() int {
	return shopSetting["server"]["config"].Workers
}

func GetServerCpuCoreConfig() int {
	return shopSetting["server"]["config"].CpuCore
}

func GetServerTempFilePathConfig() string {
	return shopSetting["server"]["config"].TempFilepath
}

func GetServerFilePathPort() string {
	return shopSetting["server"]["config"].Port
}

func GetServerIndexCountConfig() int {
	return shopSetting["server"]["config"].IndexInterval
}

func GetFileSortDivision() int {
	return shopSetting["server"]["config"].FileSortDivision
}

func GetIndexFileInterval() int {
	return shopSetting["server"]["config"].IndexInterval
}

func GetFileDeletePeriodDay() int {
	return shopSetting["server"]["config"].FileDeletePeriodDay
}

func GetLogging() model.Logging {
	fileName := shopSetting["server"]["logging"].Filename
	level := shopSetting["server"]["logging"].Level
	maxSize := shopSetting["server"]["logging"].MaxSize
	MaxBackups := shopSetting["server"]["logging"].MaxBackups
	MaxAge := shopSetting["server"]["logging"].MaxAge
	compress := shopSetting["server"]["logging"].Compress

	return model.Logging{Filename: fileName, Level: level, MaxSize: maxSize, MaxBackups: MaxBackups, MaxAge: MaxAge, Compress: compress}
}
