FROM golang:1.20.5 AS builder
WORKDIR /go/src/github.com/dedalusj/ctar/
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o ctar

FROM alpine:latest AS runner
WORKDIR /root/
COPY  --from=builder /go/src/github.com/dedalusj/ctar/ctar ./
CMD ["./ctar"]