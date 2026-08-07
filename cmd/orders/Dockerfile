FROM golang:1.25-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /app/cmd/orders/exe app/cmd/orders

FROM alpine:3.23

WORKDIR /app

COPY --from=builder /app/cmd/orders/exe /app

CMD ["/app/exe"]