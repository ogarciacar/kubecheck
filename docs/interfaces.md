---
generated_by: structure repo kubecheck
model: claude-opus-4-8
reviewed: false
title: kubecheck — Public Interfaces & Contracts
area: interfaces
audience: [agent, human]
tags: [api, public-api, entry-points, contracts, go-module]
gold_for_tasks: [T5]
last_verified: 2026-06-11
---
# Interfaces

`kubecheck` exposes a Go library API only — there is no CLI, HTTP server, or event bus. Consumers import the packages below. Module path: `github.com/ogarciacar/kubecheck`.

## Primary entry point — package `kubecheck` (`kubecheck/kubecheck.go`, `models.go`)

**Construction**
- `func NewCluster(k8sVersion K8sVersion) (*K, error)` — provisions a single-node cluster, builds a clientset, returns the handle `*K`. Logs progress (cluster name, kubeconfig path, ingress port). Requires Docker. Errors propagate from temp-dir creation, cluster create, or client build.

**Type `K` methods (the test-facing contract)**
- `func (k *K) GetClientset() *kubernetes.Clientset` — standard `client-go` clientset for the cluster.
- `func (k *K) GetKubeconfigPath() string` — path to the cluster's temp kubeconfig.
- `func (k *K) GetIngressPort() int` — host port mapped into the node for reaching workloads (e.g. `http://localhost:<port>/`).
- `func (k *K) Destroy() error` — deletes the cluster and removes the temp kubeconfig directory. Intended for `defer`/teardown.

**Version values (`models.go`)**
- `type K8sVersion struct { ... }` with `func (k K8sVersion) String() string` → `v<release>` (e.g. `v1.30.10`).
- `var K8sRelease_v1_30_10 K8sVersion` — the currently supported, pinned release passed to `NewCluster`.

**Internal contracts (not for external implementation, but part of the design):** `k8sRuntime` (Create/Delete) and `k8sKubeconfig` (CreateTempKubeconfig/DeleteTempKubeconfig) are unexported interfaces in `kubecheck.go`; their implementations are described in `architecture.md`.

## Helper package `sdk` (`sdk/free_port.go`, `sdk/unique_id.go`)
- `func GetHostFreePort() (int, error)` — returns an OS-assigned free TCP port.
- `func GenerateUniqueID() string` — 8-character unique id (UUID prefix).

## Usage contract (from `kubecheck/quick_start_test.go`)
```go
k1, err := kubecheck.NewCluster(kubecheck.K8sRelease_v1_30_10)
if err != nil { t.Fatalf("...") }
defer k1.Destroy()
clientset := k1.GetClientset()
// ... exercise the real cluster via client-go ...
```
A common pattern (see `kubecheck/shared_cluster_test.go`) is to create one cluster in `TestMain`, share it across parallel tests via a package-level `*kubecheck.K`, and `Destroy()` in teardown.

Packages under `compute/runtime/...` and `storage/persistence/...` are implementation detail and not part of the public contract.