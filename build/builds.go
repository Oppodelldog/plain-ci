package build

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/Oppodelldog/plainci/config"
	"github.com/sirupsen/logrus"
)

var buildQueue []*Build

type Build struct {
	ID            string                      `json:"id"`
	No            int                         `json:"no"`
	LogFile       *LogFile                    `json:"-"`
	RepoURL       string                      `json:"repo_url"`
	RepoPath      string                      `json:"repo_path"`
	CommitHash    string                      `json:"commit_hash"`
	Originator    string                      `json:"originator"`
	StartedAt     time.Time                   `json:"started_at"`
	FinishedAt    time.Time                   `json:"finished_at"`
	Status        Status                      `json:"status"`
	Result        Result                      `json:"result"`
	Error         string                      `json:"error"`
	Context       context.Context             `json:"-"`
	CancelFunc    context.CancelFunc          `json:"-"`
	Notifications map[string]NotificationFunc `json:"-"`
}

func (b *Build) Abort() {
	b.CancelFunc()
	b.ChangeStatus(Finished)
	b.ChangeResult(Aborted)
}

func (b *Build) ChangeStatus(status Status) {
	b.LogFile.ChangeStatus(status)
	b.Status = b.LogFile.Status
}

func (b *Build) ChangeResult(result Result) {
	b.LogFile.ChangeResult(result)
	b.Result = b.LogFile.Result
}

type Project struct {
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

func GetProjects() []Project {
	var repositories []Project
	err := filepath.Walk(config.BuildDir, func(f string, info os.FileInfo, err error) error {
		if info.IsDir() && info.Name() == ".git" {
			l, _ := filepath.Abs(f + "/../..")
			if err != nil {
				logrus.Errorf("error making absolute path of project dir %s: %v", f, err)
				return nil
			}
			repoPath := strings.Replace(l, config.BuildDir+"/", "", -1)

			repositories = append(repositories, Project{
				RepoPath: repoPath,
			})
		}
		return nil
	})

	if err != nil {
		logrus.Errorf("error reading project dirs '%s': %v", config.BuildDir, err)
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
