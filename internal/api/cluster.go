package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kunal-511/gitops-deployment-manager/internal/k8s"
	"github.com/kunal-511/gitops-deployment-manager/internal/models"
	"github.com/kunal-511/gitops-deployment-manager/internal/storage"
	"gorm.io/gorm"
)

type clusterCreateDTO struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Kubeconfig  string `json:"kubeconfig" binding:"required"`
}

func CreateCluster(c *gin.Context) {
	var req clusterCreateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cluster := models.Cluster{
		Name:        req.Name,
		Description: req.Description,
		Kubeconfig:  []byte(req.Kubeconfig),
	}
	if err := storage.DB.Create(&cluster).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": cluster.ID, "name": cluster.Name})

}

func ListClusters(c *gin.Context) {
	var items []models.Cluster
	if err := storage.DB.Select("id", "name", "description", "created_at", "updated_at").Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, items)
}

// by id
func GetCluster(c *gin.Context) {
	var cluster models.Cluster
	if err := storage.DB.First(&cluster, "id = ?", c.Param("id")).Error; err != nil {
		status := http.StatusNotFound
		if err == gorm.ErrRecordNotFound {
			c.JSON(status, gin.H{"error": "cluster not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":          cluster.ID,
		"name":        cluster.Name,
		"description": cluster.Description,
		"createdAt":   cluster.CreatedAt,
		"updatedAt":   cluster.UpdatedAt,
	})
}

func TestCluster(c *gin.Context) {
	var cluster models.Cluster
	if err := storage.DB.First(&cluster, "id = ?", c.Param("id")).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "cluster not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	cs, _, err := k8s.BuildClientsetFromKubeconfig(cluster.Kubeconfig)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "invalid_kubeconfig", "error": err.Error()})
		return
	}

	if err := k8s.QuickPing(cs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "connection_failed", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
