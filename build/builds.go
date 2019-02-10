package build

import (
	"context"
	"io/ioutil"
	"sync"
	"time"

	"github.com/Oppodelldog/plainci/config"
	"github.com/sirupsen/logrus"
)

var buildQueue []*Build

type Build struct {
	ID            string                      `json:"id"`
	No            int                         `json:"no"`
	RepoURL       string                      `json:"repo_url"`
	RepoPath      string                      `json:"repo_path"`
	CommitHash    string                      `json:"commit_hash"`
	Originator    string                      `json:"originator"`
	StartedAt     time.Time                   `json:"started_at"`
	FinishedAt    time.Time                   `json:"finished_at"`
	Status        Status                      `json:"status"`
	Error         string                      `json:"error"`
	Context       context.Context             `json:"-"`
	CancelFunc    context.CancelFunc          `json:"-"`
	Notifications map[string]NotificationFunc `json:"-"`
}

func (b *Build) Abort() {
	b.CancelFunc()
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

func removeBuildByIndex(i int) {
	buildQueue[i] = buildQueue[len(buildQueue)-1]
	buildQueue[len(buildQueue)-1] = nil
	buildQueue = buildQueue[:len(buildQueue)-1]
}

func removeBuild(build *Build) {
	m := sync.Mutex{}
	m.Lock()
	defer m.Unlock()

	for index, b := range buildQueue {
		if b.ID == build.ID {
			removeBuildByIndex(index)
			return
		}
	}
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

func GetBuildsByState(status Status) []*Build {
	var builds []*Build
	for _, build := range buildQueue {
		if build.Status == status {
			builds = append(builds, build)
		}
	}
	return builds
}

func GetNextBuild() (*Build, bool) {
	if len(GetBuildsByState(Building)) > 0 {
		return nil, false
	}

	for _, build := range buildQueue {
		if build.Status == Queued {
			return build, true
		}
	}

	return nil, false
}
