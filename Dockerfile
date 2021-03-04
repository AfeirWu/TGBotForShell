FROM golang:1.13
COPY trsh.go .
RUN go get -u "gopkg.in/telegram-bot-api.v4" \
    && CGO_ENABLED=0 go build -ldflags "-s -w" ./trsh.go

FROM scratch
COPY --from=0 /go/trsh .
COPY --from=0 /bin/sh /bin/sh
COPY --from=0 /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["./trsh"]