# Knative Serving

This repository contains an open source specification and implementation of a Kubernetes- and Istio-based container platform.

If you are interested in contributing to `Knative Serving`, see
[CONTRIBUTING.md](./CONTRIBUTING.md) and [DEVELOPMENT.md](./DEVELOPMENT.md).

## Getting Started

* [Setup Istio](https://istio.io/docs/setup/kubernetes/quick-start.html): Make sure to enable automatic sidecar injection for the default namespace (or any other namespace containing Knative Serving services).
* [Setup Knative Serving](#latest-release): See `Latest Release` below.
* [Run samples](./sample/README.md)

## Configuring Knative Serving
See [CONFIG.md](./install/CONFIG.md).

### Latest Release

You can install the latest release of Knative Serving via:

```shell
kubectl apply -f https://storage.googleapis.com/knative-releases/latest/release.yaml
```
