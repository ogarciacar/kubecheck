package kindcluster

import (
	"bytes"
	"fmt"
	"log"
	"text/template"
	"time"

	"github.com/ogarciacar/kubecheck/sdk"
	"sigs.k8s.io/kind/pkg/cluster"
)

type KindRuntime struct {
	provider *cluster.Provider
}

func New() *KindRuntime {
	return &KindRuntime{
		provider: cluster.NewProvider(cluster.ProviderWithDocker()),
	}
}

func (ke *KindRuntime) Create(name string, apiServerVersion string, kubeconfig string, options ...string) (int32, error) {

	port, err := sdk.GetHostFreePort()

	if err != nil {
		return 0, fmt.Errorf("failed to get a free port: %v", err)
	}

	kindconfig, err := generateKindConfig(port)

	if err != nil {
		return 0, fmt.Errorf("failed to create %s cluster: %v", name, err)
	}

	return int32(port), ke.provider.Create(name,
		cluster.CreateWithRawConfig(kindconfig),
		cluster.CreateWithNodeImage(getNodeImage(apiServerVersion)),
		cluster.CreateWithKubeconfigPath(kubeconfig),
		cluster.CreateWithWaitForReady(5*time.Second))
}

func (ke *KindRuntime) Delete(name string, kubeconfig string) error {
	return ke.provider.Delete(name, kubeconfig)
}

func getNodeImage(apiServerVersion string) string {
	return fmt.Sprintf("kindest/node:v%s", apiServerVersion)
}

// GenerateKindConfig generates a Kind cluster config with a dynamic hostPort
func generateKindConfig(port int) ([]byte, error) {

	log.Println("Generating cluster config file...")

	kindTemplate := `kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  extraPortMappings:
  - containerPort: {{ .HostPort }}
    hostPort: {{ .HostPort }}
    protocol: TCP`

	tmpl, err := template.New("kindConfig").Parse(kindTemplate)
	if err != nil {
		return make([]byte, 0), fmt.Errorf("failed to parse template: %v", err)
	}

	var result bytes.Buffer
	err = template.Must(tmpl, err).Execute(&result, struct {
		HostPort int
	}{HostPort: port})

	if err != nil {
		return make([]byte, 0), fmt.Errorf("failed to execute template: %v", err)
	}

	return result.Bytes(), nil
}
