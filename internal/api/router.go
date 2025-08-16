package api

import "github.com/gin-gonic/gin"

func NewRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/healthz", HealthHandler)

	api := r.Group("/api")
	{
		// Cluster routes
		api.POST("/clusters", CreateCluster)
		api.GET("/clusters", ListClusters)
		api.GET("/clusters/:id", GetCluster)
		api.POST("/clusters/:id/test", TestCluster)

		// Repository routes
		api.POST("/repos", CreateRepo)
		api.GET("/repos", ListRepos)
		api.GET("/repos/:id", GetRepo)

	}

	return r
}
