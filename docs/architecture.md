---
generated_by: structure repo kubecheck
model: claude-opus-4-8
reviewed: false
title: kubecheck — Architecture & Dependencies
area: architecture
audience: [agent, human]
tags: [architecture, components, data-flow, dependencies, kind]
gold_for_tasks: [T3, T6]
last_verified: 2026-06-11
---
# Architecture

`kubecheck` is a small layered Go module. The top-level orchestrator coordinates two pluggable collaborators (a cluster runtime and a kubeconfig manager) behind interfaces, plus a stateless `sdk` of helpers.

## Components (all paths verified against the tree)
- **`kubecheck/` (package `kubecheck`)** — public orchestrator.
  - `kubecheck.go`: the `K` struct and lifecycle API (`NewCluster`, `Destroy`, `GetClientset`, `GetKubeconfigPath`, `GetIngressPort`). Defines two collaborator interfaces — `k8sRuntime` (Create/Delete) and `k8sKubeconfig` (CreateTempKubeconfig/DeleteTempKubeconfig).
  - `models.go`: `K8sVersion` value type and the pinned release `K8sRelease_v1_30_10` (`1.30.10`).
- **`kubecheck/compute/runtime/kindcluster/` (package `kindcluster`)** — the concrete `k8sRuntime`. Wraps `sigs.k8s.io/kind`'s `cluster.Provider` (Docker provider), renders a single-control-plane Kind config from a template with a dynamic `hostPort`, and selects the node image `kindest/node:v<version>`.
- **`kubecheck/storage/persistence/tempkubeconfig/` (package `tempkubeconfig`)** — the concrete `k8sKubeconfig`. Creates a unique temp directory (`os.MkdirTemp` with a `*-kubecheck` pattern), returns the `config` path inside it, and removes the directory on teardown.
- **`sdk/` (package `sdk`)** — stateless helpers: `GetHostFreePort` (binds `tcp :0`, lets the OS assign a port) and `GenerateUniqueID` (first 8 chars of a UUID).

## Data flow (cluster creation)
1. `NewCluster(version)` builds a `kindcluster` runtime and a `tempkubeconfig` manager, and derives a cluster name `k1-<uniqueID>`.
2. The kubeconfig manager creates a temp dir and returns the kubeconfig path.
3. The runtime asks `sdk.GetHostFreePort()` for a port, renders the Kind config mapping that `hostPort`, and creates the cluster with the chosen node image and kubeconfig path.
4. `NewCluster` builds REST config from the kubeconfig via `clientcmd.BuildConfigFromFlags` and constructs a `*kubernetes.Clientset`.
5. The populated `K` (kubeconfig, name, runtime, clientset, ingress port, kubeconfig manager) is returned. `Destroy()` deletes the cluster, then the temp kubeconfig dir.

## Key design decisions
- **Interface-seamed collaborators** (`k8sRuntime`, `k8sKubeconfig`) keep the runtime and persistence pluggable and testable.
- **Pinned Kubernetes version** as a typed value (`K8sVersion`) keeps client-go and node image in lockstep.
- **Ephemeral, single-node** topology by construction.

## Upstream dependencies (from `go.mod`, Go 1.23)
- `sigs.k8s.io/kind` v0.27.0 — cluster provisioning (requires **Docker** at runtime; pulls `kindest/node` images).
- `k8s.io/client-go`, `k8s.io/api`, `k8s.io/apimachinery` v0.30.0 — Kubernetes client and types.
- `github.com/google/uuid` v1.6.0 — unique IDs.
- `github.com/stretchr/testify` v1.10.0 — test assertions.
Numerous transitive `// indirect` deps support the above; see `go.mod`/`go.sum` for the full set.

Procedures for changing this design live in `workflows.md`; the consumer-facing surface is in `interfaces.md`.