package api

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kunal-511/gitops-deployment-manager/internal/deploy"
	"github.com/kunal-511/gitops-deployment-manager/internal/git"
)

func SyncHandler(c *gin.Context) {
	repoURL := os.Getenv("GIT_REPO_URL")
	branch := os.Getenv("GIT_BRANCH")
	username := os.Getenv("GIT_USERNAME") // optional
	token := os.Getenv("GIT_TOKEN")       // optional
	repoPath := "./manifests"

	if err := git.PullLatest(repoPath, repoURL, branch, username, token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := deploy.ApplyManifests(repoPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "synced", "branch": branch})

}
