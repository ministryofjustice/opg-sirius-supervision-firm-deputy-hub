package main

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHelloWorld(t *testing.T) {
	assert := assert.New(t)
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/", nil)

	HelloServer(w, r)
	actual := w.Body.String()
	expected := "Hello world!"

	assert.Contains(actual, expected)
}