FROM golang as builder
COPY . /app
WORKDIR /app
ENV GO111MODULE=on
RUN go mod vendor && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build cmd/go-open-registry.go

FROM golang:buster-slim
COPY --from=builder /etc/ssl/certs /etc/ssl/certs
COPY --from=builder /app/go-open-registry /go-open-registry
ENTRYPOINT ["/go-open-registry"]