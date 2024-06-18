package server

import (
	"aid-server/pkg/ldb"
	"aid-server/services/user"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var UserDB ldb.DB

func setGracefulShutdown(server *http.Server) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signals
		if err := server.Close(); err != nil {
			panic(err)
		}
		if err := user.FreeDB(UserDB); err != nil {
			return
		}
	}()
}
