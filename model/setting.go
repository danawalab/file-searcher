package model

type Setting struct {
	Code                string `yaml:"code"`
	Delimiter           string `yaml:"delimiter"`
	Epfilepath          string `yaml:"epfilepath"`
	CpuCore             int    `yaml:"cpuCore"`
	TempFilepath        string `yaml:"tempfilepath"`
	Port                string `yaml:"port"`
	Server              string `yaml:"server"`
	Config              string `yaml:"config"`
	StartWord           string `yaml:"startWord"`
	EndWord             string `yaml:"endWord"`
	IdWord              string `yaml:"idWord"`
	Header              bool   `yaml:"header"`
	Custom              bool   `yaml:"custom"`
	IdPostion           int    `yaml:"idPostion"`
	IndexInterval       int    `yaml:"indexInterval"`
	FileSortDivision    int    `yaml:"fileSortDivision"`
	ColumnCount         int    `yaml:"columnCount"`
	FileDeletePeriodDay int    `yaml:"fileDeletePeriodDay"`
	Filename            string `yaml:"filename"`
	Level               string `yaml:"level"`
	MaxSize             int    `yaml:"max_size"`
	MaxBackups          int    `yaml:"max_backups"`
	MaxAge              int    `yaml:"max_age"`
	Compress            bool   `yaml:"compress"`
	Workers             int    `yaml:"workers"`
}
