// --------------------------------------------------------------------------------
// File:        serve.go
// Author:      TRAE AI
// Created:     12/30/2025 11:03:46
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
	testServeMain()
	fmt.Println("\nAll examples completed. Running for up to 5 seconds...")

	// Block for at most 5 seconds
	select {
	case <-time.After(5 * time.Second):
		fmt.Println("Program has been running for 5 seconds, exiting automatically")
	}
}

func testServeBasicHTTPServer() {
	serve := NewServe()

	// Get available port
	port, err := serve.GetAvailablePort()
	if err != nil {
		fmt.Printf("Failed to get available port: %v\n", err)
		return
	}
	addr := fmt.Sprintf(":%d", port)
	serve.On(http.MethodGet, "/", func(w http.ResponseWriter, r *http.Request) error {
		_, err := fmt.Fprintf(w, "Hello, World!\n")
		return err
	})
	fmt.Printf("Basic HTTP server started on http://localhost%s\n", addr)
	fmt.Println("Route: GET / -> Hello, World!")

	// Start server in a goroutine to avoid blocking
	go func() {
		if err := serve.Listen(addr); err != nil {
			fmt.Printf("Error starting basic HTTP server: %v\n", err)
		}
	}()
	time.Sleep(100 * time.Millisecond)
}

func testServeHTTPServerWithRoutes() {
	serve := NewServe()

	// Get available port
	port, err := serve.GetAvailablePort()
	if err != nil {
		fmt.Printf("Failed to get available port: %v\n", err)
		return
	}
	addr := fmt.Sprintf(":%d", port)
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
	fmt.Printf("HTTP server with routes started on http://localhost%s\n", addr)
	fmt.Println("Routes:")
	fmt.Println("  GET / -> Home Page")
	fmt.Println("  GET /about -> About Page")
	fmt.Println("  POST /api/data -> JSON response")
	// Start server in a goroutine to avoid blocking
	go func() {
		if err := serve.Listen(addr); err != nil {
			fmt.Printf("Error starting HTTP server with routes: %v\n", err)
		}
	}()
	time.Sleep(100 * time.Millisecond)
}

func testServeHTTPSServerWithTLS() {
	serve := NewServe()

	// Get available port
	port, err := serve.GetAvailablePort()
	if err != nil {
		fmt.Printf("Failed to get available port: %v\n", err)
		return
	}
	addr := fmt.Sprintf(":%d", port)
	serve.On(http.MethodGet, "/", func(w http.ResponseWriter, r *http.Request) error {
		_, err := fmt.Fprintf(w, "Secure HTTPS Server\n")
		return err
	})
	fmt.Printf("HTTPS server with TLS started on https://localhost%s\n", addr)
	fmt.Println("Note: Using self-signed certificate (may show security warning in browser)")

	// Start server in a goroutine to avoid blocking
	go func() {
		if err := serve.ListenTLS(addr); err != nil {
			fmt.Printf("Error starting HTTPS server with TLS: %v\n", err)
		}
	}()
	time.Sleep(100 * time.Millisecond)
}

func testServeTLSWithSelfSigned() {
	serve := NewServe()

	// Get available port
	port, err := serve.GetAvailablePort()
	if err != nil {
		fmt.Printf("Failed to get available port: %v\n", err)
		return
	}
	addr := fmt.Sprintf(":%d", port)

	// Start HTTPS server with self-signed certificate
	serve.On(http.MethodGet, "/", func(w http.ResponseWriter, r *http.Request) error {
		_, err := fmt.Fprintf(w, "HTTPS Server with Self-Signed Certificate\n")
		return err
	})
	fmt.Printf("Starting HTTPS server with self-signed certificate on https://localhost%s...\n", addr)
	// Start server with self-signed certificate
	go func() {
		if err := serve.ListenTLS(addr); err != nil {
			fmt.Printf("Error starting self-signed HTTPS server: %v\n", err)
		}
	}()
	fmt.Printf("HTTPS server with self-signed certificate started on https://localhost%s\n", addr)
	fmt.Println("Note: Using self-signed certificate (may show security warning in browser)")
	time.Sleep(100 * time.Millisecond)
}

