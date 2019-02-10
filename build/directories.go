package build

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/Oppodelldog/plainci/config"
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
	tmpPath := strings.Replace(repoURL, "http://", "", 1)
	tmpPath = strings.Replace(tmpPath, "https://", "", 1)
	tmpPath = strings.Replace(tmpPath, "/", "_", -1)
	var gitExtension = ".git"
	if path.Ext(tmpPath) == gitExtension {
		tmpPath = tmpPath[:len(tmpPath)-len(gitExtension)]
	}

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
