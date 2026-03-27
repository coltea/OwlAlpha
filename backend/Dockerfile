FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN go build -o owlalpha-server .

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/owlalpha-server /app/owlalpha-server
COPY --from=builder /app/manifest/config/config.yaml /app/config.yaml
COPY --from=builder /app/resource /app/resource

EXPOSE 8000

CMD ["/app/owlalpha-server"]
