package build

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

const preBuildNotifyName = "pre"
const postBuildNotifyName = "post"

type NotificationFunc func(build Build) error

type Queue struct {
	newBuildsChannel chan *Build
	shutdownChannel  chan bool
	ctx              context.Context
}

func NewQueue(ctx context.Context) *Queue {
	newBuildsChannel := make(chan *Build)
	q := &Queue{
		newBuildsChannel: newBuildsChannel,
		shutdownChannel:  make(chan bool),
		ctx:              ctx,
	}

	q.start()

	return q
}

func (q *Queue) GetShutDownChannel() chan bool {
	return q.shutdownChannel
}

func (q *Queue) NewBuild(repoURL string, commit string, originator string, customPreBuildNotify NotificationFunc, customPostBuildNotify NotificationFunc) {

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

	q.newBuildsChannel <- &Build{
		ID:         buildID,
		RepoPath:   normalizeUrl(repoURL),
		No:         buildNo,
		Originator: originator,
		RepoURL:    repoURL,
		CommitHash: commit,
		StartedAt:  time.Now(),
		FinishedAt: time.Time{},
		Error:      "",
		Notifications: map[string]NotificationFunc{
			preBuildNotifyName:  customPreBuildNotify,
			postBuildNotifyName: customPostBuildNotify,
		},
	}
}

func (q *Queue) AbortBuild(id string) error {
	if _, build, ok := q.findBuildIndex(id); ok {
		build.Abort()

		return nil
	} else {
		return fmt.Errorf("build %v not found", id)
	}
}

func (q *Queue) findBuildIndex(id string) (int, *Build, bool) {
	for index, build := range buildQueue {
		if build.ID == id {
			return index, build, true
		}
	}
	return 0, nil, false
}

func (q *Queue) initializeBuild(build *Build) {
	defer removeBuild(build)

	build.Context, build.CancelFunc = context.WithDeadline(q.ctx, time.Now().Add(time.Hour*2))

	triggerNotifyApi(*build, preBuildNotifyName, getPreNotificationScriptsDir())

	build.Status = Building

	err := startBuild(build, build.No, build.RepoURL, build.CommitHash)
	if err != nil {
		logrus.Errorf("error during Build: %v", err)
		build.Error = err.Error()
	}

	build.FinishedAt = time.Now()
	build.Status = Finished

	triggerNotifyApi(*build, postBuildNotifyName, getPostNotificationScriptsDir())
}

func startBuild(build *Build, buildNo int, repoURL, commit string) error {

	logrus.Infof("Starting build %v for %s on commit %s", buildNo, repoURL, commit)
	err := prepareGitRepository(build.Context, repoURL, commit)
	if err != nil {
		return err
	}

	err = executeBuild(build.Context, repoURL, commit, buildNo)
	if err != nil {
		return err
	}

	return nil
}

func (q *Queue) start() {

	logrus.Info("starting queue")
	go func() {
		for {
			select {
			case <-q.ctx.Done():
				logrus.Info("stopping queue")
				close(q.shutdownChannel)
				return

			case newBuild := <-q.newBuildsChannel:
				logrus.Infof("queue received new build (%s, %s, %v)", newBuild.RepoURL, newBuild.CommitHash, newBuild.No)
				buildQueue = append(buildQueue, newBuild)
			default:
				if build, ok := GetNextBuild(); ok {
					logrus.Infof("queue starting build (%s, %s, %v)", build.RepoURL, build.CommitHash, build.No)
					go q.initializeBuild(build)
				}
				time.Sleep(time.Second)
			}
		}
	}()
}
