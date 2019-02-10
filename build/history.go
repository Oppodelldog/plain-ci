package build

import (
	"fmt"
	"github.com/Oppodelldog/plainci/config"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"path"
	"regexp"
	"sort"
)

func GetAllBuilds() map[string][]LogFile {
	res := map[string][]LogFile{}
	repoFolders, err := ioutil.ReadDir(config.BuildDir)
	if err != nil {
		logrus.Errorf("error reading build dir %s: %v", config.BuildDir, err)
		return nil
	}

	for _, repoFolder := range repoFolders {
		if !repoFolder.IsDir() {
			continue
		}

		res[repoFolder.Name()] = GetBuildLogs(repoFolder.Name())
	}
	return res
}

func GetBuildLogs(repoDir string) []LogFile {
	var logFiles []LogFile
	logPath := getBuildLogPath(repoDir)
	files, err := ioutil.ReadDir(logPath)
	if err != nil {
		logrus.Errorf("error reading build dir for repositoryURL '%s': %v", err)
		return nil
	}

	for _, file := range files {
		logFilePath := path.Join(logPath, file.Name())
		logFile, err := FromFilePath(logFilePath)
		if err != nil {
			logrus.Errorf("error reading build log info: %v", err)
			continue
		}
		logFiles = append(logFiles, logFile)
	}

	sort.Slice(logFiles, func(i, j int) bool {
		return logFiles[i].Number > logFiles[j].Number
	})

	return logFiles
}

func GetBuildLog(repoDir string, logId int) string {
	buildPath := getBuildLogPath(repoDir)
	files, err := ioutil.ReadDir(buildPath)
	if err != nil {
		logrus.Errorf("Error reading build dir %v", buildPath)
		return ""
	}

	for _, file := range files {
		filePrefix := fmt.Sprintf("%v_", logId)
		if file.Name()[:len(filePrefix)] == filePrefix {
			filePath := path.Join(buildPath, file.Name())

			b, err := ioutil.ReadFile(filePath)
			if err != nil {
				logrus.Errorf("Error reading build log %s %v", repoDir, logId)
				return ""
			}

			return string(b)
		}
	}

	return ""
}

func readRegExGroups(regEx *regexp.Regexp, s string) (paramsMap map[string]string) {

	match := regEx.FindStringSubmatch(s)

	paramsMap = make(map[string]string)
	for i, name := range regEx.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}
	return
}
