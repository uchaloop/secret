# Changelog

All notable changes to this module are documented in this file.

Entries are grouped by version and change type, with the newest version first.

## [2.2.0] - Unreleased

### Added

- Added the `Secret.IsSensitive()` marker so integrations can recognize sensitive types
  through a structural interface without importing this module. The existing `Value`
  interface and all masking and decoding behavior remain unchanged.

## [2.1.0]

The `/v2` import path was retained despite the breaking changes in this release.
The release date is not recorded in the existing changelog.

### Added

- Added Codecov and contribution/security guidelines; updated documentation and tests.

### Changed

- **Breaking:** Made `Secret` non-comparable and unusable as a map key; use `IsZero()`
  to check emptiness. JSON `null` now clears the secret.
- Lowered the minimum Go version to 1.21.0.

### Fixed

- Rejected malformed JSON and Unicode with safe errors, preserving the previous value.
  Nil decode receivers return errors.
- Fixed the `/v2` module path in the release workflow.

### Security

- Fixed secret disclosure through `fmt` and `slog` in nested and unexported fields.

## [2.0.2] - 2026-09-01

### Changed

- Clarified masking guarantees and limitations in GoDoc; simplified the README.
- Standardized changelog version headings.

### Removed

- Removed the obsolete `koanf` tag from the README example.

## [2.0.1] - 2026-08-06

### Changed

- Simplified the README and usage documentation.

## [2.0.0] - 2026-08-04

### Added

- Added `New`, `IsZero`, `Clear` and the `Value` sensitivity marker.

### Changed

- **Breaking:** Replaced the string-based `Secret` with an opaque struct.
- **Breaking:** Changed the module path to `github.com/uchaloop/secret/v2`.
- Expanded masking for `fmt`, text/JSON/XML serialization and `slog`.

## [1.0.0] - 2026-08-04

### Added

- Introduced a string-based `Secret` with masked formatting, text/JSON output and `slog`
  support; empty secrets produce empty output.
- Added `Reveal` for explicit access and `UnmarshalText` for input.

[2.2.0]: https://github.com/uchaloop/secret/compare/v2.1.0...HEAD
[2.1.0]: https://github.com/uchaloop/secret/compare/v2.0.2...v2.1.0
[2.0.2]: https://github.com/uchaloop/secret/compare/v2.0.1...v2.0.2
[2.0.1]: https://github.com/uchaloop/secret/compare/v2.0.0...v2.0.1
[2.0.0]: https://github.com/uchaloop/secret/compare/v1.0.0...v2.0.0
[1.0.0]: https://github.com/uchaloop/secret/releases/tag/v1.0.0
