FROM golang:1.18-bullseye as builder

WORKDIR /build
COPY go.* ./
RUN go mod download
COPY *.go ./
RUN go build -o /build/app .

FROM scratch:latest
COPY --from=builder /build/app /bin/
