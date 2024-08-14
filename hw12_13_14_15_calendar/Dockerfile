FROM golang:1.22 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /app/cmd ./app/service/cmd

FROM debian:stable-slim

WORKDIR /root/

COPY --from=builder /app/cmd .

RUN mkdir -p /etc/calendar

COPY config.yaml /etc/calendar/config.yaml

COPY migrations ./migrations

CMD ["sh", "-c", "./cmd & wait"]
