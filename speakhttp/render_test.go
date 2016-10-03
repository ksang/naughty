package speakhttp

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
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

func TestOrder(t *testing.T) {
	tests := []struct {
		req map[string]string
		res map[string]string
		// request order + response order
		expected []string
	}{
		{
			req: map[string]string{
				"Client":        "Test_Client",
				"Authorization": "Test_Authorization",
				"Date":          "Test_Date",
			},
			res: map[string]string{
				"Set-Cookie": "Test_Cookie",
				"X-Reqid":    "Test_reqid",
				"Connection": "Test_Connection",
				"Status":     "Test-Status",
			},
			expected: []string{"Date", "Authorization", "Client",
				"Connection", "Status", "Set-Cookie", "X-Reqid"},
		},
	}
	for _, te := range tests {
		req, err := http.NewRequest("GET", "http://example.com", nil)
		if err != nil {
			t.Errorf("Failed to new request: %s", err)
		}
		res, err := http.NewRequest("GET", "http://example.com", nil)
		if err != nil {
			t.Errorf("Failed to new request: %s", err)
		}

		for k, v := range te.req {
			req.Header.Add(k, v)
		}
		for k, v := range te.res {
			res.Header.Add(k, v)
		}
		order := getHeaderOrder(req.Header, RequestWellKnownHeaders)
		order = append(order, getHeaderOrder(res.Header, ResponseWellKnownHeaders)...)
		t.Logf("Result order: %#v", order)
		if !reflect.DeepEqual(order, te.expected) {
			t.Errorf("Order incorrect, actual: %#v, expected: %#v", order, te.expected)
		}
	}
}
