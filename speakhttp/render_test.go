package speakhttp

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRender(t *testing.T) {
	backendServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "from backend server.")
	}))
	defer backendServer.Close()

	tests := []struct {
		speakReq func(*http.Request)
		speakRes func(*http.Response)
	}{
		{
			speakReq: nil,
			speakRes: nil,
		},
		{
			speakReq: ColorRequestWithBody,
			speakRes: ColorResponseWithBody,
		},
		{
			speakReq: CurlRequest,
			speakRes: CurlResponse,
		},
		{
			speakReq: CurlRequestWithBody,
			speakRes: CurlResponseWithBody,
		},
	}
	for _, te := range tests {
		r := NewRender(te.speakReq, te.speakRes)
		req, err := http.NewRequest("POST", "http://www.test.com", nil)
		if err != nil {
			t.Errorf("Failed to gen request: %s", err)
			return
		}
		req.Header.Set("Test-Key", "Test-Value")
		req.SetBasicAuth("Username", "Password")
		r.SpeakRequest(req)
		resp, err := http.Get(backendServer.URL)
		if err != nil {
			t.Errorf("failed connect, %s", err)
		}
		r.SpeakResponse(resp)
	}
}
