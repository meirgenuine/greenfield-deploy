FROM golang:1.18

WORKDIR /go/src/github.com/din-mukhammed/greenfield-deploy/bot
COPY . .
ENV CGO_ENABLED 0
ENV GO111MODULE on
RUN go build .



FROM alpine:latest

RUN apk --no-cache add ca-certificates && apk add tzdata
WORKDIR /root/
COPY --from=0 /go/src/github.com/din-mukhammed/greenfield-deploy/bot/bot ./
