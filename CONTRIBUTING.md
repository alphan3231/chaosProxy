# Contributing to Chaos-Proxy

First off, thanks for taking the time to contribute! 🎉

## How Can I Contribute?

### Reporting Bugs

- Use the GitHub issue tracker
- Include steps to reproduce
- Include Go/Python versions
- Include relevant logs

### Suggesting Features

- Open an issue with the `enhancement` label
- Describe the use case
- Explain how it would work

### Pull Requests

1. Fork the repo
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Commit Message Convention

We follow [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` New feature
- `fix:` Bug fix
- `docs:` Documentation changes
- `style:` Code style changes (formatting, etc.)
- `refactor:` Code refactoring
- `test:` Adding tests
- `chore:` Maintenance tasks
- `security:` Security improvements

### Code Style

**Go:**
- Run `go fmt` before committing
- Follow [Effective Go](https://golang.org/doc/effective_go.html)

**Python:**
- Use Black for formatting
- Follow PEP 8

**TypeScript:**
- Use ESLint configuration provided
- Run `npm run lint` before committing

## Development Setup

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/chaosProxy.git
cd chaosProxy

# Install Go dependencies
go mod download

# Install Python dependencies
cd brain && pip install -r requirements.txt && cd ..

# Install Dashboard dependencies
cd dashboard && npm install && cd ..

# Start Redis
docker-compose up -d

# Run tests (when available)
go test ./...
```

## Questions?

Feel free to open an issue with the `question` label.

---

Thank you for contributing! 👻
