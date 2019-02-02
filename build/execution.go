package build

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"os/exec"
	"path"
	"time"
)

func executeBuild(ctx context.Context, repoURL, commitHash string, buildNo int) error {

	logrus.Infof("Starting build %v", buildNo)

	repoDir := getRepoPath(repoURL)
	buildsDir := getBuildLogPath(repoURL)
	err := ensureDir(buildsDir)
	if err != nil {
		return err
	}

	ciScript := path.Join(repoDir, ".build", "ci.sh")

	buildLogFile := path.Join(buildsDir, fmt.Sprintf("%v.txt", buildNo))
	f, err := os.OpenFile(buildLogFile, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer func() {
		err := f.Close()
		if err != nil {
			logrus.Errorf("error closing build log file: %v", err)
		}
	}()

	_, err = fmt.Fprintf(f, "Build: %v\n", buildNo)
	_, err = fmt.Fprintf(f, "Date: %v\n", time.Now())
	_, err = fmt.Fprintf(f, "Repository: %v\n", repoURL)
	_, err = fmt.Fprintf(f, "Commit-Hash: %v\n", commitHash)
	if err != nil {
		panic(err)
	}
	writeSeperator(f)

	ciCmd := exec.CommandContext(ctx, ciScript)
	ciCmd.Stdout = f
	ciCmd.Stderr = f

	err = ciCmd.Start()
	if err != nil {
		writeBuildFinishError(f, err)
		return err
	}

	logrus.Infof("Waiting for build %v to finish", buildNo)
	err = ciCmd.Wait()
	if err != nil {
		writeBuildFinishError(f, err)
		return err
	} else {
		writeBuildFinishSuccess(f)
	}

	return nil
}

func writeBuildFinishError(f io.Writer, err error) {
	writeSeperator(f)
	_, _ = fmt.Fprintf(f, "error: %v\n", err)
	_, _ = fmt.Fprint(f, "BUILD FAILED")
}

func writeBuildFinishSuccess(f io.Writer) {
	writeSeperator(f)
	_, _ = fmt.Fprint(f, "BUILD SUCCESS")
}

func writeSeperator(f io.Writer) {
	_, _ = fmt.Fprint(f, "------------------------------------------------------------------------------------------------------------------------------\n")
}

func ensureDir(path string) error {
	return os.MkdirAll(path, 0766)
}
