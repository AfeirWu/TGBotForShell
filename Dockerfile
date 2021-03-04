FROM golang:1.13
COPY main.go .
RUN go get -u "github.com/go-telegram-bot-api/telegram-bot-api" \
    && CGO_ENABLED=0 go build -ldflags "-s -w" ./main.go

FROM scratch
COPY --from=0 /go/main .
COPY --from=0 /bin/sh /bin/sh
COPY --from=0 /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["./main"]