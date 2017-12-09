FROM alpine

ADD main.go  /go/

RUN apk add --no-cache go git libc-dev; \
    cd /go; \
    go get -v; \
    go build main.go; \
    mv main /; \
    apk del go git libc-dev;

CMD ["/main"]

