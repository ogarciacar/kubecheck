package kluster1_test

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/ogarciacar/kubecheck/sdk"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// TestCreatePod tests the functionality of creating a new Pod in a specified namespace
// using the provided pod configuration.
//
// The test performs the following steps:
// 0. Runs the test in parallel.
// 1. Defines a Pod configuration.
// 2. Creates the Pod in the specified namespace.
// 3. Verifies that the created Pod is not nil and matches the expected configuration.
// 4. Verifies that the created Pod is exposed.
func TestCreatePod(t *testing.T) {

	// 0. Runs the test in parallel.
	t.Parallel()

	safeTest(t, func(t *testing.T) {
		// 1. Defines a Pod configuration.
		uniqueID := sdk.GenerateUniqueID()
		namespace := fmt.Sprintf("test-namespace-%s", uniqueID)
		var testPod = createTestPodDefinition("nginx", namespace, uniqueID, k1.GetIngressPort())

		// 2. Creates the Pod in the specified namespace.
		createdPod, err := createPod(k1.GetClientset(), namespace, testPod)

		// 3. Verifies that the created Pod is not nil and matches the expected configuration.
		require.NoError(t, err, "should not be error when creating pod")
		require.NotNil(t, createdPod, "created pod should not be nil")
		require.Equal(t, testPod.Name, createdPod.Name, "pod name should match")
		require.Equal(t, testPod.Spec.Containers[0].Name, createdPod.Spec.Containers[0].Name, "container name should match")
		require.Equal(t, testPod.Spec.Containers[0].Image, createdPod.Spec.Containers[0].Image, "container image should match")
		require.Equal(t, k1.GetIngressPort(), int(testPod.Spec.Containers[0].Ports[0].HostPort), "should be equal")

		resp, err := retryHttpGet(fmt.Sprintf("http://localhost:%d/", k1.GetIngressPort()), 16, 5)
		require.NoError(t, err, "should not be error")
		require.Equal(t, http.StatusOK, resp.StatusCode, "should be 200 OK")
	})
}

func createTestPodDefinition(image string, namespace string, uniqueID string, ingressPort int) *corev1.Pod {

	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("test-pod-%s", uniqueID),
			Namespace: namespace,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  fmt.Sprintf("test-container-%s", uniqueID),
					Image: image,
					Ports: []corev1.ContainerPort{
						{
							ContainerPort: 80,
							HostPort:      int32(ingressPort),
						},
					},
				},
			},
		},
	}
}

// CreatePod creates a new Pod in the specified namespace using the provided pod configuration.
//
// Parameters:
//   - k8sClient: A Kubernetes clientset to interact with the Kubernetes API.
//   - namespace: The namespace in which to create the pod.
//   - pod: The pod object to be created.
//
// Returns:
//   - A pointer to the created pod object.
//   - An error if there was an issue creating the namespace, service account, or pod.
func createPod(k8sClient *kubernetes.Clientset, namespace string, pod *corev1.Pod) (*corev1.Pod, error) {

	// Create the namespace
	_, err := k8sClient.CoreV1().Namespaces().Create(context.TODO(), &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}, metav1.CreateOptions{})

	if err != nil {
		return nil, fmt.Errorf("error creating namespace %s: %v", namespace, err)
	}

	// Create the default service account
	_, err = k8sClient.CoreV1().ServiceAccounts(namespace).Create(context.TODO(), &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name: "default",
		},
	}, metav1.CreateOptions{})

	if err != nil {
		return nil, fmt.Errorf("error creating service account %s: %v", namespace, err)
	}

	// Create the pod in the namespace
	return k8sClient.CoreV1().Pods(namespace).Create(context.TODO(), pod, metav1.CreateOptions{})
}

// up to N http GET request every N seconds
func retryHttpGet(url string, times int, waitTime time.Duration) (*http.Response, error) {

	var resp *http.Response

	var err error

	for range times {
		resp, err = http.Get(url)
		if err == nil {
			log.Printf("GET %s %v", url, resp.Body)
			break
		} else {
			log.Printf("GET %s %v", url, err)
		}
		time.Sleep(waitTime * time.Second)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}
