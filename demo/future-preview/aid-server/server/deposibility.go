package server

import (
	"aid-server/pkg/ldb"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var UserDB ldb.DB
var UserMapDB ldb.DB

func setGracefulShutdown(server *http.Server) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signals
		if err := server.Close(); err != nil {
			panic(err)
		}
		if err := ldb.FreeDB(UserDB); err != nil {
			return
		}
		if err := ldb.FreeDB(UserMapDB); err != nil {
			return
		}
	}()
}
