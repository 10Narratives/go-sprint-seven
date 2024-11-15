package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMainHandlerWhenOK(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/cafe?count=2&city=moscow", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, request)

	assert.Equal(t, http.StatusOK, responseRecorder.Code)
	assert.NotEqual(t, 0, responseRecorder.Body.Len())
}

func TestMainHandleWhenCountMissing(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/cafe?city=moscow", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, request)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)

	expectedBody := `count missing`
	assert.Equal(t, expectedBody, responseRecorder.Body.String())
}

func TestMainHandleWhenWrongCount(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/cafe?count=count&city=moscow", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, request)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)

	expectedBody := `wrong count value`
	assert.Equal(t, expectedBody, responseRecorder.Body.String())
}

func TestMainHandleWhenNotSupportedCity(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/cafe?count=2&city=london", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, request)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)

	expectedBody := `wrong city value`
	assert.Equal(t, expectedBody, responseRecorder.Body.String())
}

func TestMainHandleWhenCountOverflow(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/cafe?count=10&city=moscow", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, request)

	assert.Equal(t, http.StatusOK, responseRecorder.Code)

	totalCount := 4
	slicedBody := strings.Split(responseRecorder.Body.String(), ",")
	assert.Len(t, slicedBody, totalCount)
}
