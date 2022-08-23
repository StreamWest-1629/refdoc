FROM golang:1.18-bullseye as builder

WORKDIR /build
COPY go.* ./
RUN go mod download
COPY *.go ./
RUN go build -o /bin/refdoc .
