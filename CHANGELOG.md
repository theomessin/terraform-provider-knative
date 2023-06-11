# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Ability to specify Kubernetes namespace for Service Data Source. This is optional and defaults to `default`.

## [0.1.1] - 2023-06-11

### Fixed

- Remove non-existent `kube_config_path` attribute from Provider docs.

## [0.1.0] - 2023-06-11

### Added

- Knative Service Data Source that can be queried by name only (namespace hardcoded to default).
  Partly implemented reading the status, although some fields are still missing.

[unreleased]: https://github.com/theomessin/terraform-provider-knative/compare/v0.1.1...HEAD
[0.1.1]: https://github.com/theomessin/terraform-provider-knative/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/theomessin/terraform-provider-knative/releases/tag/v0.1.0