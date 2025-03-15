FROM golang:alpine AS builder
WORKDIR /app

COPY . .
RUN go mod download
RUN go build -o bin/app .

FROM alpine
WORKDIR /app
COPY --from=builder /app/bin/app /app/app
CMD ["/app/app"]
