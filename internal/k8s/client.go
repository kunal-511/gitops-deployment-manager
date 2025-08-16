package k8s

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func decodeKubeconfig(kubeconfig []byte) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(string(kubeconfig))
	if err != nil {
		decoded = kubeconfig
	}
	return decoded, nil
}

// BuildClientsetFromKubeconfig creates a clientset using kubeconfig bytes
func BuildClientsetFromKubeconfig(kubeconfig []byte) (*kubernetes.Clientset, *rest.Config, error) {
	decoded, err := decodeKubeconfig(kubeconfig)

	loader, err := clientcmd.NewClientConfigFromBytes(decoded)
	if err != nil {
		return nil, nil, err
	}
	cfg, err := loader.ClientConfig()
	if err != nil {
		return nil, nil, err
	}
	cs, err := kubernetes.NewForConfig(cfg)
	return cs, cfg, err
	// clientset is used to interact with  Kubernetes resources like Pods, Deployments, Namespaces, etc.
	//rest config contains the configuration for connecting to the Kubernetes API server.
}

// QuickPing checks we can reach the API server by listing namespaces.
func QuickPing(cs *kubernetes.Clientset) error {
	_, err := cs.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	// Create namespace 'deploy'
	deployNS := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "deploy",
		},
	}
	_, err = cs.CoreV1().Namespaces().Create(context.Background(), deployNS, metav1.CreateOptions{})
	// If the namespace already exists, ignore the error
	if err != nil && !isAlreadyExistsError(err) {
		return err
	}
	return nil
}

func isAlreadyExistsError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "already exists")
}

// Fingerprint returns a short fingerprint for kubeconfigs (for logs only)
func Fingerprint(kubeconfig []byte) string { // just a short hash for logging
	if len(kubeconfig) == 0 {
		return "empty"
	}
	sum := sha1.Sum(kubeconfig)
	return hex.EncodeToString(sum[:])[:10]
}

func NewClientFromKubeconfig(kubeconfig []byte) (*kubernetes.Clientset, error) {
	decoded, err := decodeKubeconfig(kubeconfig)
	config, err := clientcmd.RESTConfigFromKubeConfig(decoded)
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}
func ApplyManifest(clientset *kubernetes.Clientset, manifestPath string) error {
	// For now, we’ll just print — later integrate `sigs.k8s.io/controller-runtime/pkg/client`
	fmt.Println("Would apply manifest:", manifestPath)
	return nil
}
