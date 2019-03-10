package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/acme/autocert"
)

var opts struct {
	Domain string `short:"d" long:"domain" description:"Domain for which to request certs from Let's Encrypt" required:"true"`
	To     string `short:"t" long:"to" description:"http[s]://IP:port to which traffic will be redirected" required:"true"`
}

func init() {
	_, err := flags.Parse(&opts)
	if err != nil {
		if err.(*flags.Error).Type == flags.ErrHelp {
			os.Exit(0)
		}
		os.Exit(1)
	}
}

func main() {
	mainLogger := log.WithFields(log.Fields{
		"to":     opts.To,
		"domain": opts.Domain,
	})

	url, err := url.Parse(opts.To)
	if err != nil {
		mainLogger.Fatal(err)
	}

	mainLogger.Info("Attempting to fetch TLS certificates from Let's Encrypt")
	manager := &autocert.Manager{
		Cache:      autocert.DirCache("certs"),
		HostPolicy: autocert.HostWhitelist(opts.Domain),
		Prompt:     autocert.AcceptTOS,
	}

	s := &http.Server{
		Addr:      "0.0.0.0:443",
		TLSConfig: manager.TLSConfig(),
	}

	s.Handler = httputil.NewSingleHostReverseProxy(url)
	mainLogger.Info("Starting Reverse Proxy Server")
	mainLogger.Fatal(s.ListenAndServeTLS("", ""))
}
