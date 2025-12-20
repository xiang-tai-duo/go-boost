// --------------------------------------------------------------------------------
// File:        serve.go
// Author:      TRAE AI
// Created:     2025/12/20 12:31:58
// Description: Examples of using the go-boost HTTP server
// --------------------------------------------------------------------------------
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	. "github.com/xiang-tai-duo/go-boost"
)

func main() {
	sample_serve_main()

	fmt.Println("\nAll examples completed. Press Ctrl+C to exit.")
	select {}
}

func sample_serve_basic_http_server() {
	serve := NewServe()

	// Check if default port 80 is available
	if !serve.CheckPortAvailable(80) {
		fmt.Println("Port 80 is in use, getting available port...")
		port, err := serve.GetAvailablePort()
		if err != nil {
			fmt.Printf("Failed to get available port: %v\n", err)
			return
		}
		serve.Server.Addr = fmt.Sprintf(":%d", port)
	}

	serve.On(http.MethodGet, "/", func(w http.ResponseWriter, r *http.Request) error {
		_, err := fmt.Fprintf(w, "Hello, World!\n")
		return err
	})

	fmt.Printf("Basic HTTP server started on http://localhost%s\n", serve.Server.Addr)
	fmt.Println("Route: GET / -> Hello, World!")

	time.Sleep(100 * time.Millisecond)
}

func sample_serve_http_server_with_routes() {
	serve := NewServe()

	// Check if default port 80 is available
	if !serve.CheckPortAvailable(80) {
		fmt.Println("Port 80 is in use, getting available port...")
		port, err := serve.GetAvailablePort()
		if err != nil {
			fmt.Printf("Failed to get available port: %v\n", err)
			return
		}
		serve.Server.Addr = fmt.Sprintf(":%d", port)
	}

	serve.On(http.MethodGet, "/", func(w http.ResponseWriter, r *http.Request) error {
		_, err := fmt.Fprintf(w, "Home Page\n")
		return err
	})

	serve.On(http.MethodGet, "/about", func(w http.ResponseWriter, r *http.Request) error {
		_, err := fmt.Fprintf(w, "About Page\n")
		return err
	})

	serve.On(http.MethodPost, "/api/data", func(w http.ResponseWriter, r *http.Request) error {
		w.Header().Set("Content-Type", "application/json")
		return json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Data received"})
	})

	fmt.Printf("HTTP server with routes started on http://localhost%s\n", serve.Server.Addr)
	fmt.Println("Routes:")
	fmt.Println("  GET / -> Home Page")
	fmt.Println("  GET /about -> About Page")
	fmt.Println("  POST /api/data -> JSON response")

	time.Sleep(100 * time.Millisecond)
}

func sample_serve_https_server_with_tls() {
	serve := NewServe()

	// Check if default port 80 is available
	if !serve.CheckPortAvailable(80) {
		fmt.Println("Port 80 is in use, getting available port...")
		port, err := serve.GetAvailablePort()
		if err != nil {
			fmt.Printf("Failed to get available port: %v\n", err)
			return
		}
		serve.Server.Addr = fmt.Sprintf(":%d", port)
	}

	serve.On(http.MethodGet, "/", func(w http.ResponseWriter, r *http.Request) error {
		_, err := fmt.Fprintf(w, "Secure HTTPS Server\n")
		return err
	})

	fmt.Printf("HTTPS server with TLS started on https://localhost%s\n", serve.Server.Addr)
	fmt.Println("Note: Using self-signed certificate (may show security warning in browser)")

	time.Sleep(100 * time.Millisecond)
}

func sample_serve_main() {
	fmt.Println("=== Example 1: Basic HTTP Server ===")
	sample_serve_basic_http_server()

	fmt.Println("\n=== Example 2: HTTP Server with Routes ===")
	sample_serve_http_server_with_routes()

	fmt.Println("\n=== Example 3: HTTPS Server with TLS ===")
	sample_serve_https_server_with_tls()

	fmt.Println("\n=== Example 4: Static File Server ===")
	sample_serve_static_file_server()
}

func sample_serve_static_file_server() {
	serve := NewServe()

	// Check if default port 80 is available
	if !serve.CheckPortAvailable(80) {
		fmt.Println("Port 80 is in use, getting available port...")
		port, err := serve.GetAvailablePort()
		if err != nil {
			fmt.Printf("Failed to get available port: %v\n", err)
			return
		}
		serve.Server.Addr = fmt.Sprintf(":%d", port)
	}

	serve.AddStaticDirectory("/static", "./public")

	serve.On(http.MethodGet, "/", func(w http.ResponseWriter, r *http.Request) error {
		_, err := fmt.Fprintf(w, "Static File Server\n")
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(w, "Static files available at: /static/*\n")
		return err
	})

	fmt.Printf("Static file server started on http://localhost%s\n", serve.Server.Addr)
	fmt.Printf("Static files served from: ./public -> http://localhost%s/static/\n", serve.Server.Addr)
	fmt.Println("Note: Create a 'public' directory with files to test this example")

	time.Sleep(100 * time.Millisecond)
}
