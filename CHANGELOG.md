# Changelog

All notable changes to this module are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this module adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.0.2] - 2026-09-01

### Changed

- The package documentation carries what is masked, what `Reveal` and `Clear`
  promise and what the type does not protect against; the README is a landing
  page.
- The README no longer shows a `koanf:"-"` tag beside the `env` one. There is no
  file decoder left to keep a secret out of.

### Fixed

- The changelog headings carry the version in brackets, the form the release
  workflow looks for. It would have found no section for the next tag and failed
  the job.

## [2.0.1] - 2026-08-06

- Reworked the README as concise, user-focused documentation.

## [2.0.0] - 2026-08-04

- `Secret` is now an opaque struct, preventing ordinary string conversion,
  concatenation, indexing and slicing.
- Added `New`, `IsZero` and `Clear`.
- Added the sealed `Value` marker interface for integration libraries.
- Added comprehensive masking through `fmt.Formatter`, text/JSON/XML
  serialization and `log/slog`.
- This release changes the module path to `github.com/uchaloop/secret/v2`.

## [1.0.0] - 2026-08-04

A masked string type for sensitive values.

- `Secret` string type that masks itself in `fmt` (`%v`/`%s`/`%q`/`%+v`/`%#v`/`%x`),
  `encoding/json`, `encoding/text` (via `MarshalText`) and `log/slog` (via
  `LogValue`); an empty `Secret` renders as empty.
- `Reveal` returns the real value - the only way it leaves the type, so every real
  use is greppable.
- `UnmarshalText` reads a value in (for env/text decoders) while every outward
  representation stays masked.
