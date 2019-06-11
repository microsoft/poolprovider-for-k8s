FROM golang:alpine

MAINTAINER Divyansh Manchanda <divyanshm@gmail.com>

RUN apk add --no-cache git mercurial \
    && go get github.com/garyburd/redigo/redis \
    && go get github.com/gorilla/handlers \
    && go get k8s.io/client-go/kubernetes \
    && go get k8s.io/client-go/rest \
    && go get k8s.io/api/core/v1 \
    && apk del git mercurial

RUN apk --no-cache add curl

RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go build -o main .

EXPOSE 8082
CMD ["/app/main"]
