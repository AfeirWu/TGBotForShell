FROM golang:1.13
COPY trsh.go .
RUN CGO_ENABLED=0 go build ./trsh.go

FROM alpine
COPY --from=0 /go/trsh .
COPY --from=0 /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["./trsh"]