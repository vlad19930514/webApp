# Build stage
FROM golang:1.22-alpine3.19 AS builder
WORKDIR /app
COPY . .
RUN go build -o main ./cmd/web/main.go

# Run stage
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .
COPY start.sh .
COPY wait-for.sh .
COPY db/migration ./db/migration

RUN chmod +x wait-for.sh start.sh

EXPOSE 8081 8081
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]
