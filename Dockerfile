FROM golang:1.13 AS builder
WORKDIR /
COPY ./trsh.go /
RUN go get -u "gopkg.in/telegram-bot-api.v4" \
    && CGO_ENABLED=0 go build ./trsh.go

FROM ubuntu AS production
COPY --from=builder /trsh /
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/trsh"]