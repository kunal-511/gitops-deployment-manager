package models

import "time"

type Repo struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:120" json:"name"`
	URL       string    `gorm:"size:300" json:"url"`
	Branch    string    `gorm:"size:60" json:"branch"`
	Path      string    `gorm:"size:200" json:"path"` // subdir with manifests
	ClusterID uint      `json:"clusterId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
