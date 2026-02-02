package kluster1_test

import (
	"testing"

	"github.com/ogarciacar/kubecheck/kluster1"
)

func TestMyKubernetesApp(t *testing.T) {

	t.Parallel()

	// 1. Creates a new Kubernetes cluster with a specified release version.
	k1, err := kluster1.NewCluster(kluster1.K8sRelease_v1_30_10)

	// 2. Check the cluster was successfully created
	if err != nil {
		t.Fatalf("should no be error %s", err)
	}

	// 3. Ensures the cluster is cleaned up after the test completes.
	defer k1.Destroy()

	// ... Your Kubernetes test logic here, using k1.GetClientset()
	clientset := k1.GetClientset()

	serverVersion, err := clientset.ServerVersion()
	if err != nil {
		t.Fatalf("should no be error %s", err)
	}

	expected := kluster1.K8sRelease_v1_30_10.String()
	got := serverVersion.String()

	if expected != got {
		t.Fatalf("should be the same server version: expected %s got, %s", expected, got)
	}
}
