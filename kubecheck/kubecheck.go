package kubecheck

import (
	"fmt"
	"log"

	"github.com/ogarciacar/kubecheck/kubecheck/compute/runtime/kindcluster"
	"github.com/ogarciacar/kubecheck/kubecheck/storage/persistence/tempkubeconfig"
	"github.com/ogarciacar/kubecheck/sdk"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type K struct {
	kubeconfig               string
	name                     string
	k8sRuntime               k8sRuntime
	k8sClient                *kubernetes.Clientset
	k8sIngressPort           int32
	k8sKubeconfigPathManager k8sKubeconfig
}

type k8sRuntime interface {
	Create(name string, apiServerVersion string, kubeconfig string, options ...string) (int32, error)
	Delete(name string, kubeconfig string) error
}

type k8sKubeconfig interface {
	CreateTempKubeconfig() (*string, error)
	DeleteTempKubeconfig() error
}

// NewCluster creates a new Kubernetes cluster using the specified Kubernetes version and a temporary directory for the kubeconfig file.
// It returns a pointer to a kluster1 instance and an error if the cluster creation fails.
//
// Parameters:
//   - k8sVersion: The version of Kubernetes to use for the cluster.
//   - tempDir: The directory where the temporary kubeconfig file will be stored.
//
// Returns:
//   - *kluster1: A pointer to the created kluster1 instance.
//   - error: An error if the cluster creation fails.
func NewCluster(k8sVersion K8sVersion) (*K, error) {

	// get cluster provider
	k8sRuntime := kindcluster.New()

	name := fmt.Sprintf("k1-%s", sdk.GenerateUniqueID())

	// dev feedback
	log.Println("Creating single-node Kubernetes cluster...")
	log.Printf("Kubernetes API Server %s", k8sVersion.release)
	log.Printf("Cluster name %s-control-plane", name)

	// temp kubeconfigPath file
	k8sKubeconfigPathManager := tempkubeconfig.New()
	kubeconfigPath, err := k8sKubeconfigPathManager.CreateTempKubeconfig()
	if err != nil {
		return nil, err
	}

	log.Printf("KUBECONFIG=%s", *kubeconfigPath)

	// create cluster
	k8sIngressPort, err := k8sRuntime.Create(name, k8sVersion.release, *kubeconfigPath)
	if err != nil {
		return nil, err
	}

	log.Printf("Kubernetes Ingress Port %d", k8sIngressPort)

	// create a standard kubernetes client
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfigPath)
	if err != nil {
		return nil, err
	}

	k8sClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &K{
		kubeconfig:               *kubeconfigPath,
		name:                     name,
		k8sRuntime:               k8sRuntime,
		k8sClient:                k8sClient,
		k8sIngressPort:           k8sIngressPort,
		k8sKubeconfigPathManager: k8sKubeconfigPathManager,
	}, nil
}

// Destroy deletes the cluster associated with the kluster1 instance.
// It logs the deletion process and calls the engine's Delete method
// with the cluster's name and kubeconfig.
//
// Returns an error if the deletion process fails.
func (k *K) Destroy() error {

	log.Printf("Deleting cluster %s-control-plane", k.name)

	err := k.k8sRuntime.Delete(k.name, k.kubeconfig)

	if err != nil {
		return err
	}

	k.k8sKubeconfigPathManager.DeleteTempKubeconfig()

	return nil
}

// GetClientset returns the Kubernetes clientset associated with the kluster1 instance.
// This clientset can be used to interact with the Kubernetes API server.
func (k *K) GetClientset() *kubernetes.Clientset {
	return k.k8sClient
}

// GetKubeconfigPath returns the kubeconfig string associated with the kluster1 instance.
func (k *K) GetKubeconfigPath() string {
	return k.kubeconfig
}

func (k *K) GetIngressPort() int {
	return int(k.k8sIngressPort)
}
