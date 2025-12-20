# Contributing to terraform-provider-idgen

## Feedback and Discussion

I welcome your feedback, ideas, and questions about this project.  Feel free to:

* Open discussions to share thoughts or propose changes
* Ask questions about usage or implementation
* Report issues or suggest improvements
* Share your use cases and experiences

## Requirements

All commits must be cryptographically signed. Configure commit signing by following the [GitHub documentation on signing commits](https://docs.github.com/en/authentication/managing-commit-signature-verification/signing-commits).

## Getting Started

Before contributing code:

1. Open a discussion or issue to outline your proposed changes
2. Wait for feedback to ensure alignment with project goals
3. Fork the repository and create a feature branch
4. Ensure tests pass and add new tests for your changes
5. Submit a pull request with a clear description

## Code Standards

* Follow existing code style and conventions
* Include tests for new functionality
* Update documentation as needed
* Keep commits focused and atomic
* Commit messages should follow the [conventional commits style](https://www.conventionalcommits.org/en/v1.0.0-beta.4/#specification)

To ensure code standards locally, install [pre-commit](https://pre-commit.com) and run before push

```bash
pre-commit run --all-files
```

## Development

### Running Tests


```bash
go test ./...
```
