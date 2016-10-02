package proxy

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"testing"

	"github.com/ksang/naughty/speakhttp"
)

func TestNaughtyProxy(t *testing.T) {
	// start fake backend server
	backendServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "from backend server.")
	}))
	defer backendServer.Close()

	// init naughty proxy
	tarUrl, err := url.Parse(backendServer.URL)
	if err != nil {
		t.Errorf("Failed to parse backend server URL: %s", err)
		return
	}
	rp := httputil.NewSingleHostReverseProxy(tarUrl)
	speaker := speakhttp.NewRender(
		speakhttp.ColorRequestWithBody,
		speakhttp.ColorResponseWithBody)

	np := &NaughtyProxy{rp, tarUrl, speaker, nil}

	npServer := httptest.NewServer(np)
	defer npServer.Close()

	resp, err := http.Get(npServer.URL)
	if err != nil {
		t.Errorf("Failed to send request to NaughtyProxy, err: %s", err)
	}
	t.Logf("Response: %#v", resp)

}
