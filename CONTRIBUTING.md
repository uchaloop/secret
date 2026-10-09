# Contributing

Bug fixes, documentation improvements and focused feature proposals are welcome.
Discuss new features and incompatible changes in an issue before implementing them.
For ordinary bug reports, include a minimal reproduction, the Go and secret
versions, and the expected and actual behavior. Use synthetic values only.
For suspected disclosures, follow [SECURITY.md](SECURITY.md).

## Project scope

secret protects sensitive strings from accidental disclosure through supported
formatting, logging and serialization interfaces. `Reveal` provides explicit
access to the underlying string. The library has no external dependencies.

ENV loading, configuration management, Vault clients, encryption and secure
memory management are outside its responsibility. Do not promise physical memory
erasure or protection from reflection, unsafe code, debuggers or memory dumps.
Discuss additions to supported interfaces before implementing them: an interface
can change how existing decoders or loggers treat a value.

## Behavioral guarantees

Changes must preserve the documented public contract, including:

- Masking through supported output interfaces, including nested values and
  formatting of unexported fields without disclosure.
- A usable empty zero value and independent copies when clearing or replacing
  a secret.
- Input-free errors from the library's decoding methods and preservation of the
  receiver on their errors.
- Strict JSON input without silently replacing malformed Unicode, and text
  input that preserves arbitrary bytes.
- Concurrent reads, with caller synchronization for mutation of a shared variable.

These guarantees apply within the boundaries described in the package GoDoc.
They do not extend to enclosing decoders or types that embed Secret and override
its behavior. The sensitivity marker is a hint, not proof of safe formatting.

## Branches

Start from an up-to-date `main` and open pull requests against `main`. Use a
short English description with hyphens:

- `fix/json-unicode` for bug fixes.
- `feat/output-interface` for new features.
- `docs/masking-guarantees` for documentation.
- `refactor/json-validation` for internal restructuring.
- `ci/release-checks` for workflow changes.

```sh
git switch main
git pull --ff-only
git switch -c fix/json-unicode
```

For a fork, update from this repository's `main` and target this repository in
the pull request. Maintainers use `release/<version>` for release preparation.
Branch names are recommendations, not CI requirements; issue numbers are optional.

## Commit messages

Start each subject with the full source branch name, a colon and a short English
description of the actual change:

```text
fix/json-unicode: reject unpaired surrogate escapes
refactor/json-validation: simplify Unicode checks
release/2.1.0: prepare library release
```

For a branch named `feat/2.1.0`, use the prefix `feat/2.1.0: `. Preserve the source
branch prefix in squash commit subjects. This convention is not enforced by CI.

## Development and tests

The minimum Go version is defined in [go.mod](go.mod), currently Go 1.21.0.
Check changes on the minimum and newer toolchains listed in
[CI](.github/workflows/ci.yml). Keep newer APIs behind appropriate build tags
when they are used only by tests for newer standard-library integrations.

Format changed Go files with `gofmt`, then run:

```sh
GOWORK=off go test -race ./...
GOWORK=off go vet ./...
```

Tests should verify public behavior useful to callers. Add regression tests for
bug fixes, especially disclosure and input-corruption cases. Preserve meaningful
checks for copy independence, concurrency and codec correctness. Avoid tests
that mirror the storage representation or duplicate an existing scenario.

Run relevant fuzz tests when changing input handling:

```sh
GOWORK=off go test -run '^$' -fuzz '^FuzzJSONInput$' -fuzztime=30s
GOWORK=off go test -run '^$' -fuzz '^FuzzSecretTextMasking$' -fuzztime=30s
```

Use coverage to identify missing scenarios, not as a target percentage. Support
performance claims with relevant benchmarks and before/after results. Do not
weaken masking or error guarantees for speculative optimization.

## Coverage

CI uploads coverage from the Go 1.27 job to Codecov. Both toolchains still run
all applicable tests. Project and patch coverage statuses are informational;
there is no required percentage. Upload failures do not block CI, while test
failures do.

