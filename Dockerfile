FROM golang:alpine AS builder
MAINTAINER Temesxgn Gebrehiwet, temesxgn@gmail.com

ADD https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 /usr/bin/dep
RUN chmod +x /usr/bin/dep

RUN apk update && apk add --no-cache git ca-certificates
WORKDIR $GOPATH/src/github.com/temesxgn/redeam
RUN git clone -b master --single-branch https://github.com/temesxgn/redeam.git .
RUN dep ensure --vendor-only && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go test -v ./... && CGO_ENABLED=0 GOOS=linux go build -a -o /bin/api main.go

FROM scratch
COPY --from=builder /bin/api api
EXPOSE 8080
ENTRYPOINT ["./api"]