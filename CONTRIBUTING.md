# Contributing to GoRedis-Lite

We love your input! We want to make contributing to GoRedis-Lite as easy and transparent as possible, whether it's:

- Reporting a bug
- Discussing the current state of the code
- Submitting a fix
- Proposing new features
- Becoming a maintainer

## Development Process

We use GitHub to host code, to track issues and feature requests, as well as accept pull requests.

## Pull Requests

1. Fork the repo and create your branch from `main`.
2. If you've added code that should be tested, add tests.
3. If you've changed APIs, update the documentation.
4. Ensure the test suite passes.
5. Make sure your code lints.
6. Issue that pull request!

## Setting Up Development Environment

```bash
# Clone your fork
git clone https://github.com/your-username/goredis-lite.git
cd goredis-lite

# Install dependencies
go mod download

# Run tests
make test

# Run linter
make lint

# Build the project
make build
```

## Code Style

- Use `gofmt` to format your code
- Follow standard Go conventions
- Write clear commit messages
- Add comments for complex logic

## Testing

- Write unit tests for new features
- Ensure all tests pass before submitting PR
- Include integration tests where appropriate
- Aim for good test coverage

## Reporting Bugs

We use GitHub issues to track public bugs. Report a bug by opening a new issue.

**Great Bug Reports** tend to have:

- A quick summary and/or background
- Steps to reproduce
- What you expected would happen
- What actually happens
- Notes (possibly including why you think this might be happening, or stuff you tried that didn't work)

## Feature Requests

We actively welcome feature requests. Please use GitHub issues with the "enhancement" label.

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

## Code of Conduct

This project adheres to a code of conduct. By participating, you are expected to uphold this code.
