package build

import (
	"github.com/Oppodelldog/simpleci/config"
	"path/filepath"
	"strings"
)

func getRepoPath(repoUrl string) string {
	return filepath.Join(getPath(repoUrl), "repo")
}

func getBuildsPath(repoUrl string) string {
	return filepath.Join(getPath(repoUrl), "builds")
}
func getPath(repoURL string) string {
	tmpPath := repoURL
	tmpPath = strings.Replace(tmpPath, "/", "_", -1)
	tmpPath = strings.Replace(tmpPath, ":", "_", -1)
	tmpPath = strings.Replace(tmpPath, "\\", "_", -1)
	tmpPath = strings.Replace(tmpPath, "?", "_", -1)
	tmpPath = strings.Replace(tmpPath, " ", "_", -1)

	tmpPath = filepath.Join(config.BuildDir, tmpPath)

	return tmpPath
}
