package server_test

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alam0rt/tuna/server"
	"github.com/alam0rt/tuna/vtuner"
)

func setupTestServer() http.Handler {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	config := &server.Config{
		Host: "localhost",
		Port: "9191",
	}
	return server.NewServer(logger, config)
}

func TestSetupAppHandler(t *testing.T) {

	type test struct {
		name    string
		code    int
		body    string
		request *http.Request
	}

	tests := []test{
		{
			name:    "gets a token if none is provided",
			request: httptest.NewRequest("GET", "/setupapp/", nil),
			code:    200,
			body:    string(vtuner.EncryptedToken),
		},
		{
			name: "does not get a token if one is provided as an argument",
			request: func() *http.Request {
				httpReq := httptest.NewRequest("GET", "/setupapp/loginXML.asp", nil)
				q := httpReq.URL.Query()
				q.Add("token", "1234")
				httpReq.URL.RawQuery = q.Encode()
				return httpReq
			}(),
			code: 200,
			body: "",
		},
		{
			name: "is not found",
			request: func() *http.Request {
				httpReq := httptest.NewRequest("GET", "/setupapp/missing.asp", nil)
				q := httpReq.URL.Query()
				q.Add("token", "1234")
				httpReq.URL.RawQuery = q.Encode()
				return httpReq
			}(),
			code: 404,
			body: "404 page not found\n",
		},
	}

	srv := setupTestServer()
	for _, tc := range tests {
		rr := httptest.NewRecorder()
		srv.ServeHTTP(rr, tc.request)
		if rr.Code != tc.code {
			t.Fatalf("%s: expected status code %d, got %d", tc.name, tc.code, rr.Code)
		}
		if rr.Body.String() != tc.body {
			t.Fatalf("%s: expected body to be %s, got %s", tc.name, tc.body, rr.Body.String())
		}
	}
}
