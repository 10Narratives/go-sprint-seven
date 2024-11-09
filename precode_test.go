package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestMainHandlerWhenOK(t *testing.T) {
	u, err := url.Parse("/cafe")
	if err != nil {
		t.Fatal(err)
	}

	query := u.Query()
	query.Add("city", "moscow")
	query.Add("count", "2")
	u.RawQuery = query.Encode()

	request := httptest.NewRequest(http.MethodGet, u.String(), nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("expected status code: %d, got %d", http.StatusOK, status)
	}

	expectedBody := ``
	if responseRecorder.Body.String() == expectedBody {
		t.Errorf("expected not empty response body")
	}
}

func TestMainHandleWhenCountMissing(t *testing.T) {
	u, err := url.Parse("/cafe")
	if err != nil {
		t.Fatal(err)
	}

	query := u.Query()
	query.Add("city", "moscow")
	u.RawQuery = query.Encode()

	request := httptest.NewRequest(http.MethodGet, u.String(), nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusBadRequest {
		t.Errorf("expected status code: %d, got %d", http.StatusBadRequest, status)
	}

	expectedBody := `count missing`
	if gottenBody := responseRecorder.Body.String(); gottenBody != expectedBody {
		t.Errorf("expected body: %s, got %s", expectedBody, gottenBody)
	}
}

func TestMainHandleWhenWrongCount(t *testing.T) {
	u, err := url.Parse("/cafe")
	if err != nil {
		t.Fatal(err)
	}

	query := u.Query()
	query.Add("city", "moscow")
	query.Add("count", "count")
	u.RawQuery = query.Encode()

	request := httptest.NewRequest(http.MethodGet, u.String(), nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusBadRequest {
		t.Errorf("expected status code: %d, got %d", http.StatusBadRequest, status)
	}

	expectedBody := `wrong count value`
	if gottenBody := responseRecorder.Body.String(); gottenBody != expectedBody {
		t.Errorf("expected body: %s, got %s", expectedBody, gottenBody)
	}
}

func TestMainHandleWhenNotSupportedCity(t *testing.T) {
	u, err := url.Parse("/cafe")
	if err != nil {
		t.Fatal(err)
	}

	query := u.Query()
	query.Add("city", "london")
	query.Add("count", "2")
	u.RawQuery = query.Encode()

	request := httptest.NewRequest(http.MethodGet, u.String(), nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusBadRequest {
		t.Errorf("expected status code: %d, got %d", http.StatusBadRequest, status)
	}

	expectedBody := `wrong city value`
	if gottenBody := responseRecorder.Body.String(); gottenBody != expectedBody {
		t.Errorf("expected body: %s, got %s", expectedBody, gottenBody)
	}
}

func TestMainHandleWhenCountOverflow(t *testing.T) {
	u, err := url.Parse("/cafe")
	if err != nil {
		t.Fatal(err)
	}

	query := u.Query()
	query.Add("city", "moscow")
	query.Add("count", "5")
	u.RawQuery = query.Encode()

	request := httptest.NewRequest(http.MethodGet, u.String(), nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("expected status code: %d, got %d", http.StatusOK, status)
	}

	totalCount := 4
	slicedBody := strings.Split(responseRecorder.Body.String(), ",")
	if totalCount != len(slicedBody) {
		t.Errorf("expected count of cities: %d, got %d", totalCount, len(slicedBody))
	}

}
