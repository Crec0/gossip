# gossip or go simply serve IP

A small dockerized go application that returns IPv4 of the requester.

You can run it as a standalone application or in a docker container.

## Standalone Application

#### Prerequisites

- Go

#### Steps

- Clone the repo 
- run `go build`
- It should produce an executable in the root directory of the project
- Simply run the executable in your terminal/screen/tmux instance
- Make sure to open ports or reverse proxy, if you want to publicly expose the application ~~Otherwise it's just gonna give you local ip xD~~

## Docker

#### Prerequisite

- Docker
- Docker Compose

#### Steps

- `docker compose up`