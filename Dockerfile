FROM golang:alpine

ENV GO111MODULE=on

RUN apk add --no-cache git 

RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go build -o main .

COPY agentpods/* agentpods/

EXPOSE 8082
CMD ["/app/main"]
