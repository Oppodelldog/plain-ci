package build

import (
	"context"
	"fmt"
	"time"
)

var buildQueue []*Build

type Build struct {
	ID         string             `json:"id"`
	RepoURL    string             `json:"repo_url"`
	CommitHash string             `json:"commit_hash"`
	Originator string             `json:"originator"`
	StartedAt  time.Time          `json:"started_at"`
	FinishedAt time.Time          `json:"finished_at"`
	Error      string             `json:"error"`
	Context    context.Context    `json:"-"`
	CancelFunc context.CancelFunc `json:"-"`
}

func GetBuildQueueList() []Build {
	var newBuilds []Build
	for _, build := range buildQueue {
		newBuilds = append(newBuilds, *build)
	}

	return newBuilds
}

func AbortBuild(id string) error {
	for _, build := range buildQueue {
		if build.ID == id {
			build.CancelFunc()
			return nil
		}
	}

	return fmt.Errorf("build %v not found", id)
}
