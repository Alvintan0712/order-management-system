FROM golang:1.23.2 AS builder

WORKDIR /app
COPY ./common /app/common
COPY ./order /app/order

WORKDIR /app/order

RUN go mod tidy
RUN go build -v -o /usr/local/bin/order .

EXPOSE 8001

CMD ["order"]