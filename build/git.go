package build

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

func prepareGitRepository(ctx context.Context, repoURL string, buildRev string) error {

	tmpPath := getRepoPath(repoURL)

	logrus.Infof("checking out repository into: %v", tmpPath)
	repo, err := git.PlainCloneContext(ctx, tmpPath, false, &git.CloneOptions{
		URL:      repoURL,
		Progress: os.Stdout,
	})
	if err != nil {
		if err == git.ErrRepositoryAlreadyExists {
			logrus.Infof("repository already cloned")
			repo, err = git.PlainOpen(tmpPath)
			if err != nil {
				return err
			}
			logrus.Infof("fetching")
			err = repo.FetchContext(ctx, &git.FetchOptions{})
			if err != nil && err != git.NoErrAlreadyUpToDate {
				return err
			}
		} else {
			return err
		}
	}
	workTree, err := repo.Worktree()
	if err != nil {
		return err
	}
	logrus.Infof("reset hard to build revision %v", buildRev)
	err = workTree.Reset(&git.ResetOptions{
		Commit: plumbing.NewHash(buildRev),
		Mode:   git.HardReset,
	})
	if err != nil {
		return err
	}
	logrus.Infof("done")

	return nil
}
