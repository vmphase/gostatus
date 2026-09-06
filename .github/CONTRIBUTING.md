# Contributing to gostatus

Thank you for your interest in improving gostatus. This document covers how to
set up your environment, the conventions used, and
what to expect from the review process. Please also read the
[Code of Conduct](CODE_OF_CONDUCT.md), [AI Contribution Policy](AI_POLICY.md),
and [Security Policy](SECURITY.md) before opening an issue or pull request.
All contributions are held to all three.

## Before you start

- For small fixes (typos, broken links, minor bugs), feel free to open a pull
  request directly.
- For anything larger, like new features, breaking changes, or architectural
  changes — please open an issue first to discuss the approach. This avoids
  wasted effort if the direction doesn't fit the project.
- Check existing issues and pull requests to avoid duplicating work already
  in progress.
- Fill out the pull request template. At minimum, all template sections
  must be present, at least one box under `Type of change` must be checked, and every box under
  `Checklist` must be checked before your PR can be merged. CI will flag
  the PR and block it otherwise.

## Setting up your environment

- [Go 1.26.0+](https://go.dev/dl/) is required.
- The project uses `go.mod` for dependency management; no separate tooling
  is needed beyond the standard toolchain.

## Commit and branch conventions

- Use descriptive branch names, e.g. `fix/reconnect-timeout` or
  `feat/expose-history`.
- Commit messages must follow the
  [Conventional Commits](https://www.conventionalcommits.org/) format, e.g.
  `fix(gateway): resolve token exchange error` or
  `feat(handler): support additional badge styles`.
- Your **pull request title** must also follow
  [Conventional Commits](https://www.conventionalcommits.org/) format,
  independently of your individual commit messages — this is checked
  automatically and is required for the PR to pass CI.
- Keep commits focused; unrelated changes should be split into separate
  commits or pull requests.

## Code quality checks

CI runs the following checks on every pull request — run them locally before
pushing to avoid review round-trips:

- `gofmt` and `goimports` — files must be formatted (run
  `gofmt -w .` and `goimports -w .` to fix).
- `golangci-lint run ./...` — no lint issues.
- `go build ./...` — all packages must compile.

## Issue reports

- Use the appropriate issue template (bug report or feature request).
- Include enough detail for a maintainer to understand and, where
  applicable, reproduce the problem: expected behavior, actual behavior, and
  steps to reproduce.
- One issue per report. Don't bundle multiple unrelated problems or
  requests into a single issue.

## Usage of AI

AI tools may be used when preparing contributions, but contributors are fully
responsible for the work they submit and must be able to explain and justify
it themselves. All use of AI must comply with the
[AI Contribution Policy](AI_POLICY.md).

## Code review

- Reviews focus on correctness, maintainability, and fit with the project's
  existing design. Feedback is intended to improve the contribution, not to
  discourage the contributor.
- I may request changes, ask clarifying questions, or close
  contributions that don't meet the standards described in this document.
- Please engage with feedback in good faith.

## Security issues

Do not report security vulnerabilities through public issues or pull
requests. See [SECURITY.md](SECURITY.md) for how to report them privately.

## Questions

If anything in this document is unclear, or you're unsure whether a change
is a good fit, ask in an issue or by emailing
[vmphase.dev@gmail.com](mailto:vmphase.dev@gmail.com) before investing
significant time.
