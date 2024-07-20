# Build stage
FROM golang:1.23rc2-alpine3.20 AS Builder
WORKDIR /build/app
COPY . .
RUN go build -o main main.go

# Run stage
FROM alpine:latest
WORKDIR /run/app
COPY --from=builder /build/app/app.env .
COPY --from=builder /build/app/main .

EXPOSE 8080
CMD ["./main"]