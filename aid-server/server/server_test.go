package server

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"syscall"
	"testing"
	"time"
)

func TestSetGracefulShutdown(t *testing.T) {
	// 創建一個模擬的 http.Server
	server := &http.Server{}

	// 調用 setGracefulShutdown 函數
	setGracefulShutdown(server)

	stopSignal := make(chan bool, 1)

	// start server
	go func() {
		err := server.ListenAndServe()
		assert.NotNil(t, err)
		assert.Equal(t, http.ErrServerClosed, err)
		stopSignal <- true
	}()

	// 模擬發送 SIGINT 信號
	err := syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	assert.Nil(t, err)

	// 等待一段時間,讓 Goroutine 有機會執行
	time.Sleep(100 * time.Millisecond)

	// check if the server is closed
	select {
	case <-stopSignal:
		break
	default:
		t.Error("server is not closed")
		return
	}
}