func testServeTLSWithCertKey() {
	serve := NewServe()

	// Get available port
	port, err := serve.GetAvailablePort()
	if err != nil {
		fmt.Printf("Failed to get available port: %v\n", err)
		return
	}
	addr := fmt.Sprintf(":%d", port)

	// Start HTTPS server with cert+key files
	serve.On(http.MethodGet, "/", func(w http.ResponseWriter, r *http.Request) error {
		_, err := fmt.Fprintf(w, "HTTPS Server with Cert+Key Files\n")
		return err
	})
	fmt.Println("\nStarting HTTPS server with cert+key files...")

	// Note: Replace with actual certificate and key file paths
	// Using example paths here, ensure files exist when running
	certPath := "./cert.pem"
	keyPath := "./key.pem"
	go func() {
		if err := serve.ListenTLS(addr, certPath, keyPath); err != nil {
			fmt.Printf("Error starting HTTPS server with cert+key: %v\n", err)
			fmt.Println("Note: This example will fail if cert.pem and key.pem files are not present")
		}
	}()
	fmt.Printf("HTTPS server with cert+key started on https://localhost%s (if cert.pem and key.pem exist)\n", addr)
	time.Sleep(100 * time.Millisecond)
}

func testServeTLSWithPFX() {
	serve := NewServe()

	// Get available port
	port, err := serve.GetAvailablePort()
	if err != nil {
		fmt.Printf("Failed to get available port: %v\n", err)
		return
	}
	addr := fmt.Sprintf(":%d", port)

	// Start HTTPS server with PFX file
	serve.On(http.MethodGet, "/", func(w http.ResponseWriter, r *http.Request) error {
		_, err := fmt.Fprintf(w, "HTTPS Server with PFX File\n")
		return err
	})
	fmt.Println("\nStarting HTTPS server with PFX file...")

	// Note: Replace with actual PFX file path
	// Using example path here, ensure file exists when running
	pfxPath := "./cert.pfx"
	pfxPassword := "password"
	go func() {
		if err := serve.ListenTLS(addr, pfxPath, pfxPassword); err != nil {
			fmt.Printf("Error starting HTTPS server with PFX: %v\n", err)
			fmt.Println("Note: This example will fail if cert.pfx file is not present")
		}
	}()
	fmt.Printf("HTTPS server with PFX started on https://localhost%s (if cert.pfx exists)\n", addr)
	time.Sleep(100 * time.Millisecond)
}

