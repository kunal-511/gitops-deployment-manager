package sync

import (
	"fmt"
	"time"

	"github.com/kunal-511/gitops-deployment-manager/internal/git"

	"github.com/kunal-511/gitops-deployment-manager/internal/k8s"
)

type RepoSyncJob struct {
	ID          string
	URL         string
	Branch      string
	ManifestDir string
	Kubeconfig  string
}

func StartSyncWorker(jobs []RepoSyncJob, interval time.Duration) {
	ticker := time.NewTicker(interval)
	for {
		<-ticker.C
		for _, job := range jobs {
			fmt.Println("Syncing repo:", job.URL)
			dest := git.RepoPath("/tmp/gitops-repos", job.ID)
			if err := git.CloneOrPull(job.URL, job.Branch, dest); err != nil {
				fmt.Println("Git sync failed:", err)
				continue
			}

			clientset, err := k8s.NewClientFromKubeconfig([]byte(job.Kubeconfig))
			if err != nil {
				fmt.Println("K8s client error:", err)
				continue
			}

			manifestPath := dest + "/" + job.ManifestDir
			if err := k8s.ApplyManifest(clientset, manifestPath); err != nil {
				fmt.Println("Apply failed:", err)
			} else {
				fmt.Println("Apply success for repo:", job.URL)
			}
		}
	}
}
