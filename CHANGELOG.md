# Changelog

All notable changes to this project will be documented in this file.

The format is based on [keep a changelog](https://keepachangelog.com/en/1.0.0/) and this project adheres to [semantic versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Fixed

- Do not output `ProgressMessage` to non-interactive terminals. This avoids including a progress update in the started message that is printed to non-interactive terminals.

## [v1.0.2] - 2022-11-23

### Fixed

- Normalise whitespace (`\s` → ` `) in messages to avoid newlines and tabs breaking in-progress message updating.
- Assume details message to be preformatted, if it contains newline characters (`\n`). Preformatted message details are wrapped so that newline characters are maintained. This makes, for example, stack traces and console output in message details more readable.

## [v1.0.1] - 2022-08-30

### Fixed

- Do not try to render message if terminal width is zero. This happens with some terminals on first terminal width get(s).

## [v1.0.0] - 2022-08-26

### Added

- Extract and refactor livelog functionality from [UpCloud CLI (`upctl`)](https://github.com/UpCloudLtd/upcloud-cli.git).

[Unreleased]: https://github.com/UpCloudLtd/progress/compare/v1.0.2...HEAD
[v1.0.2]: https://github.com/UpCloudLtd/progress/compare/v1.0.1...v1.0.2
[v1.0.1]: https://github.com/UpCloudLtd/progress/compare/v1.0.0...v1.0.1
[v1.0.0]: https://github.com/UpCloudLtd/progress/releases/tag/v1.0.0
