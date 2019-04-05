FROM golang:alpine AS builder
MAINTAINER Temesxgn Gebrehiwet, temesxgn@gmail.com

RUN apk update && apk add --no-cache git ca-certificates
WORKDIR $GOPATH/src/github.com/temesxgn/redeam
RUN git clone -b master --single-branch https://github.com/temesxgn/redeam.git .
RUN go mod download
RUN go test -v ./...
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o bin/api

FROM scratch
COPY --from=builder bin/api api
EXPOSE 8080
ENTRYPOINT ["api"]