```sh
GOWORK=off go test -race -covermode=atomic -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

Maintainers must enable `uchaloop/secret` in Codecov and save its upload token as
the GitHub Actions repository secret `CODECOV_TOKEN`. Never commit the token.
Fork pull requests do not receive this secret; tokenless uploads depend on the
Codecov settings. The README badge needs a successful upload for `main`.
See the [Codecov action documentation](https://github.com/codecov/codecov-action).

## Code style

- Choose names that explain purpose and avoid shadowing built-in identifiers.
- Preserve the existing public `New` constructor. Use `Make…` / `make…` for new
  constructors; do not rename public API solely to apply a style convention.
- Prefer early returns, `continue` or `break` when they reduce nesting. Do not
  introduce `goto` for speculative optimization.
- Handle errors explicitly before inspecting or returning successful results.
  Preserve partial results only when required by the function's contract.
- Separate logical steps with blank lines; keep calls and their error checks
  together. Use `len(value)` for string emptiness checks.
- Inline local struct types used only once when this improves readability.
- Write GoDoc for callers: behavior, guarantees and relevant limits. Internal
  comments should explain non-obvious reasons rather than restate code.

## Pull requests and versions

Keep changes focused. Describe the problem, resulting behavior and validation.
Update affected examples and documentation, and add a concise entry under the current unreleased version to [CHANGELOG.md](CHANGELOG.md) for user-visible changes.

Compatibility includes masking, decoding, error behavior, zero values,
comparability and copy semantics as well as exported signatures. Discuss
incompatible changes explicitly; unchanged signatures do not imply compatibility.

Release branches use `release/x.y.z`, where `x.y.z` is the target version:

- `x` (major) changes for incompatible public contract changes.
- `y` (minor) changes for compatible functionality additions.
- `z` (patch) changes for compatible bug fixes.

Reset lower components to zero when increasing a major or minor version.
Feature and fix branches normally describe the task without a version number.
Their prefixes do not determine the release version.

Maintainers publish Git tags such as `v2.1.0` or `v2.1.0-rc.1`; a branch alone
publishes nothing. Examples here do not select the next release version. Never
move or replace a published tag. The current module path ends in `/v2`; another
major version requires the corresponding module path and import changes.
See [Semantic Versioning](https://semver.org/) and
[Go module version numbering](https://go.dev/doc/modules/version-numbers).

Never include real credentials in source, examples, fixtures, logs or reports.

The 2.1.0 release was an explicit exception to the compatibility policy:
it retained `/v2` while changing comparability and JSON input behavior. These
changes are listed in the README and changelog; this exception does not redefine
Semantic Versioning or authorize future incompatible minor releases.

## Changelog style

Use the same format in confmaker, confx and secret:

- Keep newest versions first. Use `## [x.y.z] - Unreleased` once the target version
  is known, or `## [Unreleased]` before choosing it. Do not add a separate target
  release paragraph. Version headings omit `v`; Git tags include it.
- At publication, replace `Unreleased` with the actual release date in
  `YYYY-MM-DD` format. Do not infer publication dates from commit dates. Preserve
  undated historical entries when the release date cannot be verified.
- Group entries under `Added`, `Changed`, `Deprecated`, `Removed`, `Fixed`,
  `Security`, in that order. Omit empty categories, even for small releases.
- Write in English, using past-tense opening verbs such as `Added`, `Changed`,
  `Removed` and `Fixed`. Each bullet describes one user-visible change, normally
  in one or two sentences. Include internal work only when it explains a useful
  result. Preserve necessary historical compatibility and retraction notes.
- Prefix incompatible changes with `**Breaking:**` in their normal category and
  explain the replacement or required action when applicable. Do not duplicate
  them in a separate breaking-changes section.
- Define heading links at the bottom of the file. An unreleased target compares
  the last released tag with `HEAD`; a released version compares its predecessor
  with its tag. Link the first release to its release page. Preserve explicit
  historical exceptions for retracted or incorrectly tagged versions.
- Separate headings, paragraphs and lists with one blank line. Wrap continuation
  lines consistently and format API names as inline code.

Example before publication:

```markdown
## [2.0.0] - Unreleased

### Added

- Added support for application-provided configuration engines.

### Changed

- **Breaking:** Renamed `EnvOption` to `LoaderOption`; update option declarations.

[2.0.0]: https://github.com/uchaloop/confmaker/compare/v1.0.1...HEAD
```
