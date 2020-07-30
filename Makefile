APP = letsproxy
DIR = release

FLAGS=-trimpath -ldflags="-s -w -buildid="

PLATFORMS ?= linux windows darwin
os = $(word 1, $@)


all: ${PLATFORMS}

${PLATFORMS}:
	GOOS=${os} \
	     go build ${FLAGS} -o ${DIR}/${APP}.${os}

clean: 
	rm -rf ${DIR}

