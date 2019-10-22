FROM golang:alpine

ENV GO111MODULE=on

RUN apk add --no-cache git 

RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go build -o main .

COPY agentpods/* agentpods/

EXPOSE 8080
CMD ["/app/main"]
