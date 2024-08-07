# Build go image
FROM golang:1.23rc1-alpine3.20 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run image stage
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .
COPY start.sh .
COPY wait-for.sh .
COPY db/migration ./db/migration

EXPOSE 8080
CMD ["./main"]
ENTRYPOINT ["./app/start.sh"]

