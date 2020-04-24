# LetsProxy

No frills, no config, reverse proxy that automatically pulls TLS certificates from Let's Encrypt.

## Usage

```sh
./letsproxy --domain example.com --to http://localhost:3000
./letsproxy -d example.com -t http://127.0.0.1:9090

# multiple hosts allowed with csv
./letsproxy --domain corp.com,dev.com,example.com -t http://localhost:8080
```

## Build

**WebSocket proxying requires building with Go v1.12+**

```
# required for concurrent cross-compilation build process and Makefile usage
go install github.com/audibleblink/gox

# OSARCH ?= "linux/amd64 linux/386 linux/arm windows/amd64 windows/386 darwin/amd64 darwin/386"
make
```

To selectively build a binary os/arch combination, simply override the `OSARCH` variable at
make-time

```
make OSARCH=windows/amd64
```

