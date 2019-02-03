package build

import (
	"github.com/Oppodelldog/simpleci/config"
	"os"
	"path/filepath"
	"strings"
)

func getRepoPath(repoUrl string) string {
	return filepath.Join(getPath(repoUrl), "repo")
}

func getBuildLogPath(repoUrl string) string {
	return filepath.Join(getPath(repoUrl), "builds")
}

func getPath(repoURL string) string {
	tmpPath := normalizeUrl(repoURL)

	tmpPath = filepath.Join(config.BuildDir, tmpPath)

	return tmpPath
}

func normalizeUrl(repoURL string) string {
	tmpPath := repoURL
	tmpPath = strings.Replace(tmpPath, "/", "_", -1)
	tmpPath = strings.Replace(tmpPath, ":", "_", -1)
	tmpPath = strings.Replace(tmpPath, "\\", "_", -1)
	tmpPath = strings.Replace(tmpPath, "?", "_", -1)
	tmpPath = strings.Replace(tmpPath, " ", "_", -1)

	return tmpPath
}

func getPreNotificationScriptsDir() string {
	return filepath.Join(getNotificationScriptsDir(), "pre")
}

func getPostNotificationScriptsDir() string {
	return filepath.Join(getNotificationScriptsDir(), "post")
}

func getNotificationScriptsDir() string {
	workingDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return filepath.Join(workingDir, ".notify")
}
