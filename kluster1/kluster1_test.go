package kluster1_test

import (
	"testing"

	"github.com/ogarciacar/kubecheck/kluster1"
	"github.com/stretchr/testify/require"
)

// TestGetClient tests the functionality of creating a new Kubernetes cluster,
// obtaining a Kubernetes client from the cluster, and verifying the client's
// functionality by checking the server version.
//
// The test performs the following steps:
//  0. Runs the test in parallel.
//  1. Obtains a standard Kubernetes client from the cluster.
//  2. Verifies that the client is not nil.
//  3. Checks that the client can successfully retrieve the server version and
//     compares it to the expected version.
func TestGetClientset(t *testing.T) {

	//  0. Runs the test in parallel.
	t.Parallel()

	//  1. Obtains a standard Kubernetes client from the cluster.
	k8sClient := k1.GetClientset()

	//  2. Verifies that the client is not nil.
	require.NotNil(t, k8sClient, "should not be nil")

	//  3. Checks that the client can successfully retrieve the server version and
	//     compares it to the expected version.
	got, err := k8sClient.ServerVersion()
	require.NoError(t, err, "should not be error")

	// additional checks to ensure the client is functional
	want := kluster1.K8sRelease_v1_30_10.String()
	require.Equal(t, want, got.String(), "should be equal")
}

// TestGetKubeconfig tests the creation of a Kubernetes cluster and ensures that
// the kubeconfig file is generated and exists while the cluster is running.
// It performs the following steps:
// 0. Runs the test in parallel.
// 1. Retrieves the kubeconfig file for the created cluster.
// 2. Checks that the kubeconfig file exists while the cluster is up and running.
func TestGetKubeconfigPath(t *testing.T) {

	//  0. Runs the test in parallel.
	t.Parallel()

	// 1. Retrieves the the path of the kubeconfig file for the created cluster.
	kubeconfig := k1.GetKubeconfigPath()

	// 2. Checks that the kubeconfig file exists while the cluster is up and running.
	require.FileExists(t, kubeconfig, "should exist")
}

// TestGetHostPort tests the GetHostPort function of the kluster1 package.
// It creates a new Kubernetes cluster, retrieves the host port where the cluster ingress is exposed,
// and asserts that the host port is in the range of Dynamic ports.
// It performs the following steps:
// 0. Runs the test in parallel.
// 1. Retrieves the host port.
// 2. Checks that the port is within the Dynamic private ports.
func TestGetIngressPort(t *testing.T) {

	// 0. Runs the test in parallel.
	t.Parallel()

	// 1. Retrieves the host port.
	got := k1.GetIngressPort()

	// assert it is in the range of Dynamic and/or private ports
	// get host port where the cluster ingress is exposed
	wantMinPort := 0
	wantMaxPort := 65535
	require.Greaterf(t, got, wantMinPort, "%T(%v) should be greater than %T(%v)", got, got, wantMinPort, wantMinPort)
	require.LessOrEqualf(t, got, wantMaxPort, "%T(%v) should be less or equal than %T(%v)", got, got, wantMaxPort, wantMaxPort)
}
