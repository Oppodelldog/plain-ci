package build

import (
	"context"
	"github.com/sirupsen/logrus"
	"time"
)

type postBuildFunc func(build Build) error
type preBuildFunc func(build Build) error

func New(repoURL string, commit string, originator string, preBuild preBuildFunc, postBuild postBuildFunc) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Hour*2))

	buildID, err := newUUID()
	if err != nil {
		logrus.Errorf("error generating build ID: %v", err)
	}

	newBuild := &Build{
		ID:         buildID,
		Originator: originator,
		StartedAt:  time.Now(),
		RepoURL:    repoURL,
		CommitHash: commit,
		Context:    ctx,
		CancelFunc: cancel,
	}

	buildQueue = append(buildQueue, newBuild)

	if preBuild != nil {
		err = preBuild(*newBuild)
		if err != nil {
			logrus.Errorf("aborting build due to error during pre-build action: %v", err)
			newBuild.Error = err.Error()
			return
		}
	}

	err = startBuild(ctx, repoURL, commit)
	if err != nil {
		logrus.Errorf("error during Build: %v", err)
		newBuild.Error = err.Error()
	}
	newBuild.FinishedAt = time.Now()

	if postBuild != nil {
		err = postBuild(*newBuild)
		if err != nil {
			logrus.Errorf("error during post-build action: %v", err)
		}
	}
}

func startBuild(ctx context.Context, repoURL, buildRev string) error {
	buildNo, err := getNextBuildNo(repoURL)
	if err != nil {
		logrus.Errorf("Error starting build for %s on commit %s", repoURL, buildRev)
		return err
	}

	logrus.Infof("Starting build %v for %s on commit %s", buildNo, repoURL, buildRev)
	err = prepareGitRepository(ctx, repoURL, buildRev)
	if err != nil {
		return err
	}

	err = executeBuild(ctx, repoURL, buildRev, buildNo)
	if err != nil {
		return err
	}

	return nil
}
