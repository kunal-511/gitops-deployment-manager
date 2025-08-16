package git

import (
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func CloneOrPull(repoURL, branch, dest string) error {
	if _, err := os.Stat(dest); os.IsNotExist(err) {
		// Clone fresh
		_, err := git.PlainClone(dest, false, &git.CloneOptions{
			URL:           repoURL,
			ReferenceName: plumbing.ReferenceName("refs/heads/" + branch),
			SingleBranch:  true,
		})
		return err
	}
	r, err := git.PlainOpen(dest)
	if err != nil {
		return err
	}
	w, err := r.Worktree()
	if err != nil {
		return err
	}
	return w.Pull(&git.PullOptions{
		RemoteName: "origin",
	})

}
func RepoPath(baseDir, repoID string) string {
	return filepath.Join(baseDir, repoID)
}
