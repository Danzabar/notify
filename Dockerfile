FROM golang:1.8-alpine

MAINTAINER Dan Cox <danzabian@gmail.com>

LABEL name="notify"

RUN apk add --no-cache git build-base
RUN go get github.com/Masterminds/glide
WORKDIR /go/src/github.com/Danzabar/notify

ADD . . 
RUN glide install --skip-test -v
RUN go install

ENTRYPOINT /go/bin/notify -a -user=$AUTH_USER -pass=$AUTH_PASS -driver=$DATABASE_DRIVER -creds=$DATABASE_USER:$DATABASE_PASS@/$DATABASE?charsetutf8&parseTime=true&loc=Local

EXPOSE 8080
