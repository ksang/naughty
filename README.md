# naughty

[![Build Status](https://travis-ci.org/ksang/naughty.svg?branch=master)](https://travis-ci.org/ksang/naughty)[![Go Report Card](https://goreportcard.com/badge/github.com/ksang/naughty)](https://goreportcard.com/report/github.com/ksang/naughty)

Naughty is a tool acting as a man-in-the-middle reverse proxy server pretty printing http request/response for debug purposes

![screenshot](./screenshot.png)

###Build:

	make build

###Usage:

	  -a string
	    	local addr to bind (default "127.0.0.1:8080")
	  -b string
	    	backend server url, e.g http://test.com
	  -body
	    	print body content