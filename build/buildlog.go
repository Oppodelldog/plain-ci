package build

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strconv"
	"sync"
)

var logFilenameRegEx = regexp.MustCompile(`(?P<Number>\d+)_(?P<Status>\d+)_(?P<Result>\d+).txt`)

type LogFile struct {
	RepoURL  string   `json:"-"`
	FilePath string   `json:"-"`
	Number   int      `json:"number"`
	Status   Status   `json:"build_state"`
	Result   Result   `json:"result"`
	file     *os.File `json:"-"`
	fileLock sync.RWMutex
}

func getBuildLogFilename(logfile *LogFile) string {
	return fmt.Sprintf("%d_%d_%d.txt", logfile.Number, logfile.Status, logfile.Result)
}

func (l *LogFile) Close() error {
	if l.file == nil {
		err := l.file.Close()
		if err != nil {
			return err
		}
		l.file = nil

		return nil
	}

	return nil
}

func (l *LogFile) Write(p []byte) (n int, err error) {
	l.fileLock.Lock()
	defer l.fileLock.Unlock()

	if l.file == nil {
		l.file, err = os.OpenFile(l.FilePath, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			return 0, err
		}
	}

	return l.file.Write(p)
}

func (l *LogFile) openFile() (error) {
	l.FilePath = path.Join(getBuildLogPath(l.RepoURL), getBuildLogFilename(l))

	err := ensureDir(path.Dir(l.FilePath))
	if err != nil {
		return err
	}

	l.file, err = os.OpenFile(l.FilePath, os.O_WRONLY|os.O_CREATE, 0666)

	return err
}

func (l *LogFile) ChangeStatus(status Status) {
	l.fileLock.Lock()
	defer l.fileLock.Unlock()

	err := l.file.Close()
	if err != nil {
		logrus.Errorf("Could not change status to %s, error closing file: %v", status, err)
		return
	}

	originalStatus := l.Status
	originalFileName := getBuildLogFilename(l)
	l.Status = status
	newFileName := getBuildLogFilename(l)

	buildPath := getBuildLogPath(l.RepoURL)
	oldPath := path.Join(buildPath, originalFileName)
	newPath := path.Join(buildPath, newFileName)
	err = os.Rename(oldPath, newPath)
	if err != nil {
		l.Status = originalStatus
		logrus.Errorf("Could not change status to %s, error renaming file: %v", status, err)
		return
	}

	err = l.openFile()
	if err != nil {
		l.Status = originalStatus
		logrus.Errorf("Could not open new build log file after change status to %s: %v", status, err)
		return
	}
}

func (l *LogFile) ChangeResult(result Result) {
	l.fileLock.Lock()
	defer l.fileLock.Unlock()

	err := l.file.Close()
	if err != nil {
		logrus.Errorf("Could not change result to %s, error closing file: %v", result, err)
		return
	}

	originalResult := l.Result
	originalFileName := getBuildLogFilename(l)
	l.Result = result
	newFileName := getBuildLogFilename(l)

	buildPath := getBuildLogPath(l.RepoURL)
	oldPath := path.Join(buildPath, originalFileName)
	newPath := path.Join(buildPath, newFileName)
	
	err = os.Rename(oldPath, newPath)
	if err != nil {
		l.Result = originalResult
		logrus.Errorf("Could not change result to %s, error renaming file: %v", result, err)
		return
	}

	err = l.openFile()
	if err != nil {
		l.Result = originalResult
		logrus.Errorf("Could not open new build log file after change result  to %s: %v", result, err)
		return
	}
}

func NewLogFile(repoURL string, number int) (*LogFile, error) {
	logFile := &LogFile{
		RepoURL: repoURL,
		Number:  number,
	}

	err := logFile.openFile()
	if err != nil {
		return nil, err
	}

	return logFile, nil
}

func FromFilePath(filePath string) (LogFile, error) {
	fileName := path.Base(filePath)

	fileNameParts := readRegExGroups(logFilenameRegEx, fileName)
	if len(fileNameParts) != 3 {
		logrus.Errorf("error parsing build log filename: %s", fileName)
	}

	number, err := strconv.Atoi(fileNameParts["Number"])
	if err != nil {
		return LogFile{}, err
	}
	if number == 0 {
		return LogFile{}, fmt.Errorf("filename parsing error: number may not be 0 (given filename '%s')", fileName)
	}
	status, err := strconv.Atoi(fileNameParts["Status"])
	if err != nil {
		return LogFile{}, err
	}
	result, err := strconv.Atoi(fileNameParts["Result"])
	if err != nil {
		return LogFile{}, err
	}

	return LogFile{
		Number:   number,
		FilePath: filePath,
		Status:   Status(status),
		Result:   Result(result),
	}, nil
}

func getNextBuildLogfile(repoUrl string) (*LogFile, error) {
	buildsDir := getBuildLogPath(repoUrl)
	var buildNo int
	if _, err := os.Stat(buildsDir); os.IsNotExist(err) {
		buildNo = 1
	} else {
		files, err := ioutil.ReadDir(buildsDir)
		if err != nil {
			return nil, err
		}
		buildNo = len(files) + 1
	}
	return NewLogFile(repoUrl, buildNo)
}

func ensureDir(path string) error {
	return os.MkdirAll(path, 0766)
}
