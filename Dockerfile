FROM golang:latest

WORKDIR /app
COPY . /app

RUN go build

CMD ./ip-pls
