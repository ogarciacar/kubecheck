---
generated_by: structure repo kubecheck
model: claude-opus-4-8
reviewed: false
title: kubecheck — Operations, Config & Local Run
area: operations
audience: [agent, human]
tags: [operations, config, env, local-dev, runbook, ci-secrets]
gold_for_tasks: [T4, T8, T11]
last_verified: 2026-06-11
---
# Operations

This is a test-support library, not a deployed service. "Operations" here means running it locally, its configuration surface, and the CI runbook.

## How to run it locally (T4)
**Prerequisites:** Go 1.23+ and a running **Docker** daemon (the `kind` provider uses Docker). No external cluster or `kubectl` needed.

Common commands (`Makefile`):
- `make test` — `go clean -testcache` then `go test -v ./... -cover`. Provisions real ephemeral clusters; the suite manages create/teardown itself.
- `make coverage` — runs tests writing `coverage.txt`.
- `make covreport` — renders `coverage.txt` to `index.html`.

To use it from your own project: `go get github.com/ogarciacar/kubecheck`, then call `kubecheck.NewCluster(...)` from a test (see `interfaces.md`). Each `NewCluster` call pulls/uses a `kindest/node` image and starts a Docker-backed node, so the first run may be slow and a working Docker context is mandatory.

## Configuration & environment (T8)
There is no config file. Behaviour is determined by code constants and a few runtime touchpoints:
- **Kubernetes version** — set in code via `K8sRelease_v1_30_10` (`kubecheck/models.go`); not an env var.
- **`KUBECONFIG`** — `kubecheck` does **not** read a `KUBECONFIG` from the environment; it generates its own temp kubeconfig under the OS temp dir (`os.MkdirTemp`, `*-kubecheck` pattern) and logs the path. Honour the OS temp-dir conventions (e.g. `TMPDIR`) if you need to relocate it.
- **Host port** — assigned dynamically by the OS at runtime; retrieve with `GetIngressPort()`. Nothing to configure.
- **CI secrets (config for the release pipeline only, names only):** `GITHUB_TOKEN`, `GPG_PRIVATE_KEY`, `PASSPHRASE`, consumed by `.github/workflows/go.yml`. Values are provided by GitHub Actions secrets; the workflow references them by name and `persist-credentials: false` is used on checkout. Do not inline secret values anywhere.

## Runbook / on-call (T11)
There is **no on-call rotation, monitoring, alerting, or production runbook** for this module, and none is documented in the repo — it ships no running service. The operational failure surface is local/CI test execution:
- *Cluster create fails / hangs* → confirm Docker is running and the `kindest/node:v<version>` image is pullable; re-run `make test`.
- *Port/HTTP checks fail* → the suite retries HTTP GETs; ensure no host firewall blocks the mapped port and that leftover containers from a crashed run are pruned (`Destroy()` normally cleans up).
- *CI red* → inspect the `build`/`test` jobs in the Go workflow; release issues are governed by the pipeline in `workflows.md`.

Escalation for any of the above is via the maintainer in `ownership.md`.