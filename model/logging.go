package model

type Logging struct {
	Filename   string `json:"filename"`
	Level      string `json:"level"`
	MaxSize    int    `json:"max_size"`
	MaxBackups int    `json:"max_backups"`
	MaxAge     int    `json:"max_age"`
	Compress   bool   `json:"compress"`
}
