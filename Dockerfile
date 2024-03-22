FROM golang:latest

WORKDIR /go/src/app

COPY . .

RUN go build -o main ./cmd/app/main.go

CMD ["./main"]
