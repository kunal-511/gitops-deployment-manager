package api

import "github.com/gin-gonic/gin"

func NewRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/healthz", HealthHandler)

	api := r.Group("/api")
	{
		api.POST("/clusters", CreateCluster)
		api.GET("/clusters", ListClusters)
		api.GET("/clusters/:id", GetCluster)
		api.POST("/clusters/:id/test", TestCluster)
		// repos, sync, webhooks will be added next
	}

	return r
}
