package build

import (
	"context"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os/exec"
	"path"
	"strings"
)

func triggerNotifyApi(ctx context.Context, buildStepName string, postBuild NotificationFunc, build Build, scriptDir string) {
	if postBuild != nil {
		err := postBuild(build)
		if err != nil {
			logrus.Errorf("error during %s-build action: %v", buildStepName, err)
		}
	}

	err := simpleNotification(ctx, scriptDir, build)
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
		"PLAIN_CIBUILD_ID=" + build.ID,
		"PLAIN_CIREPO_URL=" + build.RepoURL,
		"PLAIN_CICOMMIT_HASH=" + build.CommitHash,
		"PLAIN_CIORIGINATOR=" + build.Originator,
		"PLAIN_CIERROR=" + build.Error,
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		logrus.Errorf("error while executing script '%s': %v", filePath, err)
	}

	logrus.Infof("script '%s': %v", filePath, string(output))
}
