# chaos-proxy-cli

CLI tool for managing [Chaos-Proxy](https://github.com/elliot/chaosProxy) - The Immortality Layer for APIs.

## Installation

```bash
npm install -g chaos-proxy-cli
```

## Usage

### Initialize a new project

```bash
# Basic setup
chaos-proxy init

# With Docker support
chaos-proxy init --docker
```

### Check service status

```bash
chaos-proxy status
```

### Start services

```bash
# Using Docker (recommended)
chaos-proxy start --docker

# Manual mode (shows instructions)
chaos-proxy start
```

## Commands

| Command | Description |
|---------|-------------|
| `init` | Initialize a new Chaos-Proxy project |
| `status` | Check the status of Chaos-Proxy services |
| `start` | Start Chaos-Proxy services |

## Options

### init
- `-d, --docker` - Setup with Docker support

### status
- `-r, --redis <url>` - Redis connection URL (default: `redis://localhost:6379`)

### start
- `--docker` - Start using Docker Compose

## License

MIT
