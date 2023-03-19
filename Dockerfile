FROM golang:1.20-alpine AS build

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 go build -o gossip

FROM alpine

COPY --from=build /app/gossip /bin/gossip

ENTRYPOINT ["/bin/gossip"]
