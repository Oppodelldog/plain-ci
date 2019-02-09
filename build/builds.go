package build

import (
	"context"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/Oppodelldog/plainci/config"
	"github.com/sirupsen/logrus"
)

var buildQueue []*Build

type Build struct {
	ID         string             `json:"id"`
	No         int                `json:"no"`
	RepoURL    string             `json:"repo_url"`
	RepoPath   string             `json:"repo_path"`
	CommitHash string             `json:"commit_hash"`
	Originator string             `json:"originator"`
	StartedAt  time.Time          `json:"started_at"`
	FinishedAt time.Time          `json:"finished_at"`
	Done       bool               `json:"done"`
	Error      string             `json:"error"`
	Context    context.Context    `json:"-"`
	CancelFunc context.CancelFunc `json:"-"`
}

type Repository struct {
	URL      string
	RepoPath string
	Builds   []int
}

func GetBuildQueueList() []Build {
	var newBuilds []Build
	for _, build := range buildQueue {
		newBuilds = append(newBuilds, *build)
	}

	return newBuilds
}

func AbortBuild(id string) error {
	for index, build := range buildQueue {
		if build.ID == id {
			build.CancelFunc()
			removeBuild(index)
			return nil
		}
	}

	return fmt.Errorf("build %v not found", id)
}

func removeBuild(i int) {
	buildQueue[i] = buildQueue[len(buildQueue)-1]
	buildQueue[len(buildQueue)-1] = nil
	buildQueue = buildQueue[:len(buildQueue)-1]
}

func GetRepositories() []Repository {
	files, err := ioutil.ReadDir(config.BuildDir)
	if err != nil {
		logrus.Errorf("error reading builds dir '%s': %v", config.BuildDir, err)
		return nil
	}

	var repositories []Repository
	for _, file := range files {
		if !file.IsDir() {
			continue
		}
		repositories = append(repositories, Repository{
			RepoPath: file.Name(),
		})
	}

	return repositories
}
