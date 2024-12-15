FROM golang:1.22.3-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy && go mod download

COPY . .

WORKDIR /app/cmd/init
RUN go build -o initDb .

WORKDIR /app/cmd/server
RUN go build -o main .

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/cmd/init/initDb .
COPY --from=builder /app/cmd/server/main .
COPY .env.production ./
COPY templates ./templates

EXPOSE 8080

CMD ["sh", "-c", "./initDb && ./main"]
