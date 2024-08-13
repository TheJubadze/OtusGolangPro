FROM golang:1.22 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /app/cmd ./app/service/cmd

RUN go build -o /app/scheduler ./app/service/scheduler

RUN go build -o /app/sender ./app/service/sender

FROM debian:stable-slim

WORKDIR /root/

COPY --from=builder /app/cmd .
COPY --from=builder /app/scheduler .
COPY --from=builder /app/sender .

RUN mkdir -p /etc/calendar

COPY config.yaml /etc/calendar/config.yaml

COPY migrations ./migrations

COPY tests ./tests

CMD ["sh", "-c", "./cmd & ./scheduler & ./sender & wait"]
