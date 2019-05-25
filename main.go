package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"golang.org/x/crypto/acme/autocert"
)

var (
	domain string
	to     string
)

func init() {
	flag.StringVar(&domain, "domain", "", "Domain for which to request certs")
	flag.StringVar(&to, "to", "", "Downstream server. ex: http://localhost:8080")
	flag.Parse()

	if domain == "" || to == "" {
		log.Fatal("Both -domain and -to are required")
	}

}

func main() {
	manager := &autocert.Manager{
		Cache:      autocert.DirCache("certs"),
		HostPolicy: autocert.HostWhitelist(domain),
		Prompt:     autocert.AcceptTOS,
	}

	s := &http.Server{
		Addr:      "0.0.0.0:443",
		TLSConfig: manager.TLSConfig(),
	}

	url, err := url.Parse(to)
	if err != nil {
		log.Fatal(err)
	}

	s.Handler = httputil.NewSingleHostReverseProxy(url)
	log.Println("Starting Reverse Proxy Server")
	log.Fatal(s.ListenAndServeTLS("", ""))
}
