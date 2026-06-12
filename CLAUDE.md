---
generated_by: structure repo kubecheck
model: claude-opus-4-8
reviewed: false
last_source_sync: 44da450ec3ed559acf715c9f9886c50ef066b3cc
---
# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

`kubecheck` (module path `github.com/ogarciacar/kubecheck`) is a Go **library** that makes a real, single-node Kubernetes cluster a first-class primitive inside Go tests. A test calls `kubecheck.NewCluster(...)`, gets a live `kind` cluster plus a standard `client-go` `*kubernetes.Clientset`, exercises real Kubernetes behaviour, and tears it down on `Destroy()`. There is no CLI, server, or deployed service — the consumable surface is a Go API.

Note: the `README.md` still uses the module's former name `kluster1`; the module was rebranded to `kubecheck` (the 1.0.0 `BREAKING CHANGE`). Trust `go.mod` and `docs/` over the README.

## Prerequisites

- Go 1.23+
- A **running Docker daemon** — the `kind` provider provisions clusters as Docker containers, and the test suite creates real clusters. No external cluster or `kubectl` is needed. The first run pulls a `kindest/node` image and may be slow.

## Commands

- `make test` — clears the test cache, then `go test -v ./... -cover`. Provisions and tears down real ephemeral clusters.
- `make coverage` — runs tests writing `coverage.txt`.
- `make covreport` — renders `coverage.txt` to `index.html`.
- Run a single test: `go test -v ./kubecheck/ -run TestName` (clear cache first with `go clean -testcache` if re-running, since cluster tests are not cache-friendly).

## Architecture

A small layered module. The top-level orchestrator coordinates two pluggable collaborators behind **unexported interfaces**, plus a stateless helper package.

- **`kubecheck/`** (package `kubecheck`) — the public orchestrator.
  - `kubecheck.go`: the `K` struct and lifecycle API (`NewCluster`, `Destroy`, `GetClientset`, `GetKubeconfigPath`, `GetIngressPort`). Defines the two collaborator interfaces: `k8sRuntime` (Create/Delete) and `k8sKubeconfig` (CreateTempKubeconfig/DeleteTempKubeconfig).
  - `models.go`: the `K8sVersion` value type and the single pinned release `K8sRelease_v1_30_10` (`1.30.10`).
- **`kubecheck/compute/runtime/kindcluster/`** — concrete `k8sRuntime`. Wraps `sigs.k8s.io/kind`'s Docker provider, renders a single-control-plane Kind config from a template with a dynamic `hostPort`, and selects node image `kindest/node:v<version>`.
- **`kubecheck/storage/persistence/tempkubeconfig/`** — concrete `k8sKubeconfig`. Creates a unique temp dir (`os.MkdirTemp`, `*-kubecheck` pattern), returns the kubeconfig path inside it, removes it on teardown.
- **`sdk/`** — stateless helpers: `GetHostFreePort` (OS-assigned free TCP port) and `GenerateUniqueID` (8-char UUID prefix).

**Cluster-creation flow:** `NewCluster` derives a name `k1-<uniqueID>` → temp kubeconfig dir created → runtime gets a free host port, renders the Kind config mapping that port, creates the cluster → REST config built from the kubeconfig via `clientcmd.BuildConfigFromFlags` → `*kubernetes.Clientset` constructed and returned in `*K`. `Destroy()` deletes the cluster then the temp dir.

**Key invariant — Kubernetes version coherence:** the supported version is one pinned constant, and `client-go`/`k8s.io/api`/`apimachinery` in `go.mod` are pinned to the matching minor (v0.30.0). To bump the supported version: change the `K8sVersion` constant in `models.go`, align all `k8s.io/*` modules to the same minor, `go mod tidy`, and verify with `make test` against the new `kindest/node` image.

## Conventions

- **Conventional Commits are required** — versioning is fully automated by semantic-release. Commit type drives the release: `fix:` → patch, `feat:` → minor, `BREAKING CHANGE:`/`!` → major.
- **Do not hand-edit `CHANGELOG.md` or create git tags manually** — the CI pipeline owns versioning. On push to `main`, semantic-release computes the version, updates the changelog, tags, and creates the GitHub release; GoReleaser runs with `builds: skip: true` (this is a library — no binaries are produced).
- CI (`.github/workflows/go.yml`) gates every PR on `build` (`go mod tidy` + `go build -v ./...`) and `test` (via `gotestsum`); both must pass. There is no `CODEOWNERS`/`CONTRIBUTING.md`, so reviews are not auto-assigned.
- Sharing a cluster across tests: create one in `TestMain`, store it in a package-level `*kubecheck.K`, and `Destroy()` in teardown (see `kubecheck/shared_cluster_test.go`).

## Project documentation (`docs/`)

Richer, task-oriented documentation lives under `docs/` — consult it before deep work:

- `docs/purpose.md` — purpose, scope, and non-goals.
- `docs/architecture.md` — components, data flow, design decisions, dependencies.
- `docs/workflows.md` — shipping changes, version migrations, release/deploy procedure.
- `docs/interfaces.md` — the public Go API surface and usage contracts.
- `docs/operations.md` — local run, configuration/environment, CI runbook.
- `docs/ownership.md` — ownership and escalation (best-effort; inferred).

## Keeping docs fresh

After source changes, run `/update-docs` in Claude Code, or see `.claude/skills/update-docs/SKILL.md`.
