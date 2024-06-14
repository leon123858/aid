package cmd

import (
	"github.com/stretchr/testify/assert"
	"net"
	"net/http"
	"os"
	"testing"
)

func TestExecute(t *testing.T) {
	os.Args = []string{"app"}
	err := Execute()
	assert.NoError(t, err)

	os.Args = []string{"app", "not-a-command"}
	err = Execute()
	assert.Error(t, err)
}

func TestServe(t *testing.T) {
	testServeFunc := func(ln net.Listener) error {
		if ln == nil {
			panic("test panic")
		}
		return http.ErrServerClosed
	}
	//should panic
	assert.Panics(t, func() {
		err := serve(nil, testServeFunc)
		assert.NoError(t, err)
	})
	// should not panic
	err := serve(&mockListener{}, testServeFunc)
	assert.NoError(t, err)

	// should return error
	err = serve(&mockListener{}, func(ln net.Listener) error {
		return assert.AnError
	})
	assert.Error(t, err)
}

// mockListener is a mock implementation of net.Listener for testing
type mockListener struct {
	err error
}

func (l *mockListener) Accept() (net.Conn, error) {
	return nil, nil
}

func (l *mockListener) Close() error {
	return nil
}

func (l *mockListener) Addr() net.Addr {
	return &net.TCPAddr{}
}
