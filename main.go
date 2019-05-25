package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/acme/autocert"
)

var opts struct {
	Domain  string `short:"d" long:"domain" description:"Domain for which to request certs from Let's Encrypt" required:"true"`
	To      string `short:"t" long:"to" description:"http[s]://IP:port to which traffic will be redirected" required:"true"`
	Verbose bool   `short:"v" long:"verbose" description:"Log request data"`
	Trace   bool   `long:"trace" description:"Log request data"`
}

func init() {
	_, err := flags.Parse(&opts)
	if err != nil {
		if err.(*flags.Error).Type == flags.ErrHelp {
			os.Exit(0)
		}
		os.Exit(1)
	}

	if opts.Verbose {
		log.SetLevel(log.DebugLevel)
	} else if opts.Trace {
		log.SetLevel(log.TraceLevel)
	}
}

func main() {
	url, err := url.Parse(opts.To)
	if err != nil {
		log.Fatal(err)
	}

	manager := &autocert.Manager{
		Cache:      autocert.DirCache("certs"),
		HostPolicy: autocert.HostWhitelist(opts.Domain),
		Prompt:     autocert.AcceptTOS,
	}

	s := &http.Server{
		Addr:      "0.0.0.0:443",
		TLSConfig: manager.TLSConfig(),
		Handler:   NewSingleHostReverseProxy(url),
	}

	log.WithFields(log.Fields{
		"to":     opts.To,
		"domain": opts.Domain,
	}).Info("Forwarding")

	log.Fatal(s.ListenAndServeTLS("", ""))
}

// NewSingleHostReverseProxy is yanked from httputils so additional logging
// could be added
func NewSingleHostReverseProxy(target *url.URL) *httputil.ReverseProxy {
	targetQuery := target.RawQuery
	director := func(req *http.Request) {
		log.WithFields(log.Fields{
			"src": req.RemoteAddr,
			"dst": req.RequestURI,
		}).Debug()

		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}

		log.WithFields(log.Fields{
			"headers": req.Header,
			"method":  req.Method,
			"length":  req.ContentLength,
		}).Trace()
	}
	return &httputil.ReverseProxy{Director: director}
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
