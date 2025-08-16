package main

import (
	"log"

	"github.com/kunal-511/gitops-deployment-manager/internal/api"
	"github.com/kunal-511/gitops-deployment-manager/internal/config"
	"github.com/kunal-511/gitops-deployment-manager/internal/storage"
)

func main() {
	cfg := config.Load()

	_, err := storage.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	// Auto-migrate initial tables
	// if err := db.AutoMigrate(
	// 	&models.Cluster{},
	// 	&models.Repo{},
	// 	&models.DeploymentRecord{},
	// ); err != nil {
	// 	log.Fatalf("migrate: %v", err)
	// }

	r := api.NewRouter()
	log.Printf("listening on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
