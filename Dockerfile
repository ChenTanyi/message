from golang:1.15-alpine as builder

add . /go/src/github.com/chentanyi/message
run apk update && apk add git && \
    cd /go/src/github.com/chentanyi/message && \
    CGO_ENABLED=0 go install

from alpine:latest
workdir /
copy --from=builder /go/bin/message /usr/bin
entrypoint ["message"]