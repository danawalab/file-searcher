package service

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"sync"

	"github.com/farmergreg/rfsnotify"
	"gitlab.danawa.com/devops/file-searcher/logging"
)

var indexMutex = sync.RWMutex{}

func Initalize(epFilePath, tmpFilePath string, workers int) {
	defer func() {
		if r := recover(); r != nil {
			logging.Error(fmt.Sprint("Recovered in Initailize : ", r))
		}
	}()

	// 기존 색인 파일은 적재하기
	LoadIndexFiles(epFilePath, tmpFilePath, workers)

	// 기존 색인 안된 파일들 인덱싱
	IndexingAllFiles(epFilePath, tmpFilePath, workers)

	// 신규 디렉토리 감지, 왓처 생성
	listenNewDicrectory(epFilePath, tmpFilePath, workers)
}

// 디렉토리를 감시하다가 파일이 들어오면 색인을 실행한다
func listenNewDicrectory(epFilePath, tmpFilePath string, workers int) {
	semaphore := make(chan int, workers)

	defer func() {
		if r := recover(); r != nil {
			logging.Error(fmt.Sprint("Recovered in Listen Directory : ", r))
		}
	}()

	watcher, err := rfsnotify.NewWatcher()
	if err != nil {
		logging.Error(fmt.Sprint("CREATE NEW WATHCER ERROR ", err))
	}
	defer watcher.Close()

	// 변화를 감지한다
	done := make(chan bool)
	go func() {
		defer close(done)

		for {
			select {

			case event := <-watcher.Events:
				switch event.Op.String() {
				case "CREATE":
					// ".txt" 파일 인덱싱 처리
					if strings.HasSuffix(event.Name, ".txt") {
						check := strings.Split(event.Name, "-")

						if len(check) > 2 {
							semaphore <- 1
							go func() {
								logging.Info(fmt.Sprint("INIT INDEXING  ... ", event.Name))
								Indexing(event.Name, tmpFilePath)
								<-semaphore
							}()
						}
					}
				default:
				}
			case err := <-watcher.Errors:
				logging.Error(fmt.Sprint("TO MANY QUEUE ERROR ", err))
			}
		}
	}()

	// 감시할 경로를 등록한다
	err = watcher.AddRecursive(epFilePath)
	if err != nil {
		logging.Error(err.Error())
	}
	<-done
}

// 기존 인덱스 파일 로드
func LoadIndexFiles(epFilePath, tmpFilePath string, workers int) {
	semaphore := make(chan int, workers)

	defer func() {
		if r := recover(); r != nil {
			logging.Error(fmt.Sprint("Recovered in Indexing : ", r))
		}
	}()

	// 타겟 디렉토리 전체 파일 조사
	filepath.WalkDir(epFilePath,
		func(path string, entry fs.DirEntry, err error) error {
			// 1개씩 넣는다. 최대 크기에 도달하면 블락된다
			semaphore <- 1
			go func() {
				info, _ := entry.Info()

				if !info.IsDir() {
					check := strings.Split(info.Name(), "-")

					// tmon-1561929118-0 꼴 형태의 파일만 가능
					if len(check) > 2 {
						if strings.Contains(info.Name(), ".index.index") {
							// 최초 시작시 index.index 파일은 모두 등록해준다
							Load(path)
						}
					}
				}
				// 고 루틴 한 개가 제거되면 해당 제거한다. 다음 순서가 진행된다
				<-semaphore
			}()
			return nil
		})
}

// 전체 대상 파일에 대해 인덱싱을 한다
func IndexingAllFiles(epFilePath, tmpFilePath string, workers int) {
	semaphore := make(chan int, workers)

	defer func() {
		if r := recover(); r != nil {
			logging.Error(fmt.Sprint("Recovered in Indexing : ", r))
		}
	}()

	// 타겟 디렉토리 전체 파일 조사
	filepath.WalkDir(epFilePath,
		func(path string, entry fs.DirEntry, err error) error {
			// 1개씩 넣는다. 최대 크기에 도달하면 블락된다
			semaphore <- 1
			go func() {
				info, _ := entry.Info()

				if !info.IsDir() {
					check := strings.Split(info.Name(), "-")

					// tmon-1561929118-0 꼴 형태의 파일만 가능
					if len(check) > 2 {
						if !strings.Contains(info.Name(), ".index") && !strings.Contains(info.Name(), ".tmp") {
							curPath := strings.Replace(path, info.Name(), "", -1)
							if !scanHasIndexFile(curPath, info.Name()) {

								logging.Info(fmt.Sprint("INIT INDEXING  ... ", path))
								Indexing(path, tmpFilePath)

							}
						}
					}
				}
				// 고 루틴 한 개가 제거되면 해당 제거한다. 다음 순서가 진행된다
				<-semaphore
			}()
			return nil
		})
}

func scanHasIndexFile(filePath, fileName string) bool {
	defer func() {
		if r := recover(); r != nil {
			logging.Error(fmt.Sprint("Recovered in scanIndexFiles : ", r))
		}
	}()

	var hasIndex bool

	// 타겟 디렉토리 전체 파일 조사
	filepath.WalkDir(filePath,
		func(path string, entry fs.DirEntry, err error) error {
			info, errEntry := entry.Info()
			if errEntry != nil {
				return errEntry
			}

			if err != nil {
				return err
			}

			// 인덱스가 있는지 조사
			if !info.IsDir() && strings.Contains(info.Name(), fileName) && (strings.Contains(info.Name(), ".index") || strings.Contains(info.Name(), ".tmp")) {
				hasIndex = true
			}

			return nil
		})

	return hasIndex
}
