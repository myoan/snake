FROM golang:1.17.2 as builder
WORKDIR /go/src/snake

COPY api /go/src/snake/api
COPY server /go/src/snake/server
COPY go.mod .
COPY go.sum .

RUN cd server; GCO_ENABLED=0 GOOS=linux go build -o gameserver .
EXPOSE 8080
CMD ["/go/src/snake/server/gameserver"]
