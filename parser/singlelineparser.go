package parser

import (
	"strings"

	"gitlab.danawa.com/devops/file-searcher/setting"
)

type SingleLineParser struct {
	Delimiter   string
	IdPosition  int
	Header      bool
	ColumnCount int
}

func NewSingleLineParser(shop string) *SingleLineParser {
	shopConfig := setting.GetShopParseConfig(shop)
	singleLineParser := &SingleLineParser{Delimiter: shopConfig.Delimiter, IdPosition: shopConfig.IdPostion, ColumnCount: shopConfig.ColumnCount}
	return singleLineParser
}

func (s SingleLineParser) SingleParse(line string) string {
	var productId string

	lineSplitArray := strings.Split(line, s.Delimiter)

	// 유효한 라인인지 검사(최소한 컬럼 갯수보다 같거나 커야한다)
	if len(lineSplitArray) >= s.ColumnCount {
		if len(lineSplitArray) > s.IdPosition {
			productId = lineSplitArray[s.IdPosition]
		}
	}

	return productId
}
