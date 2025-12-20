# go-boost

<div align="center">
  <strong>🚀 Generated with TRAE CN AI 🚀</strong>
</div>

## Overview

go-boost is a comprehensive Go toolkit designed to simplify common operations in Go development, built on top of Go's standard library.

## Key Capabilities

- Advanced logging system with multiple levels and encryption
- Cross-platform operating system utilities
- Process management and information retrieval
- File and directory operations
- Database wrapper with SQLite support
- String processing utilities with extensive methods
- HTTP server utilities for easy server setup
- WebSocket client and server functionality
- JSON operations and utilities

## Installation

Install go-boost using the standard Go get command.

## Design Philosophy

go-boost follows an object-oriented design with structs for each functionality area, supporting method chaining and a fluent API.

## WebSocket Support

go-boost provides comprehensive WebSocket functionality, including:

- WebSocket server with easy setup
- WebSocket client with reconnect support
- Message handling and event processing
- Random available port generation for server startup

### Example Usage

#### WebSocket Server with Random Port

```go
package main

import (
    "fmt"
    "github.com/xiang-tai-duo/go-boost"
    "github.com/gorilla/websocket"
)

func main() {
    // Create a new server instance
    serve := boost.NewServe()
    
    // Get a random available port greater than 1024
    port, err := serve.GetRandomAvailablePort()
    if err != nil {
        panic(err)
    }
    
    // Set server to listen on the random port
    serve.Server.Addr = fmt.Sprintf(":%d", port)
    
    // Register WebSocket handler
    serve.OnWebSocket("/ws", nil, func(conn *websocket.Conn, messageType int, data []byte) error {
        // Echo message back
        return conn.WriteMessage(messageType, data)
    })
    
    fmt.Printf("WebSocket server started on port %d\n", port)
}
```

#### WebSocket Client

```go
package main

import (
    "fmt"
    "github.com/xiang-tai-duo/go-boost"
    "github.com/gorilla/websocket"
)

func main() {
    // Create WebSocket client manager
    wsClient := boost.NewWebSocket()
    
    // Connect to WebSocket server
    clientID, err := wsClient.Connect(
        "ws://localhost:8080/ws",
        func(conn *websocket.Conn, messageType int, data []byte) error {
            fmt.Printf("Received: %s\n", string(data))
            return nil
        },
    )
    if err != nil {
        panic(err)
    }
    
    // Send a message
    wsClient.Send(clientID, []byte("Hello, WebSocket!"))
}
```

## License

MIT License

## About

This project is fully generated and maintained with **TRAE CN AI** assistance, demonstrating the power of AI-driven software development in creating high-quality, maintainable Go libraries.