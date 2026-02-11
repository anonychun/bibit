FROM docker.io/library/golang:1.26.0-trixie AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go install github.com/mattn/goreman@latest

RUN go build -o tmp/db cmd/db/main.go
RUN go build -o tmp/server cmd/server/main.go
RUN go build -o tmp/worker cmd/worker/main.go

FROM docker.io/library/debian:trixie-slim
WORKDIR /app

RUN apt-get update -qq && \
	apt-get install --no-install-recommends -y curl wget telnet htop vim tmux tini postgresql-client

COPY --from=build /go/bin/goreman /usr/bin/goreman

COPY --from=build /app/bin/docker-entrypoint bin/docker-entrypoint
COPY --from=build /app/Procfile Procfile
COPY --from=build /app/tmp/db bin/db
COPY --from=build /app/tmp/server bin/server
COPY --from=build /app/tmp/worker bin/worker

ENTRYPOINT ["/usr/bin/tini", "--", "/app/bin/docker-entrypoint"]
CMD ["/usr/bin/goreman", "start"]
