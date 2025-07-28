# Contribution Guidelines for SENTINEL

Thank you for your interest in contributing to SENTINEL! This project aims to be a simple yet effective monitoring system written in Go.

Repository: [https://github.com/0xReLogic/SENTINEL](https://github.com/0xReLogic/SENTINEL)

## How to Contribute

### Reporting Bugs

If you find a bug, please create a new issue with the following information:

1. Clear and descriptive title
2. Steps to reproduce the bug
3. Expected behavior
4. Actual behavior
5. Screenshots (if applicable)
6. Environment information (OS, Go version, etc.)

### Proposing Features

If you have an idea for a new feature, please create a new issue with the "enhancement" label and explain:

1. The feature you're proposing
2. Why this feature is useful
3. How you envision its implementation

### Pull Requests

We greatly appreciate pull requests from contributors! Here are the steps to submit a PR:

1. Fork the repository
2. Create a new branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Code Conventions

- Use `gofmt` to format your code
- Follow [Effective Go](https://golang.org/doc/effective_go.html) and [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Add comments for functions and packages you create
- Write unit tests for new code

## Project Structure

```
SENTINEL/
├── checker/       # Package for service checking
├── config/        # Package for configuration management
├── main.go        # Main program file
└── sentinel.yaml  # Configuration file
```

## Development Process

1. Choose an issue you want to work on or create a new one
2. Discuss your approach in the issue
3. Implement your solution
4. Write tests for your code
5. Submit a PR

## License

By contributing to this project, you agree that your contributions will be licensed under the same MIT license as the project.