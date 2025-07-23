
# 🚀 PulseDB

[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.20+-blue.svg)](https://golang.org/)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)]()

A high-performance, lightweight Redis clone written in Go, implementing the RESP2 protocol from scratch. 
This project serves as a deep dive into building low-level systems and understanding the inner workings of Redis.

> **⚠️ Educational Project**: This is primarily an educational implementation designed for learning purposes. For production use, consider [Redis](https://redis.io/) or [KeyDB](https://keydb.dev/).

## 🎯 Quick Start

```bash
# Using Docker (Recommended)
docker-compose up --build

# Or build from source
go build -o PulseDB ./cmd/PulseDB
./PulseDB
```

---

## ✨ Features

### Core Functionality
- 🔧 **Custom RESP2 Protocol**: Complete implementation of Redis Serialization Protocol v2
- 📊 **Redis Data Types**: Support for all core Redis data types:
  - Simple Strings (`+OK\r\n`)
  - Errors (`-ERROR message\r\n`)
  - Integers (`:100\r\n`)
  - Bulk Strings (`$6\r\nfoobar\r\n`)
  - Arrays (`*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n`)
- ⚡**High Concurrency**: Thread-safe operations with concurrent client handling
- 💾 **Persistence**: AOF (Append Only File) and memory snapshots
- ⏰ **Key Expiration**: TTL support with automatic cleanup

### Performance & Reliability
- 🔒 **Thread Safety**: Robust concurrent access with RWMutex
- 📈 **Benchmarks**: Comprehensive performance testing included
- 🧪 **Test Coverage**: Extensive test suite with >95% coverage
- 📊 **Memory Efficient**: Optimized memory usage patterns

---

## 📁 Project Structure

```
PulseDB/
├── cmd/
│   └── PulseDB/        # Main application entry point
├── internal/
│   ├── command/        # Command parsing and execution
│   ├── resp/           # RESP2 protocol implementation
│   ├── server/         # TCP server setup and client handling
│   ├── store/          # In-memory data storage with persistence
│   └── utils/          # Utility functions and helpers
├── bin/                # Compiled binaries
├── docker-compose.yml  # Docker setup for easy deployment
├── Dockerfile          # Container configuration
├── Makefile           # Build automation
├── go.mod             # Go module dependencies
└── README.md          # Project documentation
```

---

## 🛠️ Getting Started

### Prerequisites

- **Go 1.20+** - [Download & Install Go](https://golang.org/dl/)
- **Docker** (optional) - [Get Docker](https://docs.docker.com/get-docker/)
- **Make** (optional) - For using the Makefile commands

### Installation Options

#### Option 1: Using Docker (Recommended)

```bash
# Clone the repository
git clone https://github.com/DNahar74/PulseDB.git
cd PulseDB

# Run with Docker Compose
docker-compose up --build

# The server will be available at localhost:6378
```

#### Option 2: Build from Source

```bash
# Clone the repository
git clone https://github.com/DNahar74/PulseDB.git
cd PulseDB

# Build using Makefile
make build

# Or build manually
go build -o bin/PulseDB ./cmd/PulseDB

# Run the server
./bin/PulseDB
```

#### Option 3: Development Mode

```bash
# Run directly with Go
make run

# Or manually
go run ./cmd/PulseDB
```

### Configuration

PulseDB accepts the following command-line flags:

- `-addr` : Server address (default: `:6379`)
- `-v` : Show version information
- `-verbose` : Enable verbose logging

```bash
./bin/PulseDB -addr :6378 -verbose
```

---

## 🧪 Usage

### Connecting to PulseDB

Once the server is running, you can connect using various Redis clients:

#### Using Redis CLI

```bash
# If you have redis-cli installed
redis-cli -p 6379

# Using Docker with redis-cli
docker run -it --rm redis:alpine redis-cli -h host.docker.internal -p 6379
```

#### Using Telnet

```bash
telnet localhost 6379
```

#### Using Docker Compose Redis CLI

```bash
# Start the redis-cli service defined in docker-compose.yml
docker-compose run --rm redis-cli
```

### Basic Commands

```bash
# Test connection
PING

# Set and get values
SET mykey "Hello, PulseDB!"
GET mykey

# Set with expiration (100 seconds)
SET tempkey "temporary value" EX 100
GET tempkey

# Delete keys
DEL mykey tempkey

# Echo command
ECHO "Hello World"
```

---

## 🔧 Development

### Building

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Clean build artifacts
make clean
```

### Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run benchmarks
make benchmark
```

### Code Quality

```bash
# Format code
make fmt

# Run linter
make lint

# Run security checks
make security
```

---

## 🐳 Docker

The project includes Docker support for easy deployment:

```bash
# Build Docker image
docker build -t pulsedb .

# Run container
docker run -p 6379:6378 pulsedb

# Using docker-compose for full setup
docker-compose up -d
```

---

## 📊 Performance

PulseDB is designed for high performance with the following characteristics:

- **Concurrent Connections**: Handles thousands of concurrent clients
- **Memory Efficiency**: Optimized data structures for minimal memory footprint
- **Protocol Efficiency**: Full RESP2 implementation with minimal overhead
- **Persistence**: AOF and snapshot mechanisms for data durability

### Benchmarks

Run benchmarks to test performance:

```bash
make benchmark
```

---

## 🛠️ Implemented Commands

### 🟢 `PING`

- **Description**: Tests the connection with the server.
- **Usage**:  
  ```bash
  PING
  ```
- **Response**:  
  ```
  +PONG
  ```

---

### 🗣️ `ECHO`

- **Description**: Returns the input string.
- **Usage**:  
  ```bash
  ECHO "hello world"
  ```
- **Response**:  
  ```
  "hello world"
  ```

---

### 💾 `SET`

- **Description**: Stores a key with a string value. Optional expiry flag in seconds.
- **Usage**:  
  ```bash
  SET hello world
  SET hello world EX 100  # Key expires in 100 seconds
  ```
- **Response**:  
  ```
  +OK
  ```

---

### 📥 `GET`

- **Description**: Retrieves the value of the given key.
- **Usage**:  
  ```bash
  GET hello
  ```
- **Response** (if found):  
  ```
  "world"
  ```

---

### ❌ `DEL`

- **Description**: Deletes the specified key.
- **Usage**:  
  ```bash
  DEL hello
  ```
- **Response**:  
  ```
  +OK
  ```


---

## 📚 RESP2 Protocol Overview

PulseDB implements the Redis Serialization Protocol (RESP) version 2 for client-server communication.

### Data Types

| Type | Format | Example |
|------|--------|---------|
| **Simple Strings** | `+<string>\r\n` | `+OK\r\n` |
| **Errors** | `-<error>\r\n` | `-ERROR Invalid command\r\n` |
| **Integers** | `:<number>\r\n` | `:1000\r\n` |
| **Bulk Strings** | `$<length>\r\n<string>\r\n` | `$6\r\nfoobar\r\n` |
| **Arrays** | `*<count>\r\n<element1>...<elementN>` | `*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n` |

For detailed protocol specification, see the [Redis Protocol documentation](https://redis.io/docs/reference/protocol-spec/).

---

## 🚀 Roadmap

### Completed ✅
- [x] RESP2 protocol implementation
- [x] Basic Redis commands (PING, ECHO, SET, GET, DEL)
- [x] Key expiration (TTL)
- [x] AOF persistence
- [x] Memory snapshots
- [x] Docker support
- [x] Concurrent client handling

### In Progress 🚧
- [ ] More Redis commands (INCR, DECR, LPUSH, RPOP, etc.)
- [ ] Redis data structures (Lists, Sets, Hashes)
- [ ] Pub/Sub functionality
- [ ] Clustering support

### Future Plans 📋
- [ ] RESP3 protocol support
- [ ] Redis modules compatibility
- [ ] Replication
- [ ] Lua scripting support

---

## 🧱 Contributing

We welcome contributions! Here's how you can help:

### Getting Started

1. **Fork** the repository
2. **Clone** your fork: `git clone https://github.com/yourusername/PulseDB.git`
3. **Create** a feature branch: `git checkout -b feature/amazing-feature`
4. **Make** your changes
5. **Test** your changes: `make test`
6. **Commit** your changes: `git commit -m 'Add amazing feature'`
7. **Push** to the branch: `git push origin feature/amazing-feature`
8. **Open** a Pull Request

### Development Guidelines

- Follow Go best practices and conventions
- Add tests for new functionality
- Update documentation for API changes
- Run `make fmt` and `make lint` before committing
- Write clear, descriptive commit messages

### Issues

Found a bug or have a feature request? Please [open an issue](https://github.com/DNahar74/PulseDB/issues) with:

- Clear description of the problem/feature
- Steps to reproduce (for bugs)
- Expected vs actual behavior
- Go version and OS information

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 🙌 Acknowledgements

- **[Redis](https://redis.io/)** - For the inspiration and protocol specification
- **[Go Team](https://golang.org/)** - For creating an amazing programming language
- **[Docker](https://docker.com/)** - For simplifying deployment and development
- **Open Source Community** - For continuous inspiration and support

---

## 📞 Contact

**Project Maintainer**: [DNahar74](https://github.com/DNahar74)

- **Issues**: [GitHub Issues](https://github.com/DNahar74/PulseDB/issues)
- **Discussions**: [GitHub Discussions](https://github.com/DNahar74/PulseDB/discussions)

---

<div align="center">
  <strong>⭐ Star this repository if you found it helpful! ⭐</strong>
</div>
