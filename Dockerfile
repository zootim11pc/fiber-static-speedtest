FROM golang:alpine AS builder

RUN apk --no-cache add ca-certificates

RUN apk add --update \
		git

ADD . /app
WORKDIR /app

RUN GOARCH=arm64 GOOS=linux CGO_ENABLED=0 go build -o static_speedtest fiber-static-speedtest

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /app/static_speedtest /

EXPOSE 3001

CMD ["/static_speedtest"]
