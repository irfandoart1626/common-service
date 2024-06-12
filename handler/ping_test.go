package handler

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetPingHandler(t *testing.T) {
	method, path, _ := GetPingHandler()
	assert.Equal(t, "GET", method)
	assert.NotEmpty(t, path)
}
