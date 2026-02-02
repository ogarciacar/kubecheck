package tempkubeconfig

import (
	"log"
	"os"
	"path/filepath"
)

type TempKubeconfig struct {
	tempKubeconfigDir string
}

func New() *TempKubeconfig {
	return &TempKubeconfig{}
}

// Creates a temporary directory to persist the kubeconfig file
func (tk *TempKubeconfig) CreateTempKubeconfig() (*string, error) {

	log.Println("Creating temporary directory to persist the Kubeconfig file")

	uniqueTempDir, err := os.MkdirTemp(os.TempDir(), "*-kluster1")

	if err != nil {
		return nil, err
	}

	tempKubeconfigPath := filepath.Join(uniqueTempDir, "config")

	tk.tempKubeconfigDir = uniqueTempDir

	return &tempKubeconfigPath, nil
}

// Deletes the temporary directory that contains the kubeconfig file
func (tk *TempKubeconfig) DeleteTempKubeconfig() error {
	log.Printf("Deleting temporary directory %s", tk.tempKubeconfigDir)
	return os.RemoveAll(tk.tempKubeconfigDir)
}
