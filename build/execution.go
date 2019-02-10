package build

import (
	"fmt"
	"io"
	"os/exec"
	"path"
	"time"

	"github.com/sirupsen/logrus"
)

func executeBuild(build *Build) error {
	f := build.LogFile
	timeStart := time.Now()
	logrus.Infof("Starting build %v", build.LogFile.Number)

	ciScript := path.Join(getRepoPath(build.RepoURL), ".build", "ci.sh")
	defer func() {
		err := f.Close()
		if err != nil {
			logrus.Errorf("error closing build log writer: %v", err)
		}
	}()

	logWrite := createWriteLogFunc(f)

	logWrite("Build: %v\n", build.LogFile.Number)
	logWrite("Date: %v\n", time.Now())
	logWrite("Project: %v\n", build.RepoURL)
	logWrite("Commit-Hash: %v\n", build.CommitHash)
	logWrite("------------------------------------------------------------------------------------------------------------------------------")

	ciCmd := exec.CommandContext(build.Context, ciScript)
	ciCmd.Stdout = f
	ciCmd.Stderr = f

	err := ciCmd.Start()
	if err != nil {
		logWrite("error while starting: %v\n", err)
		logWrite("BUILD FAILED", nil)
		build.ChangeStatus(Finished)
		build.ChangeResult(Failure)
		return err
	}

	logrus.Infof("Waiting for build %v to finish", build.LogFile.Number)
	err = ciCmd.Wait()

	logWrite("------------------------------------------------------------------------------------------------------------------------------")
	logWrite("duration: %v\n", time.Since(timeStart))

	if err != nil {
		logrus.Infof("error during build: %v", err)
		logWrite("error during build: %v\n", err)
		logWrite("BUILD FAILED")
		build.ChangeStatus(Finished)
		build.ChangeResult(Failure)
		return err
	} else {
		
		logWrite("BUILD SUCCESS")
		build.ChangeStatus(Finished)
		build.ChangeResult(Success)
	}

	return nil
}

func createWriteLogFunc(log io.Writer) func(format string, args ...interface{}) {
	return func(format string, args ...interface{}) {
		var err error
		if len(args) > 0 {
			_, err = fmt.Fprintf(log, format, args)
		} else {
			_, err = fmt.Fprint(log, format)
		}
		if err != nil {
			logrus.Errorf("could not write to log log: %v", err)
		}
		_, err = fmt.Fprintln(log)
		if err != nil {
			logrus.Errorf("could not write line break to log: %v", err)
		}
	}
}
