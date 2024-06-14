package server

import (
	"errors"
	"fmt"
	"net"
	"net/http"
)

func Serve(ln net.Listener) error {
	// Create a new server
	server := &http.Server{
		Handler: generateRouter(),
	}

	// Graceful shutdown
	setGracefulShutdown(server)

	// print server info
	fmt.Printf("Server is running on %s://%s\n", "http", ln.Addr().String())
	fmt.Printf("Swagger docs: %s://%s/swagger/index.html\n", "http", ln.Addr().String())

	// Start the server
	err := server.Serve(ln)
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
