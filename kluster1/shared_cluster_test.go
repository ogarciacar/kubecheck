package kluster1_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/ogarciacar/kubecheck/kluster1"
)

// Represents the Kubernetes cluster
var k1 *kluster1.K

// One TestMain function per package
// The TestMain function performs the following steps:
// 1. Creates a new Kubernetes cluster with a specified release version.
// 2. Ensure teardown runs even if m.Run() panics
// 3. Runs the tests in in kluster1_test package.
// 4. Ensures the cluster is cleaned up after the test completes.
func TestMain(m *testing.M) {

	var err error

	// 1. Creates a new Kubernetes cluster with a specified release version.
	k1, err = kluster1.NewCluster(kluster1.K8sRelease_v1_30_10)

	if err != nil {
		os.Exit(1)
	}

	// 2. Ensure teardown runs even if m.Run() panics
	exitCode := 0
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
			exitCode = 1 // Indicate failure due to panic
		}
		// 4. Ensures the cluster is cleaned up after the test completes.
		k1.Destroy()
		os.Exit(exitCode)
	}()

	// 3. Runs the tests in in kluster1_test package.
	exitCode = m.Run()
}

// Wrapper to catch panics inside test functions
func safeTest(t *testing.T, testFunc func(t *testing.T)) {

	// Marks safeTest function as a test helper function
	t.Helper()

	// Catches panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Test panicked: %v", r)
		}
	}()

	// Executes test function
	testFunc(t)
}
