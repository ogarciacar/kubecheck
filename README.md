# `kluster1`: Kubernetes TDD, Simplified.

[![Go Version](https://img.shields.io/badge/Go-1.23+-blue)](https://golang.org/) [![Go Report Card](https://goreportcard.com/badge/github.com/ogarciacar/kluster1)](https://goreportcard.com/report/github.com/ogarciacar/kluster1)

**Stop managing Kubernetes infrastructure, start testing your applications.**

`kluster1` is a Go module that empowers Test-Driven Development (TDD) for Kubernetes applications. It allows you to integrate local Kubernetes cluster lifecycle management directly into your Go test code, eliminating the need for external Kubernetes tools and ensuring consistent, fast, and reliable testing environments.

## Strategy Statement

**What change is being pursued?**
Shift Kubernetes from an external, shared testing dependency to an embedded test primitive that runs the same way locally and inside CI pipelines.

**What is the promise?**
Production-like feedback during development and CI without the coordination overhead, cost, and fragility of managing shared Kubernetes test clusters.

**Who is being changed?**
Platform and DevEx engineers who own CI pipelines and believe real system behavior should be tested, not assumed.

