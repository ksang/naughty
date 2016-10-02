/*
Package speakhttp provides methods of pretty printing http packets
*/
package speakhttp

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/fatih/color"
)

// Render defines the type of printing http request and response packets
type Render interface {
	// SpeakRequest is the method to printing http request
	SpeakRequest(*http.Request)
	// SpeakResponse is the method of printing http response
	SpeakResponse(*http.Response)
}

type render struct {
	speakReq func(*http.Request)
	speakRes func(*http.Response)
}

// SpeakRequest is the wrapper method to call the actual print request
func (r *render) SpeakRequest(req *http.Request) {
	r.speakReq(req)
}

// SpeakResponse is the wrapper method to call the actual print response
func (r *render) SpeakResponse(resp *http.Response) {
	r.speakRes(resp)
}

// NewRender is the factory function to return a Render interface
// You can specify your own print request and response functions
func NewRender(sReq func(*http.Request), sRes func(*http.Response)) Render {
	if sReq == nil {
		sReq = ColorRequest
	}
	if sRes == nil {
		sRes = ColorResponse
	}
	return &render{
		speakReq: sReq,
		speakRes: sRes,
	}
}

// CurlRequestWithBody is the print function for http response in curl style with body content
func CurlRequestWithBody(req *http.Request) {
	curlRequest(req, true)
}

// CurlRequest is the print function for http response in curl style without body
func CurlRequest(req *http.Request) {
	curlRequest(req, false)
}

func curlRequest(req *http.Request, withBody bool) {
	fmt.Printf("> %s %s %s\n", req.Method, req.URL.String(), req.Proto)
	// incoming req host is promoted to req.Host
	fmt.Printf("> Host: %s\n", req.Host)
	for k := range req.Header {
		fmt.Printf("> %s: %s\n", k, req.Header.Get(k))
	}
	if withBody {
		printRequestBody(req)
	}
}

// CurlResponseWithBody is the print function for http response in curl style with body content
func CurlResponseWithBody(res *http.Response) {
	curlResponse(res, true)
}

// CurlResponse is the print function for http response in curl style without body
func CurlResponse(res *http.Response) {
	curlResponse(res, false)
}

func curlResponse(res *http.Response, withBody bool) {
	fmt.Printf("< %s %s\n", res.Proto, res.Status)
	for k := range res.Header {
		fmt.Printf("< %s: %s\n", k, res.Header.Get(k))
	}
	if withBody {
		printResponseBody(res)
	}
}

// ColorRequestWithBody is the print function for http request in colored style with body
func ColorRequestWithBody(req *http.Request) {
	colorRequest(req, true)
}

// ColorRequest is the print function for http request in colored style without body
func ColorRequest(req *http.Request) {
	colorRequest(req, false)
}

func colorRequest(req *http.Request, withBody bool) {
	magenta := color.New(color.FgHiMagenta, color.Bold).SprintFunc()
	bold := color.New(color.Bold).SprintFunc()
	fmt.Printf("%s %s %s\n", magenta(req.Method), magenta(req.URL.String()), magenta(req.Proto))
	// incoming req host is promoted to req.Host
	fmt.Printf("%s %s\n", bold("Host:"), magenta(req.Host))
	for k := range req.Header {
		fmt.Printf("%s%s %s\n", bold(k), bold(":"), magenta(req.Header.Get(k)))
		if k == "Authorization" {
			u, p, ok := req.BasicAuth()
			if ok {
				fmt.Printf("    //BasicAuth decoded: %s:%s\n", u, p)
			}
		}
	}
	if withBody {
		printRequestBody(req)
	}
}

// ColorResponseWithBody is the print function for http response in colored style with body
func ColorResponseWithBody(res *http.Response) {
	colorResponse(res, true)
}

// ColorResponse is the print function for http response in colored style without body
func ColorResponse(res *http.Response) {
	colorResponse(res, false)
}

func colorResponse(res *http.Response, withBody bool) {
	green := color.New(color.FgHiGreen, color.Bold).SprintFunc()
	bold := color.New(color.Bold).SprintFunc()
	fmt.Printf("%s %s\n", green(res.Status), green(res.Proto))
	for k := range res.Header {
		fmt.Printf("%s%s %s\n", bold(k), bold(":"), green(res.Header.Get(k)))
	}
	if withBody {
		printResponseBody(res)
	}
}

func printRequestBody(req *http.Request) {
	if req.Body == nil {
		return
	}
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return
	}
	// restore request body content
	req.Body = ioutil.NopCloser(bytes.NewBuffer(b))
	fmt.Println(string(b))
}

func printResponseBody(res *http.Response) {
	if res.Body == nil {
		return
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	// restore request body content
	res.Body = ioutil.NopCloser(bytes.NewBuffer(b))
	fmt.Println(string(b))
}
