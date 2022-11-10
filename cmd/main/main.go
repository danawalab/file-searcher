package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.danawa.com/devops/file-searcher/logging"
	"gitlab.danawa.com/devops/file-searcher/service"
	"gitlab.danawa.com/devops/file-searcher/setting"
	"golang.org/x/sync/errgroup"
)

func setupRouter() http.Handler {
	e := gin.New()
	e.GET("/", func(c *gin.Context) {
		c.JSON(200, "FileSearcher Running ...")
	})
	e.GET("/search", func(c *gin.Context) {
		result := make(map[string]interface{})
		var err error

		// 현재 날짜
		t := time.Now()
		formatted := fmt.Sprintf("%d%02d%02d", t.Year(), t.Month(), t.Day())

		// 필수값 validation check
		if len(c.Query("shopCode")) == 0 || len(c.Query("productId")) == 0 {
			result["status"] = http.StatusBadRequest
			result["message"] = "필수값이 입력되지 않았습니다. 상품 아이디와 샵 코드를 입력하세요."
			c.JSON(http.StatusOK, result)

			return
		}

		// 검색 실행
		data, err := service.SearchID(c.Query("shopCode"), c.Query("productId"), c.Query("startDate"), c.Query("endDate"), c.Query("renew"))
		for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
			data[i], data[j] = data[j], data[i]
		}

		if err == nil {
			// 검색 결과를 셋팅해서 리턴
			result["status"] = http.StatusOK

			if data != nil {
				result["data"] = data
			} else {
				result["data"] = "검색 결과가 없습니다."
			}

			result["date"] = c.DefaultQuery("date", formatted)
			c.JSON(http.StatusOK, result)
		} else {
			logging.Error(err.Error())
			c.JSON(http.StatusInternalServerError, result)
		}
	})
	return e
}

func main() {
	// 설정파일을 로드합니다.
	logging.Info(fmt.Sprint("SETTING_FILE_PATH : ", os.Args[1]))
	loadSettingError := setting.LoadSetting(os.Args[1])
	if loadSettingError != nil {
		logging.Error(fmt.Sprint("[ERROR] ", loadSettingError))
		os.Exit(1)
	}

	runtime.GOMAXPROCS(setting.GetServerCpuCoreConfig())

	logging.Info(fmt.Sprint("Used CPU Core Count: ", setting.GetServerCpuCoreConfig()))

	// 색인 서비스 시작
	go service.Initalize(setting.GetServerEpFilePathConfig(), setting.GetServerTempFilePathConfig(), setting.GetServerWorkerConfig())

	// 보관기간 제외 파일 자동 삭제
	go service.ScheduleFileDeleteInit(setting.GetServerEpFilePathConfig(), setting.GetFileDeletePeriodDay())

	// 로그 설정 합니다.
	go logging.LoadLogging(setting.GetLogging())

	// 파일 서쳐 서버를 생성합니다.
	apiServer := &http.Server{
		Addr:           fmt.Sprintf(":%s", setting.GetServerFilePathPort()),
		Handler:        setupRouter(),
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	gin.SetMode(gin.ReleaseMode)
	g := errgroup.Group{}
	g.Go(func() error {
		return apiServer.ListenAndServe()
	})
	if err := g.Wait(); err != nil {
		logging.Error(fmt.Sprint("[ERROR] ", err))
	}
}
