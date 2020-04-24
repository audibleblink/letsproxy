.PHONY: all clean

DIR = bin
OUT = ${DIR}/{{.OS}}_{{.Arch}}
GOFLAGS=-trimpath -ldflags="-s -w -buildid="
OSARCH ?= "linux/amd64 linux/386 linux/arm windows/amd64 windows/386 darwin/amd64 darwin/386"

all:
	gox -osarch ${OSARCH} ${GOFLAGS} -output ${OUT}

clean: 
	rm -rf ${DIR}