func testServeShutdown() {
	fmt.Println("\n=== Example 7: Server Shutdown Test ===")

	// Test 1: HTTP server shutdown
	serve1 := NewServe()
	port1, err := serve1.GetAvailablePort()
	if err != nil {
		fmt.Printf("Failed to get available port: %v\n", err)
		return
	}
	addr1 := fmt.Sprintf(":%d", port1)

	serve1.On(http.MethodGet, "/", func(w http.ResponseWriter, r *http.Request) error {
		_, err := fmt.Fprintf(w, "Test HTTP Server for Shutdown\n")
		return err
	})

	// Start HTTP server
	fmt.Printf("Starting HTTP server for shutdown test on http://localhost%s...\n", addr1)
	go func() {
		if err := serve1.Listen(addr1); err != nil && err.Error() != "http: Server closed" {
			fmt.Printf("Error starting HTTP server for shutdown test: %v\n", err)
		}
	}()
	time.Sleep(100 * time.Millisecond)

	// Verify server is running
	if serve1.IsRunning() {
		fmt.Println("✓ HTTP server is running (as expected)")
	} else {
		fmt.Println("✗ HTTP server is not running (unexpected)")
	}

	// Shutdown server
	fmt.Println("Shutting down HTTP server...")
	if err := serve1.Shutdown(); err != nil {
		fmt.Printf("Error shutting down HTTP server: %v\n", err)
	} else {
		fmt.Println("✓ HTTP server shutdown successful")
	}

	// Verify server is not running
	if !serve1.IsRunning() {
		fmt.Println("✓ HTTP server is not running after shutdown (as expected)")
	} else {
		fmt.Println("✗ HTTP server is still running after shutdown (unexpected)")
	}

	// Test 2: HTTPS server shutdown
	serve2 := NewServe()
	port2, err := serve2.GetAvailablePort()
	if err != nil {
		fmt.Printf("Failed to get available port: %v\n", err)
		return
	}
	addr2 := fmt.Sprintf(":%d", port2)

	serve2.On(http.MethodGet, "/", func(w http.ResponseWriter, r *http.Request) error {
		_, err := fmt.Fprintf(w, "Test HTTPS Server for Shutdown\n")
		return err
	})

	// Start HTTPS server
	fmt.Printf("\nStarting HTTPS server for shutdown test on https://localhost%s...\n", addr2)
	go func() {
		if err := serve2.ListenTLS(addr2); err != nil && err.Error() != "http: Server closed" {
			fmt.Printf("Error starting HTTPS server for shutdown test: %v\n", err)
		}
	}()
	time.Sleep(100 * time.Millisecond)

	// Verify server is running
	if serve2.IsRunning() {
		fmt.Println("✓ HTTPS server is running (as expected)")
	} else {
		fmt.Println("✗ HTTPS server is not running (unexpected)")
	}

	// Shutdown server
	fmt.Println("Shutting down HTTPS server...")
	if err := serve2.Shutdown(); err != nil {
		fmt.Printf("Error shutting down HTTPS server: %v\n", err)
	} else {
		fmt.Println("✓ HTTPS server shutdown successful")
	}

	// Verify server is not running
	if !serve2.IsRunning() {
		fmt.Println("✓ HTTPS server is not running after shutdown (as expected)")
	} else {
		fmt.Println("✗ HTTPS server is still running after shutdown (unexpected)")
	}

	fmt.Println("\nAll shutdown tests completed.")
}

func testServeMain() {
	fmt.Println("=== Example 1: Basic HTTP Server ===")
	testServeBasicHTTPServer()
	fmt.Println("\n=== Example 2: HTTP Server with Routes ===")
	testServeHTTPServerWithRoutes()
	fmt.Println("\n=== Example 3: HTTPS Server with Self-Signed Certificate ===")
	testServeTLSWithSelfSigned()
	fmt.Println("\n=== Example 4: HTTPS Server with Cert+Key Files ===")
	testServeTLSWithCertKey()
	fmt.Println("\n=== Example 5: HTTPS Server with PFX File ===")
	testServeTLSWithPFX()
	fmt.Println("\n=== Example 6: Static File Server ===")
	testServeStaticFileServer()
	fmt.Println("\n=== Example 7: Server Shutdown Test ===")
	testServeShutdown()
}

func testServeStaticFileServer() {
	serve := NewServe()

	// Get available port
	port, err := serve.GetAvailablePort()
	if err != nil {
		fmt.Printf("Failed to get available port: %v\n", err)
		return
	}
	addr := fmt.Sprintf(":%d", port)
	serve.AddStaticDirectory("/static", "./public")
	serve.On(http.MethodGet, "/", func(w http.ResponseWriter, r *http.Request) error {
		_, err := fmt.Fprintf(w, "Static File Server\n")
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(w, "Static files available at: /static/*\n")
		return err
	})
	fmt.Printf("Static file server started on http://localhost%s\n", addr)
	fmt.Printf("Static files served from: ./public -> http://localhost%s/static/\n", addr)
	fmt.Println("Note: Create a 'public' directory with files to test this example")

	// Start server in a goroutine to avoid blocking
	go func() {
		if err := serve.Listen(addr); err != nil {
			fmt.Printf("Error starting static file server: %v\n", err)
		}
	}()
	time.Sleep(100 * time.Millisecond)
}
