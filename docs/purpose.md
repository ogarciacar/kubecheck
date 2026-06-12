---
generated_by: structure repo kubecheck
model: claude-opus-4-8
reviewed: false
title: kubecheck — Purpose & Scope
area: purpose
audience: [agent, human]
tags: [purpose, scope, non-goals, tdd, kubernetes]
gold_for_tasks: [T1, T2]
last_verified: 2026-06-11
---
# Purpose

`kubecheck` is a Go module (`github.com/ogarciacar/kubecheck`) that makes a real, single-node Kubernetes cluster a first-class primitive inside Go tests. Instead of pointing tests at an external, shared cluster, a test calls `kubecheck.NewCluster(...)`, receives a live cluster plus a `client-go` clientset, exercises real Kubernetes behaviour, and tears the cluster down when the test finishes.

## What it does
- Spins up an ephemeral, single-node cluster per call using `kind` (Kubernetes IN Docker), pinned to a known Kubernetes release.
- Hands the test a standard `*kubernetes.Clientset` so existing `client-go` code works unchanged.
- Allocates an OS-assigned free host port and maps it into the node so workloads under test are reachable from the host.
- Manages a throwaway kubeconfig in a temp directory and cleans up both cluster and kubeconfig on `Destroy()`.

## Strategy
Shift Kubernetes from an external, shared testing dependency to an embedded test primitive that runs identically on a laptop and in CI. The promise is production-like feedback during development and CI without the coordination overhead, cost, and fragility of shared test clusters. The intended users are platform and DevEx engineers who own CI pipelines and want real system behaviour tested rather than mocked.

## Scope
- In scope: programmatic cluster lifecycle (create/destroy) from Go test code; exposing a clientset, kubeconfig path, and ingress port; small reusable test helpers (free-port discovery, unique-id generation).
- Supported runtime today: a single pinned Kubernetes release (`K8sRelease_v1_30_10`) on a single-node `kind` cluster backed by Docker.

## Non-goals
- **Not** a production or long-lived cluster manager — clusters are ephemeral and single-node.
- **Not** a manager of shared/external test clusters; eliminating that coordination is the whole point.
- **Not** a replacement for `kubectl`, Helm, or a CI cluster-provisioning system in production environments.
- **Not** a multi-node, HA, or cloud-provider cluster tool — topology is fixed to one control-plane node.
- **Not** a general CLI; the surface is a Go API consumed from tests (see `interfaces.md`).

For the component breakdown that delivers this, see `architecture.md`; for the public API, see `interfaces.md`.