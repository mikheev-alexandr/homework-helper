FROM golang:1.22 AS builder

ENV GOPATH=/

WORKDIR /app

COPY ./ ./

RUN apt-get update
RUN apt-get -y install postgresql-client
RUN apt-get -y install dos2unix

RUN dos2unix /app/entrypoint.sh
RUN dos2unix /app/wait-for-postgres.sh

RUN chmod +x wait-for-postgres.sh
RUN chmod +x entrypoint.sh

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
RUN go mod download
RUN go build -o main ./cmd/main.go

EXPOSE 8000

ENTRYPOINT ["./entrypoint.sh"]