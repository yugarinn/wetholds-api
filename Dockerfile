FROM golang:1.22

ARG GOPATH
ENV GOPATH=/go

RUN mkdir /wetholds-api
WORKDIR /wetholds-api

ADD go.mod ./go.mod
ADD go.sum ./go.sum
ADD . .

RUN go mod download && go mod verify

ENTRYPOINT ./setup.sh
