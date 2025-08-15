package models

import "time"

type Cluster struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"uniqueIndex;size:100" json:"name"`
	Description string `gorm:"size:255" json:"description"`
	// Kubeconfig in plain text for MVP — replace with encrypted at rest in prod
	Kubeconfig []byte    `json:"-" gorm:"type:bytea"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}
