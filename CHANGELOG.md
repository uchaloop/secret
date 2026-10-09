# Changelog

Notable changes by version.

## Unreleased

- Add the `Secret.IsSensitive()` marker so integrations can recognize sensitive
  types through a structural interface without importing this module. Keep the
  existing `Value` interface and all masking and decoding behavior unchanged.

## [2.1.0]

- Fixed secret disclosure through `fmt` and `slog` in nested and unexported fields.
- **Breaking:** `Secret` is no longer comparable or usable as a map key; use
  `IsZero()` to check emptiness. JSON `null` now clears the secret.
- Reject malformed JSON and Unicode with safe errors, preserving the previous value.
  Nil decode receivers return errors.
- Lowered the minimum Go version to 1.21.0.
- Added Codecov and contribution/security guidelines; updated documentation and tests.
- Fixed the `/v2` module path in the release workflow.

The `/v2` import path is unchanged despite the breaking changes above.

## [2.0.2] - 2026-09-01

- Clarified masking guarantees and limitations in GoDoc; simplified the README.
- Removed the obsolete `koanf` tag from the README example.
- Standardized changelog version headings.

## [2.0.1] - 2026-08-06

- Simplified the README and usage documentation.

## [2.0.0] - 2026-08-04

- **Breaking:** Replaced the string-based `Secret` with an opaque struct.
- **Breaking:** Changed the module path to `github.com/uchaloop/secret/v2`.
- Added `New`, `IsZero`, `Clear` and the `Value` sensitivity marker.
- Expanded masking for `fmt`, text/JSON/XML serialization and `slog`.

## [1.0.0] - 2026-08-04

- Introduced a string-based `Secret` with masked formatting, text/JSON output
  and `slog` support; empty secrets produce empty output.
- Added `Reveal` for explicit access and `UnmarshalText` for input.
