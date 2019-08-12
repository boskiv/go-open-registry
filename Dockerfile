FROM golang as builder
COPY . /app
WORKDIR /app
ENV GO111MODULE=on
RUN go mod vendor && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build cmd/go-open-registry.go

FROM busybox
COPY --from=builder /app/go-open-registry /go-open-registry
ENTRYPOINT ["/go-open-registry"]