package models

import "time"

type DeploymentRecord struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	RepoID      uint       `json:"repoId"`
	ClusterID   uint       `json:"clusterId"`
	Commit      string     `gorm:"size:80" json:"commit"`
	Status      string     `gorm:"size:32" json:"status"` // success|failed|in_progress
	Message     string     `gorm:"size:500" json:"message"`
	StartedAt   time.Time  `json:"startedAt"`
	CompletedAt *time.Time `json:"completedAt"`
}
