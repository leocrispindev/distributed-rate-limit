FROM golang:1.25.1-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o main cmd/main.go

RUN apk --no-cache add ca-certificates

#COPY ./app/main .

EXPOSE 8080

CMD ["./main"]