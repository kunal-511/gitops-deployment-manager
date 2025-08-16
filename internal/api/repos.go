package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kunal-511/gitops-deployment-manager/internal/models"
	"github.com/kunal-511/gitops-deployment-manager/internal/storage"
	"gorm.io/gorm"
)

type repoCreateDTO struct {
	Name      string `json:"name" binding:"required"`
	URL       string `json:"url" binding:"required"`
	Branch    string `json:"branch" binding:"required"`
	Path      string `json:"path"`
	ClusterID uint   `json:"clusterId" binding:"required"`
}

func CreateRepo(c *gin.Context) {
	var req repoCreateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify cluster exists
	var cluster models.Cluster
	if err := storage.DB.First(&cluster, req.ClusterID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "cluster not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	repo := models.Repo{
		Name:      req.Name,
		URL:       req.URL,
		Branch:    req.Branch,
		Path:      req.Path,
		ClusterID: req.ClusterID,
	}

	if err := storage.DB.Create(&repo).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": repo.ID, "name": repo.Name, "url": repo.URL})
}

func ListRepos(c *gin.Context) {
	var repos []models.Repo
	if err := storage.DB.Find(&repos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	c.JSON(http.StatusOK, repos)
}

func GetRepo(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid repo id"})
		return
	}

	var repo models.Repo
	if err := storage.DB.First(&repo, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "repo not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	c.JSON(http.StatusOK, repo)
}
