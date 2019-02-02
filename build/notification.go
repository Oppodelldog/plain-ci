package build

import (
	"context"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os/exec"
	"path"
	"strings"
)

func triggerNotifyApi(ctx context.Context, buildStepName string, postBuild NotificationFunc, build Build) {
	if postBuild != nil {
		err := postBuild(build)
		if err != nil {
			logrus.Errorf("error during %s-build action: %v", buildStepName, err)
		}
	}

	err := simpleNotification(ctx, getPostNotificationScriptsDir(), build)
	if err != nil {
		logrus.Errorf("error during %s-build action: %v", buildStepName, err)
	}

}

func simpleNotification(ctx context.Context, dir string, build Build) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		extension := file.Name()[strings.LastIndex(file.Name(), ".")+1:]
		filePath := path.Join(dir, file.Name())
		switch strings.ToLower(extension) {
		case "sh":
			executeScript(ctx, filePath, build)
		}
	}

	return nil
}

func executeScript(ctx context.Context, filePath string, build Build) {
	cmd := exec.CommandContext(ctx, filePath)

	cmd.Env = []string{
		"SIMPLE_CI_BUILD_ID=" + build.ID,
		"SIMPLE_CI_REPO_URL=" + build.RepoURL,
		"SIMPLE_CI_COMMIT_HASH=" + build.CommitHash,
		"SIMPLE_CI_ORIGINATOR=" + build.Originator,
		"SIMPLE_CI_ERROR=" + build.Error,
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		logrus.Errorf("error while executing script '%s': %v", filePath, err)
	}

	logrus.Infof("script '%s': %v", filePath, string(output))
}
