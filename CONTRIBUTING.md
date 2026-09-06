# Contributing

Thank you for your interest in contributing to the SocialKit Go SDK! This
document outlines the process for submitting changes.

## Getting Started

1. Fork the repository and clone your fork.
2. Ensure you have Go 1.21 or later installed.
3. Run `go test ./...` to verify the existing test suite passes.

## Development Workflow

1. Create a branch for your change:

   ```sh
   git checkout -b my-feature
   ```

2. Make your changes, keeping the following in mind:
   - Follow the existing code style and conventions.
   - Add or update tests for any new or changed behaviour.
   - Ensure all tests pass: `go test ./...`
   - Ensure `go vet ./...` reports no issues.
   - Ensure code is formatted: `gofmt -l .` should produce no output.

3. Commit your changes with a clear, descriptive commit message.

4. Push your branch and open a pull request against the `main` branch.

## Pull Request Guidelines

- Provide a clear description of what the PR does and why.
- Reference any related issues.
- Keep changes focused; avoid unrelated refactoring in the same PR.
- All CI checks must pass before merge.

## Reporting Issues

Use GitHub Issues to report bugs or request features. Please include:

- A clear description of the problem or request.
- Steps to reproduce (for bugs).
- Expected and actual behaviour.
- Go version and OS.

## Code of Conduct

All contributors are expected to abide by the [Code of Conduct](CODE_OF_CONDUCT.md).
