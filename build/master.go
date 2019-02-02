package build

import (
	"context"
	"github.com/sirupsen/logrus"
	"time"
)

const preBuildNotifyName = "pre"
const postBuildNotifyName = "post"

type NotificationFunc func(build Build) error

func New(repoURL string, commit string, originator string, customPreBuildNotify NotificationFunc, customPostBuildNotify NotificationFunc) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Hour*2))

	buildID, err := newUUID()
	if err != nil {
		logrus.Errorf("error generating build ID: %v", err)
	}

	newBuild := &Build{
		ID:         buildID,
		Originator: originator,
		RepoURL:    repoURL,
		CommitHash: commit,
		StartedAt:  time.Now(),
		FinishedAt: time.Time{},
		Error:      "",
		Context:    ctx,
		CancelFunc: cancel,
	}

	buildQueue = append(buildQueue, newBuild)

	triggerNotifyApi(ctx, preBuildNotifyName, customPreBuildNotify, *newBuild, getPreNotificationScriptsDir())

	err = startBuild(ctx, repoURL, commit)
	if err != nil {
		logrus.Errorf("error during Build: %v", err)
		newBuild.Error = err.Error()
	}
	newBuild.FinishedAt = time.Now()

	triggerNotifyApi(ctx, postBuildNotifyName, customPostBuildNotify, *newBuild, getPostNotificationScriptsDir())
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
