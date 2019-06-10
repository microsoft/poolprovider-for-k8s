FROM golang:alpine

MAINTAINER Divyansh Manchanda <divyanshm@gmail.com>

RUN apk add --no-cache git mercurial \
    && go get github.com/garyburd/redigo/redis \
    && go get github.com/gorilla/handlers \
    && apk del git mercurial

RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go build -o main .

EXPOSE 8082
CMD ["/app/main"]
