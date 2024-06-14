package server

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func setGracefulShutdown(server *http.Server) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signals
		err := server.Close()
		if err != nil {
			panic(err)
		}
	}()
}
