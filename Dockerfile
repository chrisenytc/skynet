FROM golang:1.8

MAINTAINER Christopher EnyTC <chris@enytc.com>

RUN mkdir -p /go/src/github.com/chrisenytc/skynet

WORKDIR /go/src/github.com/chrisenytc/skynet

COPY . $APP_HOME

RUN script/install
