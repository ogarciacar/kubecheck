package tempkubeconfig_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ogarciacar/kubecheck/kubecheck/storage/persistence/tempkubeconfig"
	"github.com/stretchr/testify/require"
)

func TestCreateTempKubeconfig(t *testing.T) {

	require := require.New(t)

	tk := tempkubeconfig.New()

	tmpKubeconfigPath, err := tk.CreateTempKubeconfig()

	require.NoErrorf(err, "expected no error, got %v", err)

	require.NotNil(tmpKubeconfigPath, "expected a valid path, got nil")

	expectedDir := filepath.Dir(*tmpKubeconfigPath)

	require.Containsf(expectedDir, os.TempDir(), "expected temp directory to be within os.TempDir(), got %v", expectedDir)

	// Clean up
	err = os.RemoveAll(expectedDir)
	require.NoError(err)
}
