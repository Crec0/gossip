FROM golang:latest

WORKDIR /app
COPY . /app

RUN go build && ./ip-pls
