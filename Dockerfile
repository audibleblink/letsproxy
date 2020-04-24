FROM golang:1.14-alpine as builder
MAINTAINER github.com/audibleblink <@4lex>
WORKDIR /app

RUN apk update --no-cache && \
    apk add --no-cache \
    make

ENV CGO_ENABLED=0
COPY go.mod go.sum ./

# TODO: replace when `go mod` eventually also track `main` deps
# https://github.com/golang/go/issues/32504
RUN go mod download && go install github.com/audibleblink/gox 

COPY . .
RUN make OSARCH=linux/amd64

FROM alpine:latest
EXPOSE 443
COPY --from=builder /app/bin/linux_amd64 /app/letsproxy
