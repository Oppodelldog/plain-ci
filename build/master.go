package build

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
)

const preBuildNotifyName = "pre"
const postBuildNotifyName = "post"

type NotificationFunc func(build Build) error

func New(repoURL string, commit string, originator string, customPreBuildNotify NotificationFunc, customPostBuildNotify NotificationFunc) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Hour*2))

	buildID, err := newUUID()
	if err != nil {
		logrus.Errorf("error generating build ID: %v", err)
		return
	}
	buildNo, err := getNextBuildNo(repoURL)
	if err != nil {
		logrus.Errorf("Error starting build for %s on commit %s", repoURL, commit)
		return
	}

	newBuild := &Build{
		ID:         buildID,
		RepoPath:   normalizeUrl(repoURL),
		No:         buildNo,
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

	err = startBuild(ctx, buildNo, repoURL, commit)
	if err != nil {
		logrus.Errorf("error during Build: %v", err)
		newBuild.Error = err.Error()
	}
	newBuild.FinishedAt = time.Now()
	newBuild.Done = true

	triggerNotifyApi(ctx, postBuildNotifyName, customPostBuildNotify, *newBuild, getPostNotificationScriptsDir())
}

func startBuild(ctx context.Context, buildNo int, repoURL, commit string) error {

	logrus.Infof("Starting build %v for %s on commit %s", buildNo, repoURL, commit)
	err := prepareGitRepository(ctx, repoURL, commit)
	if err != nil {
		return err
	}

	err = executeBuild(ctx, repoURL, commit, buildNo)
	if err != nil {
		return err
	}

	return nil
}
