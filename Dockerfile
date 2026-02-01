FROM golang:1.24-alpine as builder

WORKDIR /app
COPY . .
RUN go build -o main main.go

FROM alpine:3.21
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .

EXPOSE 8080

CMD ["./main"]
