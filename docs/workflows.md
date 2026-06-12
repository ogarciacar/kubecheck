---
generated_by: structure repo kubecheck
model: claude-opus-4-8
reviewed: false
title: kubecheck — Change, Migration & Release Workflows
area: workflows
audience: [agent, human]
tags: [contributing, release, semantic-release, ci, migrations, versioning]
gold_for_tasks: [T7, T9, T10]
last_verified: 2026-06-11
---
# Workflows

## How to ship a safe change (T7)
1. Branch from `main`. The repo has no `CONTRIBUTING.md`; the de-facto process is observable from CI and git history (PRs merged to `main`, e.g. #41, #42).
2. Write/adjust tests. Tests provision real `kind` clusters, so **Docker must be running locally** (see `operations.md` for the run command).
3. Run the suite locally before pushing: `make test` (`go test -v ./... -cover`).
4. Open a PR targeting `main`. The **Go** workflow (`.github/workflows/go.yml`) runs three gates on every PR: `build` (`go mod tidy` + `go build -v ./...`) and `test` (`gotestsum ... -v ./...` with a test report). Both must pass.
5. **Use Conventional Commits.** Versioning is fully automated by semantic-release (`commit-analyzer`), so commit type drives the release: `fix:` → patch, `feat:` → minor, `BREAKING CHANGE:`/`!` → major. The 1.0.0 release was triggered by a `BREAKING CHANGE` (kluster1 → kubecheck rename).
6. Merge to `main` after review (see `ownership.md` for reviewers).

## How architecture changes are done (T9)
There is **no documented RFC/ADR or design-review process** in the repo. In practice, structural changes land as ordinary reviewed PRs to `main` under the same CI gates, and breaking ones are signalled through Conventional-Commit `BREAKING CHANGE` footers that force a major release and a `CHANGELOG.md` entry (the kluster1→kubecheck rebrand is the canonical example). Treat any need for a formal approval ceremony as currently undocumented.

## "Schema" / version-compatibility migrations (T10)
This module has **no database and no persistent schema**, so there are no schema migrations in the conventional sense. The analogous migration is keeping the supported Kubernetes contract coherent:
- The supported release is a single pinned constant, `K8sRelease_v1_30_10` in `kubecheck/models.go`.
- The runtime derives the node image as `kindest/node:v<release>` (`kindcluster.go`), and `client-go`/`k8s.io/api` are pinned in `go.mod` to the matching minor (v0.30.0).
- To "migrate" to a new Kubernetes version: bump the `K8sVersion` constant, align the `k8s.io/*` module versions in `go.mod` to the same minor (prior commit `refactor(version): match client-go with the kubernetes api`), run `go mod tidy`, and verify with `make test` against the new `kindest/node` image. Ship as a `feat:`/`BREAKING CHANGE:` commit as appropriate.

## Release & deploy procedure
Releases are CI-driven, not manual:
1. On push to `main`, the `semantic-release` job runs (after `build`+`test`) using `.releaserc.json` plugins: `commit-analyzer`, `release-notes-generator`, `changelog`, `github`, `git`. It computes the next version, updates `CHANGELOG.md`, commits, tags, and creates the GitHub release.
2. The pushed tag triggers the `release` job: it imports a GPG key and runs **GoReleaser** (`.goreleaser.yaml`, `builds: skip: true` — no binaries; this is a library, so the "release" is the tagged module + GitHub release).
Do not hand-edit `CHANGELOG.md` or create tags manually; let the pipeline own versioning. Required CI secrets are listed in `operations.md`.