package build

import (
	"context"
	"io/ioutil"
	"os/exec"
	"path"
	"strings"

	"github.com/sirupsen/logrus"
)

func triggerNotifyApi(build *Build, buildStepName string, scriptDir string) {
	if notify, ok := build.Notifications[buildStepName]; ok && notify != nil {
		err := notify(*build)
		if err != nil {
			logrus.Errorf("error during %s-build action: %v", buildStepName, err)
		}
	}

	err := simpleNotification(build.Context, scriptDir, *build)
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
		"PLAIN_CI_BUILD_ID=" + build.ID,
		"PLAIN_CI_REPO_URL=" + build.RepoURL,
		"PLAIN_CI_COMMIT_HASH=" + build.CommitHash,
		"PLAIN_CI_ORIGINATOR=" + build.Originator,
		"PLAIN_CI_ERROR=" + build.Error,
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		logrus.Errorf("error while executing script '%s': %v", filePath, err)
	}

	logrus.Infof("script '%s': %v", filePath, string(output))
}
