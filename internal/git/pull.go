package git

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

func PullLatest(repoPath, repoUrl, branch, username, token string) error {
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		fmt.Println("Cloning repository...")
		_, err := git.PlainClone(repoPath, false, &git.CloneOptions{
			URL:           repoUrl,
			Progress:      os.Stdout,
			ReferenceName: plumbing.NewBranchReferenceName(branch),
			Auth: &http.BasicAuth{
				Username: username,
				Password: token,
			},
		})
		return err
	}

	fmt.Println("Pulling latest changes...")
	repo, err := git.PlainOpen(repoPath) // returns a *repo and PlainOpen just opens the existing repository
	if err != nil {
		return err
	}
	w, err := repo.Worktree() // returns a *Worktree which is used to perform operations on the repository
	if err != nil {
		return err
	}
	err = w.Pull(&git.PullOptions{
		RemoteName:    "origin",
		ReferenceName: plumbing.NewBranchReferenceName(branch),
		Auth: &http.BasicAuth{
			Username: username,
			Password: token,
		},
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return err
	}

	return nil
}
