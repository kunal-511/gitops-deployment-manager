package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kunal-511/gitops-deployment-manager/internal/deploy"
	"github.com/kunal-511/gitops-deployment-manager/internal/git"
)

type GitHubWebhookPayload struct {
	Ref string `json:"ref"`
}

func WebhookHandler(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to read request body"})
		return
	}
	var payload GitHubWebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil { // unmarshal JSON payload converts JSON to go data structure
		c.JSON(400, gin.H{"error": "Invalid JSON payload"})
		return
	}

	fmt.Println("Received webhook for ref:", payload.Ref)

	// Extract branch from ref ("refs/heads/master")
	branch := ""
	if len(payload.Ref) > 11 {
		branch = payload.Ref[11:]
	}
	expectedBranch := os.Getenv("GIT_BRANCH")
	if branch != expectedBranch {
		fmt.Printf("Ignoring push to branch: %s\n", branch)
		c.JSON(http.StatusOK, gin.H{"status": "ignored", "branch": branch})
		return
	}

	repoURL := os.Getenv("GIT_REPO_URL")
	username := os.Getenv("GIT_USERNAME")
	token := os.Getenv("GIT_TOKEN")
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
