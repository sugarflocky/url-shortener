FROM golang:1.26-alpine AS builder

WORKDIR /build

COPY go.* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /shortener ./cmd/shortener

FROM alpine:3.21

COPY --from=builder /shortener /shortener

EXPOSE 8080
ENTRYPOINT ["/shortener"]
CMD ["-storage=memory"]