
# 🚀 GoRedis-Lite

[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)

A high-performance, lightweight Redis clone written in Go, implementing the RESP2 protocol from scratch. 
This project serves as a deep dive into building low-level systems and understanding the inner workings of Redis.

> **⚠️ Educational Project**: This is primarily an educational implementation. For production use, consider [Redis](https://redis.io/) or [KeyDB](https://keydb.dev/).

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
my-redis/
├── command/        # Command parsing and execution
├── resp/           # RESP2 protocol implementation
├── server/         # TCP server setup and client handling
├── store/          # In-memory data storage
├── utils/          # Utility functions
├── main.go         # Entry point of the application
├── go.mod          # Go module file
└── README.md       # Project documentation
```

---

## 🛠️ Getting Started

### Prerequisites

- Go 1.20 or higher installed on your machine.

### Installation

```bash
# Clone the repository
git clone https://github.com/DNahar74/my-redis.git

# Navigate to the project directory
cd my-redis

# Build the application
go build -o my-redis.exe

# Run the application
./my-redis
```

---

## 🧪 Usage

Once the server is running, you can interact with it using `telnet` or any Redis client:

```bash
telnet localhost 6379
```

Example commands:

```
SET key value
GET key
DEL key
```

Note: Command support is currently limited as this is a work in progress.

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

The Redis Serialization Protocol (RESP) is a simple protocol used by Redis for client-server communication.  
This project implements RESP version 2.

### Simple Strings

- **Format**: `+<string>\r\n`

### Errors

- **Format**: `-<error message>\r\n`

### Integers

- **Format**: `:<number>\r\n`

### Bulk Strings

- **Format**: `$<length>\r\n<string>\r\n`

### Arrays

- **Format**: `*<number of elements>\r\n<element1>\r\n<element2>\r\n...`

For a more detailed explanation, refer to the [Redis Protocol specification](https://redis.io/docs/reference/protocol-spec/).

---

## 🧱 Contributing

Contributions are welcome!  
If you'd like to add features, fix bugs, or improve documentation, please fork the repository and submit a pull request.

---

## 🙌 Acknowledgements

- [Redis](https://redis.io/) for the inspiration and protocol specification.
- [Go Programming Language](https://golang.org/) for its simplicity and performance.
