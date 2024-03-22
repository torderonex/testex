FROM golang:latest

WORKDIR /go/src/app

COPY . .

RUN apt-get update
RUN apt-get -y install postgresql-client

# make wait-for-postgres.sh executable
RUN chmod +x wait-for-postgres.sh


RUN go build -o main ./cmd/app/main.go

CMD ["./main"]
