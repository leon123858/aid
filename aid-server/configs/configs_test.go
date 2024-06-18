package configs

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfigs(t *testing.T) {
	assert.NotNil(t, Configs)
}
