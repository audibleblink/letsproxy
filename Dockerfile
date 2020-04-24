FROM golang:1.14-alpine as builder
MAINTAINER github.com/audibleblink <@4lex>
WORKDIR /app

RUN apk update --no-cache && \
    apk add --no-cache \
    make

ENV CGO_ENABLED=0
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN make linux

FROM alpine:latest
EXPOSE 443
COPY --from=builder /app/bin/linux_amd64 /app/letsproxy

