# Security policy

secret helps prevent accidental disclosure through supported formatting,
logging and serialization interfaces. It is not secure memory or encryption.
See the package GoDoc and [README](README.md) for its guarantees and limits.

## Reporting a vulnerability

Do not publish suspected disclosure details in a public issue or pull request
before discussing them privately with the maintainer.

Use GitHub's **Report a vulnerability** action on the repository's Security tab
when available:
[private vulnerability report](https://github.com/uchaloop/secret/security/advisories/new).
If that action is unavailable, open an issue asking the maintainer to enable
private reporting, without exploit details, sensitive output or credentials.

Include the affected library and Go versions, a minimal reproduction using
synthetic values, the affected interface, and expected versus actual behavior.
Never include real passwords, tokens, production configuration or memory dumps
containing sensitive information.

The maintainer will assess the report and coordinate any fix and disclosure.
No response-time or historical-version support commitment is currently defined.
