package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	_ "net/http/pprof"
	"net/url"

	np "github.com/ksang/naughty/proxy"
	"github.com/ksang/naughty/speakhttp"
)

var (
	backend  string
	bindaddr string
	withbody bool
)

func init() {
	flag.StringVar(&backend, "b", "", "backend server url, e.g http://test.com")
	flag.StringVar(&bindaddr, "a", "127.0.0.1:8080", "local addr to bind")
	flag.BoolVar(&withbody, "body", false, "print body content")
}

func main() {
	flag.Parse()
	if !checkArg() {
		flag.PrintDefaults()
		return
	}
	go func() {
		log.Println(http.ListenAndServe("127.0.0.1:6060", nil))
	}()
	log.Fatal(startProxy())
}

func startProxy() error {
	tarUrl, err := url.Parse(backend)
	if err != nil {
		return err
	}
	rp := httputil.NewSingleHostReverseProxy(tarUrl)

	var speaker speakhttp.Render
	if withbody {
		speaker = speakhttp.NewRender(
			speakhttp.ColorRequestWithBody,
			speakhttp.ColorResponseWithBody)
	} else {
		speaker = speakhttp.NewRender(nil, nil)
	}

	log.Printf("Naughty starting at: %s, backend: %s", bindaddr, backend)
	return http.ListenAndServe(bindaddr, &np.NaughtyProxy{
		rp,
		tarUrl,
		speaker,
		nil,
	})
}

func checkArg() bool {
	if len(backend) == 0 {
		return false
	}
	if _, err := url.Parse(backend); err != nil {
		return false
	}
	return true
}
