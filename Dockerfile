From golang:1.22.1-alpine3.19

WORKDIR /go/src/app

COPY . .

RUN go build -o main main.go

EXPOSE 8080

CMD ["./main"]
