from golang:1.13-alpine as builder

add . /go/src/github.com/chentanyi/message
run apk update && apk add git && \
    go get -v github.com/go-bindata/go-bindata/go-bindata && \
    cd /go/src/github.com/chentanyi/message && \
    go-bindata template/ && \
    CGO_ENABLED=0 go install

from alpine:latest
workdir /app
copy --from=builder /go/bin/message .
entrypoint ["./message"]