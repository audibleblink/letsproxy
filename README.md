# LetsProxy

Dead simple reverse proxy that automatically pulls TLS certificates from Let's Encrypt.

## Usage

```sh
./letsproxy --domain example.com --to http://localhost:3000
./letsproxy -d example.com -t http://127.0.0.1:9090

# multiple hosts allowed with csv
./letsproxy --domain corp.com,dev.com,example.com -t http://localhost:8080
```

## Build

**WebSocket proxying requires building with Go v1.12+**

This utility uses Go1.11+ built-in dependency system `go mod`. No need to be in $GOPATH, just clone the repo
and `go build`

