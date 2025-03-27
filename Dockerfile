FROM golang:1.21 as builder
WORKDIR /app
COPY . .
RUN go build -o /bot ./cmd/bot

FROM alpine:latest
COPY --from=builder /bot /bot
CMD ["/bot"]