package build

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/sirupsen/logrus"
)

func executeBuild(ctx context.Context, repoURL, commitHash string, buildNo int) error {

	timeStart := time.Now()
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

	logWrite := createWriteLogFunc(f)

	logWrite("Build: %v\n", buildNo)
	logWrite("Date: %v\n", time.Now())
	logWrite("Project: %v\n", repoURL)
	logWrite("Commit-Hash: %v\n", commitHash)
	if err != nil {
		panic(err)
	}
	logWrite("------------------------------------------------------------------------------------------------------------------------------")

	ciCmd := exec.CommandContext(ctx, ciScript)
	ciCmd.Stdout = f
	ciCmd.Stderr = f

	err = ciCmd.Start()
	if err != nil {
		logWrite("error while starting: %v\n", err)
		logWrite("BUILD FAILED", nil)
		return err
	}

	logrus.Infof("Waiting for build %v to finish", buildNo)
	err = ciCmd.Wait()

	logWrite("------------------------------------------------------------------------------------------------------------------------------")
	logWrite("duration: %v\n", time.Since(timeStart))

	if err != nil {
		logWrite("error during build: %v\n", err)
		logWrite("BUILD FAILED")
		return err
	} else {
		logWrite("BUILD SUCCESS")
	}

	return nil
}

func createWriteLogFunc(file *os.File) func(format string, args ...interface{}) {
	return func(format string, args ...interface{}) {
		var err error
		if len(args) > 0 {
			_, err = fmt.Fprintf(file, format, args)
		} else {
			_, err = fmt.Fprint(file, format)
		}
		if err != nil {
			logrus.Errorf("could not write to log file '%s': %v", file.Name(), err)
		}
		_, err = fmt.Fprintln(file)
		if err != nil {
			logrus.Errorf("could not write line break to log file '%s': %v", file.Name(), err)
		}
	}
}

func ensureDir(path string) error {
	return os.MkdirAll(path, 0766)
}
