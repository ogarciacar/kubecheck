package kubecheck

import "fmt"

type K8sVersion struct {
	release string
}

func (k K8sVersion) String() string {
	return fmt.Sprintf("v%s", k.release)
}

var K8sRelease_v1_30_10 = K8sVersion{
	release: "1.30.10",
}
