# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] - 2026-05-28

### Added
- Initial release of `runpod-top`, a terminal UI for monitoring RunPod GPU pods.
- Split-pane dashboard: pod list on the left, detailed metrics on the right.
- Live polling of the RunPod GraphQL API with configurable interval.
- Multi-select to compare metrics across pods.
- Start, stop, and restart pods from the TUI with confirmation dialogs.
- TOML config file at `~/.config/runpod-top/config.toml`, plus `--init-config` to generate a default.
- Precedence: CLI flags > environment variables > config file > defaults.
- RAM displayed in the pod detail header.
- `-version` flag prints the build version, commit, and build date.

[Unreleased]: https://github.com/Wyoc/runpod-top/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/Wyoc/runpod-top/releases/tag/v0.1.0
